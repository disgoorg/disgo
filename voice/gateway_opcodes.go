package voice

type GatewayOpcode int

const (
	GatewayOpcodeIdentify GatewayOpcode = iota
	GatewayOpcodeSelectProtocol
	GatewayOpcodeReady
	GatewayOpcodeHeartbeat
	GatewayOpcodeSessionDescription
	GatewayOpcodeSpeaking
	GatewayOpcodeHeartbeatACK
	GatewayOpcodeResume
	GatewayOpcodeHello
	GatewayOpcodeResumed
	_
	_
	GatewayOpcodeClientConnect
	GatewayOpcodeClientDisconnect
)

type GatewayCloseEventCode int

const (
	GatewayCloseEventCodeUnknownOpcode GatewayCloseEventCode = iota + 4001
	GatewayCloseEventCodeDecodeError
	GatewayCloseEventCodeNotAuthenticated
	GatewayCloseEventCodeAuthenticationFailed
	GatewayCloseEventCodeAlreadyAuthenticated
	GatewayCloseEventCodeSessionNoLongerValid
	_
	_
	GatewayCloseEventCodeSessionTimedOut
	_
	GatewayCloseEventCodeServerNotFound
	GatewayCloseEventCodeUnknownProtocol
	_
	GatewayCloseEventCodeDisconnected
	GatewayCloseEventCodeGatewayServerCrashed
	GatewayCloseEventCodeUnknownEncryptionMode
)
