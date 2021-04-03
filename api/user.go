package api

var _ Mentionable = (User)(nil)


// User is a struct for interacting with discord's users
type User interface {
	Disgo() Disgo
	ID() Snowflake
	Username() string
	Discriminator() int
	Tag() string
	AvatarURL() *string
	EffectiveAvatarURL() string
	Bot() bool
	Flags() UserFlags
	String() string
	Mention() string
	OpenDMChannel() (*DMChannel, error)
}


type UserFlags int64

const (
	UserFlagsNone           UserFlags = 0
	UserFlagDiscordEmployee UserFlags = 1 << iota
	UserFlagPartneredServerOwner
	UserFlagHypeSquadEvents
	UserFlagBugHunterLevel1
	UserFlagHouseBravery
	UserFlagHouseBrilliance
	UserFlagHouseBalance
	UserFlagEarlySupporter
	UserFlagTeamUser
	UserFlagSystem
	UserFlagBugHunterLevel2
	UserFlagVerifiedBot
	UserFlagEarlyVerifiedBotDeveloper
)
