package discord

// GatewayOpcode are opcodes used by discord
type GatewayOpcode int

// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
//goland:noinspection GoUnusedConst
const (
	GatewayOpcodeDispatch GatewayOpcode = iota
	GatewayOpcodeHeartbeat
	GatewayOpcodeIdentify
	GatewayOpcodePresenceUpdate
	GatewayOpcodeVoiceStateUpdate
	_
	GatewayOpcodeResume
	GatewayOpcodeReconnect
	GatewayOpcodeRequestGuildMembers
	GatewayOpcodeInvalidSession
	GatewayOpcodeHello
	GatewayOpcodeHeartbeatACK
)

type GatewayCloseEventCode int

//goland:noinspection GoUnusedConst
const (
	GatewayCloseEventCodeUnknownError GatewayCloseEventCode = iota + 4000
	GatewayCloseEventCodeUnknownOpcode
	GatewayCloseEventCodeDecodeError
	GatewayCloseEventCodeNotAuthenticated
	GatewayCloseEventCodeAuthenticationFailed
	GatewayCloseEventCodeAlreadyAuthenticated
	_
	GatewayCloseEventCodeInvalidSeq
	GatewayCloseEventCodeRateLimited
	GatewayCloseEventCodeSessionTimedOut
	GatewayCloseEventCodeInvalidShard
	GatewayCloseEventCodeShardingRequired
	GatewayCloseEventCodeInvalidAPIVersion
	GatewayCloseEventCodeInvalidIntents
	GatewayCloseEventCodeDisallowedIntents
)

type VoiceOpcode int

//goland:noinspection GoUnusedConst
const (
	VoiceOpcodeIdentify VoiceOpcode = iota
	VoiceOpcodeSelectProtocol
	VoiceOpcodeReady
	VoiceOpcodeHeartbeat
	VoiceOpcodeSessionDescription
	VoiceOpcodeSpeaking
	VoiceOpcodeHeartbeatACK
	VoiceOpcodeResume
	VoiceOpcodeHello
	VoiceOpcodeResumed
	_
	_
	_
	VoiceOpcodeClientDisconnect
)

type VoiceCloseEventCode int

//goland:noinspection GoUnusedConst
const (
	VoiceCloseEventCodeUnknownOpcode VoiceCloseEventCode = iota + 4001
	VoiceCloseEventCodeDecodeError
	VoiceCloseEventCodeNotAuthenticated
	VoiceCloseEventCodeAuthenticationFailed
	VoiceCloseEventCodeAlreadyAuthenticated
	VoiceCloseEventCodeSessionNoLongerValid
	_
	_
	VoiceCloseEventCodeSessionTimedOut
	_
	VoiceCloseEventCodeServerNotFound
	VoiceCloseEventCodeUnknownProtocol
	_
	VoiceCloseEventCodeDisconnected
	VoiceCloseEventCodeVoiceServerCrashed
	VoiceCloseEventCodeUnknownEncryptionMode
)
