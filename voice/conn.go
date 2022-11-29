package voice

import (
	"context"
	"sync"
	"time"

	botgateway "github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice/gateway"
	"github.com/disgoorg/disgo/voice/udp"
	"github.com/disgoorg/snowflake/v2"
)

// ConnCreateFunc is a type alias for a function that creates a new Conn.
type ConnCreateFunc func(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnConfigOpt) Conn

// Conn is a complete voice conn to discord. It holds the voice gateway.Gateway and udp.Conn conn and combines them.
type Conn interface {
	// Gateway returns the voice Gateway used by the voice Conn.
	Gateway() gateway.Gateway

	// Conn returns the udp.Conn conn used by the voice Conn.
	Conn() udp.Conn

	// ChannelID returns the ID of the voice channel the voice Conn is openedChan to.
	ChannelID() snowflake.ID

	// GuildID returns the ID of the guild the voice Conn is openedChan to.
	GuildID() snowflake.ID

	// UserIDBySSRC returns the ID of the user for the given SSRC.
	UserIDBySSRC(ssrc uint32) snowflake.ID

	// SetSpeaking sends a speaking packet to the Conn socket discord.
	SetSpeaking(ctx context.Context, flags gateway.SpeakingFlags) error

	// SetOpusFrameProvider lets you inject your own OpusFrameProvider.
	SetOpusFrameProvider(handler OpusFrameProvider)

	// SetOpusFrameReceiver lets you inject your own OpusFrameReceiver.
	SetOpusFrameReceiver(handler OpusFrameReceiver)

	// SetEventHandlerFunc lets listen for voice gateway events.
	SetEventHandlerFunc(eventHandlerFunc gateway.EventHandlerFunc)

	// Open opens the voice conn. It will connect to the voice gateway and start the Conn conn after it receives the Gateway events.
	Open(ctx context.Context) error

	// Close closes the voice conn. It will close the Conn conn and disconnect from the voice gateway.
	Close()

	// HandleVoiceStateUpdate provides the gateway.EventVoiceStateUpdate to the voice conn. Which is needed to connect to the voice Gateway.
	HandleVoiceStateUpdate(update botgateway.EventVoiceStateUpdate)

	// HandleVoiceServerUpdate provides the gateway.EventVoiceServerUpdate to the voice conn. Which is needed to connect to the voice Gateway.
	HandleVoiceServerUpdate(update botgateway.EventVoiceServerUpdate)

	// WaitUntilOpened blocks the current goroutine until the voice conn is openedChan. Make sure you call this method in its own goroutine, or it may block the gateway goroutine.
	WaitUntilOpened(ctx context.Context) error
	// WaitUntilClosed blocks the current goroutine until the voice conn is closedChan. Make sure you call this method in its own goroutine, or it may block the gateway goroutine.
	WaitUntilClosed(ctx context.Context) error
}

// NewConn returns a new default voice conn.
func NewConn(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, opts ...ConnConfigOpt) Conn {
	config := DefaultConnConfig()
	config.Apply(opts)

	conn := &connImpl{
		config: *config,
		state: gateway.State{
			GuildID:   guildID,
			UserID:    userID,
			ChannelID: channelID,
		},
		openedChan: make(chan struct{}, 1),
		closedChan: make(chan struct{}, 1),
		ssrcs:      map[uint32]snowflake.ID{},
	}

	conn.gateway = config.GatewayCreateFunc(conn.handleMessage, conn.handleGatewayClose, append([]gateway.ConfigOpt{gateway.WithLogger(config.Logger)}, config.GatewayConfigOpts...)...)
	conn.udp = config.UDPConnCreateFunc(append([]udp.ConnConfigOpt{udp.WithLogger(config.Logger)}, config.UDPConnConfigOpts...)...)

	return conn
}

type connImpl struct {
	config ConnConfig

	state   gateway.State
	stateMu sync.Mutex

	gateway gateway.Gateway
	udp     udp.Conn

	audioSender   AudioSender
	audioReceiver AudioReceiver

	openedChan chan struct{}
	closedChan chan struct{}

	ssrcs   map[uint32]snowflake.ID
	ssrcsMu sync.Mutex
}

func (c *connImpl) ChannelID() snowflake.ID {
	return c.state.ChannelID
}

func (c *connImpl) GuildID() snowflake.ID {
	return c.state.GuildID
}

func (c *connImpl) UserIDBySSRC(ssrc uint32) snowflake.ID {
	c.ssrcsMu.Lock()
	defer c.ssrcsMu.Unlock()
	return c.ssrcs[ssrc]
}

func (c *connImpl) Gateway() gateway.Gateway {
	return c.gateway
}

