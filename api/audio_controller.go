package api

import "errors"

// errors returned when no gateway or ws conn exists
var (
	ErrNoGateway = errors.New("no gateway initialized")
	ErrNoGatewayConn = errors.New("no active gateway connection found")
)

// AudioController lets you Connect / Disconnect from a VoiceChannel
type AudioController interface {
	Disgo() Disgo
	Connect(guildID Snowflake, channelID Snowflake) error
	Disconnect(guildID Snowflake) error
}
