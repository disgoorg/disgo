package route

// Discord Endpoint Constants
const (
	APIVersion = "10"
	Base       = "https://discord.com/"
	CDN        = "https://cdn.discordapp.com"
	API        = Base + "api/v" + APIVersion
)

// Misc
var (
	GetGateway      = NewAPIRoute(GET, "/gateway")
	GetGatewayBot   = NewAPIRoute(GET, "/gateway/bot")
	GetVoiceRegions = NewAPIRoute(GET, "/voice/regions")
)

// OAuth2
var (
	GetBotApplicationInfo = NewAPIRoute(GET, "/oauth2/applications/@me")
	GetAuthorizationInfo  = NewAPIRouteNoAuth(GET, "/oauth2/@me")
	Authorize             = NewRoute("/oauth2/authorize", "client_id", "permissions", "redirect_uri", "response_type", "scope", "state", "guild_id", "disable_guild_select")
	Token                 = NewAPIRouteNoAuth(POST, "/oauth2/token")
)

// Users
var (
	GetUser                   = NewAPIRoute(GET, "/users/{user.id}")
	GetCurrentUser            = NewAPIRoute(GET, "/users/@me")
	GetCurrentMember          = NewAPIRoute(GET, "/users/@me/guilds/{guild.id}/member")
	UpdateSelfUser            = NewAPIRoute(PATCH, "/users/@me")
	GetCurrentUserConnections = NewAPIRouteNoAuth(GET, "/users/@me/connections")
	GetCurrentUserGuilds      = NewAPIRouteNoAuth(GET, "/users/@me/guilds", "before", "after", "limit")
	LeaveGuild                = NewAPIRoute(DELETE, "/users/@me/guilds/{guild.id}")
	GetDMChannels             = NewAPIRoute(GET, "/users/@me/channels")
	CreateDMChannel           = NewAPIRoute(POST, "/users/@me/channels")
)

// Guilds
var (
	GetGuild          = NewAPIRoute(GET, "/guilds/{guild.id}", "with_counts")
	GetGuildPreview   = NewAPIRoute(GET, "/guilds/{guild.id}/preview")
	CreateGuild       = NewAPIRoute(POST, "/guilds")
	UpdateGuild       = NewAPIRoute(PATCH, "/guilds/{guild.id}")
	DeleteGuild       = NewAPIRoute(DELETE, "/guilds/{guild.id}")
	GetGuildVanityURL = NewAPIRoute(GET, "/guilds/{guild.id}/vanity-url")

	CreateGuildChannel     = NewAPIRoute(POST, "/guilds/{guild.id}/channels")
	GetGuildChannels       = NewAPIRoute(GET, "/guilds/{guild.id}/channels")
	UpdateChannelPositions = NewAPIRoute(PATCH, "/guilds/{guild.id}/channels")

	GetBans   = NewAPIRoute(GET, "/guilds/{guild.id}/bans", "before", "after", "limit")
	GetBan    = NewAPIRoute(GET, "/guilds/{guild.id}/bans/{user.id}")
	AddBan    = NewAPIRoute(PUT, "/guilds/{guild.id}/bans/{user.id}")
	DeleteBan = NewAPIRoute(DELETE, "/guilds/{guild.id}/bans/{user.id}")

	GetMember        = NewAPIRoute(GET, "/guilds/{guild.id}/members/{user.id}")
	GetMembers       = NewAPIRoute(GET, "/guilds/{guild.id}/members")
	SearchMembers    = NewAPIRoute(GET, "/guilds/{guild.id}/members/search", "query", "limit")
	AddMember        = NewAPIRoute(PUT, "/guilds/{guild.id}/members/{user.id}")
	UpdateMember     = NewAPIRoute(PATCH, "/guilds/{guild.id}/members/{user.id}")
	RemoveMember     = NewAPIRoute(DELETE, "/guilds/{guild.id}/members/{user.id}", "reason")
	AddMemberRole    = NewAPIRoute(PUT, "/guilds/{guild.id}/members/{user.id}/roles/{role.id}")
	RemoveMemberRole = NewAPIRoute(DELETE, "/guilds/{guild.id}/members/{user.id}/roles/{role.id}")

	UpdateSelfNick = NewAPIRoute(PATCH, "/guilds/{guild.id}/members/@me/nick")

	GetPruneMembersCount = NewAPIRoute(GET, "/guilds/{guild.id}/prune")
	PruneMembers         = NewAPIRoute(POST, "/guilds/{guild.id}/prune")

	GetGuildWebhooks = NewAPIRoute(GET, "/guilds/{guild.id}/webhooks")

	GetAuditLogs = NewAPIRoute(GET, "/guilds/{guild.id}/audit-logs", "user_id", "action_type", "before", "limit")

	GetGuildVoiceRegions = NewAPIRoute(GET, "/guilds/{guild.id}/regions")

	UpdateCurrentUserVoiceState = NewAPIRoute(PATCH, "/guilds/{guild.id}/voice-states/@me")
	UpdateUserVoiceState        = NewAPIRoute(PATCH, "/guilds/{guild.id}/voice-states/{user.id}")
)

