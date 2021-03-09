package constants

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
const (
	GuildTextChannel ChannelType = iota
	DMChannel
	GuildVoiceChannel
	GroupDMChannel
	GuildCategoryChannel
	GuildNewsChannel
	GuildStoreChannel
)
