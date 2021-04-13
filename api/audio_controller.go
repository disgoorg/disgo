package api

type AudioController interface {
	Disgo() Disgo
	Connect(guildID Snowflake, channelID Snowflake) error
	Disconnect(guildID Snowflake) error
}
