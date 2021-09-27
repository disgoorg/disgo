package core

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

type Presence struct {
	discord.Presence
	Bot *Bot
}

func (p *Presence) User() *User {
	return p.Bot.Caches.UserCache().Get(p.PresenceUser.ID)
}

func (p *Presence) Member() *Member {
	return p.Bot.Caches.MemberCache().Get(p.GuildID, p.PresenceUser.ID)
}

func (p *Presence) Guild() *Guild {
	return p.Bot.Caches.GuildCache().Get(p.GuildID)
}

func NewPresence(activityType discord.ActivityType, name string, url string, status discord.OnlineStatus, afk bool) discord.PresenceUpdate {
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

	return discord.PresenceUpdate{
		Since:      since,
		Activities: activities,
		Status:     status,
		AFK:        afk,
	}
}

func NewGamePresence(name string, status discord.OnlineStatus, afk bool) discord.PresenceUpdate {
	return NewPresence(discord.ActivityTypeGame, name, "", status, afk)
}

func NewStreamingPresence(name string, url string, status discord.OnlineStatus, afk bool) discord.PresenceUpdate {
	return NewPresence(discord.ActivityTypeStreaming, name, url, status, afk)
}

func NewListeningPresence(name string, status discord.OnlineStatus, afk bool) discord.PresenceUpdate {
	return NewPresence(discord.ActivityTypeListening, name, "", status, afk)
}

func NewWatchingPresence(name string, status discord.OnlineStatus, afk bool) discord.PresenceUpdate {
	return NewPresence(discord.ActivityTypeWatching, name, "", status, afk)
}

func NewCompetingPresence(name string, status discord.OnlineStatus, afk bool) discord.PresenceUpdate {
	return NewPresence(discord.ActivityTypeCompeting, name, "", status, afk)
}
