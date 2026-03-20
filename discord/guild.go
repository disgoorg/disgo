package discord

import (
	"strings"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/internal/flags"
)

// PremiumTier tells you the boost level of a Guild
type PremiumTier int

// Constants for PremiumTier
const (
	PremiumTierNone PremiumTier = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

// SystemChannelFlags contains the settings for the Guild(s) system channel
type SystemChannelFlags int

// Constants for SystemChannelFlags
const (
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlags = 1 << iota
	SystemChannelFlagSuppressPremiumSubscriptions
	SystemChannelFlagSuppressGuildReminderNotifications
	SystemChannelFlagSuppressJoinNotificationReplies
	SystemChannelFlagSuppressRoleSubscriptionPurchaseNotifications
	SystemChannelFlagSuppressRoleSubscriptionPurchaseNotificationReplies
)

// Add allows you to add multiple bits together, producing a new bit
func (f SystemChannelFlags) Add(bits ...SystemChannelFlags) SystemChannelFlags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f SystemChannelFlags) Remove(bits ...SystemChannelFlags) SystemChannelFlags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f SystemChannelFlags) Has(bits ...SystemChannelFlags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f SystemChannelFlags) Missing(bits ...SystemChannelFlags) bool {
	return flags.Missing(f, bits...)
}

// The VerificationLevel of a Guild that members must be to send messages
type VerificationLevel int

// Constants for VerificationLevel
const (
	VerificationLevelNone VerificationLevel = iota
	VerificationLevelLow
	VerificationLevelMedium
	VerificationLevelHigh
	VerificationLevelVeryHigh
)

// MessageNotificationsLevel indicates whether users receive @ mentions on a new message
type MessageNotificationsLevel int

// Constants for MessageNotificationsLevel
const (
	MessageNotificationsLevelAllMessages MessageNotificationsLevel = iota
	MessageNotificationsLevelOnlyMentions
)

// The ExplicitContentFilterLevel of a Guild
type ExplicitContentFilterLevel int

// Constants for ExplicitContentFilterLevel
const (
	ExplicitContentFilterLevelDisabled ExplicitContentFilterLevel = iota
	ExplicitContentFilterLevelMembersWithoutRoles
	ExplicitContentFilterLevelAllMembers
)

// The MFALevel of a Guild
type MFALevel int

// Constants for MFALevel
const (
	MFALevelNone MFALevel = iota
	MFALevelElevated
)

// The GuildFeature (s) that a Guild contains
type GuildFeature string

// Constants for GuildFeature
const (
	GuildFeatureAnimatedBanner                        GuildFeature = "ANIMATED_BANNER"
	GuildFeatureAnimatedIcon                          GuildFeature = "ANIMATED_ICON"
	GuildFeatureAutoModeration                        GuildFeature = "AUTO_MODERATION"
	GuildFeatureBanner                                GuildFeature = "BANNER"
	GuildFeatureCommunity                             GuildFeature = "COMMUNITY"
	GuildFeatureCreatorMonetizableProvisional         GuildFeature = "CREATOR_MONETIZABLE_PROVISIONAL"
	GuildFeatureCreatorStorePage                      GuildFeature = "CREATOR_STORE_PAGE"
	GuildFeatureDeveloperSupportServer                GuildFeature = "DEVELOPER_SUPPORT_SERVER"
	GuildFeatureDiscoverable                          GuildFeature = "DISCOVERABLE"
	GuildFeatureEnhancedRoleColors                    GuildFeature = "ENHANCED_ROLE_COLORS"
	GuildFeatureFeaturable                            GuildFeature = "FEATURABLE"
	GuildFeatureGuestsEnabled                         GuildFeature = "GUESTS_ENABLED"
	GuildFeatureGuildTags                             GuildFeature = "GUILD_TAGS"
	GuildFeatureInvitesDisabled                       GuildFeature = "INVITES_DISABLED"
	GuildFeatureInviteSplash                          GuildFeature = "INVITE_SPLASH"
	GuildFeatureMemberVerificationGateEnabled         GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeatureMoreSoundboard                        GuildFeature = "MORE_SOUNDBOARD"
	GuildFeatureMoreStickers                          GuildFeature = "MORE_STICKERS"
	GuildFeatureNews                                  GuildFeature = "NEWS"
	GuildFeaturePartnered                             GuildFeature = "PARTNERED"
	GuildFeaturePreviewEnabled                        GuildFeature = "PREVIEW_ENABLED"
	GuildFeatureRaidAlertsDisabled                    GuildFeature = "RAID_ALERTS_DISABLED"
	GuildFeatureRoleIcons                             GuildFeature = "ROLE_ICONS"
	GuildFeatureRoleSubscriptionsAvailableForPurchase GuildFeature = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"
	GuildFeatureRoleSubscriptionsEnabled              GuildFeature = "ROLE_SUBSCRIPTIONS_ENABLED"
	GuildFeatureSoundboard                            GuildFeature = "SOUNDBOARD"
	GuildFeatureTicketedEventsEnabled                 GuildFeature = "TICKETED_EVENTS_ENABLED"
	GuildFeatureVanityURL                             GuildFeature = "VANITY_URL"
	GuildFeatureVerified                              GuildFeature = "VERIFIED"
	GuildFeatureVipRegions                            GuildFeature = "VIP_REGIONS"
	GuildFeatureWelcomeScreenEnabled                  GuildFeature = "WELCOME_SCREEN_ENABLED"
)

