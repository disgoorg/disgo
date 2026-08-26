package voice

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/godave"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	botgateway "github.com/disgoorg/disgo/gateway"
)

type (
	// ConnCreateFunc is a type alias for a function that creates a new Conn.
	ConnCreateFunc func(guildID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn

	// Conn is a complete voice conn to discord. It holds the Gateway and voiceudp.UDPConn conn and combines them.
	Conn interface {
		// Gateway returns the voice Gateway used by the voice Conn.
		Gateway() Gateway

		// UDP returns the voice UDPConn conn used by the voice Conn.
		UDP() UDPConn

		DAVE() godave.Session

		// ChannelID returns the ID of the voice channel the voice Conn is openedChan to.
		ChannelID() *snowflake.ID

		// GuildID returns the ID of the guild the voice Conn is openedChan to.
		GuildID() snowflake.ID

		// UserIDBySSRC returns the ID of the user for the given SSRC.
		UserIDBySSRC(ssrc uint32) snowflake.ID

		// SetSpeaking sends a speaking packet to the Conn socket discord.
		SetSpeaking(ctx context.Context, flags SpeakingFlags) error

		// SetOpusFrameProvider lets you inject your own OpusFrameProvider.
		SetOpusFrameProvider(handler OpusFrameProvider)

		// SetOpusFrameReceiver lets you inject your own OpusFrameReceiver.
		SetOpusFrameReceiver(handler OpusFrameReceiver)

		// Open opens the voice conn. It will connect to the voice gateway and start the Conn conn after it receives the Gateway events.
		Open(ctx context.Context, channelID snowflake.ID, selfMute bool, selfDeaf bool) error

		// Close closes the voice conn. It will close the Conn conn and disconnect from the voice gateway.
		Close(ctx context.Context)

		// HandleVoiceStateUpdate provides the gateway.EventVoiceStateUpdate to the voice conn. Which is needed to connect to the voice Gateway.
		HandleVoiceStateUpdate(update botgateway.EventVoiceStateUpdate)

		// HandleVoiceServerUpdate provides the gateway.EventVoiceServerUpdate to the voice conn. Which is needed to connect to the voice Gateway.
		HandleVoiceServerUpdate(update botgateway.EventVoiceServerUpdate)
	}
)

const voiceReconnectTimeout = 30 * time.Second

// NewConn returns a new default voice conn.
func NewConn(guildID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn {
	cfg := defaultConnConfig()
	cfg.apply(opts)

	conn := &connImpl{
		config:               cfg,
		voiceStateUpdateFunc: voiceStateUpdateFunc,
		removeConnFunc:       removeConnFunc,
		state: State{
			GuildID: guildID,
			UserID:  userID,
		},
		ssrcs: map[uint32]snowflake.ID{},
	}

	conn.initTransports()

	return conn
}

func (c *connImpl) initTransports() {
	daveSession := c.config.DaveSessionCreate(c.config.DaveSessionLogger, godave.UserID(c.state.UserID.String()), c)
	c.dave = daveSession
	c.gateway = c.config.GatewayCreateFunc(daveSession, c.handleMessage, c.handleGatewayClose, append([]GatewayConfigOpt{WithGatewayLogger(c.config.Logger)}, c.config.GatewayConfigOpts...)...)
	c.udp = c.config.UDPConnCreateFunc(daveSession, c.UserIDBySSRC, append([]UDPConnConfigOpt{WithUDPConnLogger(c.config.Logger)}, c.config.UDPConnConfigOpts...)...)
}

// resetTransportsLocked replaces the voice transport generation. stateMu must be held.
func (c *connImpl) resetTransportsLocked() {
	oldGateway, oldUDP, oldDAVE := c.gateway, c.udp, c.dave
	c.initTransports()
	oldGateway.Close()
	_ = oldUDP.Close()
	_ = oldDAVE.Close()

	c.ssrcsMu.Lock()
	clear(c.ssrcs)
	c.ssrcsMu.Unlock()
}

type connImpl struct {
	config               connConfig
	voiceStateUpdateFunc StateUpdateFunc
	removeConnFunc       func()

	state   State
	stateMu sync.Mutex

	gateway Gateway
	udp     UDPConn
	dave    godave.Session

	audioSender   AudioSender
	audioReceiver AudioReceiver

	openedFunc context.CancelFunc

	// Discord completes a voice handshake with one event from each gateway.
	// Only open after receiving a fresh pair so old state is never reused.
	voiceStateReceived  bool
	voiceServerReceived bool

	speakingFlags SpeakingFlags
	speakingSet   bool

	ssrcs   map[uint32]snowflake.ID
	ssrcsMu sync.Mutex
}

func (c *connImpl) SendMLSKeyPackage(mlsKeyPackage []byte) error {
	return c.gateway.Send(context.Background(), OpcodeDaveMLSKeyPackage, GatewayMessageDataDaveMLSKeyPackage(mlsKeyPackage))
}

func (c *connImpl) SendMLSCommitWelcome(mlsCommitWelcome []byte) error {
	return c.gateway.Send(context.Background(), OpcodeDaveMLSCommitWelcome, GatewayMessageDataDaveMLSCommitWelcome(mlsCommitWelcome))
}

func (c *connImpl) SendReadyForTransition(transitionID uint16) error {
	return c.gateway.Send(context.Background(), OpcodeDaveTransitionReady, GatewayMessageDataDaveProtocolReadyForTransition{TransitionID: transitionID})
}

func (c *connImpl) SendInvalidCommitWelcome(transitionID uint16) error {
	return c.gateway.Send(context.Background(), OpcodeDaveMLSInvalidCommitWelcome, GatewayMessageDataDaveInvalidCommitWelcome{TransitionID: transitionID})
}

func (c *connImpl) ChannelID() *snowflake.ID {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if c.state.ChannelID == 0 {
		return nil
	}
	channelID := c.state.ChannelID
	return &channelID
}

func (c *connImpl) GuildID() snowflake.ID {
	return c.state.GuildID
}

func (c *connImpl) UserIDBySSRC(ssrc uint32) snowflake.ID {
	c.ssrcsMu.Lock()
	defer c.ssrcsMu.Unlock()
	return c.ssrcs[ssrc]
}

func (c *connImpl) Gateway() Gateway {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.gateway
}

func (c *connImpl) SetSpeaking(ctx context.Context, flags SpeakingFlags) error {
	c.stateMu.Lock()
	c.speakingFlags = flags
	c.speakingSet = true
	gateway := c.gateway
	c.stateMu.Unlock()
	return gateway.Send(ctx, OpcodeSpeaking, GatewayMessageDataSpeaking{
		SSRC:     gateway.SSRC(),
		Speaking: flags,
	})
}

func (c *connImpl) UDP() UDPConn {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.udp
}

func (c *connImpl) DAVE() godave.Session {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.dave
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

func (c *connImpl) HandleVoiceStateUpdate(update botgateway.EventVoiceStateUpdate) {
	if update.GuildID != c.state.GuildID || update.UserID != c.state.UserID {
		return
	}

	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if update.ChannelID == nil {
		c.state.ChannelID = 0
		c.resetVoiceEventsLocked()
		if c.audioSender != nil {
			c.audioSender.Close()
			c.audioSender = nil
		}
		if c.audioReceiver != nil {
			c.audioReceiver.Close()
			c.audioReceiver = nil
		}
		c.speakingFlags = 0
		c.speakingSet = false
		c.resetTransportsLocked()
	} else {
		moved := c.state.ChannelID != 0 && c.state.ChannelID != *update.ChannelID
		c.state.ChannelID = *update.ChannelID
		if moved {
			// A move needs a new voice-server handshake
			c.resetTransportsLocked()
		}
		c.voiceStateReceived = true
	}
	c.state.SessionID = update.SessionID
	c.state.SelfMute = update.SelfMute
	c.state.SelfDeaf = update.SelfDeaf

	c.tryOpenGateway()
}

func (c *connImpl) HandleVoiceServerUpdate(update botgateway.EventVoiceServerUpdate) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if update.GuildID != c.state.GuildID || update.Endpoint == nil {
		return
	}

	c.state.Token = update.Token
	c.state.Endpoint = *update.Endpoint
	c.voiceServerReceived = true
	c.tryOpenGateway()
}

func (c *connImpl) tryOpenGateway() {
	if c.state.SessionID == "" || c.state.Token == "" || c.state.Endpoint == "" || c.state.ChannelID == 0 {
		return
	}
	if !c.voiceStateReceived || !c.voiceServerReceived {
		return
	}
	state := c.state
	c.resetVoiceEventsLocked()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), voiceReconnectTimeout)
		defer cancel()
		if err := c.gateway.Open(ctx, state); err != nil {
			c.config.Logger.Error("error opening voice gateway", slog.Any("err", err))
		}
	}()
}

