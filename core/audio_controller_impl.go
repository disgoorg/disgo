package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

func NewAudioController(disgo Disgo) AudioController {
	return &AudioControllerImpl{disgo: disgo}
}

// AudioControllerImpl lets you Connect / Disconnect from an api.VoiceChannel
type AudioControllerImpl struct {
	disgo Disgo
}

// Disgo returns the api.Disgo instance
func (c *AudioControllerImpl) Disgo() Disgo {
	return c.disgo
}

// Connect sends an api.GatewayCommand to connect to an api.VoiceChannel
func (c *AudioControllerImpl) Connect(guildID discord.Snowflake, channelID discord.Snowflake) error {
	gw, err := c.getGateway()
	if err != nil {
		return err
	}
	return gw.Send(gateway.NewGatewayCommand(gateway.OpVoiceStateUpdate, gateway.UpdateVoiceStateCommand{
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
	return gw.Send(gateway.NewGatewayCommand(gateway.OpVoiceStateUpdate, gateway.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *AudioControllerImpl) getGateway() (gateway.Gateway, error) {
	gw := c.Disgo().Gateway()
	if gw == nil {
		return nil, discord.ErrNoGateway
	}
	if !gw.Status().IsConnected() {
		return nil, discord.ErrNoGatewayConn
	}
	return gw, nil
}
