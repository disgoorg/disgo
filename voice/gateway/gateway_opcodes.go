package gateway

type Opcode int

const (
	OpcodeIdentify Opcode = iota
	OpcodeSelectProtocol
	OpcodeReady
	OpcodeHeartbeat
	OpcodeSessionDescription
	OpcodeSpeaking
	OpcodeHeartbeatACK
	OpcodeResume
	OpcodeHello
	OpcodeResumed
	_
	_
	_
	OpcodeClientDisconnect
	OpcodeGuildSync
)

type CloseEventCode struct {
	Code        int
	Description string
	Explanation string
	Reconnect   bool
}

var (
	CloseEventCodeHeartbeatTimeout = CloseEventCode{
		Code:        1000,
		Description: "Heartbeat timeout",
		Explanation: "We did not heartbeat in time.",
		Reconnect:   true,
	}

	CloseEventCodeUnknownError = CloseEventCode{
		Code:        4000,
		Description: "Unknown error",
		Explanation: "We're not sure what went wrong. Try reconnecting?",
		Reconnect:   true,
	}

	CloseEventCodeUnknownOpcode = CloseEventCode{
		Code:        4001,
		Description: "Unknown opcode",
		Explanation: "You sent an invalid opcode.",
		Reconnect:   true,
	}

	CloseEventCodeFailedDecode = CloseEventCode{
		Code:        4002,
		Description: "Failed to decode payload",
		Explanation: "You sent a invalid payload in your identifying to the Gateway.",
		Reconnect:   true,
	}

	CloseEventCodeNotAuthenticated = CloseEventCode{
		Code:        4003,
		Description: "Not authenticated",
		Explanation: "You sent a payload before identifying with the Gateway.",
		Reconnect:   false,
	}

	CloseEventCodeAuthenticationFailed = CloseEventCode{
		Code:        4004,
		Description: "Authentication failed",
		Explanation: "The token you sent in your identify payload is incorrect.",
		Reconnect:   false,
	}

	CloseEventCodeAlreadyAuthenticated = CloseEventCode{
		Code:        4005,
		Description: "Already authenticated",
		Explanation: "You sent more than one identify payload. Stahp.",
		Reconnect:   true,
	}

	CloseEventCodeSessionNoLongerValid = CloseEventCode{
		Code:        4006,
		Description: "Session no longer valid",
		Explanation: "Your session is no longer valid.",
		Reconnect:   false,
	}

	CloseEventCodeSessionTimeout = CloseEventCode{
		Code:        4009,
		Description: "Session timeout",
		Explanation: "Your session has timed out.",
		Reconnect:   false,
	}

	CloseEventCodeServerNotFound = CloseEventCode{
		Code:        4011,
		Description: "Server not found",
		Explanation: "We can't find the server you're trying to connect to.",
		Reconnect:   false,
	}

	CloseEventCodeUnknownProtocol = CloseEventCode{
		Code:        4012,
		Description: "Unknown protocol",
		Explanation: "We didn't recognize the protocol you sent.",
		Reconnect:   false,
	}

	CloseEventCodeDisconnected = CloseEventCode{
		Code:        4014,
		Description: "Disconnected",
		Explanation: "Channel was deleted, you were kicked, voice server changed, or the main gateway session was dropped. Don't reconnect.",
		Reconnect:   false,
	}

	CloseEventCodeVoiceServerCrash = CloseEventCode{
		Code:        4015,
		Description: "Voice server crashed",
		Explanation: "The server crashed. Our bad! Try resuming.",
		Reconnect:   true,
	}

	CloseEventCodeUnknownEncryptionMode = CloseEventCode{
		Code:        4016,
		Description: "Unknown encryption mode",
		Explanation: "We didn't recognize your encryption.",
		Reconnect:   false,
	}

	CloseEventCodeUnknown = CloseEventCode{
		Code:        0,
		Description: "Unknown",
		Explanation: "Unknown Voice Close Event Code",
		Reconnect:   true,
	}

	CloseEventCodes = map[int]CloseEventCode{
		CloseEventCodeHeartbeatTimeout.Code:      CloseEventCodeHeartbeatTimeout,
		CloseEventCodeUnknownError.Code:          CloseEventCodeUnknownError,
		CloseEventCodeUnknownOpcode.Code:         CloseEventCodeUnknownOpcode,
		CloseEventCodeFailedDecode.Code:          CloseEventCodeFailedDecode,
		CloseEventCodeNotAuthenticated.Code:      CloseEventCodeNotAuthenticated,
		CloseEventCodeAuthenticationFailed.Code:  CloseEventCodeAuthenticationFailed,
		CloseEventCodeAlreadyAuthenticated.Code:  CloseEventCodeAlreadyAuthenticated,
		CloseEventCodeSessionNoLongerValid.Code:  CloseEventCodeSessionNoLongerValid,
		CloseEventCodeSessionTimeout.Code:        CloseEventCodeSessionTimeout,
		CloseEventCodeServerNotFound.Code:        CloseEventCodeServerNotFound,
		CloseEventCodeUnknownProtocol.Code:       CloseEventCodeUnknownProtocol,
		CloseEventCodeDisconnected.Code:          CloseEventCodeDisconnected,
		CloseEventCodeVoiceServerCrash.Code:      CloseEventCodeVoiceServerCrash,
		CloseEventCodeUnknownEncryptionMode.Code: CloseEventCodeUnknownEncryptionMode,
	}
)

func CloseEventCodeByCode(code int) CloseEventCode {
	closeCode, ok := CloseEventCodes[code]
	if !ok {
		return CloseEventCodeUnknown
	}
	return closeCode
}
