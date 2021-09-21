package discord

// Op are opcodes used by discord
type Op int

// Constants for the gateway opcodes
//goland:noinspection GoUnusedConst
const (
	OpDispatch Op = iota
	OpHeartbeat
	OpIdentify
	OpPresenceUpdate
	OpVoiceStateUpdate
	_
	OpResume
	OpReconnect
	OpRequestGuildMembers
	OpInvalidSession
	OpHello
	OpHeartbeatACK
)