// Guild represents a discord Guild
type Guild struct {
	ID                          snowflake.ID               `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        *string                    `json:"icon"`
	Splash                      *string                    `json:"splash"`
	DiscoverySplash             *string                    `json:"discovery_splash"`
	OwnerID                     snowflake.ID               `json:"owner_id"`
	AfkChannelID                *snowflake.ID              `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             snowflake.ID               `json:"widget_channel_id"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	DefaultMessageNotifications MessageNotificationsLevel  `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Features                    []GuildFeature             `json:"features"`
	MFALevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               *snowflake.ID              `json:"application_id"`
	SystemChannelID             *snowflake.ID              `json:"system_channel_id"`
	SystemChannelFlags          SystemChannelFlags         `json:"system_channel_flags"`
	RulesChannelID              *snowflake.ID              `json:"rules_channel_id"`
	MemberCount                 int                        `json:"member_count"`
	MaxPresences                *int                       `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	VanityURLCode               *string                    `json:"vanity_url_code"`
	Description                 *string                    `json:"description"`
	Banner                      *string                    `json:"banner"`
	PremiumTier                 PremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      *snowflake.ID              `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	MaxStageVideoChannelUsers   int                        `json:"max_stage_video_channel_users"`
	WelcomeScreen               GuildWelcomeScreen         `json:"welcome_screen"`
	NSFWLevel                   NSFWLevel                  `json:"nsfw_level"`
	PremiumProgressBarEnabled   bool                       `json:"premium_progress_bar_enabled"`
	JoinedAt                    time.Time                  `json:"joined_at"`
	SafetyAlertsChannelID       *snowflake.ID              `json:"safety_alerts_channel_id"`
	IncidentsData               *GuildIncidentsData        `json:"incidents_data"`

	// only over GET /guilds/{guild.id}
	ApproximateMemberCount   int `json:"approximate_member_count"`
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

func (g Guild) IconURL(opts ...CDNOpt) *string {
	if g.Icon == nil {
		return nil
	}
	url := formatAssetURL(GuildIcon, opts, g.ID, *g.Icon)
	return &url
}

func (g Guild) SplashURL(opts ...CDNOpt) *string {
	if g.Splash == nil {
		return nil
	}
	url := formatAssetURL(GuildSplash, opts, g.ID, *g.Splash)
	return &url
}

func (g Guild) DiscoverySplashURL(opts ...CDNOpt) *string {
	if g.DiscoverySplash == nil {
		return nil
	}
	url := formatAssetURL(GuildDiscoverySplash, opts, g.ID, *g.DiscoverySplash)
	return &url
}

func (g Guild) BannerURL(opts ...CDNOpt) *string {
	if g.Banner == nil {
		return nil
	}
	url := formatAssetURL(GuildBanner, opts, g.ID, *g.Banner)
	return &url
}

func (g Guild) CreatedAt() time.Time {
	return g.ID.Time()
}

type RestGuild struct {
	Guild
	Stickers []Sticker `json:"stickers"`
	Roles    []Role    `json:"roles"`
	Emojis   []Emoji   `json:"emojis"`
}

type GatewayGuild struct {
	RestGuild
	Large                bool                  `json:"large"`
	Unavailable          bool                  `json:"unavailable"`
	VoiceStates          []VoiceState          `json:"voice_states"`
	Members              []Member              `json:"members"`
	Channels             []GuildChannel        `json:"channels"`
	Threads              []GuildThread         `json:"threads"`
	Presences            []Presence            `json:"presences"`
	StageInstances       []StageInstance       `json:"stage_instances"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	SoundboardSounds     []SoundboardSound     `json:"soundboard_sounds"`
}

func (g *GatewayGuild) UnmarshalJSON(data []byte) error {
	type gatewayGuild GatewayGuild
	var v struct {
		Channels []UnmarshalChannel `json:"channels"`
		gatewayGuild
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*g = GatewayGuild(v.gatewayGuild)

	g.Channels = make([]GuildChannel, len(v.Channels))
	for i := range v.Channels {
		g.Channels[i] = v.Channels[i].Channel.(GuildChannel)
	}

	return nil
}

type UnavailableGuild struct {
	ID          snowflake.ID `json:"id"`
	Unavailable bool         `json:"unavailable"`
}

// OAuth2Guild is returned on the GetGuilds route
type OAuth2Guild struct {
	ID                       snowflake.ID   `json:"id"`
	Name                     string         `json:"name"`
	Icon                     *string        `json:"icon"`
	Banner                   *string        `json:"banner"`
	Owner                    bool           `json:"owner"`
	Permissions              Permissions    `json:"permissions"`
	Features                 []GuildFeature `json:"features"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
}

func (g OAuth2Guild) IconURL(opts ...CDNOpt) *string {
	if g.Icon == nil {
		return nil
	}
	url := formatAssetURL(GuildIcon, opts, g.ID, *g.Icon)
	return &url
}

func (g OAuth2Guild) BannerURL(opts ...CDNOpt) *string {
	if g.Banner == nil {
		return nil
	}
	url := formatAssetURL(GuildBanner, opts, g.ID, *g.Banner)
	return &url
}

// GuildWelcomeScreen is the Welcome Screen of a Guild
type GuildWelcomeScreen struct {
	Description     *string               `json:"description,omitempty"`
	WelcomeChannels []GuildWelcomeChannel `json:"welcome_channels"`
}

// GuildWelcomeChannel is one of the channels in a GuildWelcomeScreen
type GuildWelcomeChannel struct {
	ChannelID   snowflake.ID  `json:"channel_id"`
	Description string        `json:"description"`
	EmojiID     *snowflake.ID `json:"emoji_id,omitempty"`
	EmojiName   *string       `json:"emoji_name,omitempty"`
}

// GuildWelcomeScreenUpdate is used to update the GuildWelcomeScreen of a Guild
type GuildWelcomeScreenUpdate struct {
	Enabled         *bool                  `json:"enabled,omitempty"`
	WelcomeChannels *[]GuildWelcomeChannel `json:"welcome_channels,omitempty"`
	Description     *string                `json:"description,omitempty"`
}

// GuildPreview is used for previewing public Guild(s) before joining them
type GuildPreview struct {
	ID                       snowflake.ID   `json:"id"`
	Name                     string         `json:"name"`
	Icon                     *string        `json:"icon"`
	DiscoverySplash          *string        `json:"discovery_splash"`
	Splash                   *string        `json:"splash"`
	Features                 []GuildFeature `json:"features"`
	Description              *string        `json:"description"`
	ApproximateMemberCount   *int           `json:"approximate_member_count"`
	ApproximatePresenceCount *int           `json:"approximate_presence_count"`
	Emojis                   []Emoji        `json:"emojis"`
	Stickers                 []Sticker      `json:"stickers"`
}

type GuildIncidentsData struct {
	InvitesDisabledUntil *time.Time `json:"invites_disabled_until"`
	DMsDisabledUntil     *time.Time `json:"dms_disabled_until"`
	DMSpamDetectedAt     *time.Time `json:"dm_spam_detected_at"`
	RaidDetectedAt       *time.Time `json:"raid_detected_at"`
}

type GuildIncidentActionsUpdate struct {
	InvitesDisabledUntil omit.Omit[*time.Time] `json:"invites_disabled_until,omitzero"`
	DMsDisabledUntil     omit.Omit[*time.Time] `json:"dms_disabled_until,omitzero"`
}

// GuildUpdate is the payload used to update a Guild
type GuildUpdate struct {
	Name                        *string                                `json:"name,omitempty"`
	VerificationLevel           omit.Omit[*VerificationLevel]          `json:"verification_level,omitzero"`
	DefaultMessageNotifications omit.Omit[*MessageNotificationsLevel]  `json:"default_message_notification,omitzero"`
	ExplicitContentFilter       omit.Omit[*ExplicitContentFilterLevel] `json:"explicit_content_filter,omitzero"`
	AFKChannelID                *snowflake.ID                          `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int                                   `json:"afk_timeout,omitempty"`
	Icon                        omit.Omit[*Icon]                       `json:"icon,omitzero"`
	Splash                      omit.Omit[*Icon]                       `json:"splash,omitzero"`
	DiscoverySplash             omit.Omit[*Icon]                       `json:"discovery_splash,omitzero"`
	Banner                      omit.Omit[*Icon]                       `json:"banner,omitzero"`
	SystemChannelID             *snowflake.ID                          `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *SystemChannelFlags                    `json:"system_channel_flags,omitempty"`
	RulesChannelID              *snowflake.ID                          `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *snowflake.ID                          `json:"public_updates_channel_id,omitempty"`
	SafetyAlertsChannelID       *snowflake.ID                          `json:"safety_alerts_channel_id,omitempty"`
	PreferredLocale             *string                                `json:"preferred_locale,omitempty"`
	Features                    *[]GuildFeature                        `json:"features,omitempty"`
	Description                 *string                                `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                                  `json:"premium_progress_bar_enabled,omitempty"`
}

type NSFWLevel int

const (
	NSFWLevelDefault NSFWLevel = iota
	NSFWLevelExplicit
	NSFWLevelSafe
	NSFWLevelAgeRestricted
)

type GuildCreateRole struct {
	RoleCreate
	ID int `json:"id,omitempty"`
}

type GuildCreateChannel struct {
	ChannelCreate
	ID       int `json:"id,omitempty"`
	ParentID int `json:"parent_id,omitempty"`
}

type GuildPrune struct {
	Days              int            `json:"days"`
	ComputePruneCount bool           `json:"compute_prune_count"`
	IncludeRoles      []snowflake.ID `json:"include_roles"`
}

type GuildPruneResult struct {
	Pruned *int `json:"pruned"`
}

type GuildActiveThreads struct {
	Threads []GuildThread  `json:"threads"`
	Members []ThreadMember `json:"members"`
}

// MessageSearchAuthorType is the author type filter for guild message search.
type MessageSearchAuthorType string

// Constants for MessageSearchAuthorType.
const (
	MessageSearchAuthorTypeUser    MessageSearchAuthorType = "user"
	MessageSearchAuthorTypeBot     MessageSearchAuthorType = "bot"
	MessageSearchAuthorTypeWebhook MessageSearchAuthorType = "webhook"

	// Negated types, results will NOT include messages from these author types.
	MessageSearchAuthorTypeNotUser    MessageSearchAuthorType = "-user"
	MessageSearchAuthorTypeNotBot     MessageSearchAuthorType = "-bot"
	MessageSearchAuthorTypeNotWebhook MessageSearchAuthorType = "-webhook"
)

// MessageSearchHasType filters messages by whether they contain a specific type of content.
type MessageSearchHasType string

// Constants for MessageSearchHasType.
const (
	MessageSearchHasTypeImage    MessageSearchHasType = "image"
	MessageSearchHasTypeSound    MessageSearchHasType = "sound"
	MessageSearchHasTypeVideo    MessageSearchHasType = "video"
	MessageSearchHasTypeFile     MessageSearchHasType = "file"
	MessageSearchHasTypeSticker  MessageSearchHasType = "sticker"
	MessageSearchHasTypeEmbed    MessageSearchHasType = "embed"
	MessageSearchHasTypeLink     MessageSearchHasType = "link"
	MessageSearchHasTypePoll     MessageSearchHasType = "poll"
	MessageSearchHasTypeSnapshot MessageSearchHasType = "snapshot"

	// Negated types, results will NOT include messages that have these content types.
	MessageSearchHasTypeNotImage    MessageSearchHasType = "-image"
	MessageSearchHasTypeNotSound    MessageSearchHasType = "-sound"
	MessageSearchHasTypeNotVideo    MessageSearchHasType = "-video"
	MessageSearchHasTypeNotFile     MessageSearchHasType = "-file"
	MessageSearchHasTypeNotSticker  MessageSearchHasType = "-sticker"
	MessageSearchHasTypeNotEmbed    MessageSearchHasType = "-embed"
	MessageSearchHasTypeNotLink     MessageSearchHasType = "-link"
	MessageSearchHasTypeNotPoll     MessageSearchHasType = "-poll"
	MessageSearchHasTypeNotSnapshot MessageSearchHasType = "-snapshot"
)

// MessageSearchEmbedType filters messages by the type of embed they contain.
type MessageSearchEmbedType string

// Constants for MessageSearchEmbedType.
const (
	MessageSearchEmbedTypeImage   MessageSearchEmbedType = "image"
	MessageSearchEmbedTypeVideo   MessageSearchEmbedType = "video"
	MessageSearchEmbedTypeGif     MessageSearchEmbedType = "gif"
	MessageSearchEmbedTypeSound   MessageSearchEmbedType = "sound"
	MessageSearchEmbedTypeArticle MessageSearchEmbedType = "article"

	// Negated types, results will NOT include messages that have these embed types.
	MessageSearchEmbedTypeNotImage   MessageSearchEmbedType = "-image"
	MessageSearchEmbedTypeNotVideo   MessageSearchEmbedType = "-video"
	MessageSearchEmbedTypeNotGif     MessageSearchEmbedType = "-gif"
	MessageSearchEmbedTypeNotSound   MessageSearchEmbedType = "-sound"
	MessageSearchEmbedTypeNotArticle MessageSearchEmbedType = "-article"
)

// MessageSearchSortBy is the field to sort search results by.
type MessageSearchSortBy string

// Constants for MessageSearchSortBy.
const (
	MessageSearchSortByTimestamp MessageSearchSortBy = "timestamp"
	MessageSearchSortByRelevance MessageSearchSortBy = "relevance"
)

// MessageSearchSortOrder is the direction to sort search results.
type MessageSearchSortOrder string

// Constants for MessageSearchSortOrder.
const (
	MessageSearchSortOrderAsc  MessageSearchSortOrder = "asc"
	MessageSearchSortOrderDesc MessageSearchSortOrder = "desc"
)

// GuildMessagesSearch holds the query parameters for searching guild messages.
type GuildMessagesSearch struct {
	// Limit is the max number of messages to return (1-25, default 25).
	Limit                int
	// Offset is the number to offset the returned messages by (max 9975).
	Offset               int
	// MaxID filters messages before this message ID.
	MaxID                snowflake.ID
	// MinID filters messages after this message ID.
	MinID                snowflake.ID
	// Slop is the max number of words to skip between matching tokens in the search Content (max 100, default 2).
	Slop                 int
	// Content filters messages by content (max 1024 characters).
	Content              string
	// ChannelIDs filters messages by these channel IDs (max 500).
	ChannelIDs           []snowflake.ID
	// AuthorTypes filters messages by author type.
	AuthorTypes          []MessageSearchAuthorType
	// AuthorIDs filters messages by these author IDs (max 100).
	AuthorIDs            []snowflake.ID
	// Mentions filters messages that mention these user IDs (max 100).
	Mentions             []snowflake.ID
	// MentionsRoleIDs filters messages that mention these role IDs (max 100).
	MentionsRoleIDs      []snowflake.ID
	// MentionEveryone filters messages that do or do not mention @everyone.
	MentionEveryone      *bool
	// RepliedToUserIDs filters messages that reply to these user IDs (max 100).
	RepliedToUserIDs     []snowflake.ID
	// RepliedToMessageIDs filters messages that reply to these message IDs (max 100).
	RepliedToMessageIDs  []snowflake.ID
	// Pinned filters messages by whether they are pinned.
	Pinned               *bool
	// Has filters messages by whether they contain specific content types.
	Has                  []MessageSearchHasType
	// EmbedTypes filters messages by embed type.
	EmbedTypes           []MessageSearchEmbedType
	// EmbedProviders filters messages by embed provider (case-sensitive, e.g. "Tenor").
	EmbedProviders       []string
	// LinkHostnames filters messages by link hostname (e.g. "discordapp.com").
	LinkHostnames        []string
	// AttachmentFilenames filters messages by attachment filename.
	AttachmentFilenames  []string
	// AttachmentExtensions filters messages by attachment extension (e.g. "txt").
	AttachmentExtensions []string
	// SortBy is the sorting algorithm to use.
	SortBy               MessageSearchSortBy
	// SortOrder is the direction to sort ("asc" or "desc", default "desc").
	SortOrder            MessageSearchSortOrder
	// IncludeNSFW controls whether to include results from age-restricted channels (default false).
	IncludeNSFW          *bool
}

// joinSnowflakes joins a slice of snowflake IDs into a comma separated string.
func joinSnowflakes(ids []snowflake.ID) string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = id.String()
	}
	return strings.Join(strs, ",")
}

