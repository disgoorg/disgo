package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/disgoorg/disgo/discord"
)

var (
	// Version is the Discord API version DisGo should use
	Version = 10

	// API is the base path of the Discord API
	API = "https://discord.com/api/"
)

// MajorParameters is a list of url parameters which decide in which bucket a route belongs (https://discord.com/developers/docs/topics/rate-limits#rate-limits)
const MajorParameters = "guild.id:channel.id:webhook.id:interaction.token"

// Misc
var (
	GetGateway      = NewNoBotAuthEndpoint(http.MethodGet, "/gateway")
	GetGatewayBot   = NewEndpoint(http.MethodGet, "/gateway/bot")
	GetVoiceRegions = NewNoBotAuthEndpoint(http.MethodGet, "/voice/regions")
)

// OAuth2
var (
	GetBotApplicationInfo = NewEndpoint(http.MethodGet, "/oauth2/applications/@me")
	GetAuthorizationInfo  = NewNoBotAuthEndpoint(http.MethodGet, "/oauth2/@me")
	Token                 = NewEndpoint(http.MethodPost, "/oauth2/token")
)

// Users
var (
	GetUser                                    = NewEndpoint(http.MethodGet, "/users/{user.id}")
	GetCurrentUser                             = NewEndpoint(http.MethodGet, "/users/@me")
	UpdateCurrentUser                          = NewEndpoint(http.MethodPatch, "/users/@me")
	GetCurrentUserGuilds                       = NewEndpoint(http.MethodGet, "/users/@me/guilds")
	GetCurrentMember                           = NewNoBotAuthEndpoint(http.MethodGet, "/users/@me/guilds/{guild.id}/member")
	GetCurrentUserConnections                  = NewNoBotAuthEndpoint(http.MethodGet, "/users/@me/connections")
	GetCurrentUserApplicationRoleConnection    = NewNoBotAuthEndpoint(http.MethodGet, "/users/@me/applications/{application.id}/role-connection")
	UpdateCurrentUserApplicationRoleConnection = NewNoBotAuthEndpoint(http.MethodPut, "/users/@me/applications/{application.id}/role-connection")
	LeaveGuild                                 = NewEndpoint(http.MethodDelete, "/users/@me/guilds/{guild.id}")
	CreateDMChannel                            = NewEndpoint(http.MethodPost, "/users/@me/channels")
)

// Guilds
var (
	GetGuild          = NewEndpoint(http.MethodGet, "/guilds/{guild.id}")
	GetGuildPreview   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/preview")
	CreateGuild       = NewEndpoint(http.MethodPost, "/guilds")
	UpdateGuild       = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}")
	DeleteGuild       = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}")
	GetGuildVanityURL = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/vanity-url")

	CreateGuildChannel     = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/channels")
	GetGuildChannels       = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/channels")
	UpdateChannelPositions = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/channels")

	GetBans   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/bans")
	GetBan    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/bans/{user.id}")
	AddBan    = NewEndpoint(http.MethodPut, "/guilds/{guild.id}/bans/{user.id}")
	DeleteBan = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/bans/{user.id}")
	BulkBan   = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/bulk-ban")

	GetMember        = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/members/{user.id}")
	GetMembers       = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/members")
	SearchMembers    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/members/search")
	AddMember        = NewEndpoint(http.MethodPut, "/guilds/{guild.id}/members/{user.id}")
	UpdateMember     = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/members/{user.id}")
	RemoveMember     = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/members/{user.id}")
	AddMemberRole    = NewEndpoint(http.MethodPut, "/guilds/{guild.id}/members/{user.id}/roles/{role.id}")
	RemoveMemberRole = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/members/{user.id}/roles/{role.id}")

	UpdateCurrentMember = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/members/@me")

	GetGuildPruneCount = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/prune")
	BeginGuildPrune    = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/prune")

	GetGuildWebhooks = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/webhooks")

	GetAuditLogs = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/audit-logs")

	GetGuildVoiceRegions = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/regions")

	GetGuildWelcomeScreen    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/welcome-screen")
	UpdateGuildWelcomeScreen = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/welcome-screen")

	GetGuildOnboarding    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/onboarding")
	UpdateGuildOnboarding = NewEndpoint(http.MethodPut, "/guilds/{guild.id}/onboarding")

	UpdateGuildIncidentActions = NewEndpoint(http.MethodPut, "/guilds/{guild.id}/incident-actions")

	GetCurrentUserVoiceState    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/voice-states/@me")
	GetUserVoiceState           = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/voice-states/{user.id}")
	UpdateCurrentUserVoiceState = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/voice-states/@me")
	UpdateUserVoiceState        = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/voice-states/{user.id}")
)

