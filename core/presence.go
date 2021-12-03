package core

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

type Presence struct {
	discord.Presence
	Bot *Bot
}

// User returns the User this Presence belongs to.
// This will only check cached users!
func (p *Presence) User() *User {
	return p.Bot.Caches.Users().Get(p.PresenceUser.ID)
}

// Member returns the Member this Presence belongs to.
// This will only check cached members!
func (p *Presence) Member() *Member {
	return p.Bot.Caches.Members().Get(p.GuildID, p.PresenceUser.ID)
}

// Guild returns the Guild this Presence belongs to.
// This will only check cached guilds!
func (p *Presence) Guild() *Guild {
	return p.Bot.Caches.Guilds().Get(p.GuildID)
}

// NewPresence creates a new Presence with the provided properties
//goland:noinspection GoUnusedExportedFunction
func NewPresence(activityType discord.ActivityType, name string, url string, status discord.OnlineStatus, afk bool) discord.UpdatePresenceCommandData {
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

	return discord.UpdatePresenceCommandData{
		Since:      since,
		Activities: activities,
		Status:     status,
		AFK:        afk,
	}
}

// NewGamePresence creates a new Presence of type discord.ActivityTypeGame
//goland:noinspection GoUnusedExportedFunction
func NewGamePresence(name string, status discord.OnlineStatus, afk bool) discord.UpdatePresenceCommandData {
	return NewPresence(discord.ActivityTypeGame, name, "", status, afk)
}

// NewStreamingPresence creates a new Presence of type discord.ActivityTypeStreaming
//goland:noinspection GoUnusedExportedFunction
func NewStreamingPresence(name string, url string, status discord.OnlineStatus, afk bool) discord.UpdatePresenceCommandData {
	return NewPresence(discord.ActivityTypeStreaming, name, url, status, afk)
}

// NewListeningPresence creates a new Presence of type discord.ActivityTypeListening
//goland:noinspection GoUnusedExportedFunction
func NewListeningPresence(name string, status discord.OnlineStatus, afk bool) discord.UpdatePresenceCommandData {
	return NewPresence(discord.ActivityTypeListening, name, "", status, afk)
}

// NewWatchingPresence creates a new Presence of type discord.ActivityTypeWatching
//goland:noinspection GoUnusedExportedFunction
func NewWatchingPresence(name string, status discord.OnlineStatus, afk bool) discord.UpdatePresenceCommandData {
	return NewPresence(discord.ActivityTypeWatching, name, "", status, afk)
}

// NewCompetingPresence creates a new Presence of type discord.ActivityTypeCompeting
//goland:noinspection GoUnusedExportedFunction
func NewCompetingPresence(name string, status discord.OnlineStatus, afk bool) discord.UpdatePresenceCommandData {
	return NewPresence(discord.ActivityTypeCompeting, name, "", status, afk)
}
