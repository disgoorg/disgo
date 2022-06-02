package voice

import (
	"context"
	"errors"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var (
	ErrAlreadyConnected = errors.New("already connected")
)

type ConnectionCreateFunc func(guildID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) *Connection

func NewConnection(guildID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) *Connection {
	config := DefaultConnectionConfig()
	config.Apply(opts)

	return &Connection{
		config:  *config,
		guildID: guildID,
		userID:  userID,
	}
}

type Connection struct {
	config  ConnectionConfig
	guildID snowflake.ID
	userID  snowflake.ID

	channelID snowflake.ID
	sessionID string
	token     string
	endpoint  string
	ssrc      uint32

	sendHandler    SendHandler
	receiveHandler ReceiveHandler

	gateway *Gateway
	conn    *UDPConn
}

func (c *Connection) Gateway() *Gateway {
	return c.gateway
}

func (c *Connection) UDPConn() *UDPConn {
	return c.conn
}

func (c *Connection) SetSendHandler(handler SendHandler) {
	c.sendHandler = handler
}

func (c *Connection) SetReceiveHandler(handler ReceiveHandler) {
	c.receiveHandler = handler
}

func (c *Connection) HandleVoiceStateUpdate(update discord.VoiceState) {
	if update.GuildID != c.guildID || update.UserID != c.userID {
		return
	}
	if update.ChannelID == nil {
		c.Close()
		return
	}
	c.channelID = *update.ChannelID
	c.sessionID = update.SessionID
}

func (c *Connection) HandleVoiceServerUpdate(update discord.VoiceServerUpdate) {
	if update.GuildID != c.guildID || update.Endpoint == nil {
		return
	}
	c.token = update.Token
	c.endpoint = *update.Endpoint
	go c.reconnect()
}

func (c *Connection) handleGatewayMessage(opCode GatewayOpcode, data GatewayMessageData) {
	switch d := data.(type) {
	case GatewayMessageDataReady:
		c.ssrc = d.SSRC
		udpConn := c.config.UDPConnCreateFunc(d.IP, d.Port, d.SSRC)
		address, port, err := udpConn.Open(context.Background())
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
			return
		}
	case GatewayMessageDataSessionDescription:
		c.conn.HandleGatewayMessageSessionDescription(d)
	}
}

func (c *Connection) handleGatewayClose(gateway *Gateway, err error) {
	c.config.Logger.Error("voice gateway closed. error: ", err)
}

func (c *Connection) Open(ctx context.Context) error {
	c.config.Logger.Debug("voice: opening connection")
	if c.gateway != nil {
		return ErrAlreadyConnected
	}
	c.gateway = c.config.GatewayCreateFunc(c.guildID, c.userID, c.sessionID, c.token, c.endpoint, c.handleGatewayMessage, c.handleGatewayClose, c.config.GatewayConfigOpts...)
	return c.gateway.Open(ctx)
}

func (c *Connection) reconnect() {
	c.Open(context.Background())
}

func (c *Connection) Close() {
	c.conn.Close()
	c.conn = nil
	c.gateway.Close()
	c.gateway = nil
}
