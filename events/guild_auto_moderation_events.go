package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

type GenericAutoModerationRule struct {
	*GenericEvent
	discord.AutoModerationRule
}

// Guild returns the discord.Guild the event happened in.
// This will only check cached guilds!
func (e *GenericAutoModerationRule) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guild(e.GuildID)
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
	gateway.EventAutoModerationActionExecution
}

// Guild returns the discord.Guild the event happened in.
// This will only check cached guilds!
func (e *AutoModerationActionExecution) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guild(e.GuildID)
}
