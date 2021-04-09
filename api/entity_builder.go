package api

type EntityBuilder interface {
	Disgo() Disgo

	CreateGlobalCommand(command *Command, updateCache bool) *Command

	CreateUser(user *User, updateCache bool) *User

	CreateGuild(guild *Guild, updateCache bool) *Guild
	CreateMember(member *Member, updateCache bool) *Member
	CreateGuildCommand(guildID Snowflake, command *Command, updateCache bool) *Command

	CreateTextChannel(channel *Channel, updateCache bool) *TextChannel
	CreateVoiceChannel(channel *Channel, updateCache bool) *VoiceChannel
	CreateStoreChannel(channel *Channel, updateCache bool) *StoreChannel
	CreateCategory(channel *Channel, updateCache bool) *Category
	CreateDMChannel(channel *Channel, updateCache bool) *DMChannel

	CreateEmote(emote *Emote, updateCache bool) *Emote
}
