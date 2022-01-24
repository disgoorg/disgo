package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/snowflake"
)

type Integration interface {
	discord.Integration
}

type TwitchIntegration struct {
	discord.TwitchIntegration
	Bot     *Bot
	GuildID snowflake.Snowflake
	User    *User
}

// Guild returns the Guild the Integration belongs to
func (i *TwitchIntegration) Guild() *Guild {
	return i.Bot.Caches.Guilds().Get(i.GuildID)
}

// Role returns the Subscriber Role the Integration uses
func (i *TwitchIntegration) Role() *Role {
	return i.Bot.Caches.Roles().Get(i.GuildID, i.RoleID)
}

// Delete deletes the Integration from the Guild
func (i *TwitchIntegration) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.GuildService().DeleteIntegration(i.GuildID, i.ID(), opts...)
}

type YouTubeIntegration struct {
	discord.YouTubeIntegration
	Bot     *Bot
	GuildID snowflake.Snowflake
	User    *User
}

// Guild returns the Guild the Integration belongs to
func (i *YouTubeIntegration) Guild() *Guild {
	return i.Bot.Caches.Guilds().Get(i.GuildID)
}

// Role returns the Subscriber Role the Integration uses
func (i *YouTubeIntegration) Role() *Role {
	return i.Bot.Caches.Roles().Get(i.GuildID, i.RoleID)
}

// Delete deletes the Integration from the Guild
func (i *YouTubeIntegration) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.GuildService().DeleteIntegration(i.GuildID, i.ID(), opts...)
}

type BotIntegration struct {
	discord.BotIntegration
	Bot         *Bot
	GuildID     snowflake.Snowflake
	Application *IntegrationApplication
}

// Guild returns the Guild the Integration belongs to
func (i *BotIntegration) Guild() *Guild {
	return i.Bot.Caches.Guilds().Get(i.GuildID)
}

// Delete deletes the Integration from the Guild
func (i *BotIntegration) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.GuildService().DeleteIntegration(i.GuildID, i.ID(), opts...)
}

type IntegrationApplication struct {
	discord.IntegrationApplication
	Bot *User
}