// AutoModeration
var (
	GetAutoModerationRules   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/auto-moderation/rules")
	GetAutoModerationRule    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}")
	CreateAutoModerationRule = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/auto-moderation/rules")
	UpdateAutoModerationRule = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}")
	DeleteAutoModerationRule = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}")
)

// GuildIntegrations
var (
	GetIntegrations   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/integrations")
	CreateIntegration = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/integrations")
	UpdateIntegration = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/integrations/{integration.id}")
	DeleteIntegration = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/integrations/{integration.id}")
	SyncIntegration   = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/integrations/{integration.id}/sync")
)

// GuildTemplates
var (
	GetGuildTemplate        = NewEndpoint(http.MethodGet, "/guilds/templates/{template.code}")
	GetGuildTemplates       = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/templates")
	CreateGuildTemplate     = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/templates")
	SyncGuildTemplate       = NewEndpoint(http.MethodPut, "/guilds/{guild.id}/templates/{template.code}")
	UpdateGuildTemplate     = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/templates/{template.code}")
	DeleteGuildTemplate     = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/templates/{template.code}")
	CreateGuildFromTemplate = NewEndpoint(http.MethodPost, "/guilds/templates/{template.code}")
)

// GuildScheduledEvents
var (
	GetGuildScheduledEvents   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/scheduled-events")
	GetGuildScheduledEvent    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}")
	CreateGuildScheduledEvent = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/scheduled-events")
	UpdateGuildScheduledEvent = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}")
	DeleteGuildScheduledEvent = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}")

	GetGuildScheduledEventUsers = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users")
)

// Sounds
var (
	GetSoundboardDefaultSounds = NewEndpoint(http.MethodGet, "/soundboard-default-sounds")
	GetGuildSoundboardSounds   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/soundboard-sounds")
	CreateGuildSoundboardSound = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/soundboard-sounds")
	GetGuildSoundboardSound    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/soundboard-sounds/{sound.id}")
	UpdateGuildSoundboardSound = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/soundboard-sounds/{sound.id}")
	DeleteGuildSoundboardSound = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/soundboard-sounds/{sound.id}")
)

// StageInstance
var (
	GetStageInstance    = NewEndpoint(http.MethodGet, "/stage-instances/{channel.id}")
	CreateStageInstance = NewEndpoint(http.MethodPost, "/stage-instances")
	UpdateStageInstance = NewEndpoint(http.MethodPatch, "/stage-instances/{channel.id}")
	DeleteStageInstance = NewEndpoint(http.MethodDelete, "/stage-instances/{channel.id}")
)

// Roles
var (
	GetRoles            = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/roles")
	GetRole             = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/roles/{role.id}")
	CreateRole          = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/roles")
	UpdateRole          = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/roles/{role.id}")
	UpdateRolePositions = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/roles")
	DeleteRole          = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/roles/{role.id}")
)

// Channels
var (
	GetChannel    = NewEndpoint(http.MethodGet, "/channels/{channel.id}")
	UpdateChannel = NewEndpoint(http.MethodPatch, "/channels/{channel.id}")
	DeleteChannel = NewEndpoint(http.MethodDelete, "/channels/{channel.id}")

	GetChannelWebhooks = NewEndpoint(http.MethodGet, "/channels/{channel.id}/webhooks")
	CreateWebhook      = NewEndpoint(http.MethodPost, "/channels/{channel.id}/webhooks")

	UpdatePermissionOverwrite = NewEndpoint(http.MethodPut, "/channels/{channel.id}/permissions/{overwrite.id}")
	DeletePermissionOverwrite = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/permissions/{overwrite.id}")

	SendTyping    = NewEndpoint(http.MethodPost, "/channels/{channel.id}/typing")
	FollowChannel = NewEndpoint(http.MethodPost, "/channels/{channel.id}/followers")

	GetPollAnswerVotes = NewEndpoint(http.MethodGet, "/channels/{channel.id}/polls/{message.id}/answers/{answer.id}")
	ExpirePoll         = NewEndpoint(http.MethodPost, "/channels/{channel.id}/polls/{message.id}/expire")

	SendSoundboardSound = NewEndpoint(http.MethodPost, "/channels/{channel.id}/send-soundboard-sound")
)

