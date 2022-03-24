package discord

import (
	"time"

	"github.com/disgoorg/snowflake"
)

// GatewayCommand object is used when sending data to discord's websocket, it's recommended that you don't use these
//goland:noinspection GoNameStartsWithPackageName
type GatewayCommand struct {
	GatewayPayload
	D GatewayCommandData `json:"d"`
}

type GatewayCommandData interface {
	gatewayCommandData()
}

var _ GatewayCommandData = (*GatewaySequence)(nil)

type GatewaySequence int

func (GatewaySequence) gatewayCommandData() {}

var _ GatewayCommandData = (*IdentifyCommandData)(nil)

// IdentifyCommandData is the data used in IdentifyCommandData
type IdentifyCommandData struct {
	Token          string                        `json:"token"`
	Properties     IdentifyCommandDataProperties `json:"properties"`
	Compress       bool                          `json:"compress,omitempty"`
	LargeThreshold int                           `json:"large_threshold,omitempty"`
	Shard          []int                         `json:"shard,omitempty"`
	GatewayIntents GatewayIntents                `json:"intents"`
	Presence       *UpdatePresenceCommandData    `json:"presence,omitempty"`
}

func (IdentifyCommandData) gatewayCommandData() {}

// IdentifyCommandDataProperties is used for specifying to discord which library and OS the bot is using, is
// automatically handled by the library and should rarely be used.
type IdentifyCommandDataProperties struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
}

var _ GatewayCommandData = (*IdentifyCommandData)(nil)

// ResumeCommandData is used to resume a connection to discord in the case that you are disconnected. Is automatically
// handled by the library and should rarely be used.
type ResumeCommandData struct {
	Token     string          `json:"token"`
	SessionID string          `json:"session_id"`
	Seq       GatewaySequence `json:"seq"`
}

func (ResumeCommandData) gatewayCommandData() {}

var _ GatewayCommandData = (*HeartbeatCommandData)(nil)

// HeartbeatCommandData is used to ensure the websocket connection remains open, and disconnect if not.
type HeartbeatCommandData int

func (HeartbeatCommandData) gatewayCommandData() {}

// RequestGuildMembersCommandData is used for fetching all the members of a guild_events. It is recommended you have a strict
// member caching policy when using this.
type RequestGuildMembersCommandData struct {
	GuildID   snowflake.Snowflake   `json:"guild_id"`
	Query     *string               `json:"query,omitempty"` //If specified, user_ids must not be entered
	Limit     *int                  `json:"limit,omitempty"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool                  `json:"presences,omitempty"`
	UserIDs   []snowflake.Snowflake `json:"user_ids,omitempty"` //If specified, query must not be entered
	Nonce     string                `json:"nonce,omitempty"`    //All responses are hashed with this nonce, optional
}

func (RequestGuildMembersCommandData) gatewayCommandData() {}

// UpdateVoiceStateCommandData is used for updating the bots voice state in a guild_events
type UpdateVoiceStateCommandData struct {
	GuildID   snowflake.Snowflake  `json:"guild_id"`
	ChannelID *snowflake.Snowflake `json:"channel_id"`
	SelfMute  bool                 `json:"self_mute"`
	SelfDeaf  bool                 `json:"self_deaf"`
}

func (UpdateVoiceStateCommandData) gatewayCommandData() {}

// UpdatePresenceCommandData is used for updating Client's presence
type UpdatePresenceCommandData struct {
	Since      *int64       `json:"since"`
	Activities []Activity   `json:"activities"`
	Status     OnlineStatus `json:"status"`
	AFK        bool         `json:"afk"`
}

func (UpdatePresenceCommandData) gatewayCommandData() {}

// NewPresence creates a new Presence with the provided properties
func NewPresence(activityType ActivityType, name string, url string, status OnlineStatus, afk bool) UpdatePresenceCommandData {
	var since *int64
	if status == OnlineStatusIdle {
		unix := time.Now().Unix()
		since = &unix
	}

	var activities []Activity
	if name != "" {
		activity := Activity{
			Name: name,
			Type: activityType,
		}
		if activityType == ActivityTypeStreaming && url != "" {
			activity.URL = &url
		}
		activities = append(activities, activity)
	}

	return UpdatePresenceCommandData{
		Since:      since,
		Activities: activities,
		Status:     status,
		AFK:        afk,
	}
}

// NewGamePresence creates a new Presence of type ActivityTypeGame
func NewGamePresence(name string, status OnlineStatus, afk bool) UpdatePresenceCommandData {
	return NewPresence(ActivityTypeGame, name, "", status, afk)
}

// NewStreamingPresence creates a new Presence of type ActivityTypeStreaming
func NewStreamingPresence(name string, url string, status OnlineStatus, afk bool) UpdatePresenceCommandData {
	return NewPresence(ActivityTypeStreaming, name, url, status, afk)
}

// NewListeningPresence creates a new Presence of type ActivityTypeListening
func NewListeningPresence(name string, status OnlineStatus, afk bool) UpdatePresenceCommandData {
	return NewPresence(ActivityTypeListening, name, "", status, afk)
}

// NewWatchingPresence creates a new Presence of type ActivityTypeWatching
func NewWatchingPresence(name string, status OnlineStatus, afk bool) UpdatePresenceCommandData {
	return NewPresence(ActivityTypeWatching, name, "", status, afk)
}

// NewCompetingPresence creates a new Presence of type ActivityTypeCompeting
func NewCompetingPresence(name string, status OnlineStatus, afk bool) UpdatePresenceCommandData {
	return NewPresence(ActivityTypeCompeting, name, "", status, afk)
}
