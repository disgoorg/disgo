package api

type EntityBuilder interface {
	Disgo() Disgo

	CreateGlobalCommand(command *Command, updateCache bool) *Command

	CreateUser(user *User, updateCache bool) *User

	CreateMessage(message *Message, updateCache bool) *Message

	CreateGuild(guild *Guild, updateCache bool) *Guild
	CreateMember(guildID Snowflake, member *Member, updateCache bool) *Member
	CreateGuildCommand(guildID Snowflake, command *Command, updateCache bool) *Command
	CreateRole(guildID Snowflake, role *Role, updateCache bool) *Role
	CreateVoiceState(role *VoiceState, updateCache bool) *VoiceState

	CreateTextChannel(channel *Channel, updateCache bool) *TextChannel
	CreateVoiceChannel(channel *Channel, updateCache bool) *VoiceChannel
	CreateStoreChannel(channel *Channel, updateCache bool) *StoreChannel
	CreateCategory(channel *Channel, updateCache bool) *Category
	CreateDMChannel(channel *Channel, updateCache bool) *DMChannel

	CreateEmote(guildID Snowflake, emote *Emote, updateCache bool) *Emote
}
