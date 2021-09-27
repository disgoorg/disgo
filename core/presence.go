package core

import "github.com/DisgoOrg/disgo/discord"

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
