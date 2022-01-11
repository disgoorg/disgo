package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Integration interface {
	Type() discord.IntegrationType
	ID() discord.Snowflake
}

type TwitchIntegration struct {
	IntegrationID     discord.Snowflake
	Name              string
	Enabled           bool
	Syncing           bool
	RoleID            discord.Snowflake
	EnableEmoticons   bool
	ExpireBehavior    int
	ExpireGracePeriod int
	Account           discord.IntegrationAccount
	SyncedAt          string
	SubscriberCount   int
	Revoked           bool
	Bot               Bot
	GuildID           discord.Snowflake
	User              *User
}

// Guild returns the Guild the Integration belongs to
func (i *TwitchIntegration) Guild() (Guild, bool) {
	return i.Bot.Caches().Guilds().Get(i.GuildID)
}

// Role returns the Subscriber Role the Integration uses
func (i *TwitchIntegration) Role() (Role, bool) {
	return i.Bot.Caches().Roles().Get(i.GuildID, i.RoleID)
}

// Delete deletes the Integration from the Guild
func (i *TwitchIntegration) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices().GuildService().DeleteIntegration(i.GuildID, i.ID(), opts...)
}

type YouTubeIntegration struct {
	discord.YouTubeIntegration
	Bot     Bot
	GuildID discord.Snowflake
	User    User
}

// Guild returns the Guild the Integration belongs to
func (i *YouTubeIntegration) Guild() (Guild, bool) {
	return i.Bot.Caches().Guilds().Get(i.GuildID)
}

// Role returns the Subscriber Role the Integration uses
func (i *YouTubeIntegration) Role() (Role, bool) {
	return i.Bot.Caches().Roles().Get(i.GuildID, i.RoleID)
}

// Delete deletes the Integration from the Guild
func (i *YouTubeIntegration) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices().GuildService().DeleteIntegration(i.GuildID, i.ID(), opts...)
}

type BotIntegration struct {
	IntegrationID discord.Snowflake
	Name          string
	Enabled       bool
	Account       discord.IntegrationAccount
	Bot           Bot
	GuildID       discord.Snowflake
	Application   IntegrationApplication
}

// Guild returns the Guild the Integration belongs to
func (i BotIntegration) Guild() (Guild, bool) {
	return i.Bot.Caches().Guilds().Get(i.GuildID)
}

// Delete deletes the Integration from the Guild
func (i BotIntegration) Delete(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices().GuildService().DeleteIntegration(i.GuildID, i.ID(), opts...)
}

type IntegrationApplication struct {
	ID          discord.Snowflake
	Name        string
	Icon        string
	Description string
	Summary     string
	Bot         User
}
