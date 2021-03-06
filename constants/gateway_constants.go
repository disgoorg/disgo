package constants

// Gateway opcodes used by discord
type GatewayOp int

// Constants for the gateway opcodes
const (
	OpDispatch GatewayOp = iota
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
