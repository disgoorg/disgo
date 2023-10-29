package discord

type Connection struct {
	ID           string         `json:"id"`
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
	ConnectionTypeBattleNet          ConnectionType = "battlenet"
	ConnectionTypeEbay               ConnectionType = "ebay"
	ConnectionTypeEpicGames          ConnectionType = "epicgames"
	ConnectionTypeFacebook           ConnectionType = "facebook"
	ConnectionTypeGitHub             ConnectionType = "github"
	ConnectionTypeInstagram          ConnectionType = "instagram"
	ConnectionTypeLeagueOfLegends    ConnectionType = "leagueoflegends"
	ConnectionTypePayPal             ConnectionType = "paypal"
	ConnectionTypePlayStationNetwork ConnectionType = "playstation"
	ConnectionTypeReddit             ConnectionType = "reddit"
	ConnectionTypeRiotGames          ConnectionType = "riotgames"
	ConnectionTypeSpotify            ConnectionType = "spotify"
	ConnectionTypeSkype              ConnectionType = "skype"
	ConnectionTypeSteam              ConnectionType = "steam"
	ConnectionTypeTikTok             ConnectionType = "tiktok"
	ConnectionTypeTwitch             ConnectionType = "twitch"
	ConnectionTypeTwitter            ConnectionType = "twitter"
	ConnectionTypeXbox               ConnectionType = "xbox"
	ConnectionTypeYouTube            ConnectionType = "youtube"
)

type VisibilityType int

const (
	VisibilityTypeNone VisibilityType = iota
	VisibilityTypeEveryone
)