func (c *connImpl) SetSpeaking(ctx context.Context, flags gateway.SpeakingFlags) error {
	return c.gateway.Send(ctx, gateway.OpcodeSpeaking, gateway.MessageDataSpeaking{
		SSRC:     c.Gateway().SSRC(),
		Speaking: flags,
	})
}

func (c *connImpl) Conn() udp.Conn {
	return c.udp
}

func (c *connImpl) SetOpusFrameProvider(provider OpusFrameProvider) {
	if c.audioSender != nil {
		c.audioSender.Close()
	}
	c.audioSender = c.config.AudioSenderCreateFunc(c.config.Logger, provider, c)
	c.audioSender.Open()
}

func (c *connImpl) SetOpusFrameReceiver(handler OpusFrameReceiver) {
	if c.audioReceiver != nil {
		c.audioReceiver.Close()
	}
	c.audioReceiver = c.config.AudioReceiverCreateFunc(c.config.Logger, handler, c)
	c.audioReceiver.Open()
}

func (c *connImpl) SetEventHandlerFunc(eventHandlerFunc gateway.EventHandlerFunc) {
	c.config.EventHandlerFunc = eventHandlerFunc
}

func (c *connImpl) HandleVoiceStateUpdate(update botgateway.EventVoiceStateUpdate) {
	if update.GuildID != c.state.GuildID || update.UserID != c.state.UserID {
		return
	}

	if update.ChannelID == nil {
		c.state.ChannelID = 0
		c.closedChan <- struct{}{}
	} else {
		c.state.ChannelID = *update.ChannelID
	}
	c.state.SessionID = update.SessionID
}

func (c *connImpl) HandleVoiceServerUpdate(update botgateway.EventVoiceServerUpdate) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if update.GuildID != c.state.GuildID || update.Endpoint == nil {
		return
	}

	c.state.Token = update.Token
	c.state.Endpoint = *update.Endpoint
	go c.reconnect()
}

func (c *connImpl) handleMessage(op gateway.Opcode, data gateway.MessageData) {
	switch d := data.(type) {
	case gateway.MessageDataReady:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ourAddress, ourPort, err := c.udp.Open(ctx, d.IP, d.Port, d.SSRC)
		if err != nil {
			c.config.Logger.Error("voice: failed to open udp conn. error: ", err)
			break
		}
		if err = c.Gateway().Send(ctx, gateway.OpcodeSelectProtocol, gateway.MessageDataSelectProtocol{
			Protocol: gateway.VoiceProtocolUDP,
			Data: gateway.MessageDataSelectProtocolData{
				Address: ourAddress,
				Port:    ourPort,
				Mode:    gateway.EncryptionModeNormal,
			},
		}); err != nil {
			c.config.Logger.Error("voice: failed to send select protocol. error: ", err)
		}

	case gateway.MessageDataSessionDescription:
		c.udp.SetSecretKey(d.SecretKey)
		c.openedChan <- struct{}{}

	case gateway.MessageDataSpeaking:
		c.ssrcsMu.Lock()
		defer c.ssrcsMu.Unlock()
		c.ssrcs[d.SSRC] = d.UserID

	case gateway.MessageDataClientDisconnect:
		c.ssrcsMu.Lock()
		defer c.ssrcsMu.Unlock()
		for ssrc, userID := range c.ssrcs {
			if userID == d.UserID {
				delete(c.ssrcs, ssrc)
				break
			}
		}
		c.audioReceiver.CleanupUser(d.UserID)
	}
	if c.config.EventHandlerFunc != nil {
		c.config.EventHandlerFunc(op, data)
	}
}

func (c *connImpl) handleGatewayClose(gateway gateway.Gateway, err error) {
	c.Close()
}

func (c *connImpl) Open(ctx context.Context) error {
	c.config.Logger.Debug("opening voice conn")

	c.stateMu.Lock()
	state := c.state
	c.stateMu.Unlock()
	return c.gateway.Open(ctx, state)
}

func (c *connImpl) reconnect() {
	c.stateMu.Lock()
	if c.state.Endpoint == "" || c.state.Token == "" {
		c.stateMu.Unlock()
		return
	}
	c.stateMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.config.ReconnectTimeout)
	defer cancel()
	if err := c.Open(ctx); err != nil {
		c.config.Logger.Error("failed to reconnect to voice gateway. error: ", err)
	}
}

func (c *connImpl) Close() {
	_ = c.udp.Close()
	c.gateway.Close()
}

func (c *connImpl) WaitUntilOpened(ctx context.Context) error {
	select {
	case <-c.openedChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connImpl) WaitUntilClosed(ctx context.Context) error {
	select {
	case <-c.closedChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