// AutoModeration
var (
	GetAutoModerationRules   = NewAPIRoute(GET, "/guilds/{guild.id}/auto-moderation/rules")
	GetAutoModerationRule    = NewAPIRoute(GET, "/guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}")
	CreateAutoModerationRule = NewAPIRoute(POST, "/guilds/{guild.id}/auto-moderation/rules")
	UpdateAutoModerationRule = NewAPIRoute(PATCH, "/guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}")
	DeleteAutoModerationRule = NewAPIRoute(DELETE, "/guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}")
)

// GuildIntegrations
var (
	GetIntegrations   = NewAPIRoute(GET, "/guilds/{guild.id}/integrations")
	CreateIntegration = NewAPIRoute(POST, "/guilds/{guild.id}/integrations")
	UpdateIntegration = NewAPIRoute(PATCH, "/guilds/{guild.id}/integrations/{integration.id}")
	DeleteIntegration = NewAPIRoute(DELETE, "/guilds/{guild.id}/integrations/{integration.id}")
	SyncIntegration   = NewAPIRoute(POST, "/guilds/{guild.id}/integrations/{integration.id}/sync")
)

// GuildTemplates
var (
	GetGuildTemplate        = NewAPIRoute(GET, "/guilds/templates/{template.code}")
	GetGuildTemplates       = NewAPIRoute(GET, "/guilds/{guild.id}/templates")
	CreateGuildTemplate     = NewAPIRoute(POST, "/guilds/{guild.id}/templates")
	SyncGuildTemplate       = NewAPIRoute(PUT, "/guilds/{guild.id}/templates/{template.code}")
	UpdateGuildTemplate     = NewAPIRoute(PATCH, "/guilds/{guild.id}/templates/{template.code}")
	DeleteGuildTemplate     = NewAPIRoute(DELETE, "/guilds/{guild.id}/templates/{template.code}")
	CreateGuildFromTemplate = NewAPIRoute(POST, "/guilds/templates/{template.code}")
)