func (c *connImpl) handleMessage(gateway Gateway, op Opcode, sequenceNumber int, data GatewayMessageData) {
	c.stateMu.Lock()
	if gateway != c.gateway {
		c.stateMu.Unlock()
		return
	}
	udp := c.udp
	c.stateMu.Unlock()

	switch d := data.(type) {
	case GatewayMessageDataReady:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ourAddress, ourPort, err := udp.Open(ctx, d.IP, d.Port, d.SSRC)
		if err != nil {
			c.config.Logger.Error("voice: failed to open voiceudp conn", slog.Any("err", err))
			break
		}

		encryptionMode, err := ChooseEncryptionMode(d.Modes)
		if err != nil {
			c.config.Logger.Error("voice: failed to choose encryption mode", slog.Any("err", err))
			break
		}

		if err = gateway.Send(ctx, OpcodeSelectProtocol, GatewayMessageDataSelectProtocol{
			Protocol: ProtocolUDP,
			Data: GatewayMessageDataSelectProtocolData{
				Address: ourAddress,
				Port:    ourPort,
				Mode:    encryptionMode,
			},
		}); err != nil {
			c.config.Logger.Error("voice: failed to send select protocol", slog.Any("err", err))
		}

	case GatewayMessageDataSessionDescription:
		if err := udp.SetSecretKey(d.Mode, d.SecretKey); err != nil {
			c.config.Logger.Error("voice: failed to set secret key", slog.Any("err", err))
		}
		c.stateMu.Lock()
		speakingFlags, speakingSet := c.speakingFlags, c.speakingSet
		current := gateway == c.gateway
		c.stateMu.Unlock()
		if current && speakingSet {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := gateway.Send(ctx, OpcodeSpeaking, GatewayMessageDataSpeaking{
				SSRC:     gateway.SSRC(),
				Speaking: speakingFlags,
			}); err != nil {
				c.config.Logger.Error("voice: failed to restore speaking state", slog.Any("err", err))
			}
			cancel()
		}
		if openedFunc := c.openedFunc; openedFunc != nil {
			openedFunc()
		}

	case GatewayMessageDataSpeaking:
		c.ssrcsMu.Lock()
		c.ssrcs[d.SSRC] = d.UserID
		c.ssrcsMu.Unlock()

	case GatewayMessageDataClientDisconnect:
		c.ssrcsMu.Lock()
		for ssrc, userID := range c.ssrcs {
			if userID == d.UserID {
				delete(c.ssrcs, ssrc)
				break
			}
		}
		c.ssrcsMu.Unlock()
		if c.audioReceiver != nil {
			c.audioReceiver.CleanupUser(d.UserID)
		}
	}
	if c.config.EventHandlerFunc != nil {
		c.config.EventHandlerFunc(gateway, op, sequenceNumber, data)
	}
}

