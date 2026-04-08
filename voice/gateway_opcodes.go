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
	OpcodeClientsConnect
	_
	OpcodeClientDisconnect
	OpcodeGuildSync
	_
	_
	_
	_
	_
	_
	OpcodeDavePrepareTransition
	OpcodeDaveExecuteTransition
	OpcodeDaveTransitionReady
	OpcodeDavePrepareEpoch
	OpcodeDaveMLSExternalSenderPackage
	OpcodeDaveMLSKeyPackage
	OpcodeDaveMLSProposals
	OpcodeDaveMLSCommitWelcome
	OpcodeDaveMLSPrepareCommitTransition
	OpcodeDaveMLSWelcome
	OpcodeDaveMLSInvalidCommitWelcome
)

func (o Opcode) IsBinary() bool {
	switch o {
	case OpcodeDavePrepareTransition,
		OpcodeDaveMLSKeyPackage,
		OpcodeDaveMLSProposals,
		OpcodeDaveMLSCommitWelcome,
		OpcodeDaveMLSPrepareCommitTransition,
		OpcodeDaveMLSWelcome:
		return true
	default:
		return false
	}
}

type GatewayCloseEventCode struct {
	Code        int
	Description string
	Explanation string
	// Reconnect indicated if we should reconnect to the voice gateway with the same session id, endpoint, token & channel id.
	Reconnect bool
	// NewConnection indicates if we should start a new connection when the error is received.
	NewConnection bool
}

