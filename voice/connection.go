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

func NewConnection(guildID snowflake.ID, userID snowflake.ID, opts ...ConfigOpt) *Connection {
	config := DefaultConfig()
	config.Apply(opts)

	return &Connection{
		config:  *config,
		guildID: guildID,
		userID:  userID,
	}
}

type Connection struct {
	config  Config
	guildID snowflake.ID
	userID  snowflake.ID

	channelID snowflake.ID
	sessionID string
	token     string
	endpoint  string

	gateway *Gateway
	conn    *UDPConn
}

func (c *Connection) HandleVoiceStateUpdate(update discord.VoiceStateUpdate) {
	if update.GuildID != c.guildID || update.UserID != c.userID {
		return
	}
	if update.ChannelID == nil {
		c.Close()
		return
	}
	c.channelID = *update.ChannelID
	c.sessionID = update.SessionID

	c.reconnect()
}

func (c *Connection) HandleVoiceServerUpdate(update discord.VoiceServerUpdate) {
	if update.GuildID != c.guildID || update.Endpoint == nil {
		return
	}
	c.token = update.Token
	c.endpoint = *update.Endpoint
	c.reconnect()
}

func (c *Connection) handleGatewayMessage(opCode GatewayOpcode, data GatewayMessageData) {
	switch d := data.(type) {
	case GatewayMessageDataReady:
		c.conn
	}
}

func (c *Connection) handleGatewayClose(gateway *Gateway, err error) {

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

}

func (c *Connection) Close() {
	c.conn.Close()
	c.conn = nil
	c.gateway.Close()
	c.gateway = nil
}
