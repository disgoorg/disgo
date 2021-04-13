package internal

import "github.com/DisgoOrg/disgo/api"

func newAudioControllerImpl(disgo api.Disgo) api.AudioController {
	return &AudioControllerImpl{disgo: disgo}
}

type AudioControllerImpl struct {
	disgo api.Disgo
}

func (c *AudioControllerImpl) Disgo() api.Disgo {
	return c.disgo
}

func (c *AudioControllerImpl) Connect(guildID api.Snowflake, channelID api.Snowflake) error {
	return c.Disgo().Gateway().Conn().WriteJSON(api.NewGatewayCommand(api.OpVoiceStateUpdate, api.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}
func (c *AudioControllerImpl) Disconnect(guildID api.Snowflake) error {
	return c.Disgo().Gateway().Conn().WriteJSON(api.NewGatewayCommand(api.OpVoiceStateUpdate, api.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}
