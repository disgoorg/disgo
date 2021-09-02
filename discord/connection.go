package discord

type Connection struct {
	ID           Snowflake      `json:"id"`
	Name         string         `json:"name"`
	Type         ConnectionType `json:"type"`
	Revoked      bool           `json:"revoked,omitempty"`
	Integrations []Integration  `json:"integrations,omitempty"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	Visibility   VisibilityType `json:"visibility"`
}

type ConnectionType string

const (
	ConnectionTypeYouTube   CompressType = "youtube"
	ConnectionTypeBattleNet CompressType = "battlenet"
	ConnectionTypeGitHub    CompressType = "github"
	ConnectionTypeReddit    CompressType = "reddit"
	ConnectionTypeSpotify   CompressType = "spotify"
	ConnectionTypeSteam     CompressType = "steam"
	ConnectionTypeTwitch    CompressType = "twitch"
	ConnectionTypeTwitter   CompressType = "twitter"
	ConnectionTypeXBox      CompressType = "xbox"
	ConnectionTypeFacebook  CompressType = "facebook"
)

type VisibilityType int

const (
	VisibilityTypeNone VisibilityType = iota
	VisibilityTypeEveryone
)
