package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// context support?

// AudioController lets you Connect / Disconnect from a VoiceChannel
type AudioController interface {
	Bot() *Bot
	Connect(guildID discord.Snowflake, channelID discord.Snowflake) error
	Disconnect(guildID discord.Snowflake) error
}

func NewAudioController(bot *Bot) AudioController {
	return &AudioControllerImpl{bot: bot}
}

// AudioControllerImpl lets you Connect / Disconnect from an api.VoiceChannel
type AudioControllerImpl struct {
	bot *Bot
}

// Bot returns the api.Bot instance
func (c *AudioControllerImpl) Bot() *Bot {
	return c.bot
}

// Connect sends an api.GatewayCommand to connect to an api.VoiceChannel
func (c *AudioControllerImpl) Connect(guildID discord.Snowflake, channelID discord.Snowflake) error {
	gw, err := c.getGateway()
	if err != nil {
		return err
	}
	return gw.Send(discord.NewGatewayCommand(discord.OpVoiceStateUpdate, discord.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}

// Disconnect sends an api.GatewayCommand to disconnect from an api.VoiceChannel
func (c *AudioControllerImpl) Disconnect(guildID discord.Snowflake) error {
	gw, err := c.getGateway()
	if err != nil {
		return err
	}
	return gw.Send(discord.NewGatewayCommand(discord.OpVoiceStateUpdate, discord.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *AudioControllerImpl) getGateway() (gateway.Gateway, error) {
	if c.Bot().Gateway == nil {
		return nil, discord.ErrNoGateway
	}
	if !c.Bot().Gateway.Status().IsConnected() {
		return nil, discord.ErrNoGatewayConn
	}
	return c.Bot().Gateway, nil
}
