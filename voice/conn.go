package voice

import (
	"context"
	"sync"
	"time"

	botgateway "github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice/voicegateway"
	"github.com/disgoorg/disgo/voice/voiceudp"
	"github.com/disgoorg/snowflake/v2"
)

// ConnCreateFunc is a type alias for a function that creates a new Conn.
type ConnCreateFunc func(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn

// Conn is a complete voice conn to discord. It holds the voicegateway.Gateway and voiceudp.Conn conn and combines them.
type Conn interface {
	// Gateway returns the voice Gateway used by the voice Conn.
	Gateway() voicegateway.Gateway

	// Conn returns the voiceudp.Conn conn used by the voice Conn.
	Conn() voiceudp.Conn

	// ChannelID returns the ID of the voice channel the voice Conn is openedChan to.
	ChannelID() snowflake.ID

	// GuildID returns the ID of the guild the voice Conn is openedChan to.
	GuildID() snowflake.ID

	// UserIDBySSRC returns the ID of the user for the given SSRC.
	UserIDBySSRC(ssrc uint32) snowflake.ID

	// SetSpeaking sends a speaking packet to the Conn socket discord.
	SetSpeaking(ctx context.Context, flags voicegateway.SpeakingFlags) error

	// SetOpusFrameProvider lets you inject your own OpusFrameProvider.
	SetOpusFrameProvider(handler OpusFrameProvider)

	// SetOpusFrameReceiver lets you inject your own OpusFrameReceiver.
	SetOpusFrameReceiver(handler OpusFrameReceiver)

	// SetEventHandlerFunc lets listen for voice gateway events.
	SetEventHandlerFunc(eventHandlerFunc voicegateway.EventHandlerFunc)

	// Open opens the voice conn. It will connect to the voice gateway and start the Conn conn after it receives the Gateway events.
	Open(ctx context.Context, selfMute bool, selfDeaf bool) error

	// Close closes the voice conn. It will close the Conn conn and disconnect from the voice gateway.
	Close(ctx context.Context)

	// HandleVoiceStateUpdate provides the gateway.EventVoiceStateUpdate to the voice conn. Which is needed to connect to the voice Gateway.
	HandleVoiceStateUpdate(update botgateway.EventVoiceStateUpdate)

	// HandleVoiceServerUpdate provides the gateway.EventVoiceServerUpdate to the voice conn. Which is needed to connect to the voice Gateway.
	HandleVoiceServerUpdate(update botgateway.EventVoiceServerUpdate)
}

// NewConn returns a new default voice conn.
func NewConn(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn {
	config := DefaultConnConfig()
	config.Apply(opts)

	conn := &connImpl{
		config:               *config,
		voiceStateUpdateFunc: voiceStateUpdateFunc,
		removeConnFunc:       removeConnFunc,
		state: voicegateway.State{
			GuildID:   guildID,
			UserID:    userID,
			ChannelID: channelID,
		},
		openedChan: make(chan struct{}, 1),
		closedChan: make(chan struct{}, 1),
		ssrcs:      map[uint32]snowflake.ID{},
	}

	conn.gateway = config.GatewayCreateFunc(conn.handleMessage, conn.handleGatewayClose, append([]voicegateway.ConfigOpt{voicegateway.WithLogger(config.Logger)}, config.GatewayConfigOpts...)...)
	conn.udp = config.UDPConnCreateFunc(append([]voiceudp.ConnConfigOpt{voiceudp.WithLogger(config.Logger)}, config.UDPConnConfigOpts...)...)

	return conn
}

type connImpl struct {
	config               ConnConfig
	voiceStateUpdateFunc StateUpdateFunc
	removeConnFunc       func()

	state   voicegateway.State
	stateMu sync.Mutex

	selfMute bool
	selfDeaf bool

	gateway voicegateway.Gateway
	udp     voiceudp.Conn

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

func (c *connImpl) Gateway() voicegateway.Gateway {
	return c.gateway
}

func (c *connImpl) SetSpeaking(ctx context.Context, flags voicegateway.SpeakingFlags) error {
	return c.gateway.Send(ctx, voicegateway.OpcodeSpeaking, voicegateway.MessageDataSpeaking{
		SSRC:     c.Gateway().SSRC(),
		Speaking: flags,
	})
}

func (c *connImpl) Conn() voiceudp.Conn {
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

func (c *connImpl) SetEventHandlerFunc(eventHandlerFunc voicegateway.EventHandlerFunc) {
	c.config.EventHandlerFunc = eventHandlerFunc
}

func (c *connImpl) HandleVoiceStateUpdate(update botgateway.EventVoiceStateUpdate) {
	if update.GuildID != c.state.GuildID || update.UserID != c.state.UserID {
		return
	}

	if update.ChannelID == nil {
		c.state.ChannelID = 0
		if c.audioSender != nil {
			c.audioSender.Close()
			c.audioSender = nil
		}
		if c.audioReceiver != nil {
			c.audioReceiver.Close()
			c.audioReceiver = nil
		}
		c.udp.Close()
		c.gateway.Close()
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
	go func() {
		if err := c.gateway.Open(context.Background(), c.state); err != nil {
			c.config.Logger.Error("error opening voice gateway. error: ", err)
		}
	}()
}

func (c *connImpl) handleMessage(op voicegateway.Opcode, data voicegateway.MessageData) {
	switch d := data.(type) {
	case voicegateway.MessageDataReady:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ourAddress, ourPort, err := c.udp.Open(ctx, d.IP, d.Port, d.SSRC)
		if err != nil {
			c.config.Logger.Error("voice: failed to open voiceudp conn. error: ", err)
			break
		}
		if err = c.Gateway().Send(ctx, voicegateway.OpcodeSelectProtocol, voicegateway.MessageDataSelectProtocol{
			Protocol: voicegateway.VoiceProtocolUDP,
			Data: voicegateway.MessageDataSelectProtocolData{
				Address: ourAddress,
				Port:    ourPort,
				Mode:    voicegateway.EncryptionModeNormal,
			},
		}); err != nil {
			c.config.Logger.Error("voice: failed to send select protocol. error: ", err)
		}

	case voicegateway.MessageDataSessionDescription:
		c.udp.SetSecretKey(d.SecretKey)
		c.openedChan <- struct{}{}

	case voicegateway.MessageDataSpeaking:
		c.ssrcsMu.Lock()
		defer c.ssrcsMu.Unlock()
		c.ssrcs[d.SSRC] = d.UserID

	case voicegateway.MessageDataClientDisconnect:
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

func (c *connImpl) handleGatewayClose(gateway voicegateway.Gateway, err error) {
	c.Close(context.Background())
}

func (c *connImpl) Open(ctx context.Context, selfMute bool, selfDeaf bool) error {
	c.config.Logger.Debug("opening voice conn")

	c.selfMute = selfMute
	c.selfDeaf = selfDeaf
	c.stateMu.Lock()

	if err := c.voiceStateUpdateFunc(ctx, c.state.GuildID, &c.state.ChannelID, selfMute, selfDeaf); err != nil {
		c.stateMu.Unlock()
		return err
	}
	c.stateMu.Unlock()

	select {
	case <-c.openedChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connImpl) Close(ctx context.Context) {
	_ = c.voiceStateUpdateFunc(ctx, c.state.GuildID, nil, c.selfMute, c.selfDeaf)
	defer c.gateway.Close()
	defer c.udp.Close()

	select {
	case _, ok := <-c.closedChan:
		if ok {
			close(c.closedChan)
		}
	case <-ctx.Done():
	}
	c.removeConnFunc()
}
