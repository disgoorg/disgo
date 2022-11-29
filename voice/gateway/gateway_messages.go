package gateway

import (
	"errors"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type Message struct {
	Op Opcode      `json:"op"`
	D  MessageData `json:"d,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var v struct {
		Op Opcode          `json:"op"`
		D  json.RawMessage `json:"d"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var (
		messageData MessageData
		err         error
	)

	switch v.Op {
	case OpcodeIdentify:
		var d MessageDataIdentify
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeSelectProtocol:
		var d MessageDataSelectProtocol
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeReady:
		var d MessageDataReady
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHeartbeat:
		var d MessageDataHeartbeat
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeSessionDescription:
		var d MessageDataSessionDescription
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeSpeaking:
		var d MessageDataSpeaking
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHeartbeatACK:
		var d MessageDataHeartbeatACK
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeResume:
		var d MessageDataResume
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHello:
		var d MessageDataHello
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeResumed:
		// no data

	case OpcodeClientDisconnect:
		var d MessageDataClientDisconnect
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeGuildSync:
		// ignore this opcode

	default:
		err = errors.New("unknown gateway event type")
	}
	if err != nil {
		return err
	}
	m.Op = v.Op
	m.D = messageData
	return nil
}

type MessageData interface {
	voiceMessageData()
}

type MessageDataIdentify struct {
	GuildID   snowflake.ID `json:"server_id"`
	UserID    snowflake.ID `json:"user_id"`
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
}

func (MessageDataIdentify) voiceMessageData() {}

type MessageDataReady struct {
	SSRC  uint32   `json:"ssrc"`
	IP    string   `json:"ip"`
	Port  int      `json:"port"`
	Modes []string `json:"modes"`
}

func (MessageDataReady) voiceMessageData() {}

type MessageDataHello struct {
	HeartbeatInterval float64 `json:"heartbeat_interval"`
}

func (MessageDataHello) voiceMessageData() {}

type MessageDataHeartbeat int64

func (MessageDataHeartbeat) voiceMessageData() {}

type MessageDataSessionDescription struct {
	Mode      string   `json:"mode"`
	SecretKey [32]byte `json:"secret_key"`
}

func (MessageDataSessionDescription) voiceMessageData() {}

type VoiceProtocol string

const (
	VoiceProtocolUDP VoiceProtocol = "udp"
)

type MessageDataSelectProtocol struct {
	Protocol VoiceProtocol                 `json:"protocol"`
	Data     MessageDataSelectProtocolData `json:"data"`
}

func (MessageDataSelectProtocol) voiceMessageData() {}

type MessageDataSelectProtocolData struct {
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

type MessageDataSpeaking struct {
	Speaking SpeakingFlags `json:"speaking"`
	Delay    int           `json:"delay"`
	SSRC     uint32        `json:"ssrc"`
	UserID   snowflake.ID  `json:"user_id,omitempty"`
}

func (MessageDataSpeaking) voiceMessageData() {}

type SpeakingFlags int

const (
	SpeakingFlagMicrophone SpeakingFlags = 1 << iota
	SpeakingFlagSoundshare
	SpeakingFlagPriority
	SpeakingFlagNone SpeakingFlags = 0
)

type MessageDataResume struct {
	GuildID   snowflake.ID `json:"server_id"` // wtf is this?
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
}

func (MessageDataResume) voiceMessageData() {}

type MessageDataHeartbeatACK int64

func (MessageDataHeartbeatACK) voiceMessageData() {}

type MessageDataClientConnect struct {
	UserID     snowflake.ID `json:"user_id"`
	AudioCodec string       `json:"audio_codec"`
	VideoCodec string       `json:"video_codec"`
}

func (MessageDataClientConnect) voiceMessageData() {}

type MessageDataClientDisconnect struct {
	UserID snowflake.ID `json:"user_id"`
}

func (MessageDataClientDisconnect) voiceMessageData() {}
