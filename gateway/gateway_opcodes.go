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

type CloseEventCode struct {
	Code        int
	Description string
	Explanation string
	Reconnect   bool
}

var (
	CloseEventCodeUnknownError = CloseEventCode{
		Code:        4000,
		Description: "Unknown error",
		Explanation: "We're not sure what went wrong. Try reconnecting?",
		Reconnect:   true,
	}

	CloseEventCodeUnknownOpcode = CloseEventCode{
		Code:        4001,
		Description: "Unknown opcode",
		Explanation: "You sent an invalid Gateway opcode or an invalid payload for an opcode. Don't do that!",
		Reconnect:   true,
	}

	CloseEventCodeDecodeError = CloseEventCode{
		Code:        4002,
		Description: "Decode error",
		Explanation: "You sent an invalid payload to Discord. Don't do that!",
		Reconnect:   true,
	}

	CloseEventCodeNotAuthenticated = CloseEventCode{
		Code:        4003,
		Description: "Not authenticated",
		Explanation: "You sent us a payload prior to identifying.",
		Reconnect:   true,
	}

	CloseEventCodeAuthenticationFailed = CloseEventCode{
		Code:        4004,
		Description: "Authentication failed",
		Explanation: "The account token sent with your identify payload is incorrect.",
		Reconnect:   false,
	}

	CloseEventCodeAlreadyAuthenticated = CloseEventCode{
		Code:        4005,
		Description: "Already authenticated",
		Explanation: "You sent more than one identify payload. Don't do that!",
		Reconnect:   true,
	}

	CloseEventCodeInvalidSeq = CloseEventCode{
		Code:        4007,
		Description: "Invalid seq",
		Explanation: "The sequence sent when resuming the session was invalid. Reconnect and start a new session.",
		Reconnect:   true,
	}

	CloseEventCodeRateLimited = CloseEventCode{
		Code:        4008,
		Description: "Rate limited.",
		Explanation: "You're sending payloads to us too quickly. Slow it down! You will be disconnected on receiving this.",
		Reconnect:   true,
	}

	CloseEventCodeSessionTimed = CloseEventCode{
		Code:        4009,
		Description: "Session timed out",
		Explanation: "Your session timed out. Reconnect and start a new one.",
		Reconnect:   true,
	}

	CloseEventCodeInvalidShard = CloseEventCode{
		Code:        4010,
		Description: "Invalid shard",
		Explanation: "You sent us an invalid shard when identifying.",
		Reconnect:   false,
	}

	CloseEventCodeShardingRequired = CloseEventCode{
		Code:        4011,
		Description: "Sharding required",
		Explanation: "The session would have handled too many guilds - you are required to shard your connection in order to connect.",
		Reconnect:   false,
	}

	CloseEventCodeInvalidAPIVersion = CloseEventCode{
		Code:        4012,
		Description: "Invalid API version",
		Explanation: "You sent an invalid version for the gateway.",
		Reconnect:   false,
	}

	CloseEventCodeInvalidIntent = CloseEventCode{
		Code:        4013,
		Description: "Invalid intent(s)",
		Explanation: "You sent an invalid intent for a Gateway Intent. You may have incorrectly calculated the bitwise value.",
		Reconnect:   false,
	}

	CloseEventCodeDisallowedIntent = CloseEventCode{
		Code:        4014,
		Description: "Disallowed intent(s)",
		Explanation: "You sent a disallowed intent for a Gateway Intent. You may have tried to specify an intent that you have not enabled or are not approved for.",
		Reconnect:   false,
	}

	CloseEventCodeUnknown = CloseEventCode{
		Code:        0,
		Description: "Unknown",
		Explanation: "Unknown Gateway Close Event Code",
		Reconnect:   true,
	}

	CloseEventCodes = map[int]CloseEventCode{
		CloseEventCodeUnknownError.Code:         CloseEventCodeUnknownError,
		CloseEventCodeUnknownOpcode.Code:        CloseEventCodeUnknownOpcode,
		CloseEventCodeDecodeError.Code:          CloseEventCodeDecodeError,
		CloseEventCodeNotAuthenticated.Code:     CloseEventCodeNotAuthenticated,
		CloseEventCodeAuthenticationFailed.Code: CloseEventCodeAuthenticationFailed,
		CloseEventCodeAlreadyAuthenticated.Code: CloseEventCodeAlreadyAuthenticated,
		CloseEventCodeInvalidSeq.Code:           CloseEventCodeInvalidSeq,
		CloseEventCodeRateLimited.Code:          CloseEventCodeRateLimited,
		CloseEventCodeSessionTimed.Code:         CloseEventCodeSessionTimed,
		CloseEventCodeInvalidShard.Code:         CloseEventCodeInvalidShard,
		CloseEventCodeInvalidAPIVersion.Code:    CloseEventCodeInvalidAPIVersion,
		CloseEventCodeInvalidIntent.Code:        CloseEventCodeInvalidIntent,
		CloseEventCodeDisallowedIntent.Code:     CloseEventCodeDisallowedIntent,
	}
)

func CloseEventCodeByCode(code int) CloseEventCode {
	closeCode, ok := CloseEventCodes[code]
	if !ok {
		return CloseEventCodeUnknown
	}
	return closeCode
}
