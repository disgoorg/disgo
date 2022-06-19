package events

import "github.com/disgoorg/disgo/discord"

type GenericAutoModerationRule struct {
	*GenericEvent
	discord.AutoModerationRule
}

// Guild returns the discord.Guild the event happened in.
// This will only check cached guilds!
func (e *GenericAutoModerationRule) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildID)
}

type AutoModerationRuleCreate struct {
	*GenericAutoModerationRule
}

type AutoModerationRuleUpdate struct {
	*GenericAutoModerationRule
}

type AutoModerationRuleDelete struct {
	*GenericAutoModerationRule
}

type AutoModerationActionExecution struct {
	*GenericEvent
	discord.GatewayEventAutoModerationActionExecution
}

// Guild returns the discord.Guild the event happened in.
// This will only check cached guilds!
func (e *AutoModerationActionExecution) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildID)
}
