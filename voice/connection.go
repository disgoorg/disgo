package voice

import (
	"context"
	"errors"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var ErrAlreadyConnected = errors.New("already connected")

type ConnectionCreateFunc func(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) Connection

// Connection is a complete voice connection to discord. It holds the voice Gateway and UDP connection and combines them.
type Connection interface {
	// Gateway returns the voice Gateway used by the voice Connection.
	Gateway() Gateway

	// UDP returns the UDP connection used by the voice Connection.
	UDP() UDP

	// ChannelID returns the ID of the voice channel the voice Connection is connected to.
	ChannelID() snowflake.ID

	// GuildID returns the ID of the guild the voice Connection is connected to.
	GuildID() snowflake.ID

	// UserIDBySSRC returns the ID of the user for the given SSRC.
	UserIDBySSRC(ssrc uint32) snowflake.ID

	// Speaking sends a speaking packet to the UDP socket discord.
	Speaking(flags SpeakingFlags) error

	// SetAudioSendSystem lets you inject your own AudioSendSystem. This is useful if you want to handle audio sending yourself.
	SetAudioSendSystem(sendSystem AudioSendSystem)

	// SetOpusFrameProvider lets you inject your own OpusFrameProvider.
	SetOpusFrameProvider(handler OpusFrameProvider)

	// SetAudioReceiveSystem lets you inject your own AudioReceiveSystem. This is useful if you want to handle audio receiving yourself.
	SetAudioReceiveSystem(receiveSystem AudioReceiveSystem)

	// SetOpusFrameReceiver lets you inject your own OpusFrameReceiver.
	SetOpusFrameReceiver(handler OpusFrameReceiver)

	// SetEventHandlerFunc lets listen for voice gateway events.
	SetEventHandlerFunc(eventHandlerFunc EventHandlerFunc)

	// Open opens the voice connection. It will connect to the voice gateway and start the UDP connection after it receives the Gateway events.
	Open(ctx context.Context) error

	// Close closes the voice connection. It will close the UDP connection and disconnect from the voice gateway.
	Close()

	// HandleVoiceStateUpdate provides the discord.VoiceStateUpdate to the voice connection. Which is needed to connect to the voice Gateway.
	HandleVoiceStateUpdate(update discord.VoiceStateUpdate)

	// HandleVoiceServerUpdate provides the discord.VoiceServerUpdate to the voice connection. Which is needed to connect to the voice Gateway.
	HandleVoiceServerUpdate(update discord.VoiceServerUpdate)

	// WaitUntilConnected blocks the current goroutine until the voice connection is connected. Make sure you call this method in its own goroutine, or it may block the gateway goroutine.
	WaitUntilConnected(ctx context.Context) error
	// WaitUntilDisconnected blocks the current goroutine until the voice connection is disconnected. Make sure you call this method in its own goroutine, or it may block the gateway goroutine.
	WaitUntilDisconnected(ctx context.Context) error
}

// NewConnection returns a new default voice connection.
func NewConnection(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) Connection {
	config := DefaultConnectionConfig()
	config.Apply(opts)

	return &connectionImpl{
		config: *config,
		state: state{
			guildID:   guildID,
			userID:    userID,
			channelID: channelID,
		},
		connected:    make(chan struct{}, 1),
		disconnected: make(chan struct{}, 1),
		ssrcs:        map[uint32]snowflake.ID{},
	}
}

type state struct {
	guildID snowflake.ID
	userID  snowflake.ID

	channelID snowflake.ID
	sessionID string
	token     string
	endpoint  string
}

type connectionImpl struct {
	config ConnectionConfig

	state   state
	gateway Gateway
	udp     UDP
	mu      sync.Mutex

	audioSendSystem    AudioSendSystem
	audioReceiveSystem AudioReceiveSystem

	connected    chan struct{}
	disconnected chan struct{}

	ssrcs   map[uint32]snowflake.ID
	ssrcsMu sync.Mutex
}

func (c *connectionImpl) ChannelID() snowflake.ID {
	return c.state.channelID
}

func (c *connectionImpl) GuildID() snowflake.ID {
	return c.state.guildID
}

func (c *connectionImpl) UserIDBySSRC(ssrc uint32) snowflake.ID {
	c.ssrcsMu.Lock()
	defer c.ssrcsMu.Unlock()
	return c.ssrcs[ssrc]
}

func (c *connectionImpl) Gateway() Gateway {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.gateway
}

func (c *connectionImpl) Speaking(flags SpeakingFlags) error {
	c.mu.Lock()
	gateway := c.gateway
	c.mu.Unlock()
	if gateway == nil {
		return ErrGatewayNotConnected
	}
	return gateway.Send(GatewayOpcodeSpeaking, GatewayMessageDataSpeaking{
		SSRC:     c.Gateway().SSRC(),
		Speaking: flags,
	})
}

func (c *connectionImpl) UDP() UDP {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.udp
}

func (c *connectionImpl) SetAudioSendSystem(sendSystem AudioSendSystem) {
	if c.audioSendSystem != nil {
		c.audioSendSystem.Close()
	}
	c.audioSendSystem = sendSystem
	if c.audioSendSystem != nil {
		go c.audioSendSystem.Open()
	}
}

func (c *connectionImpl) SetOpusFrameProvider(handler OpusFrameProvider) {
	c.SetAudioSendSystem(NewAudioSendSystem(c.config.Logger, handler, c))
}

func (c *connectionImpl) SetAudioReceiveSystem(receiveSystem AudioReceiveSystem) {
	if c.audioReceiveSystem != nil {
		c.audioReceiveSystem.Close()
	}
	c.audioReceiveSystem = receiveSystem
	if c.audioReceiveSystem != nil {
		go c.audioReceiveSystem.Open()
	}
}

func (c *connectionImpl) SetOpusFrameReceiver(handler OpusFrameReceiver) {
	c.SetAudioReceiveSystem(NewAudioReceiveSystem(c.config.Logger, handler, c))
}

func (c *connectionImpl) SetEventHandlerFunc(eventHandlerFunc EventHandlerFunc) {
	c.config.EventHandlerFunc = eventHandlerFunc
}

func (c *connectionImpl) HandleVoiceStateUpdate(update discord.VoiceStateUpdate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if update.GuildID != c.state.guildID || update.UserID != c.state.userID {
		return
	}

	if update.ChannelID == nil {
		c.state.channelID = 0
		c.disconnected <- struct{}{}
	} else {
		c.state.channelID = *update.ChannelID
	}
	c.state.sessionID = update.SessionID
}

func (c *connectionImpl) HandleVoiceServerUpdate(update discord.VoiceServerUpdate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if update.GuildID != c.state.guildID || update.Endpoint == nil {
		return
	}

	c.state.token = update.Token
	c.state.endpoint = *update.Endpoint
	go c.reconnect(context.Background())
}

func (c *connectionImpl) handleGatewayMessage(op GatewayOpcode, data GatewayMessageData) {
	switch d := data.(type) {
	case GatewayMessageDataReady:
		c.mu.Lock()
		c.udp = c.config.UDPConnCreateFunc(d.IP, d.Port, d.SSRC, append([]UDPConfigOpt{WithUDPLogger(c.config.Logger)}, c.config.UDPConnConfigOpts...)...)
		c.mu.Unlock()
		address, port, err := c.udp.Open(context.Background())
		if err != nil {
			c.config.Logger.Error("voice: failed to open udp connection. error: ", err)
			break
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
		if c.udp != nil {
			c.udp.SetSecretKey(d.SecretKey)
		}
		c.mu.Unlock()
		c.connected <- struct{}{}

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
		c.audioReceiveSystem.CleanupUser(d.UserID)
	}
	if c.config.EventHandlerFunc != nil {
		c.config.EventHandlerFunc(op, data)
	}
}

func (c *connectionImpl) handleGatewayClose(gateway Gateway, err error) {
	c.Close()
}

func (c *connectionImpl) Open(ctx context.Context) error {
	c.config.Logger.Debug("opening voice connection")
	c.mu.Lock()
	c.gateway = c.config.GatewayCreateFunc(c.state, c.handleGatewayMessage, c.handleGatewayClose, append([]GatewayConfigOpt{WithGatewayLogger(c.config.Logger)}, c.config.GatewayConfigOpts...)...)
	c.mu.Unlock()
	return c.gateway.Open(ctx)
}

func (c *connectionImpl) reconnect(ctx context.Context) {
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

func (c *connectionImpl) Close() {
	c.mu.Lock()
	if c.udp != nil {
		c.udp.Close()
	}

	if c.gateway != nil {
		c.gateway.Close()
	}
	c.mu.Unlock()
}

func (c *connectionImpl) WaitUntilConnected(ctx context.Context) error {
	select {
	case <-c.connected:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connectionImpl) WaitUntilDisconnected(ctx context.Context) error {
	select {
	case <-c.disconnected:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
