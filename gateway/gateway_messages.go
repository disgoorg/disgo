package gateway

import (
	"errors"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake/v2"
)

// Message raw Message type
type Message struct {
	Op Opcode      `json:"op"`
	S  int         `json:"s,omitempty"`
	T  EventType   `json:"t,omitempty"`
	D  MessageData `json:"d,omitempty"`
}

func (e *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v struct {
		D json.RawMessage `json:"d,omitempty"`
		message
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var (
		messageData MessageData
		err         error
	)

	switch v.Op {
	case OpcodeDispatch:
		messageData = MessageDataDispatch(v.D)

	case OpcodeHeartbeat:
		var d MessageDataHeartbeat
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeIdentify:
		var d MessageDataIdentify
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodePresenceUpdate:
		var d MessageDataPresenceUpdate
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeVoiceStateUpdate:
		var d MessageDataVoiceStateUpdate
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeResume:
		var d MessageDataResume
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeReconnect:
		// no data

	case OpcodeRequestGuildMembers:
		var d MessageDataRequestGuildMembers
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeInvalidSession:
		var d MessageDataInvalidSession
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHello:
		var d MessageDataHello
		err = json.Unmarshal(v.D, &d)
		messageData = d

	case OpcodeHeartbeatACK:
		// no data

	default:
		err = errors.New("unknown gateway event type")
	}
	if err != nil {
		return err
	}
	*e = Message(v.message)
	e.D = messageData
	return nil
}

type MessageData interface {
	gatewayMessageData()
}

type MessageDataDispatch json.RawMessage

func (MessageDataDispatch) gatewayMessageData() {}

// MessageDataHeartbeat is used to ensure the websocket connection remains open, and disconnect if not.
type MessageDataHeartbeat int

func (MessageDataHeartbeat) gatewayMessageData() {}

// MessageDataIdentify is the data used in IdentifyCommandData
type MessageDataIdentify struct {
	Token          string                        `json:"token"`
	Properties     IdentifyCommandDataProperties `json:"properties"`
	Compress       bool                          `json:"compress,omitempty"`
	LargeThreshold int                           `json:"large_threshold,omitempty"`
	Shard          *[2]int                       `json:"shard,omitempty"`
	Intents        Intents                       `json:"intents"`
	Presence       *MessageDataPresenceUpdate    `json:"presence,omitempty"`
}

func (MessageDataIdentify) gatewayMessageData() {}

// IdentifyCommandDataProperties is used for specifying to discord which library and OS the bot is using, is
// automatically handled by the library and should rarely be used.
type IdentifyCommandDataProperties struct {
	OS      string `json:"os"`      // user OS
	Browser string `json:"browser"` // library name
	Device  string `json:"device"`  // library name
}

// NewPresence creates a new Presence with the provided properties
func NewPresence(activityType discord.ActivityType, name string, url string, status discord.OnlineStatus, afk bool) MessageDataPresenceUpdate {
	var since *int64
	if status == discord.OnlineStatusIdle {
		unix := time.Now().Unix()
		since = &unix
	}

	var activities []discord.Activity
	if name != "" {
		activity := discord.Activity{
			Name: name,
			Type: activityType,
		}
		if activityType == discord.ActivityTypeStreaming && url != "" {
			activity.URL = &url
		}
		activities = append(activities, activity)
	}

	return MessageDataPresenceUpdate{
		Since:      since,
		Activities: activities,
		Status:     status,
		AFK:        afk,
	}
}

// NewGamePresence creates a new Presence of type ActivityTypeGame
func NewGamePresence(name string, status discord.OnlineStatus, afk bool) MessageDataPresenceUpdate {
	return NewPresence(discord.ActivityTypeGame, name, "", status, afk)
}

// NewStreamingPresence creates a new Presence of type ActivityTypeStreaming
func NewStreamingPresence(name string, url string, status discord.OnlineStatus, afk bool) MessageDataPresenceUpdate {
	return NewPresence(discord.ActivityTypeStreaming, name, url, status, afk)
}

// NewListeningPresence creates a new Presence of type ActivityTypeListening
func NewListeningPresence(name string, status discord.OnlineStatus, afk bool) MessageDataPresenceUpdate {
	return NewPresence(discord.ActivityTypeListening, name, "", status, afk)
}

// NewWatchingPresence creates a new Presence of type ActivityTypeWatching
func NewWatchingPresence(name string, status discord.OnlineStatus, afk bool) MessageDataPresenceUpdate {
	return NewPresence(discord.ActivityTypeWatching, name, "", status, afk)
}

// NewCompetingPresence creates a new Presence of type ActivityTypeCompeting
func NewCompetingPresence(name string, status discord.OnlineStatus, afk bool) MessageDataPresenceUpdate {
	return NewPresence(discord.ActivityTypeCompeting, name, "", status, afk)
}

// MessageDataPresenceUpdate is used for updating Client's presence
type MessageDataPresenceUpdate struct {
	Since      *int64               `json:"since"`
	Activities []discord.Activity   `json:"activities"`
	Status     discord.OnlineStatus `json:"status"`
	AFK        bool                 `json:"afk"`
}

func (MessageDataPresenceUpdate) gatewayMessageData() {}

// MessageDataVoiceStateUpdate is used for updating the bots voice state in a guild
type MessageDataVoiceStateUpdate struct {
	GuildID   snowflake.ID  `json:"guild_id"`
	ChannelID *snowflake.ID `json:"channel_id"`
	SelfMute  bool          `json:"self_mute"`
	SelfDeaf  bool          `json:"self_deaf"`
}

func (MessageDataVoiceStateUpdate) gatewayMessageData() {}

// MessageDataResume is used to resume a connection to discord in the case that you are disconnected. Is automatically
// handled by the library and should rarely be used.
type MessageDataResume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

func (MessageDataResume) gatewayMessageData() {}

// MessageDataRequestGuildMembers is used for fetching all the members of a guild_events. It is recommended you have a strict
// member caching policy when using this.
type MessageDataRequestGuildMembers struct {
	GuildID   snowflake.ID   `json:"guild_id"`
	Query     *string        `json:"query,omitempty"` //If specified, user_ids must not be entered
	Limit     *int           `json:"limit,omitempty"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool           `json:"presences,omitempty"`
	UserIDs   []snowflake.ID `json:"user_ids,omitempty"` //If specified, query must not be entered
	Nonce     string         `json:"nonce,omitempty"`    //All responses are hashed with this nonce, optional
}

func (MessageDataRequestGuildMembers) gatewayMessageData() {}

type MessageDataInvalidSession bool

func (MessageDataInvalidSession) gatewayMessageData() {}

type MessageDataHello struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

func (MessageDataHello) gatewayMessageData() {}
