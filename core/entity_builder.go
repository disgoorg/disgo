package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// CacheStrategy is used to determine whether something should be cached when making an api request. When using the
// gateway, you'll receive the event shortly afterwards if you have the correct GatewayIntents.
type CacheStrategy func(disgo Disgo) bool

// Default cache strategy choices
var (
	CacheStrategyYes  CacheStrategy = func(_ Disgo) bool { return true }
	CacheStrategyNo   CacheStrategy = func(_ Disgo) bool { return true }
	CacheStrategyNoWs CacheStrategy = func(disgo Disgo) bool { return disgo.HasGateway() }
)

// EntityBuilder is used to create structs for disgo's cache
type EntityBuilder interface {
	Disgo() Disgo

	CreateInteraction(unmarshalInteraction discord.UnmarshalInteraction, responseChannel chan discord.InteractionResponse, updateCache CacheStrategy) *Interaction

	CreateApplicationCommandInteraction(interaction *Interaction, updateCache CacheStrategy) *ApplicationCommandInteraction

	CreateSlashCommandInteraction(applicationInteraction *ApplicationCommandInteraction) *SlashCommandInteraction

	CreateContextCommandInteraction(applicationInteraction *ApplicationCommandInteraction) *ContextCommandInteraction
	CreateUserCommandInteraction(contextCommandInteraction *ContextCommandInteraction) *UserCommandInteraction
	CreateMessageCommandInteraction(contextCommandInteraction *ContextCommandInteraction) *MessageCommandInteraction

	CreateComponentInteraction(interaction *Interaction, updateCache CacheStrategy) *ComponentInteraction
	CreateButtonInteraction(componentInteraction *ComponentInteraction) *ButtonInteraction
	CreateSelectMenuInteraction(componentInteraction *ComponentInteraction) *SelectMenuInteraction

	CreateUser(user discord.User, updateCache CacheStrategy) *User
	CreateSelfUser(selfUser discord.SelfUser, updateCache CacheStrategy) *SelfUser

	CreateMessage(message discord.Message, updateCache CacheStrategy) *Message
	CreateComponents(components []discord.Component, updateCache CacheStrategy) []Component

	CreateGuild(guild discord.Guild, updateCache CacheStrategy) *Guild
	CreateGuildTemplate(guildTemplate discord.GuildTemplate, updateCache CacheStrategy) *GuildTemplate
	CreateStageInstance(stageInstance discord.StageInstance, updateCache CacheStrategy) *StageInstance

	CreateRole(guildID discord.Snowflake, role discord.Role, updateCache CacheStrategy) *Role
	CreateMember(guildID discord.Snowflake, member discord.Member, updateCache CacheStrategy) *Member
	CreateBan(guildID discord.Snowflake, ban discord.Ban, updateCache CacheStrategy) *Ban
	CreateVoiceState(guildID discord.Snowflake, voiceState discord.VoiceState, updateCache CacheStrategy) *VoiceState

	CreateCommand(command discord.ApplicationCommand, updateCache CacheStrategy) *ApplicationCommand
	CreateCommandPermissions(guildCommandPermissions discord.GuildCommandPermissions, updateCache CacheStrategy) *GuildCommandPermissions

	CreateAuditLog(guildID discord.Snowflake, auditLog discord.AuditLog, filterOptions AuditLogFilterOptions, updateCache CacheStrategy) *AuditLog
	CreateIntegration(guildID discord.Snowflake, integration discord.Integration, updateCache CacheStrategy) *Integration

	CreateChannel(channel discord.Channel, updateCache CacheStrategy) Channel

	CreateInvite(invite discord.Invite, updateCache CacheStrategy) *Invite

	CreateEmoji(guildID discord.Snowflake, emoji discord.Emoji, updateCache CacheStrategy) *Emoji
}
