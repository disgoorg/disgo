package endpoints

// Discord Endpoint Constants
const (
	APIVersion = "8"
	Base       = "https://discord.com/"
	CDN        = "https://cdn.discordapp.com"
	API        = Base + "api/v" + APIVersion
	WS         = "wss://gateway.discord.gg/"
)

// All of the Discord endpoints used by disgo
var (
	GetGateway    = NewAPIRoute(GET, "/gateway")
	GetGatewayBot = NewAPIRoute(GET, "/gateway/bot")

	GetSelfUser     = NewAPIRoute(GET, "/users/@me")
	UpdateSelfNick  = NewAPIRoute(PATCH, "/guilds/{guild.id}/members/@me/nick")
	CreateMessage   = NewAPIRoute(POST, "/channels/{channel.id}/messages")
	CreateDMChannel = NewAPIRoute(POST, "/users/@me/channels")

	GetUser            = NewAPIRoute(GET, "/users/{user.id}")
	GetMember          = NewAPIRoute(GET, "/guilds/{guild_events.id}/members/{user.id}")
	GetMembers         = NewAPIRoute(GET, "/guilds/{guild.id}/members")
	AddMember          = NewAPIRoute(PUT, "/guilds/{guild.id}/members/{user.id}")
	RemoveMember       = NewAPIRoute(DELETE, "/guilds/{guild.id}/members/{user.id}")
	RemoveMemberReason = NewAPIRoute(DELETE, "/guilds/{guild.id}/members/{user.id}?reason={reason}")
	UpdateMember       = NewAPIRoute(PATCH, "/guilds/{guild.id}/members/{user.id}")
	AddMemberRole      = NewAPIRoute(PUT, "/guilds/{guild.id}/members/{user.id}/roles/{role.id}")
	RemoveMemberRole   = NewAPIRoute(DELETE, "/guilds/{guild.id}/members/{user.id}/roles/{role.id}")

	PutReaction        = NewAPIRoute(PUT, "/channels/{channel.id}/messages/{message_events.id}/reactions/{emoji}/@me")
	DeleteOwnReaction  = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message_events.id}/reactions/{emoji}/@me")
	DeleteUserReaction = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message_events.id}/reactions/{emoji}/{user.id}")
	GetReactions       = NewAPIRoute(GET, "/channels/{channel.id}/messages/{message_events.id}/reactions/{emoji}/{user.id}")

	GetRoles            = NewAPIRoute(GET, "/guilds/{guild.id}/roles")
	CreateRole          = NewAPIRoute(POST, "/guilds/{guild.id}/roles")
	UpdateRole          = NewAPIRoute(PATCH, "/guilds/{guild.id}/roles/{role.id}")
	UpdateRolePositions = NewAPIRoute(PATCH, "/guilds/{guild.id}/roles")
	DelteRole           = NewAPIRoute(DELETE, "/guilds/{guild.id}/roles/{role.id}")

	GetGlobalApplicationCommands   = NewAPIRoute(GET, "/applications/{application.id}/commands")
	CreateGlobalApplicationCommand = NewAPIRoute(POST, "/applications/{application.id}/commands")
	SetGlobalApplicationCommands   = NewAPIRoute(PUT, "/applications/{application.id}/commands")
	GetGlobalApplicationCommand    = NewAPIRoute(GET, "/applications/{application.id}/command/{command.id}")
	EditGlobalApplicationCommand   = NewAPIRoute(PATCH, "/applications/{application.id}/commands/{command.id}")
	DeleteGlobalApplicationCommand = NewAPIRoute(DELETE, "/applications/{application.id}/commands")

	GetGuildApplicationCommands   = NewAPIRoute(GET, "/applications/{application.id}/guilds/{guild_events.id}/commands")
	CreateGuildApplicationCommand = NewAPIRoute(POST, "/applications/{application.id}/guilds/{guild_events.id}/commands")
	SetGuildApplicationCommands   = NewAPIRoute(PUT, "/applications/{application.id}/guilds/{guild_events.id}/commands")
	GetGuildApplicationCommand    = NewAPIRoute(GET, "/applications/{application.id}/guilds/{guild_events.id}/command/{command.id}")
	EditGuildApplicationCommand   = NewAPIRoute(PATCH, "/applications/{application.id}/guilds/{guild_events.id}/commands/{command.id}")
	DeleteGuildApplicationCommand = NewAPIRoute(DELETE, "/applications/{application.id}/guilds/{guild_events.id}/commands")

	CreateInteractionResponse = NewAPIRoute(POST, "/interactions/{interaction.id}/{interaction.token}/callback")
	EditInteractionResponse   = NewAPIRoute(PATCH, "/webhooks/{application.id}/{interaction.token}/messages/@original")
	DeleteInteractionResponse = NewAPIRoute(DELETE, "/webhooks/{application.id}/{interaction.token}/messages/@original")
	CreateFollowupMessage     = NewAPIRoute(POST, "/webhooks/{application.id}/{interaction.token}")
	EditFollowupMessage       = NewAPIRoute(PATCH, "/webhooks/{application.id}/{interaction.token}/messages/{message_events.id}")
	DeleteFollowupMessage     = NewAPIRoute(PATCH, "/webhooks/{application.id}/{interaction.token}/messages/{message_events.id}")

	Emote                = NewCDNRoute("/emojis/{emote.id}.", PNG, GIF)
	GuildIcon            = NewCDNRoute("/icons/{guild_events.id}/{icon.hash}.", PNG, JPEG, WEBP, GIF)
	GuildSplash          = NewCDNRoute("/splashes/{guild_events.id}/guild_splash.", PNG, JPEG, WEBP)
	GuildDiscoverySplash = NewCDNRoute("/discovery-splashes/{guild_events.id}/guild_discovery_splash.", PNG, JPEG, WEBP)
	GuildBanner          = NewCDNRoute("/banners/{guild_events.id}/guild_banner.", PNG, JPEG, WEBP)
	DefaultUserAvatar    = NewCDNRoute("/embed/avatars/{user.discriminator%5}.", PNG)
	UserAvatar           = NewCDNRoute("/avatars/{user.id}/user_avatar.", PNG, JPEG, WEBP, GIF)
	ApplicationIcon      = NewCDNRoute("/app-icons/{application.id}/icon.", PNG, JPEG, WEBP)
	ApplicationAsset     = NewCDNRoute("/app-assets/{application.id}/{asset.id}.", PNG, JPEG, WEBP)
	AchievementIcon      = NewCDNRoute("/app-assets/{application.id}/achievements/{achievement.id}/icons/{icon.hash}.", PNG, JPEG, WEBP)
	TeamIcon             = NewCDNRoute("/team-icons/{team.id}/team_icon.", PNG, JPEG, WEBP)
)
