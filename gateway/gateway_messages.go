package gateway

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake/v2"
)

// Message raw Message type
type Message struct {
	Op   Opcode          `json:"op"`
	S    int             `json:"s,omitempty"`
	T    EventType       `json:"t,omitempty"`
	D    MessageData     `json:"d,omitempty"`
	RawD json.RawMessage `json:"-"`
}

func (e *Message) UnmarshalJSON(data []byte) error {
	var v struct {
		Op Opcode    `json:"op"`
		S  int       `json:"s,omitempty"`
		T  EventType `json:"t,omitempty"`
		D  json.RawMessage
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
		switch v.T {
		case EventTypeReady:
			var d EventReady
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeResumed:
			// no data

		case EventTypeApplicationCommandPermissionsUpdate:
			var d EventApplicationCommandPermissionsUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeAutoModerationRuleCreate:
			var d EventAutoModerationRuleCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeAutoModerationRuleUpdate:
			var d EventAutoModerationRuleUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeAutoModerationRuleDelete:
			var d EventAutoModerationRuleDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeAutoModerationActionExecution:
			var d EventAutoModerationActionExecution
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeChannelCreate:
			var d EventChannelCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeChannelUpdate:
			var d EventChannelUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeChannelDelete:
			var d EventChannelDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeChannelPinsUpdate:
			var d EventChannelPinsUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeThreadCreate:
			var d EventThreadCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeThreadUpdate:
			var d EventThreadUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeThreadDelete:
			var d EventThreadDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeThreadListSync:
			var d EventThreadListSync
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeThreadMemberUpdate:
			var d EventThreadMemberUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeThreadMembersUpdate:
			var d EventThreadMembersUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildCreate:
			var d EventGuildCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildUpdate:
			var d EventGuildUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildDelete:
			var d EventGuildDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildBanAdd:
			var d EventGuildBanAdd
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildBanRemove:
			var d EventGuildBanRemove
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildEmojisUpdate:
			var d EventGuildEmojisUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildStickersUpdate:
			var d EventGuildStickersUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildIntegrationsUpdate:
			var d EventGuildIntegrationsUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildMemberAdd:
			var d EventGuildMemberAdd
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildMemberUpdate:
			var d EventGuildMemberUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildMembersChunk:
			var d EventGuildMembersChunk
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildRoleCreate:
			var d EventGuildRoleCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildRoleUpdate:
			var d EventGuildRoleUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildRoleDelete:
			var d EventGuildRoleDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildScheduledEventCreate:
			var d EventGuildScheduledEventCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildScheduledEventUpdate:
			var d EventGuildScheduledEventUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildScheduledEventDelete:
			var d EventGuildScheduledEventDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildScheduledEventUserAdd:
			var d EventGuildScheduledEventUserAdd
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeGuildScheduledEventUserRemove:
			var d EventGuildScheduledEventUserRemove
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeIntegrationCreate:
			var d EventIntegrationCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeIntegrationUpdate:
			var d EventIntegrationUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeIntegrationDelete:
			var d EventIntegrationDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeInteractionCreate:
			var d EventInteractionCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeInviteCreate:
			var d EventInviteCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeInviteDelete:
			var d EventInviteDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageCreate:
			var d EventMessageCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageUpdate:
			var d EventMessageUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageDelete:
			var d EventMessageDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageDeleteBulk:
			var d EventMessageDeleteBulk
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageReactionAdd:
			var d EventMessageReactionAdd
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageReactionRemove:
			var d EventMessageReactionRemove
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageReactionRemoveAll:
			var d EventMessageReactionRemoveAll
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeMessageReactionRemoveEmoji:
			var d EventMessageReactionRemoveEmoji
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypePresenceUpdate:
			var d EventPresenceUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeStageInstanceCreate:
			var d EventStageInstanceCreate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeStageInstanceUpdate:
			var d EventStageInstanceUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeStageInstanceDelete:
			var d EventStageInstanceDelete
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeTypingStart:
			var d EventTypingStart
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeUserUpdate:
			var d EventUserUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeVoiceStateUpdate:
			var d EventVoiceStateUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeVoiceServerUpdate:
			var d EventVoiceServerUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d

		case EventTypeWebhooksUpdate:
			var d EventWebhooksUpdate
			err = json.Unmarshal(v.D, &d)
			messageData = d
		}

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
		err = fmt.Errorf("unknown opcode %d", v.Op)
	}
	if err != nil {
		return err
	}
	e.Op = v.Op
	e.S = v.S
	e.T = v.T
	e.D = messageData
	e.RawD = v.D
	return nil
}

type MessageData interface {
	messageData()
}

type MessageDataDispatch json.RawMessage

func (MessageDataDispatch) messageData() {}

// MessageDataHeartbeat is used to ensure the websocket connection remains open, and disconnect if not.
type MessageDataHeartbeat int

func (MessageDataHeartbeat) messageData() {}

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

func (MessageDataIdentify) messageData() {}

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

func (MessageDataPresenceUpdate) messageData() {}

// MessageDataVoiceStateUpdate is used for updating the bots voice state in a guild
type MessageDataVoiceStateUpdate struct {
	GuildID   snowflake.ID  `json:"guild_id"`
	ChannelID *snowflake.ID `json:"channel_id"`
	SelfMute  bool          `json:"self_mute"`
	SelfDeaf  bool          `json:"self_deaf"`
}

func (MessageDataVoiceStateUpdate) messageData() {}

// MessageDataResume is used to resume a connection to discord in the case that you are disconnected. Is automatically
// handled by the library and should rarely be used.
type MessageDataResume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

func (MessageDataResume) messageData() {}

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

func (MessageDataRequestGuildMembers) messageData() {}

type MessageDataInvalidSession bool

func (MessageDataInvalidSession) messageData() {}

type MessageDataHello struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

func (MessageDataHello) messageData() {}
