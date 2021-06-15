package internal

import (
	"github.com/DisgoOrg/disgo/api"
)

func newAudioControllerImpl(disgo api.Disgo) api.AudioController {
	return &AudioControllerImpl{disgo: disgo}
}

// AudioControllerImpl lets you Connect / Disconnect from a api.VoiceChannel
type AudioControllerImpl struct {
	disgo api.Disgo
}

// Disgo returns the api.Disgo instance
func (c *AudioControllerImpl) Disgo() api.Disgo {
	return c.disgo
}

// Connect sends a api.GatewayCommand to connect to a api.VoiceChannel
func (c *AudioControllerImpl) Connect(guildID api.Snowflake, channelID api.Snowflake) error {
	gateway, err := c.getGateway()
	if err != nil {
		return err
	}
	return gateway.Send(api.NewGatewayCommand(api.OpVoiceStateUpdate, api.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}

// Disconnect sends a api.GatewayCommand to disconnect from a api.VoiceChannel
func (c *AudioControllerImpl) Disconnect(guildID api.Snowflake) error {
	gateway, err := c.getGateway()
	if err != nil {
		return err
	}
	return gateway.Send(api.NewGatewayCommand(api.OpVoiceStateUpdate, api.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *AudioControllerImpl) getGateway() (api.Gateway, error) {
	gateway := c.Disgo().Gateway()
	if gateway == nil {
		return nil, api.ErrNoGateway
	}
	if !gateway.Status().IsConnected() {
		return nil, api.ErrNoGatewayConn
	}
	return gateway, nil
}
