package voice

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

type GatewayCloseEventCode struct {
	Code        int
	Description string
	Explanation string
	Reconnect   bool
}

var (
	GatewayCloseEventCodeHeartbeatTimeout = GatewayCloseEventCode{
		Code:        1000,
		Description: "Heartbeat timeout",
		Explanation: "We did not heartbeat in time.",
		Reconnect:   true,
	}

	GatewayCloseEventCodeUnknownError = GatewayCloseEventCode{
		Code:        4000,
		Description: "Unknown error",
		Explanation: "We're not sure what went wrong. Try reconnecting?",
		Reconnect:   true,
	}

	GatewayCloseEventCodeUnknownOpcode = GatewayCloseEventCode{
		Code:        4001,
		Description: "Unknown opcode",
		Explanation: "You sent an invalid opcode.",
		Reconnect:   true,
	}

	GatewayCloseEventCodeFailedDecode = GatewayCloseEventCode{
		Code:        4002,
		Description: "Failed to decode payload",
		Explanation: "You sent a invalid payload in your identifying to the Gateway.",
		Reconnect:   true,
	}

	GatewayCloseEventCodeNotAuthenticated = GatewayCloseEventCode{
		Code:        4003,
		Description: "Not authenticated",
		Explanation: "You sent a payload before identifying with the Gateway.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeAuthenticationFailed = GatewayCloseEventCode{
		Code:        4004,
		Description: "Authentication failed",
		Explanation: "The token you sent in your identify payload is incorrect.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeAlreadyAuthenticated = GatewayCloseEventCode{
		Code:        4005,
		Description: "Already authenticated",
		Explanation: "You sent more than one identify payload. Stahp.",
		Reconnect:   true,
	}

	GatewayCloseEventCodeSessionNoLongerValid = GatewayCloseEventCode{
		Code:        4006,
		Description: "Session no longer valid",
		Explanation: "Your session is no longer valid.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeSessionTimeout = GatewayCloseEventCode{
		Code:        4009,
		Description: "Session timeout",
		Explanation: "Your session has timed out.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeServerNotFound = GatewayCloseEventCode{
		Code:        4011,
		Description: "Server not found",
		Explanation: "We can't find the server you're trying to connect to.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeUnknownProtocol = GatewayCloseEventCode{
		Code:        4012,
		Description: "Unknown protocol",
		Explanation: "We didn't recognize the protocol you sent.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeDisconnected = GatewayCloseEventCode{
		Code:        4014,
		Description: "Disconnected",
		Explanation: "Channel was deleted, you were kicked, voice server changed, or the main voicegateway session was dropped. Don't reconnect.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeVoiceServerCrash = GatewayCloseEventCode{
		Code:        4015,
		Description: "Voice server crashed",
		Explanation: "The server crashed. Our bad! Try resuming.",
		Reconnect:   true,
	}

	GatewayCloseEventCodeUnknownEncryptionMode = GatewayCloseEventCode{
		Code:        4016,
		Description: "Unknown encryption mode",
		Explanation: "We didn't recognize your encryption.",
		Reconnect:   false,
	}

	GatewayCloseEventCodeUnknown = GatewayCloseEventCode{
		Code:        0,
		Description: "Unknown",
		Explanation: "Unknown Voice Close Event Code",
		Reconnect:   true,
	}

	GatewayCloseEventCodes = map[int]GatewayCloseEventCode{
		GatewayCloseEventCodeHeartbeatTimeout.Code:      GatewayCloseEventCodeHeartbeatTimeout,
		GatewayCloseEventCodeUnknownError.Code:          GatewayCloseEventCodeUnknownError,
		GatewayCloseEventCodeUnknownOpcode.Code:         GatewayCloseEventCodeUnknownOpcode,
		GatewayCloseEventCodeFailedDecode.Code:          GatewayCloseEventCodeFailedDecode,
		GatewayCloseEventCodeNotAuthenticated.Code:      GatewayCloseEventCodeNotAuthenticated,
		GatewayCloseEventCodeAuthenticationFailed.Code:  GatewayCloseEventCodeAuthenticationFailed,
		GatewayCloseEventCodeAlreadyAuthenticated.Code:  GatewayCloseEventCodeAlreadyAuthenticated,
		GatewayCloseEventCodeSessionNoLongerValid.Code:  GatewayCloseEventCodeSessionNoLongerValid,
		GatewayCloseEventCodeSessionTimeout.Code:        GatewayCloseEventCodeSessionTimeout,
		GatewayCloseEventCodeServerNotFound.Code:        GatewayCloseEventCodeServerNotFound,
		GatewayCloseEventCodeUnknownProtocol.Code:       GatewayCloseEventCodeUnknownProtocol,
		GatewayCloseEventCodeDisconnected.Code:          GatewayCloseEventCodeDisconnected,
		GatewayCloseEventCodeVoiceServerCrash.Code:      GatewayCloseEventCodeVoiceServerCrash,
		GatewayCloseEventCodeUnknownEncryptionMode.Code: GatewayCloseEventCodeUnknownEncryptionMode,
	}
)

func GatewayCloseEventCodeByCode(code int) GatewayCloseEventCode {
	closeCode, ok := GatewayCloseEventCodes[code]
	if !ok {
		return GatewayCloseEventCodeUnknown
	}
	return closeCode
}