// Threads
var (
	CreateThreadWithMessage = NewEndpoint(http.MethodPost, "/channels/{channel.id}/messages/{message.id}/threads")
	CreateThread            = NewEndpoint(http.MethodPost, "/channels/{channel.id}/threads")
	JoinThread              = NewEndpoint(http.MethodPut, "/channels/{channel.id}/thread-members/@me")
	LeaveThread             = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/thread-members/@me")
	AddThreadMember         = NewEndpoint(http.MethodPut, "/channels/{channel.id}/thread-members/{user.id}")
	RemoveThreadMember      = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/thread-members/{user.id}")
	GetThreadMember         = NewEndpoint(http.MethodGet, "/channels/{channel.id}/thread-members/{user.id}")
	GetThreadMembers        = NewEndpoint(http.MethodGet, "/channels/{channel.id}/thread-members")

	GetPublicArchivedThreads        = NewEndpoint(http.MethodGet, "/channels/{channel.id}/threads/archived/public")
	GetPrivateArchivedThreads       = NewEndpoint(http.MethodGet, "/channels/{channel.id}/threads/archived/private")
	GetJoinedPrivateArchivedThreads = NewEndpoint(http.MethodGet, "/channels/{channel.id}/users/@me/threads/archived/private")
	GetActiveGuildThreads           = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/threads/active")
)

// Messages
var (
	GetMessages        = NewEndpoint(http.MethodGet, "/channels/{channel.id}/messages")
	GetMessage         = NewEndpoint(http.MethodGet, "/channels/{channel.id}/messages/{message.id}")
	CreateMessage      = NewEndpoint(http.MethodPost, "/channels/{channel.id}/messages")
	UpdateMessage      = NewEndpoint(http.MethodPatch, "/channels/{channel.id}/messages/{message.id}")
	DeleteMessage      = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/messages/{message.id}")
	BulkDeleteMessages = NewEndpoint(http.MethodPost, "/channels/{channel.id}/messages/bulk-delete")

	GetPinnedMessages = NewEndpoint(http.MethodGet, "/channels/{channel.id}/pins")
	PinMessage        = NewEndpoint(http.MethodPut, "/channels/{channel.id}/pins/{message.id}")
	UnpinMessage      = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/pins/{message.id}")

	CrosspostMessage = NewEndpoint(http.MethodPost, "/channels/{channel.id}/messages/{message.id}/crosspost")

	GetReactions               = NewEndpoint(http.MethodGet, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}")
	AddReaction                = NewEndpoint(http.MethodPut, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	RemoveOwnReaction          = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me")
	RemoveUserReaction         = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}")
	RemoveAllReactions         = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/messages/{message.id}/reactions")
	RemoveAllReactionsForEmoji = NewEndpoint(http.MethodDelete, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}")
)

// Emojis
var (
	GetEmojis   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/emojis")
	GetEmoji    = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/emojis/{emoji.id}")
	CreateEmoji = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/emojis")
	UpdateEmoji = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/emojis/{emote.id}")
	DeleteEmoji = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/emojis/{emote.id}")
)

// Stickers
var (
	GetNitroStickerPacks = NewEndpoint(http.MethodGet, "/sticker-packs")
	GetNitroStickerPack  = NewEndpoint(http.MethodGet, "/sticker-packs/{pack.id}")
	GetSticker           = NewEndpoint(http.MethodGet, "/stickers/{sticker.id}")
	GetGuildStickers     = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/stickers")
	CreateGuildSticker   = NewEndpoint(http.MethodPost, "/guilds/{guild.id}/stickers")
	UpdateGuildSticker   = NewEndpoint(http.MethodPatch, "/guilds/{guild.id}/stickers/{sticker.id}")
	DeleteGuildSticker   = NewEndpoint(http.MethodDelete, "/guilds/{guild.id}/stickers/{sticker.id}")
)