var (
	GatewayCloseEventCodeHeartbeatTimeout = GatewayCloseEventCode{
		Code:          1000,
		Description:   "Heartbeat timeout",
		Explanation:   "We did not heartbeat in time.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeUnknownError = GatewayCloseEventCode{
		Code:          4000,
		Description:   "Unknown error",
		Explanation:   "We're not sure what went wrong. Try reconnecting?",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeUnknownOpcode = GatewayCloseEventCode{
		Code:          4001,
		Description:   "Unknown opcode",
		Explanation:   "You sent an invalid opcode.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeFailedDecode = GatewayCloseEventCode{
		Code:          4002,
		Description:   "Failed to decode payload",
		Explanation:   "You sent a invalid payload in your identifying to the Gateway.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeNotAuthenticated = GatewayCloseEventCode{
		Code:          4003,
		Description:   "Not authenticated",
		Explanation:   "You sent a payload before identifying with the Gateway.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeAuthenticationFailed = GatewayCloseEventCode{
		Code:          4004,
		Description:   "Authentication failed",
		Explanation:   "The token you sent in your identify payload is incorrect.",
		Reconnect:     false,
		NewConnection: false,
	}

	GatewayCloseEventCodeAlreadyAuthenticated = GatewayCloseEventCode{
		Code:          4005,
		Description:   "Already authenticated",
		Explanation:   "You sent more than one identify payload. Stahp.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeSessionNoLongerValid = GatewayCloseEventCode{
		Code:          4006,
		Description:   "Session no longer valid",
		Explanation:   "Your session is no longer valid.",
		Reconnect:     false,
		NewConnection: true,
	}

	GatewayCloseEventCodeSessionTimeout = GatewayCloseEventCode{
		Code:          4009,
		Description:   "Session timeout",
		Explanation:   "Your session has timed out.",
		Reconnect:     false,
		NewConnection: true,
	}

	GatewayCloseEventCodeServerNotFound = GatewayCloseEventCode{
		Code:          4011,
		Description:   "Server not found",
		Explanation:   "We can't find the server you're trying to connect to.",
		Reconnect:     false,
		NewConnection: true,
	}

	GatewayCloseEventCodeUnknownProtocol = GatewayCloseEventCode{
		Code:          4012,
		Description:   "Unknown protocol",
		Explanation:   "We didn't recognize the protocol you sent.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeDisconnected = GatewayCloseEventCode{
		Code:          4014,
		Description:   "Disconnected",
		Explanation:   "Disconnect individual client (you were kicked, the main gateway session was dropped, etc.). Should not reconnect.",
		Reconnect:     false,
		NewConnection: false,
	}

	GatewayCloseEventCodeVoiceServerCrash = GatewayCloseEventCode{
		Code:          4015,
		Description:   "Voice server crashed",
		Explanation:   "The server crashed. Our bad! Try resuming.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeUnknownEncryptionMode = GatewayCloseEventCode{
		Code:          4016,
		Description:   "Unknown encryption mode",
		Explanation:   "We didn't recognize your encryption.",
		Reconnect:     false,
		NewConnection: false,
	}

	GatewayCloseEventCodeEndtoEndEncryptionDAVEProtocolRequired = GatewayCloseEventCode{
		Code:          4017,
		Description:   "E2EE/DAVE protocol required",
		Explanation:   "This channel requires a client supporting E2EE via the DAVE Protocol.",
		Reconnect:     false,
		NewConnection: false,
	}

	GatewayCloseEventCodeBadRequest = GatewayCloseEventCode{
		Code:          4020,
		Description:   "Bad request",
		Explanation:   "You sent a malformed request.",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodeRateLimited = GatewayCloseEventCode{
		Code:          4021,
		Description:   "Disconnected: Rate Limited",
		Explanation:   "Disconnect due to rate limit exceeded. Should not reconnect.",
		Reconnect:     false,
		NewConnection: true,
	}

	GatewayCloseEventCodeCallTerminated = GatewayCloseEventCode{
		Code:          4022,
		Description:   "Disconnected: Call Terminated",
		Explanation:   "Disconnect all clients due to call terminated (channel deleted, voice server changed, etc.). Should not reconnect.",
		Reconnect:     false,
		NewConnection: true,
	}

	GatewayCloseEventCodeUnknown = GatewayCloseEventCode{
		Code:          0,
		Description:   "Unknown",
		Explanation:   "Unknown Voice Close Event Code",
		Reconnect:     true,
		NewConnection: true,
	}

	GatewayCloseEventCodes = map[int]GatewayCloseEventCode{
		GatewayCloseEventCodeHeartbeatTimeout.Code:                       GatewayCloseEventCodeHeartbeatTimeout,
		GatewayCloseEventCodeUnknownError.Code:                           GatewayCloseEventCodeUnknownError,
		GatewayCloseEventCodeUnknownOpcode.Code:                          GatewayCloseEventCodeUnknownOpcode,
		GatewayCloseEventCodeFailedDecode.Code:                           GatewayCloseEventCodeFailedDecode,
		GatewayCloseEventCodeNotAuthenticated.Code:                       GatewayCloseEventCodeNotAuthenticated,
		GatewayCloseEventCodeAuthenticationFailed.Code:                   GatewayCloseEventCodeAuthenticationFailed,
		GatewayCloseEventCodeAlreadyAuthenticated.Code:                   GatewayCloseEventCodeAlreadyAuthenticated,
		GatewayCloseEventCodeSessionNoLongerValid.Code:                   GatewayCloseEventCodeSessionNoLongerValid,
		GatewayCloseEventCodeSessionTimeout.Code:                         GatewayCloseEventCodeSessionTimeout,
		GatewayCloseEventCodeServerNotFound.Code:                         GatewayCloseEventCodeServerNotFound,
		GatewayCloseEventCodeUnknownProtocol.Code:                        GatewayCloseEventCodeUnknownProtocol,
		GatewayCloseEventCodeDisconnected.Code:                           GatewayCloseEventCodeDisconnected,
		GatewayCloseEventCodeVoiceServerCrash.Code:                       GatewayCloseEventCodeVoiceServerCrash,
		GatewayCloseEventCodeUnknownEncryptionMode.Code:                  GatewayCloseEventCodeUnknownEncryptionMode,
		GatewayCloseEventCodeEndtoEndEncryptionDAVEProtocolRequired.Code: GatewayCloseEventCodeEndtoEndEncryptionDAVEProtocolRequired,
		GatewayCloseEventCodeBadRequest.Code:                             GatewayCloseEventCodeBadRequest,
		GatewayCloseEventCodeRateLimited.Code:                            GatewayCloseEventCodeRateLimited,
		GatewayCloseEventCodeCallTerminated.Code:                         GatewayCloseEventCodeCallTerminated,
	}
)

func GatewayCloseEventCodeByCode(code int) GatewayCloseEventCode {
	closeCode, ok := GatewayCloseEventCodes[code]
	if !ok {
		return GatewayCloseEventCodeUnknown
	}
	return closeCode
}