// joinStrings joins a slice of values with a Stringer or ~string underlying type into a comma separated string.
func joinStrings[T ~string](vals []T) string {
	strs := make([]string, len(vals))
	for i, v := range vals {
		strs[i] = string(v)
	}
	return strings.Join(strs, ",")
}

// ToQueryValues converts the GuildMessageSearch into QueryValues for use in HTTP requests.
func (s GuildMessagesSearch) ToQueryValues() QueryValues {
	q := QueryValues{}
	if s.Limit != 0 {
		q["limit"] = s.Limit
	}
	if s.Offset != 0 {
		q["offset"] = s.Offset
	}
	if s.MaxID != 0 {
		q["max_id"] = s.MaxID
	}
	if s.MinID != 0 {
		q["min_id"] = s.MinID
	}
	if s.Slop != 0 {
		q["slop"] = s.Slop
	}
	if s.Content != "" {
		q["content"] = s.Content
	}
	if len(s.ChannelIDs) > 0 {
		q["channel_id"] = joinSnowflakes(s.ChannelIDs)
	}
	if len(s.AuthorTypes) > 0 {
		q["author_type"] = joinStrings(s.AuthorTypes)
	}
	if len(s.AuthorIDs) > 0 {
		q["author_id"] = joinSnowflakes(s.AuthorIDs)
	}
	if len(s.Mentions) > 0 {
		q["mentions"] = joinSnowflakes(s.Mentions)
	}
	if len(s.MentionsRoleIDs) > 0 {
		q["mentions_role_id"] = joinSnowflakes(s.MentionsRoleIDs)
	}
	if s.MentionEveryone != nil {
		q["mention_everyone"] = *s.MentionEveryone
	}
	if len(s.RepliedToUserIDs) > 0 {
		q["replied_to_user_id"] = joinSnowflakes(s.RepliedToUserIDs)
	}
	if len(s.RepliedToMessageIDs) > 0 {
		q["replied_to_message_id"] = joinSnowflakes(s.RepliedToMessageIDs)
	}
	if s.Pinned != nil {
		q["pinned"] = *s.Pinned
	}
	if len(s.Has) > 0 {
		q["has"] = joinStrings(s.Has)
	}
	if len(s.EmbedTypes) > 0 {
		q["embed_type"] = joinStrings(s.EmbedTypes)
	}
	if len(s.EmbedProviders) > 0 {
		q["embed_provider"] = strings.Join(s.EmbedProviders, ",")
	}
	if len(s.LinkHostnames) > 0 {
		q["link_hostname"] = strings.Join(s.LinkHostnames, ",")
	}
	if len(s.AttachmentFilenames) > 0 {
		q["attachment_filename"] = strings.Join(s.AttachmentFilenames, ",")
	}
	if len(s.AttachmentExtensions) > 0 {
		q["attachment_extension"] = strings.Join(s.AttachmentExtensions, ",")
	}
	if s.SortBy != "" {
		q["sort_by"] = s.SortBy
	}
	if s.SortOrder != "" {
		q["sort_order"] = s.SortOrder
	}
	if s.IncludeNSFW != nil {
		q["include_nsfw"] = *s.IncludeNSFW
	}
	return q
}

