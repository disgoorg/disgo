package discord

import (
	"errors"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake/v2"
)

// GatewayMessage raw GatewayMessage type
type GatewayMessage struct {
	Op GatewayOpcode      `json:"op"`
	S  int                `json:"s,omitempty"`
	T  GatewayEventType   `json:"t,omitempty"`
	D  GatewayMessageData `json:"d,omitempty"`
}

func (e *GatewayMessage) UnmarshalJSON(data []byte) error {
	type gatewayMessage GatewayMessage
	var v struct {
		D json.RawMessage `json:"d,omitempty"`
		gatewayMessage
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var (
		messageData GatewayMessageData
		err         error
	)

	switch v.Op {
	case GatewayOpcodeDispatch:
		messageData = GatewayMessageDataDispatch(v.D)

	case GatewayOpcodeHeartbeat:
		var d GatewayMessageDataHeartbeat
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeIdentify:
		var d GatewayMessageDataIdentify
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodePresenceUpdate:
		var d GatewayMessageDataPresenceUpdate
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeVoiceStateUpdate:
		var d GatewayMessageDataVoiceStateUpdate
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeResume:
		var d GatewayMessageDataResume
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeReconnect:
		// no data

	case GatewayOpcodeRequestGuildMembers:
		var d GatewayMessageDataRequestGuildMembers
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeInvalidSession:
		var d GatewayMessageDataInvalidSession
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeHello:
		var d GatewayMessageDataHello
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case GatewayOpcodeHeartbeatACK:
		// no data

	default:
		err = errors.New("unknown gateway event type")
	}
	if err != nil {
		return err
	}
	*e = GatewayMessage(v.gatewayMessage)
	e.D = messageData
	return nil
}

type GatewayMessageData interface {
	gatewayMessageData()
}

type GatewayMessageDataDispatch json.RawMessage

func (GatewayMessageDataDispatch) gatewayMessageData() {}

// GatewayMessageDataHeartbeat is used to ensure the websocket connection remains open, and disconnect if not.
type GatewayMessageDataHeartbeat int

func (GatewayMessageDataHeartbeat) gatewayMessageData() {}

// GatewayMessageDataIdentify is the data used in IdentifyCommandData
type GatewayMessageDataIdentify struct {
	Token          string                            `json:"token"`
	Properties     IdentifyCommandDataProperties     `json:"properties"`
	Compress       bool                              `json:"compress,omitempty"`
	LargeThreshold int                               `json:"large_threshold,omitempty"`
	Shard          *[2]int                           `json:"shard,omitempty"`
	GatewayIntents GatewayIntents                    `json:"intents"`
	Presence       *GatewayMessageDataPresenceUpdate `json:"presence,omitempty"`
}

func (GatewayMessageDataIdentify) gatewayMessageData() {}

// IdentifyCommandDataProperties is used for specifying to discord which library and OS the bot is using, is
// automatically handled by the library and should rarely be used.
type IdentifyCommandDataProperties struct {
	OS      string `json:"os"`      // user OS
	Browser string `json:"browser"` // library name
	Device  string `json:"device"`  // library name
}

// GatewayMessageDataPresenceUpdate is used for updating Client's presence
type GatewayMessageDataPresenceUpdate struct {
	Since      *int64       `json:"since"`
	Activities []Activity   `json:"activities"`
	Status     OnlineStatus `json:"status"`
	AFK        bool         `json:"afk"`
}

func (GatewayMessageDataPresenceUpdate) gatewayMessageData() {}

// GatewayMessageDataVoiceStateUpdate is used for updating the bots voice state in a guild
type GatewayMessageDataVoiceStateUpdate struct {
	GuildID   snowflake.ID  `json:"guild_id"`
	ChannelID *snowflake.ID `json:"channel_id"`
	SelfMute  bool          `json:"self_mute"`
	SelfDeaf  bool          `json:"self_deaf"`
}

func (GatewayMessageDataVoiceStateUpdate) gatewayMessageData() {}

// GatewayMessageDataResume is used to resume a connection to discord in the case that you are disconnected. Is automatically
// handled by the library and should rarely be used.
type GatewayMessageDataResume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

func (GatewayMessageDataResume) gatewayMessageData() {}

// GatewayMessageDataRequestGuildMembers is used for fetching all the members of a guild_events. It is recommended you have a strict
// member caching policy when using this.
type GatewayMessageDataRequestGuildMembers struct {
	GuildID   snowflake.ID   `json:"guild_id"`
	Query     *string        `json:"query,omitempty"` //If specified, user_ids must not be entered
	Limit     *int           `json:"limit,omitempty"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool           `json:"presences,omitempty"`
	UserIDs   []snowflake.ID `json:"user_ids,omitempty"` //If specified, query must not be entered
	Nonce     string         `json:"nonce,omitempty"`    //All responses are hashed with this nonce, optional
}

func (GatewayMessageDataRequestGuildMembers) gatewayMessageData() {}

type GatewayMessageDataInvalidSession bool

func (GatewayMessageDataInvalidSession) gatewayMessageData() {}

type GatewayMessageDataHello struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

func (GatewayMessageDataHello) gatewayMessageData() {}
