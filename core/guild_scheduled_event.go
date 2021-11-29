package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type GuildScheduledEvent struct {
	discord.GuildScheduledEvent
	Creator *User
	Bot     *Bot
}

func (e *GuildScheduledEvent) Update(guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...rest.RequestOpt) (*GuildScheduledEvent, error) {
	guildScheduledEvent, err := e.Bot.RestServices.GuildScheduledEventService().UpdateGuildScheduledEvent(e.GuildID, e.ID, guildScheduledEventUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return e.Bot.EntityBuilder.CreateGuildScheduledEvent(*guildScheduledEvent, CacheStrategyNoWs), nil
}

func (e *GuildScheduledEvent) Delete(opts ...rest.RequestOpt) error {
	return e.Bot.RestServices.GuildScheduledEventService().DeleteGuildScheduledEvent(e.GuildID, e.ID, opts...)
}

func (e *GuildScheduledEvent) VoiceChannelEntity() *GuildVoiceChannel {
	if e.EntityType != discord.ScheduledEventEntityTypeVoice {
		return nil
	}
	if ch := e.Bot.Caches.Channels().Get(*e.EntityID); ch != nil {
		return ch.(*GuildVoiceChannel)
	}
	return nil
}

func (e *GuildScheduledEvent) StageInstanceEntity() *StageInstance {
	if e.EntityType != discord.ScheduledEventEntityTypeStageInstance {
		return nil
	}
	if e.EntityID == nil {
		return nil
	}
	return e.Bot.Caches.StageInstances().Get(*e.EntityID)
}

func (e *GuildScheduledEvent) Guild() *Guild {
	return e.Bot.Caches.Guilds().Get(e.GuildID)
}

func (e *GuildScheduledEvent) Channel() GuildChannel {
	if e.ChannelID == nil {
		return nil
	}
	if channel := e.Bot.Caches.Channels().Get(*e.ChannelID); channel != nil {
		return channel.(GuildChannel)
	}
	return nil
}

func (e *GuildScheduledEvent) GetUsers(limit int, withMember bool, before discord.Snowflake, after discord.Snowflake, opts ...rest.RequestOpt) ([]*GuildScheduledEventUser, error) {
	users, err := e.Bot.RestServices.GuildScheduledEventService().GetGuildScheduledEventUsers(e.GuildID, e.ID, limit, withMember, before, after, opts...)
	if err != nil {
		return nil, err
	}
	eventUsers := make([]*GuildScheduledEventUser, len(users))
	for i := range users {
		eventUsers[i] = e.Bot.EntityBuilder.CreateGuildScheduledEventUser(e.GuildID, users[i], CacheStrategyNoWs)
	}
	return eventUsers, nil
}

type GuildScheduledEventUser struct {
	discord.GuildScheduledEventUser
	User   *User
	Member *Member
	Bot    *Bot
}
