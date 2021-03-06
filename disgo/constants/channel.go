package constants

type ChannelType int

const (
	GuildText ChannelType = iota
	DM
	GuildVoice
	GroupDM
	GuildCategory
	GuildNews
	GuildStore
)