// GuildScheduledEvents
var (
	GetGuildScheduledEvents   = NewAPIRoute(GET, "/guilds/{guild.id}/scheduled-events", "with_user_count")
	GetGuildScheduledEvent    = NewAPIRoute(GET, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}", "with_user_count")
	CreateGuildScheduledEvent = NewAPIRoute(POST, "/guilds/{guild.id}/scheduled-events")
	UpdateGuildScheduledEvent = NewAPIRoute(PATCH, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}")
	DeleteGuildScheduledEvent = NewAPIRoute(DELETE, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}")

	GetGuildScheduledEventUsers = NewAPIRoute(GET, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users", "limit", "with_member", "before", "after")
)

// StageInstance
var (
	GetStageInstance    = NewAPIRoute(GET, "/stage-instances/{channel.id}")
	CreateStageInstance = NewAPIRoute(POST, "/stage-instances")
	UpdateStageInstance = NewAPIRoute(PATCH, "/stage-instances/{channel.id}")
	DeleteStageInstance = NewAPIRoute(DELETE, "/stage-instances/{channel.id}")
)

// Roles
var (
	GetRoles            = NewAPIRoute(GET, "/guilds/{guild.id}/roles")
	GetRole             = NewAPIRoute(GET, "/guilds/{guild.id}/roles/{role.id}")
	CreateRole          = NewAPIRoute(POST, "/guilds/{guild.id}/roles")
	UpdateRole          = NewAPIRoute(PATCH, "/guilds/{guild.id}/roles/{role.id}")
	UpdateRolePositions = NewAPIRoute(PATCH, "/guilds/{guild.id}/roles")
	DeleteRole          = NewAPIRoute(DELETE, "/guilds/{guild.id}/roles/{role.id}")
)

// Channels
var (
	GetChannel    = NewAPIRoute(GET, "/channels/{channel.id}")
	UpdateChannel = NewAPIRoute(PATCH, "/channels/{channel.id}")
	DeleteChannel = NewAPIRoute(DELETE, "/channels/{channel.id}")

	GetChannelWebhooks = NewAPIRoute(GET, "/channels/{channel.id}/webhooks")
	CreateWebhook      = NewAPIRoute(POST, "/channels/{channel.id}/webhooks")

	GetPermissionOverwrites   = NewAPIRoute(GET, "/channels/{channel.id}/permissions")
	GetPermissionOverwrite    = NewAPIRoute(GET, "/channels/{channel.id}/permissions/{overwrite.id}")
	UpdatePermissionOverwrite = NewAPIRoute(PUT, "/channels/{channel.id}/permissions/{overwrite.id}")
	DeletePermissionOverwrite = NewAPIRoute(DELETE, "/channels/{channel.id}/permissions/{overwrite.id}")

	SendTyping    = NewAPIRoute(POST, "/channels/{channel.id}/typing")
	FollowChannel = NewAPIRoute(POST, "/channels/{channel.id}/followers")
)

// Threads
var (
	CreateThreadWithMessage = NewAPIRoute(POST, "/channels/{channel.id}/messages/{message.id}/threads")
	CreateThread            = NewAPIRoute(POST, "/channels/{channel.id}/threads")
	JoinThread              = NewAPIRoute(PUT, "/channels/{channel.id}/thread-members/@me")
	LeaveThread             = NewAPIRoute(DELETE, "/channels/{channel.id}/thread-members/@me")
	AddThreadMember         = NewAPIRoute(PUT, "/channels/{channel.id}/thread-members/{user.id}")
	RemoveThreadMember      = NewAPIRoute(DELETE, "/channels/{channel.id}/thread-members/{user.id}")
	GetThreadMember         = NewAPIRoute(GET, "/channels/{channel.id}/thread-members/{user.id}")
	GetThreadMembers        = NewAPIRoute(GET, "/channels/{channel.id}/thread-members")

	GetArchivedPublicThreads        = NewAPIRoute(GET, "/channels/{channel.id}/threads/archived/public", "before", "limit")
	GetArchivedPrivateThreads       = NewAPIRoute(GET, "/channels/{channel.id}/threads/archived/private", "before", "limit")
	GetJoinedAchievedPrivateThreads = NewAPIRoute(GET, "/channels/{channel.id}/users/@me/threads/archived/private", "before", "limit")
)

// Messages
var (
	GetMessages        = NewAPIRoute(GET, "/channels/{channel.id}/messages")
	GetMessage         = NewAPIRoute(GET, "/channels/{channel.id}/messages/{message.id}")
	CreateMessage      = NewAPIRoute(POST, "/channels/{channel.id}/messages")
	UpdateMessage      = NewAPIRoute(PATCH, "/channels/{channel.id}/messages/{message.id}")
	DeleteMessage      = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message.id}")
	BulkDeleteMessages = NewAPIRoute(POST, "/channels/{channel.id}/messages/bulk-delete")

	GetPinnedMessages = NewAPIRoute(GET, "/channels/{channel.id}/pins")
	PinMessage        = NewAPIRoute(PUT, "/channels/{channel.id}/pins/{message.id}")
	UnpinMessage      = NewAPIRoute(DELETE, "/channels/{channel.id}/pins/{message.id}")

	CrosspostMessage = NewAPIRoute(POST, "/channels/{channel.id}/messages/{message.id}/crosspost")

	GetReactions               = NewAPIRoute(GET, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}", "after", "limit")
	AddReaction                = NewAPIRoute(PUT, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	RemoveOwnReaction          = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	RemoveUserReaction         = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}")
	RemoveAllReactions         = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message.id}/reactions")
	RemoveAllReactionsForEmoji = NewAPIRoute(DELETE, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}")
)

// Emojis
var (
	GetEmojis   = NewAPIRoute(GET, "/guilds/{guild.id}/emojis")
	GetEmoji    = NewAPIRoute(GET, "/guilds/{guild.id}/emojis/{emoji.id}")
	CreateEmoji = NewAPIRoute(POST, "/guilds/{guild.id}/emojis")
	UpdateEmoji = NewAPIRoute(PATCH, "/guilds/{guild.id}/emojis/{emote.id}")
	DeleteEmoji = NewAPIRoute(DELETE, "/guilds/{guild.id}/emojis/{emote.id}")
)

