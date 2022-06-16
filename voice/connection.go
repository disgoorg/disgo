package voice

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var ErrAlreadyConnected = errors.New("already connected")

type ConnectionCreateFunc func(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) Connection

type Connection interface {
	Gateway() Gateway
	UDP() UDP

	ChannelID() snowflake.ID
	GuildID() snowflake.ID

	UserIDBySSRC(ssrc uint32) snowflake.ID
	Speaking(flags SpeakingFlags) error

	SetAudioSendSystem(sendSystem AudioSendSystem)
	SetOpusFrameProvider(handler OpusFrameProvider)
	SetAudioReceiveSystem(receiveSystem AudioReceiveSystem)
	SetOpusFrameReceiver(handler OpusFrameReceiver)
	SetEventHandlerFunc(eventHandlerFunc EventHandlerFunc)

	Open(ctx context.Context) error
	Close()

	HandleVoiceStateUpdate(update discord.VoiceStateUpdate)
	HandleVoiceServerUpdate(update discord.VoiceServerUpdate)

	WaitUntilConnected(ctx context.Context) error
	WaitUntilDisconnected(ctx context.Context) error
}

func NewConnection(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) Connection {
	config := DefaultConnectionConfig()
	config.Apply(opts)

	return &connectionImpl{
		config: *config,
		state: State{
			guildID:   guildID,
			userID:    userID,
			channelID: channelID,
		},
		connected:    make(chan struct{}),
		disconnected: make(chan struct{}),
		ssrcs:        map[uint32]snowflake.ID{},
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

type connectionImpl struct {
	config ConnectionConfig

	state   State
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
	return c.gateway.Send(GatewayOpcodeSpeaking, GatewayMessageDataSpeaking{
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
	c.SetAudioSendSystem(NewAudioSendSystem(handler, c))
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
	c.SetAudioReceiveSystem(NewAudioReceiveSystem(handler, c))
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
	if !errors.Is(err, net.ErrClosed) {
		c.config.Logger.Error("voice gateway closed. error: ", err)
	}
	gateway.Close()
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
		c.udp = nil
	}

	if c.gateway != nil {
		c.gateway.Close()
		c.gateway = nil
	}
	c.mu.Unlock()
	c.disconnected <- struct{}{}
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
