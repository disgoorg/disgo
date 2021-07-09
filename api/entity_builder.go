package api

// CacheStrategy is used to determine whether something should be cached when making an api request. When using the
// gateway, you'll receive the event shortly afterwards if you have the correct GatewayIntents.
type CacheStrategy func(disgo Disgo) bool

// Default cache strategy choices
var (
	CacheStrategyYes  CacheStrategy = func(disgo Disgo) bool { return true }
	CacheStrategyNo   CacheStrategy = func(disgo Disgo) bool { return true }
	CacheStrategyNoWs CacheStrategy = func(disgo Disgo) bool { return disgo.HasGateway() }
)

// EntityBuilder is used to create structs for disgo's cache
type EntityBuilder interface {
	Disgo() Disgo

	CreateInteraction(fullInteraction *FullInteraction, c chan InteractionResponse, updateCache CacheStrategy) *Interaction
	CreateCommandInteraction(fullInteraction *FullInteraction, interaction *Interaction, updateCache CacheStrategy) *CommandInteraction
	CreateComponentInteraction(fullInteraction *FullInteraction, interaction *Interaction, updateCache CacheStrategy) *ComponentInteraction
	CreateButtonInteraction(fullInteraction *FullInteraction, componentInteraction *ComponentInteraction) *ButtonInteraction
	CreateSelectMenuInteraction(fullInteraction *FullInteraction, componentInteraction *ComponentInteraction) *SelectMenuInteraction

	CreateGlobalCommand(command *Command, updateCache CacheStrategy) *Command

	CreateUser(user *User, updateCache CacheStrategy) *User

	CreateMessage(message *Message, updateCache CacheStrategy) *Message

	CreateGuildTemplate(guildTemplate *GuildTemplate, updateCache CacheStrategy) *GuildTemplate
	CreateGuild(fullGuild *FullGuild, updateCache CacheStrategy) *Guild
	CreateMember(guildID Snowflake, member *Member, updateCache CacheStrategy) *Member
	CreateThreadMember(guildID Snowflake, member *ThreadMember, updateCache CacheStrategy) *ThreadMember
	CreateGuildCommand(guildID Snowflake, command *Command, updateCache CacheStrategy) *Command
	CreateGuildCommandPermissions(guildCommandPermissions *GuildCommandPermissions, updateCache CacheStrategy) *GuildCommandPermissions
	CreateRole(guildID Snowflake, role *Role, updateCache CacheStrategy) *Role
	CreateVoiceState(guildID Snowflake, voiceState *VoiceState, updateCache CacheStrategy) *VoiceState

	CreateTextChannel(channel *ChannelImpl, updateCache CacheStrategy) TextChannel
	CreateThread(channel *ChannelImpl, updateCache CacheStrategy) Thread
	CreateVoiceChannel(channel *ChannelImpl, updateCache CacheStrategy) VoiceChannel
	CreateStoreChannel(channel *ChannelImpl, updateCache CacheStrategy) StoreChannel
	CreateCategory(channel *ChannelImpl, updateCache CacheStrategy) Category
	CreateDMChannel(channel *ChannelImpl, updateCache CacheStrategy) DMChannel

	CreateEmoji(guildID Snowflake, emoji *Emoji, updateCache CacheStrategy) *Emoji
}
