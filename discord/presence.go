package discord

import (
	"time"

	"github.com/disgoorg/snowflake"
)

// Presence (https://discord.com/developers/docs/topics/gateway#presence-update)
type Presence struct {
	PresenceUser PresenceUser        `json:"user"`
	GuildID      snowflake.Snowflake `json:"guild_id"`
	Status       OnlineStatus        `json:"status"`
	Activities   []Activity          `json:"activities"`
	ClientStatus ClientStatus        `json:"client_status"`
}

type PresenceUser struct {
	ID snowflake.Snowflake `json:"id"`
}

// OnlineStatus (https://discord.com/developers/docs/topics/gateway#update-presence-status-types)
type OnlineStatus string

const (
	OnlineStatusOnline    OnlineStatus = "online"
	OnlineStatusDND       OnlineStatus = "dnd"
	OnlineStatusIdle      OnlineStatus = "idle"
	OnlineStatusInvisible OnlineStatus = "invisible"
	OnlineStatusOffline   OnlineStatus = "offline"
)

// ClientStatus (https://discord.com/developers/docs/topics/gateway#client-status-object)
type ClientStatus struct {
	Desktop OnlineStatus `json:"desktop,omitempty"`
	Mobile  OnlineStatus `json:"mobile,omitempty"`
	Web     OnlineStatus `json:"web,omitempty"`
}

// NewPresence creates a new Presence with the provided properties
func NewPresence(activityType ActivityType, name string, url string, status OnlineStatus, afk bool) GatewayMessageDataPresenceUpdate {
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

	return GatewayMessageDataPresenceUpdate{
		Since:      since,
		Activities: activities,
		Status:     status,
		AFK:        afk,
	}
}

// NewGamePresence creates a new Presence of type ActivityTypeGame
func NewGamePresence(name string, status OnlineStatus, afk bool) GatewayMessageDataPresenceUpdate {
	return NewPresence(ActivityTypeGame, name, "", status, afk)
}

// NewStreamingPresence creates a new Presence of type ActivityTypeStreaming
func NewStreamingPresence(name string, url string, status OnlineStatus, afk bool) GatewayMessageDataPresenceUpdate {
	return NewPresence(ActivityTypeStreaming, name, url, status, afk)
}

// NewListeningPresence creates a new Presence of type ActivityTypeListening
func NewListeningPresence(name string, status OnlineStatus, afk bool) GatewayMessageDataPresenceUpdate {
	return NewPresence(ActivityTypeListening, name, "", status, afk)
}

// NewWatchingPresence creates a new Presence of type ActivityTypeWatching
func NewWatchingPresence(name string, status OnlineStatus, afk bool) GatewayMessageDataPresenceUpdate {
	return NewPresence(ActivityTypeWatching, name, "", status, afk)
}

// NewCompetingPresence creates a new Presence of type ActivityTypeCompeting
func NewCompetingPresence(name string, status OnlineStatus, afk bool) GatewayMessageDataPresenceUpdate {
	return NewPresence(ActivityTypeCompeting, name, "", status, afk)
}
