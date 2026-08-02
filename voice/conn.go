package voice

import (
	"context"
	"errors"
	"log/slog"
	"net"
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

	openedChan        chan struct{}
	gatewayOpenCancel context.CancelFunc
	closed            bool

	voiceStateReceived  bool
	voiceServerReceived bool
	speakingFlags       SpeakingFlags
	speakingSet         bool

	// gatewayOpenMu serializes calls to Gateway.Open.
	gatewayOpenMu sync.Mutex

	ssrcs   map[uint32]snowflake.ID
	ssrcsMu sync.Mutex
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

func (c *connImpl) SendMLSKeyPackage(mlsKeyPackage []byte) error {
	return c.Gateway().Send(context.Background(), OpcodeDaveMLSKeyPackage, GatewayMessageDataDaveMLSKeyPackage(mlsKeyPackage))
}

func (c *connImpl) SendMLSCommitWelcome(mlsCommitWelcome []byte) error {
	return c.Gateway().Send(context.Background(), OpcodeDaveMLSCommitWelcome, GatewayMessageDataDaveMLSCommitWelcome(mlsCommitWelcome))
}

func (c *connImpl) SendReadyForTransition(transitionID uint16) error {
	return c.Gateway().Send(context.Background(), OpcodeDaveTransitionReady, GatewayMessageDataDaveProtocolReadyForTransition{TransitionID: transitionID})
}

func (c *connImpl) SendInvalidCommitWelcome(transitionID uint16) error {
	return c.Gateway().Send(context.Background(), OpcodeDaveMLSInvalidCommitWelcome, GatewayMessageDataDaveInvalidCommitWelcome{TransitionID: transitionID})
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

	if c.closed {
		return
	}

	if update.ChannelID == nil {
		c.state.ChannelID = 0
		c.voiceStateReceived = false
		c.voiceServerReceived = false
		c.cancelGatewayOpenLocked()
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
			// A move needs a new voice-server handshake. Do not leave the old
			// websocket, UDP socket, or DAVE epoch around to race the server update.
			c.cancelGatewayOpenLocked()
			c.resetTransportsLocked()
		}
	}
	c.state.SessionID = update.SessionID
	c.state.SelfMute = update.SelfMute
	c.state.SelfDeaf = update.SelfDeaf

	if update.ChannelID != nil {
		c.voiceStateReceived = true
		c.tryOpenGateway()
	}
}

func (c *connImpl) HandleVoiceServerUpdate(update botgateway.EventVoiceServerUpdate) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if c.closed || update.GuildID != c.state.GuildID || update.Endpoint == nil {
		return
	}

	c.state.Token = update.Token
	c.state.Endpoint = *update.Endpoint
	c.voiceServerReceived = true
	c.tryOpenGateway()
}

func (c *connImpl) tryOpenGateway() {
	go func() {
		c.gatewayOpenMu.Lock()
		defer c.gatewayOpenMu.Unlock()

		c.stateMu.Lock()
		gateway := c.gateway
		status := gateway.Status()
		canOpen := !c.closed && c.voiceStateReceived && c.voiceServerReceived &&
			c.state.Token != "" && c.state.Endpoint != "" && c.state.ChannelID != 0 &&
			(status == StatusUnconnected || status == StatusDisconnected)
		if !canOpen {
			c.stateMu.Unlock()
			return
		}

		state := c.state
		c.voiceStateReceived = false
		c.voiceServerReceived = false
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		c.gatewayOpenCancel = cancel
		c.stateMu.Unlock()
		defer cancel()

		err := gateway.Open(ctx, state)

		c.stateMu.Lock()
		c.gatewayOpenCancel = nil
		c.stateMu.Unlock()

		if err != nil && !errors.Is(err, context.Canceled) {
			c.config.Logger.Error("error opening voice gateway", slog.Any("err", err))
		}
	}()
}

// cancelGatewayOpenLocked cancels the in-flight open. stateMu must be held.
func (c *connImpl) cancelGatewayOpenLocked() {
	if c.gatewayOpenCancel != nil {
		c.gatewayOpenCancel()
		c.gatewayOpenCancel = nil
	}
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

		c.stateMu.Lock()
		if gateway == c.gateway && c.openedChan != nil {
			close(c.openedChan)
			c.openedChan = nil
		}
		c.stateMu.Unlock()

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
	if gateway != c.gateway || c.closed {
		c.stateMu.Unlock()
		return
	}
	state := c.state
	c.stateMu.Unlock()

	var closeError *websocket.CloseError
	if errors.As(err, &closeError) {
		closeCode := GatewayCloseEventCodeByCode(closeError.Code)
		if closeCode == GatewayCloseEventCodeDisconnected || closeCode == GatewayCloseEventCodeCallTerminated {
			// These closes can precede the main gateway updates that distinguish
			// a move from a disconnect. Acting here races those definitive updates.
			return
		}
		if closeCode.NewConnection {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err = c.Open(ctx, state.ChannelID, state.SelfMute, state.SelfDeaf); err != nil {
				c.config.Logger.Error("voice: failed to reopen voice conn after full reconnect close code", slog.Any("err", err))
			} else {
				return
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c.Close(ctx)
}

func (c *connImpl) Open(ctx context.Context, channelID snowflake.ID, selfMute bool, selfDeaf bool) error {
	c.config.Logger.Debug("opening voice conn")

	openedChan := make(chan struct{})
	c.stateMu.Lock()
	if c.closed {
		c.stateMu.Unlock()
		return net.ErrClosed
	}
	c.voiceStateReceived = false
	c.voiceServerReceived = false
	c.openedChan = openedChan
	guildID := c.state.GuildID
	c.stateMu.Unlock()
	defer func() {
		c.stateMu.Lock()
		if c.openedChan == openedChan {
			c.openedChan = nil
		}
		c.stateMu.Unlock()
	}()

	if err := c.voiceStateUpdateFunc(ctx, guildID, &channelID, selfMute, selfDeaf); err != nil {
		return err
	}

	select {
	case <-openedChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connImpl) Close(ctx context.Context) {
	c.stateMu.Lock()
	if c.closed {
		c.stateMu.Unlock()
		return
	}
	c.closed = true
	c.state.ChannelID = 0
	c.voiceStateReceived = false
	c.voiceServerReceived = false
	c.cancelGatewayOpenLocked()
	guildID := c.state.GuildID
	gateway, udp, dave := c.gateway, c.udp, c.dave
	c.stateMu.Unlock()

	_ = c.voiceStateUpdateFunc(ctx, guildID, nil, false, false)

	gateway.Close()
	_ = udp.Close()
	_ = dave.Close()

	c.removeConnFunc()
}
