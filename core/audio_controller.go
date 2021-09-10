package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// context support?

// AudioController lets you Connect / Disconnect from a Channel
type AudioController interface {
	// Bot returns the core.Bot instance
	Bot() *Bot

	// Connect sends an core.GatewayCommand to connect to an core.Channel
	Connect(guildID discord.Snowflake, channelID discord.Snowflake) error

	// Disconnect sends an core.GatewayCommand to disconnect from an core.Channel
	Disconnect(guildID discord.Snowflake) error
}

func NewAudioController(bot *Bot) AudioController {
	return &audioControllerImpl{bot: bot}
}

type audioControllerImpl struct {
	bot *Bot
}

func (c *audioControllerImpl) Bot() *Bot {
	return c.bot
}

func (c *audioControllerImpl) Connect(guildID discord.Snowflake, channelID discord.Snowflake) error {
	gw, err := c.getGateway()
	if err != nil {
		return err
	}
	return gw.Send(discord.NewGatewayCommand(discord.OpVoiceStateUpdate, discord.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}

func (c *audioControllerImpl) Disconnect(guildID discord.Snowflake) error {
	gw, err := c.getGateway()
	if err != nil {
		return err
	}
	return gw.Send(discord.NewGatewayCommand(discord.OpVoiceStateUpdate, discord.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *audioControllerImpl) getGateway() (gateway.Gateway, error) {
	if c.Bot().Gateway == nil {
		return nil, discord.ErrNoGateway
	}
	if !c.Bot().Gateway.Status().IsConnected() {
		return nil, discord.ErrNoGatewayConn
	}
	return c.Bot().Gateway, nil
}