// Webhooks
var (
	GetWebhook    = NewEndpoint(http.MethodGet, "/webhooks/{webhook.id}")
	UpdateWebhook = NewEndpoint(http.MethodPatch, "/webhooks/{webhook.id}")
	DeleteWebhook = NewEndpoint(http.MethodDelete, "/webhooks/{webhook.id}")

	GetWebhookWithToken    = NewNoBotAuthEndpoint(http.MethodGet, "/webhooks/{webhook.id}/{webhook.token}")
	UpdateWebhookWithToken = NewNoBotAuthEndpoint(http.MethodPatch, "/webhooks/{webhook.id}/{webhook.token}")
	DeleteWebhookWithToken = NewNoBotAuthEndpoint(http.MethodDelete, "/webhooks/{webhook.id}/{webhook.token}")

	GetWebhookMessage          = NewNoBotAuthEndpoint(http.MethodGet, "/webhooks/{webhook.id}/{webhook.token}/messages/{message.id}")
	CreateWebhookMessage       = NewNoBotAuthEndpoint(http.MethodPost, "/webhooks/{webhook.id}/{webhook.token}")
	CreateWebhookMessageSlack  = NewNoBotAuthEndpoint(http.MethodPost, "/webhooks/{webhook.id}/{webhook.token}/slack")
	CreateWebhookMessageGitHub = NewNoBotAuthEndpoint(http.MethodPost, "/webhooks/{webhook.id}/{webhook.token}/github")
	UpdateWebhookMessage       = NewNoBotAuthEndpoint(http.MethodPatch, "/webhooks/{webhook.id}/{webhook.token}/messages/{message.id}")
	DeleteWebhookMessage       = NewNoBotAuthEndpoint(http.MethodDelete, "/webhooks/{webhook.id}/{webhook.token}/messages/{message.id}")
)

// Invites
var (
	GetInvite    = NewEndpoint(http.MethodGet, "/invites/{code}")
	CreateInvite = NewEndpoint(http.MethodPost, "/channels/{channel.id}/invites")
	DeleteInvite = NewEndpoint(http.MethodDelete, "/invites/{code}")

	GetGuildInvites   = NewEndpoint(http.MethodGet, "/guilds/{guild.id}/invites")
	GetChannelInvites = NewEndpoint(http.MethodGet, "/channels/{channel.id}/invites")
)

