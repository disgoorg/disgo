package api

type CacheStrategy func(disgo Disgo) bool

var (
	CacheStrategyYes  CacheStrategy = func(disgo Disgo) bool { return true }
	CacheStrategyNo   CacheStrategy = func(disgo Disgo) bool { return true }
	CacheStrategyNoWs CacheStrategy = func(disgo Disgo) bool { return disgo.HasGateway() }
)

type EntityBuilder interface {
	Disgo() Disgo

	CreateGlobalCommand(command *Command, updateCache CacheStrategy) *Command

	CreateUser(user *User, updateCache CacheStrategy) *User

	CreateMessage(message *Message, updateCache CacheStrategy) *Message

	CreateGuild(guild *Guild, updateCache CacheStrategy) *Guild
	CreateMember(guildID Snowflake, member *Member, updateCache CacheStrategy) *Member
	CreateGuildCommand(guildID Snowflake, command *Command, updateCache CacheStrategy) *Command
	CreateGuildCommandPermissions(guildCommandPermissions *GuildCommandPermissions, updateCache CacheStrategy) *GuildCommandPermissions
	CreateRole(guildID Snowflake, role *Role, updateCache CacheStrategy) *Role
	CreateVoiceState(role *VoiceState, updateCache CacheStrategy) *VoiceState

	CreateTextChannel(channel *Channel, updateCache CacheStrategy) *TextChannel
	CreateVoiceChannel(channel *Channel, updateCache CacheStrategy) *VoiceChannel
	CreateStoreChannel(channel *Channel, updateCache CacheStrategy) *StoreChannel
	CreateCategory(channel *Channel, updateCache CacheStrategy) *Category
	CreateDMChannel(channel *Channel, updateCache CacheStrategy) *DMChannel

	CreateEmote(guildID Snowflake, emote *Emote, updateCache CacheStrategy) *Emote
}
