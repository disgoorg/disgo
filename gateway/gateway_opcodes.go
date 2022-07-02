package gateway

// Opcode are opcodes used by discord
type Opcode int

// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
const (
	OpcodeDispatch Opcode = iota
	OpcodeHeartbeat
	OpcodeIdentify
	OpcodePresenceUpdate
	OpcodeVoiceStateUpdate
	_
	OpcodeResume
	OpcodeReconnect
	OpcodeRequestGuildMembers
	OpcodeInvalidSession
	OpcodeHello
	OpcodeHeartbeatACK
)

type CloseEventCode int

const (
	CloseEventCodeUnknownError CloseEventCode = iota + 4000
	CloseEventCodeUnknownOpcode
	CloseEventCodeDecodeError
	CloseEventCodeNotAuthenticated
	CloseEventCodeAuthenticationFailed
	CloseEventCodeAlreadyAuthenticated
	_
	CloseEventCodeInvalidSeq
	CloseEventCodeRateLimited
	CloseEventCodeSessionTimedOut
	CloseEventCodeInvalidShard
	CloseEventCodeShardingRequired
	CloseEventCodeInvalidAPIVersion
	CloseEventCodeInvalidIntents
	CloseEventCodeDisallowedIntents
)

func (c CloseEventCode) ShouldReconnect() bool {
	switch c {
	case CloseEventCodeAuthenticationFailed,
		CloseEventCodeInvalidShard,
		CloseEventCodeShardingRequired,
		CloseEventCodeInvalidAPIVersion,
		CloseEventCodeInvalidIntents,
		CloseEventCodeDisallowedIntents:
		return false

	default:
		return true
	}
}
