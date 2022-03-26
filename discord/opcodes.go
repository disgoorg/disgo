package discord

import "github.com/gorilla/websocket"

// GatewayOpcode are opcodes used by discord
type GatewayOpcode int

// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
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

func (c GatewayCloseEventCode) ShouldReconnect() bool {
	switch c {
	case websocket.CloseNormalClosure,
		GatewayCloseEventCodeAuthenticationFailed,
		GatewayCloseEventCodeInvalidShard,
		GatewayCloseEventCodeShardingRequired,
		GatewayCloseEventCodeInvalidAPIVersion,
		GatewayCloseEventCodeInvalidIntents,
		GatewayCloseEventCodeDisallowedIntents:
		return false

	default:
		return true
	}
}

type VoiceOpcode int

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
