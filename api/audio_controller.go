package api

// AudioController lets you Connect / Disconnect from a VoiceChannel
type AudioController interface {
	Disgo() Disgo
	Connect(guildID Snowflake, channelID Snowflake) error
	Disconnect(guildID Snowflake) error
}
