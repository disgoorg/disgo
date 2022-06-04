package voice

import (
	"context"
	"errors"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var ErrAlreadyConnected = errors.New("already connected")

type ConnectionCreateFunc func(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) *Connection

func NewConnection(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) *Connection {
	config := DefaultConnectionConfig()
	config.Apply(opts)

	return &Connection{
		config: *config,
		state: State{
			guildID:   guildID,
			userID:    userID,
			channelID: channelID,
		},
		connected: make(chan struct{}),
		ssrcs:     map[uint32]snowflake.ID{},
	}
}

type State struct {
	guildID snowflake.ID
	userID  snowflake.ID

	channelID snowflake.ID
	sessionID string
	token     string
	endpoint  string
}

type Connection struct {
	config ConnectionConfig

	state   State
	gateway *Gateway
	conn    *UDPConn
	mu      sync.Mutex

	connected    chan struct{}
	disconnected chan struct{}

	ssrcs   map[uint32]snowflake.ID
	ssrcsMu sync.Mutex
}

func (c *Connection) UserIDBySSRC(ssrc uint32) snowflake.ID {
	c.ssrcsMu.Lock()
	defer c.ssrcsMu.Unlock()
	return c.ssrcs[ssrc]
}

func (c *Connection) Gateway() *Gateway {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.gateway
}

func (c *Connection) Speaking(flags SpeakingFlags) error {
	return c.gateway.Send(GatewayOpcodeSpeaking, GatewayMessageDataSpeaking{
		SSRC:     c.Gateway().SSRC(),
		Speaking: flags,
	})
}

func (c *Connection) UDPConn() *UDPConn {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn
}

func (c *Connection) SetSendHandler(handler SendHandler) {
	NewSendSystem(handler, c).Start()
}

func (c *Connection) SetReceiveHandler(handler ReceiveHandler) {
	NewReceiveSystem(handler, c).Start()
}

func (c *Connection) HandleVoiceStateUpdate(update discord.VoiceStateUpdate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if update.GuildID != c.state.guildID || update.UserID != c.state.userID {
		return
	}

	if update.ChannelID == nil {
		c.state.channelID = 0
	} else {
		c.state.channelID = *update.ChannelID
	}
	c.state.sessionID = update.SessionID
	go c.reconnect(context.Background())
}

func (c *Connection) HandleVoiceServerUpdate(update discord.VoiceServerUpdate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if update.GuildID != c.state.guildID || update.Endpoint == nil {
		return
	}
	c.state.token = update.Token
	c.state.endpoint = *update.Endpoint
	go c.reconnect(context.Background())
}

func (c *Connection) handleGatewayMessage(_ GatewayOpcode, data GatewayMessageData) {
	switch d := data.(type) {
	case GatewayMessageDataReady:
		c.mu.Lock()
		println("voice: ready")
		conn := c.config.UDPConnCreateFunc(d.IP, d.Port, d.SSRC)
		c.conn = conn
		c.mu.Unlock()
		address, port, err := conn.Open(context.Background())
		if err != nil {
			c.config.Logger.Error("voice: failed to open udp connection. error: ", err)
			return
		}
		if err = c.Gateway().Send(GatewayOpcodeSelectProtocol, GatewayMessageDataSelectProtocol{
			Protocol: "udp",
			Data: GatewayMessageDataSelectProtocolData{
				Address: address,
				Port:    port,
				Mode:    EncryptionModeNormal,
			},
		}); err != nil {
			c.config.Logger.Error("voice: failed to send select protocol. error: ", err)
		}

	case GatewayMessageDataSessionDescription:
		c.mu.Lock()
		defer c.mu.Unlock()
		println("voice: session description")
		if c.conn != nil {
			c.conn.SetSecretKey(d.SecretKey)
		}
		close(c.connected)
		c.disconnected = make(chan struct{})

	case GatewayMessageDataSpeaking:
		c.ssrcsMu.Lock()
		defer c.ssrcsMu.Unlock()
		c.ssrcs[d.SSRC] = d.UserID

	case GatewayMessageDataClientDisconnect:
		c.ssrcsMu.Lock()
		defer c.ssrcsMu.Unlock()
		for ssrc, userID := range c.ssrcs {
			if userID == d.UserID {
				delete(c.ssrcs, ssrc)
				break
			}
		}
	}
}

func (c *Connection) handleGatewayClose(gateway *Gateway, err error) {
	c.config.Logger.Error("voice gateway closed. error: ", err)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gateway = nil
	close(c.disconnected)
}

func (c *Connection) Open(ctx context.Context) error {
	c.config.Logger.Debug("voice: opening connection")
	c.mu.Lock()
	if c.gateway != nil {
		c.mu.Unlock()
		return ErrAlreadyConnected
	}
	c.gateway = c.config.GatewayCreateFunc(c.state, c.handleGatewayMessage, c.handleGatewayClose, c.config.GatewayConfigOpts...)
	c.mu.Unlock()
	return c.gateway.Open(ctx)
}

func (c *Connection) reconnect(ctx context.Context) {
	c.mu.Lock()
	if c.state.endpoint == "" || c.state.token == "" {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	if err := c.Open(ctx); err != nil {
		c.config.Logger.Error("failed to reconnect to voice gateway. error: ", err)
	}
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.gateway == nil && c.conn == nil {
		return
	}

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	if c.gateway != nil {
		c.gateway.Close()
		c.gateway = nil
	}
}

func (c *Connection) WaitUntilConnected(ctx context.Context) error {
	select {
	case <-c.connected:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Connection) WaitUntilDisconnected(ctx context.Context) error {
	select {
	case <-c.disconnected:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