func (c *connImpl) handleGatewayClose(gateway Gateway, err error) {
	c.stateMu.Lock()
	if gateway != c.gateway {
		c.stateMu.Unlock()
		return
	}
	c.stateMu.Unlock()

	newConnection := true
	var closeError *websocket.CloseError
	if errors.As(err, &closeError) {
		closeCode := GatewayCloseEventCodeByCode(closeError.Code)
		// Wait for `ServerUpdate` and `StateUpdate` to close the gateway (a channel move could trigger this handler)
		if closeCode == GatewayCloseEventCodeDisconnected || closeCode == GatewayCloseEventCodeCallTerminated {
			return
		}
		newConnection = closeCode.NewConnection
	}

	if newConnection {
		c.stateMu.Lock()
		channelID := c.state.ChannelID
		selfMute := c.state.SelfMute
		selfDeaf := c.state.SelfDeaf
		c.state.SessionID = ""
		c.state.Token = ""
		c.state.Endpoint = ""
		c.stateMu.Unlock()

		udp := c.UDP()
		_ = udp.Close()
		if channelID == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), voiceReconnectTimeout)
		defer cancel()
		if err = c.Open(ctx, channelID, selfMute, selfDeaf); err != nil {
			c.config.Logger.Error("voice: failed to reopen voice conn with a fresh session", slog.Any("err", err))
		} else {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c.Close(ctx)
}

func (c *connImpl) Open(ctx context.Context, channelID snowflake.ID, selfMute bool, selfDeaf bool) error {
	c.config.Logger.Debug("opening voice conn")

	openedCtx, cancel := context.WithCancel(context.Background())
	c.openedFunc = cancel
	defer cancel()

	c.stateMu.Lock()
	c.resetVoiceEventsLocked()
	guildID := c.state.GuildID
	c.stateMu.Unlock()

	if err := c.voiceStateUpdateFunc(ctx, guildID, &channelID, selfMute, selfDeaf); err != nil {
		return err
	}

	select {
	case <-openedCtx.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connImpl) Close(ctx context.Context) {
	_ = c.voiceStateUpdateFunc(ctx, c.state.GuildID, nil, false, false)

	c.gateway.Close()
	_ = c.udp.Close()
	_ = c.dave.Close()

	c.removeConnFunc()
}

func (c *connImpl) resetVoiceEventsLocked() {
	c.voiceStateReceived = false
	c.voiceServerReceived = false
}
