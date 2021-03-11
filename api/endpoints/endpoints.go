package endpoints

// Discord Endpoint Constants
const (
	APIVersion = "8"
	Base       = "https://discord.com/"
	CDN        = "https://cdn.discordapp.com/"
	API        = Base + "api/v" + APIVersion + "/"
	WS         = "wss://gateway.discord.gg/"
)

// All of the Discord endpoints used by disgo
var (
	GetGateway          = NewAPIRoute(GET, "gateway")
	GetGatewayBot       = NewAPIRoute(GET, "gateway/bot")
	GetUsersMe          = NewAPIRoute(GET, "users/@me")
	GetUser             = NewAPIRoute(GET, "users/{user.id}")
	PostUsersMeChannels = NewAPIRoute(POST, "users/@me/channels")
	PutReaction         = NewAPIRoute(PUT, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	DeleteOwnReaction   = NewAPIRoute(DELETE, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	DeleteUserReaction  = NewAPIRoute(DELETE, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}")
	GetReactions        = NewAPIRoute(GET, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}")
	GetMember           = NewAPIRoute(GET, "guilds/{guild.id}/members/{user.id}")
	CreateMessage       = NewAPIRoute(POST, "channels/{channel.id}/messages")

	Emote                = NewCDNRoute("emojis/{emote.id}.", PNG, GIF)
	GuildIcon            = NewCDNRoute("icons/{guild.id}/{icon.hash}.", PNG, JPEG, WEBP, GIF)
	GuildSplash          = NewCDNRoute("splashes/{guild.id}/guild_splash.", PNG, JPEG, WEBP)
	GuildDiscoverySplash = NewCDNRoute("discovery-splashes/{guild.id}/guild_discovery_splash.", PNG, JPEG, WEBP)
	GuildBanner          = NewCDNRoute("banners/{guild.id}/guild_banner.", PNG, JPEG, WEBP)
	DefaultUserAvatar    = NewCDNRoute("embed/avatars/{user.discriminator%5}.", PNG)
	UserAvatar           = NewCDNRoute("avatars/{user.id}/user_avatar.", PNG, JPEG, WEBP, GIF)
	ApplicationIcon      = NewCDNRoute("app-icons/{application.id}/icon.", PNG, JPEG, WEBP)
	ApplicationAsset     = NewCDNRoute("app-assets/{application.id}/{asset.id}.", PNG, JPEG, WEBP)
	AchievementIcon      = NewCDNRoute("app-assets/{application.id}/achievements/{achievement.id}/icons/{icon.hash}.", PNG, JPEG, WEBP)
	TeamIcon             = NewCDNRoute("team-icons/{team.id}/team_icon.", PNG, JPEG, WEBP)
)