// Applications
var (
	GetCurrentApplication    = NewEndpoint(http.MethodGet, "/applications/@me")
	UpdateCurrentApplication = NewEndpoint(http.MethodPatch, "/applications/@me")

	GetGlobalCommands   = NewEndpoint(http.MethodGet, "/applications/{application.id}/commands")
	GetGlobalCommand    = NewEndpoint(http.MethodGet, "/applications/{application.id}/command/{command.id}")
	CreateGlobalCommand = NewEndpoint(http.MethodPost, "/applications/{application.id}/commands")
	SetGlobalCommands   = NewEndpoint(http.MethodPut, "/applications/{application.id}/commands")
	UpdateGlobalCommand = NewEndpoint(http.MethodPatch, "/applications/{application.id}/commands/{command.id}")
	DeleteGlobalCommand = NewEndpoint(http.MethodDelete, "/applications/{application.id}/commands/{command.id}")

	GetGuildCommands   = NewEndpoint(http.MethodGet, "/applications/{application.id}/guilds/{guild.id}/commands")
	GetGuildCommand    = NewEndpoint(http.MethodGet, "/applications/{application.id}/guilds/{guild.id}/command/{command.id}")
	CreateGuildCommand = NewEndpoint(http.MethodPost, "/applications/{application.id}/guilds/{guild.id}/commands")
	SetGuildCommands   = NewEndpoint(http.MethodPut, "/applications/{application.id}/guilds/{guild.id}/commands")
	UpdateGuildCommand = NewEndpoint(http.MethodPatch, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}")
	DeleteGuildCommand = NewEndpoint(http.MethodDelete, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}")

	GetGuildCommandsPermissions = NewEndpoint(http.MethodGet, "/applications/{application.id}/guilds/{guild.id}/commands/permissions")
	GetGuildCommandPermissions  = NewEndpoint(http.MethodGet, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions")
	SetGuildCommandPermissions  = NewNoBotAuthEndpoint(http.MethodPut, "/applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions")

	GetInteractionResponse    = NewNoBotAuthEndpoint(http.MethodGet, "/webhooks/{application.id}/{interaction.token}/messages/@original")
	CreateInteractionResponse = NewNoBotAuthEndpoint(http.MethodPost, "/interactions/{interaction.id}/{interaction.token}/callback")
	UpdateInteractionResponse = NewNoBotAuthEndpoint(http.MethodPatch, "/webhooks/{application.id}/{interaction.token}/messages/@original")
	DeleteInteractionResponse = NewNoBotAuthEndpoint(http.MethodDelete, "/webhooks/{application.id}/{interaction.token}/messages/@original")

	GetFollowupMessage    = NewNoBotAuthEndpoint(http.MethodGet, "/webhooks/{application.id}/{interaction.token}/messages/{message.id}")
	CreateFollowupMessage = NewNoBotAuthEndpoint(http.MethodPost, "/webhooks/{application.id}/{interaction.token}")
	UpdateFollowupMessage = NewNoBotAuthEndpoint(http.MethodPatch, "/webhooks/{application.id}/{interaction.token}/messages/{message.id}")
	DeleteFollowupMessage = NewNoBotAuthEndpoint(http.MethodDelete, "/webhooks/{application.id}/{interaction.token}/messages/{message.id}")

	GetApplicationRoleConnectionMetadata    = NewEndpoint(http.MethodGet, "/applications/{application.id}/role-connections/metadata")
	UpdateApplicationRoleConnectionMetadata = NewEndpoint(http.MethodPut, "/applications/{application.id}/role-connections/metadata")

	GetEntitlements       = NewEndpoint(http.MethodGet, "/applications/{application.id}/entitlements")
	GetEntitlement        = NewEndpoint(http.MethodGet, "/applications/{application.id}/entitlements/{entitlement.id}")
	CreateTestEntitlement = NewEndpoint(http.MethodPost, "/applications/{application.id}/entitlements")
	DeleteTestEntitlement = NewEndpoint(http.MethodDelete, "/applications/{application.id}/entitlements/{entitlement.id}")
	ConsumeEntitlement    = NewEndpoint(http.MethodPost, "/applications/{application.id}/entitlements/{entitlement.id}/consume")

	GetApplicationEmojis   = NewEndpoint(http.MethodGet, "/applications/{application.id}/emojis")
	GetApplicationEmoji    = NewEndpoint(http.MethodGet, "/applications/{application.id}/emojis/{emoji.id}")
	CreateApplicationEmoji = NewEndpoint(http.MethodPost, "/applications/{application.id}/emojis")
	UpdateApplicationEmoji = NewEndpoint(http.MethodPatch, "/applications/{application.id}/emojis/{emoji.id}")
	DeleteApplicationEmoji = NewEndpoint(http.MethodDelete, "/applications/{application.id}/emojis/{emoji.id}")

	GetActivityInstance = NewEndpoint(http.MethodGet, "/applications/{application.id}/activity-instances/{instance.id}")
)

// SKUs
var (
	GetSKUs = NewEndpoint(http.MethodGet, "/applications/{application.id}/skus")

	GetSKUSubscriptions = NewEndpoint(http.MethodGet, "/skus/{sku.id}/subscriptions")
	GetSKUSubscription  = NewEndpoint(http.MethodGet, "/skus/{sku.id}/subscriptions/{subscription.id}")
)

// NewEndpoint returns a new Endpoint which requires bot auth with the given http method & route.
func NewEndpoint(method string, route string) *Endpoint {
	return &Endpoint{
		Method:  method,
		Route:   route,
		BotAuth: true,
	}
}

// NewNoBotAuthEndpoint returns a new Endpoint which does not require bot auth with the given http method & route.
func NewNoBotAuthEndpoint(method string, route string) *Endpoint {
	return &Endpoint{
		Method:  method,
		Route:   route,
		BotAuth: false,
	}
}

// Endpoint represents a Discord Rest API endpoint.
type Endpoint struct {
	Method  string
	Route   string
	BotAuth bool
}

// CompiledEndpoint represents a Discord Rest API endpoint with applied url params & query values.
type CompiledEndpoint struct {
	Endpoint *Endpoint

	URL         string
	MajorParams string
}

// Compile compiles an Endpoint to a CompiledEndpoint with the given url params & query values
func (e *Endpoint) Compile(values discord.QueryValues, params ...any) *CompiledEndpoint {
	var majorParams []string
	path := e.Route
	for _, param := range params {
		start := strings.Index(path, "{")
		end := strings.Index(path, "}")
		if start == -1 || end == -1 {
			break
		}
		paramName := path[start+1 : end]
		paramValue := fmt.Sprint(param)
		if strings.Contains(MajorParameters, paramName) {
			majorParams = append(majorParams, paramName+"="+paramValue)
		}
		path = path[:start] + paramValue + path[end+1:]
	}

	query := values.Encode()
	if query != "" {
		query = "?" + query
	}

	return &CompiledEndpoint{
		Endpoint:    e,
		URL:         path + query,
		MajorParams: strings.Join(majorParams, ":"),
	}
}
