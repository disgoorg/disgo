package constants

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
const (
	GuildText ChannelType = iota
	DM
	GuildVoice
	GroupDM
	GuildCategory
	GuildNews
	GuildStore
)
