package voice

import (
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

// GatewayMessage represents a voice gateway message
type GatewayMessage struct {
	Op Opcode             `json:"op"`
	D  GatewayMessageData `json:"d,omitempty"`
}

// UnmarshalJSON unmarshalls the GatewayMessage from json
func (m *GatewayMessage) UnmarshalJSON(data []byte) error {
	var v struct {
		Op Opcode          `json:"op"`
		D  json.RawMessage `json:"d"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var (
		messageData GatewayMessageData
		err         error
	)

	switch v.Op {
	case OpcodeIdentify:
		var d GatewayMessageDataIdentify
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeSelectProtocol:
		var d GatewayMessageDataSelectProtocol
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeReady:
		var d GatewayMessageDataReady
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHeartbeat:
		var d GatewayMessageDataHeartbeat
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeSessionDescription:
		var d GatewayMessageDataSessionDescription
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeSpeaking:
		var d GatewayMessageDataSpeaking
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHeartbeatACK:
		var d GatewayMessageDataHeartbeatACK
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeResume:
		var d GatewayMessageDataResume
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHello:
		var d GatewayMessageDataHello
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeResumed:
		// no data

	case OpcodeClientDisconnect:
		var d GatewayMessageDataClientDisconnect
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeGuildSync:
		// ignore this opcode

	default:
		var d GatewayMessageDataUnknown
		err = json.Unmarshal(v.D, &d)
		messageData = d
	}
	if err != nil {
		return err
	}
	m.Op = v.Op
	m.D = messageData
	return nil
}

// GatewayMessageData represents a voice gateway message data.
type GatewayMessageData interface {
	voiceGatewayMessageData()
}

type GatewayMessageDataIdentify struct {
	GuildID   snowflake.ID `json:"server_id"`
	UserID    snowflake.ID `json:"user_id"`
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
}

func (GatewayMessageDataIdentify) voiceGatewayMessageData() {}

type GatewayMessageDataReady struct {
	SSRC  uint32   `json:"ssrc"`
	IP    string   `json:"ip"`
	Port  int      `json:"port"`
	Modes []string `json:"modes"`
}

func (GatewayMessageDataReady) voiceGatewayMessageData() {}

type GatewayMessageDataHello struct {
	HeartbeatInterval float64 `json:"heartbeat_interval"`
}

func (GatewayMessageDataHello) voiceGatewayMessageData() {}

type GatewayMessageDataHeartbeat int64

func (GatewayMessageDataHeartbeat) voiceGatewayMessageData() {}

type GatewayMessageDataSessionDescription struct {
	Mode      string   `json:"mode"`
	SecretKey [32]byte `json:"secret_key"`
}

func (GatewayMessageDataSessionDescription) voiceGatewayMessageData() {}

type Protocol string

const (
	ProtocolUDP Protocol = "udp"
)

type GatewayMessageDataSelectProtocol struct {
	Protocol Protocol                             `json:"protocol"`
	Data     GatewayMessageDataSelectProtocolData `json:"data"`
}

func (GatewayMessageDataSelectProtocol) voiceGatewayMessageData() {}

type GatewayMessageDataSelectProtocolData struct {
	Address string         `json:"address"`
	Port    int            `json:"port"`
	Mode    EncryptionMode `json:"mode"`
}

// EncryptionMode is the encryption mode used for voice data.
type EncryptionMode string

// All possible EncryptionMode(s) https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection-encryption-modes.
const (
	EncryptionModeNormal EncryptionMode = "xsalsa20_poly1305"
	EncryptionModeSuffix EncryptionMode = "xsalsa20_poly1305_suffix"
	EncryptionModeLite   EncryptionMode = "xsalsa20_poly1305_lite"
)

type GatewayMessageDataSpeaking struct {
	Speaking SpeakingFlags `json:"speaking"`
	Delay    int           `json:"delay"`
	SSRC     uint32        `json:"ssrc"`
	UserID   snowflake.ID  `json:"user_id,omitempty"`
}

func (GatewayMessageDataSpeaking) voiceGatewayMessageData() {}

type SpeakingFlags int

const (
	SpeakingFlagMicrophone SpeakingFlags = 1 << iota
	SpeakingFlagSoundshare
	SpeakingFlagPriority
	SpeakingFlagNone SpeakingFlags = 0
)

type GatewayMessageDataResume struct {
	GuildID   snowflake.ID `json:"server_id"` // wtf is this?
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
}

func (GatewayMessageDataResume) voiceGatewayMessageData() {}

type GatewayMessageDataHeartbeatACK int64

func (GatewayMessageDataHeartbeatACK) voiceGatewayMessageData() {}

type GatewayMessageDataClientConnect struct {
	UserID     snowflake.ID `json:"user_id"`
	AudioCodec string       `json:"audio_codec"`
	VideoCodec string       `json:"video_codec"`
}

func (GatewayMessageDataClientConnect) voiceGatewayMessageData() {}

type GatewayMessageDataClientDisconnect struct {
	UserID snowflake.ID `json:"user_id"`
}

func (GatewayMessageDataClientDisconnect) voiceGatewayMessageData() {}

type GatewayMessageDataUnknown json.RawMessage

func (GatewayMessageDataUnknown) voiceGatewayMessageData() {}

func (m GatewayMessageDataUnknown) MarshalJSON() ([]byte, error) {
	return json.RawMessage(m).MarshalJSON()
}

func (m *GatewayMessageDataUnknown) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(m).UnmarshalJSON(data)
}