// GuildMessagesSearchResult is the response body from the Search Guild Messages endpoint.
type GuildMessagesSearchResult struct {
	// DoingDeepHistoricalIndex indicates whether the guild is undergoing a deep historical indexing operation.
	DoingDeepHistoricalIndex bool `json:"doing_deep_historical_index"`
	// DocumentsIndexed is the number of documents indexed during the current index operation, if any.
	DocumentsIndexed         *int `json:"documents_indexed,omitempty"`
	// TotalResults is the total number of results that match the query.
	TotalResults             int `json:"total_results"`
	// Messages is the flat list of messages matching the query.
	Messages                 []Message
	// Threads contains the threads that contain the returned messages.
	Threads                  []GuildThread `json:"threads,omitempty"`
	// Members contains thread member objects for each returned thread the current user has joined.
	Members                  []ThreadMember `json:"members,omitempty"`
}

func (r *GuildMessagesSearchResult) UnmarshalJSON(data []byte) error {
	type guildMessagesSearchResult GuildMessagesSearchResult
	var v struct {
		Messages [][]Message `json:"messages"`
		guildMessagesSearchResult
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = GuildMessagesSearchResult(v.guildMessagesSearchResult)
	for _, inner := range v.Messages {
		r.Messages = append(r.Messages, inner...)
	}
	return nil
}