package discord

import "github.com/disgoorg/snowflake"

type Connection struct {
	ID           snowflake.Snowflake `json:"id"`
	Name         string              `json:"name"`
	Type         ConnectionType      `json:"type"`
	Revoked      bool                `json:"revoked,omitempty"`
	Integrations []Integration       `json:"integrations,omitempty"`
	Verified     bool                `json:"verified"`
	FriendSync   bool                `json:"friend_sync"`
	ShowActivity bool                `json:"show_activity"`
	Visibility   VisibilityType      `json:"visibility"`
}

type ConnectionType string

//goland:noinspection GoUnusedConst
const (
	ConnectionTypeYouTube   ConnectionType = "youtube"
	ConnectionTypeBattleNet ConnectionType = "battlenet"
	ConnectionTypeGitHub    ConnectionType = "github"
	ConnectionTypeReddit    ConnectionType = "reddit"
	ConnectionTypeSpotify   ConnectionType = "spotify"
	ConnectionTypeSteam     ConnectionType = "steam"
	ConnectionTypeTwitch    ConnectionType = "twitch"
	ConnectionTypeTwitter   ConnectionType = "twitter"
	ConnectionTypeXBox      ConnectionType = "xbox"
	ConnectionTypeFacebook  ConnectionType = "facebook"
)

type VisibilityType int

//goland:noinspection GoUnusedConst
const (
	VisibilityTypeNone VisibilityType = iota
	VisibilityTypeEveryone
)