// Stickers
var (
	GetNitroStickerPacks = NewAPIRoute(GET, "/sticker-packs")
	GetSticker           = NewAPIRoute(GET, "/stickers/{sticker.id}")
	GetGuildStickers     = NewAPIRoute(GET, "/guilds/{guild.id}/stickers")
	CreateGuildSticker   = NewAPIRoute(POST, "/guilds/{guild.id}/stickers")
	UpdateGuildSticker   = NewAPIRoute(PATCH, "/guilds/{guild.id}/stickers/{sticker.id}")
	DeleteGuildSticker   = NewAPIRoute(DELETE, "/guilds/{guild.id}/stickers/{sticker.id}")
)

// Webhooks
var (
	GetWebhook    = NewAPIRoute(GET, "/webhooks/{webhook.id}")
	UpdateWebhook = NewAPIRoute(PATCH, "/webhooks/{webhook.id}")
	DeleteWebhook = NewAPIRoute(DELETE, "/webhooks/{webhook.id}")

	GetWebhookWithToken    = NewAPIRouteNoAuth(GET, "/webhooks/{webhook.id}/{webhook.token}")
	UpdateWebhookWithToken = NewAPIRouteNoAuth(PATCH, "/webhooks/{webhook.id}/{webhook.token}")
	DeleteWebhookWithToken = NewAPIRouteNoAuth(DELETE, "/webhooks/{webhook.id}/{webhook.token}")

	CreateWebhookMessage       = NewAPIRouteNoAuth(POST, "/webhooks/{webhook.id}/{webhook.token}", "wait", "thread_id")
	CreateWebhookMessageSlack  = NewAPIRouteNoAuth(POST, "/webhooks/{webhook.id}/{webhook.token}/slack", "wait", "thread_id")
	CreateWebhookMessageGitHub = NewAPIRouteNoAuth(POST, "/webhooks/{webhook.id}/{webhook.token}/github", "wait", "thread_id")
	UpdateWebhookMessage       = NewAPIRouteNoAuth(PATCH, "/webhooks/{webhook.id}/{webhook.token}/messages/{message.id}", "thread_id")
	DeleteWebhookMessage       = NewAPIRouteNoAuth(DELETE, "/webhooks/{webhook.id}/{webhook.token}/messages/{message.id}", "thread_id")
)

// Invites
var (
	GetInvite    = NewAPIRoute(GET, "/invites/{code}", "with_counts", "with_expiration", "guild_scheduled_event_id")
	CreateInvite = NewAPIRoute(POST, "/channels/{channel.id}/invites")
	DeleteInvite = NewAPIRoute(DELETE, "/invites/{code}")

	GetGuildInvites   = NewAPIRoute(GET, "/guilds/{guild.id}/invites")
	GetChannelInvites = NewAPIRoute(GET, "/channels/{channel.id}/invites")
)

