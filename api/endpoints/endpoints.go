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
	GetGateway    = NewAPIRoute(GET, "gateway")
	GetGatewayBot = NewAPIRoute(GET, "gateway/bot")

	GetUsersMe          = NewAPIRoute(GET, "users/@me")
	GetUser             = NewAPIRoute(GET, "users/{user.id}")
	PostUsersMeChannels = NewAPIRoute(POST, "users/@me/channels")

	GetMember     = NewAPIRoute(GET, "guilds/{guild.id}/members/{user.id}")
	CreateMessage = NewAPIRoute(POST, "channels/{channel.id}/messages")

	PutReaction        = NewAPIRoute(PUT, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	DeleteOwnReaction  = NewAPIRoute(DELETE, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	DeleteUserReaction = NewAPIRoute(DELETE, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}")
	GetReactions       = NewAPIRoute(GET, "channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}")

	GetGlobalApplicationCommands    = NewAPIRoute(GET, "applications/{application.id}/commands")
	CreateGlobalApplicationCommands = NewAPIRoute(POST, "applications/{application.id}/commands")
	SetGlobalApplicationCommands    = NewAPIRoute(PUT, "applications/{application.id}/commands")
	GetGlobalApplicationCommand     = NewAPIRoute(GET, "applications/{application.id}/command/{command.id}")
	EditGlobalApplicationCommand    = NewAPIRoute(PATCH, "applications/{application.id}/commands/{command.id}")
	DeleteGlobalApplicationCommand  = NewAPIRoute(DELETE, "applications/{application.id}/commands")

	GetGuildApplicationCommands    = NewAPIRoute(GET, "applications/{application.id}/{guild.id}/commands")
	CreateGuildApplicationCommands = NewAPIRoute(POST, "applications/{application.id}/{guild.id}/commands")
	SetGuildApplicationCommands    = NewAPIRoute(PUT, "applications/{application.id}/{guild.id}/commands")
	GetGuildApplicationCommand     = NewAPIRoute(GET, "applications/{application.id}/{guild.id}/command/{command.id}")
	EditGuildApplicationCommand    = NewAPIRoute(PATCH, "applications/{application.id}/{guild.id}/commands/{command.id}")
	DeleteGuildApplicationCommand  = NewAPIRoute(DELETE, "applications/{application.id}/{guild.id}/commands")

	CreateInteractionResponse = NewAPIRoute(POST, "interactions/{interaction.id}/{interaction.token}/callback")
	EditInteractionResponse   = NewAPIRoute(PATCH, "webhooks/{application.id}/{interaction.token}/messages/@original")
	DeleteInteractionResponse = NewAPIRoute(DELETE, "webhooks/{application.id}/{interaction.token}/messages/@original")
	CreateFollowupMessage     = NewAPIRoute(POST, "webhooks/{application.id}/{interaction.token}")
	EditFollowupMessage       = NewAPIRoute(PATCH, "webhooks/{application.id}/{interaction.token}/messages/{message.id}")
	DeleteFollowupMessage     = NewAPIRoute(PATCH, "webhooks/{application.id}/{interaction.token}/messages/{message.id}")

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