// Interactions
var (
	GetGlobalCommands   = NewAPIRoute(GET, "/applications/{application.id}/commands", "with_localizations")
	GetGlobalCommand    = NewAPIRoute(GET, "/applications/{application.id}/command/{command.id}")
	CreateGlobalCommand = NewAPIRoute(POST, "/applications/{application.id}/commands")
	SetGlobalCommands   = NewAPIRoute(PUT, "/applications/{application.id}/commands")
	UpdateGlobalCommand = NewAPIRoute(PATCH, "/applications/{application.id}/commands/{command.id}")
	DeleteGlobalCommand = NewAPIRoute(DELETE, "/applications/{application.id}/commands/{command.id}")

	GetGuildCommands   = NewAPIRoute(GET, "/applications/{application.id}/guilds/{guild.id}/commands", "with_localizations")
	GetGuildCommand    = NewAPIRoute(GET, "/applications/{application.id}/guilds/{guild.id}/command/{command.id}")
	CreateGuildCommand = NewAPIRoute(POST, "/applications/{application.id}/guilds/{guild.id}/commands")
	SetGuildCommands   = NewAPIRoute(PUT, "/applications/{application.id}/guilds/{guild.id}/commands")
	UpdateGuildCommand = NewAPIRoute(PATCH, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}")
	DeleteGuildCommand = NewAPIRoute(DELETE, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}")

	GetGuildCommandsPermissions = NewAPIRoute(GET, "/applications/{application.id}/guilds/{guild.id}/commands/permissions")
	GetGuildCommandPermissions  = NewAPIRoute(GET, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions")
	SetGuildCommandsPermissions = NewAPIRoute(PUT, "/applications/{application.id}/guilds/{guild.id}/commands/permissions")
	SetGuildCommandPermissions  = NewAPIRoute(PUT, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions")

	GetInteractionResponse    = NewAPIRouteNoAuth(GET, "/webhooks/{application.id}/{interaction.token}/messages/@original")
	CreateInteractionResponse = NewAPIRouteNoAuth(POST, "/interactions/{interaction.id}/{interaction.token}/callback")
	UpdateInteractionResponse = NewAPIRouteNoAuth(PATCH, "/webhooks/{application.id}/{interaction.token}/messages/@original")
	DeleteInteractionResponse = NewAPIRouteNoAuth(DELETE, "/webhooks/{application.id}/{interaction.token}/messages/@original")

	GetFollowupMessage    = NewAPIRouteNoAuth(GET, "/webhooks/{application.id}/{interaction.token}")
	CreateFollowupMessage = NewAPIRouteNoAuth(POST, "/webhooks/{application.id}/{interaction.token}")
	UpdateFollowupMessage = NewAPIRouteNoAuth(PATCH, "/webhooks/{application.id}/{interaction.token}/messages/{message.id}")
	DeleteFollowupMessage = NewAPIRouteNoAuth(DELETE, "/webhooks/{application.id}/{interaction.token}/messages/{message.id}")
)

// CDN
var (
	CustomEmoji = NewCDNRoute("/emojis/{emote.id}", PNG, GIF)

	GuildIcon            = NewCDNRoute("/icons/{guild.id}/{guild.icon.hash}", PNG, JPEG, WebP, GIF)
	GuildSplash          = NewCDNRoute("/splashes/{guild.id}/{guild.splash.hash}", PNG, JPEG, WebP)
	GuildDiscoverySplash = NewCDNRoute("/discovery-splashes/{guild.id}/{guild.discovery.splash.hash}", PNG, JPEG, WebP)
	GuildBanner          = NewCDNRoute("/banners/{guild.id}/{guild.banner.hash}", PNG, JPEG, WebP, GIF)

	RoleIcon = NewCDNRoute("/role-icons/{role.id}/{role.icon.hash}", PNG, JPEG)

	UserBanner        = NewCDNRoute("/banners/{user.id}/{user.banner.hash}", PNG, JPEG, WebP, GIF)
	UserAvatar        = NewCDNRoute("/avatars/{user.id}/{user.avatar.hash}", PNG, JPEG, WebP, GIF)
	DefaultUserAvatar = NewCDNRoute("/embed/avatars/{user.discriminator%5}", PNG)

	ChannelIcon = NewCDNRoute("/channel-icons/{channel.id}/{channel.icon.hash}", PNG, JPEG, WebP)

	MemberAvatar = NewCDNRoute("/guilds/{guild.id}/users/{user.id}/avatars/{member.avatar.hash}", PNG, JPEG, WebP, GIF)

	ApplicationIcon  = NewCDNRoute("/app-icons/{application.id}/{icon.hash}", PNG, JPEG, WebP)
	ApplicationCover = NewCDNRoute("/app-assets/{application.id}/{cover.image.hash}", PNG, JPEG, WebP)
	ApplicationAsset = NewCDNRoute("/app-assets/{application.id}/{asset.id}", PNG, JPEG, WebP)

	AchievementIcon = NewCDNRoute("/app-assets/{application.id}/achievements/{achievement.id}/icons/{icon.hash}", PNG, JPEG, WebP)

	TeamIcon = NewCDNRoute("/team-icons/{team.id}/{team.icon.hash}", PNG, JPEG, WebP)

	StickerPackBanner = NewCDNRoute("app-assets/710982414301790216/store/{banner.asset.id}", PNG, JPEG, WebP)
	CustomSticker     = NewCDNRoute("/stickers/{sticker.id}", PNG, Lottie)

	Attachment = NewCDNRoute("/attachments/{channel.id}/{attachment.id}/{file.name}", BLANK)
)

// Other
var (
	InviteURL  = NewCustomRoute("https://discord.gg", "/{code}")
	WebhookURL = NewCustomRoute("https://discord.com", "/api/webhooks/{webhook.id}/{webhook.token}")
)
