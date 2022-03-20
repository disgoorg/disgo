package discord

import (
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

// PremiumTier tells you the boost level of a Guild
type PremiumTier int

// Constants for PremiumTier
//goland:noinspection GoUnusedConst
const (
	PremiumTierNone PremiumTier = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

// SystemChannelFlags contains the settings for the Guild(s) system channel
type SystemChannelFlags int

// Constants for SystemChannelFlags
//goland:noinspection GoUnusedConst
const (
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlags = 1 << iota
	SystemChannelFlagSuppressPremiumSubscriptions
)

// The VerificationLevel of a Guild that members must be to send messages
type VerificationLevel int

// Constants for VerificationLevel
//goland:noinspection GoUnusedConst
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
//goland:noinspection GoUnusedConst
const (
	MessageNotificationsLevelAllMessages MessageNotificationsLevel = iota
	MessageNotificationsLevelOnlyMentions
)

// The ExplicitContentFilterLevel of a Guild
type ExplicitContentFilterLevel int

// Constants for ExplicitContentFilterLevel
//goland:noinspection GoUnusedConst
const (
	ExplicitContentFilterLevelDisabled ExplicitContentFilterLevel = iota
	ExplicitContentFilterLevelMembersWithoutRoles
	ExplicitContentFilterLevelAllMembers
)

// The MFALevel of a Guild
type MFALevel int

// Constants for MFALevel
//goland:noinspection GoUnusedConst
const (
	MFALevelNone MFALevel = iota
	MFALevelElevated
)

// The GuildFeature (s) that a Guild contains
type GuildFeature string

// Constants for GuildFeature
//goland:noinspection GoUnusedConst
const (
	GuildFeatureAnimatedIcon                  GuildFeature = "ANIMATED_ICON"
	GuildFeatureBanner                        GuildFeature = "BANNER"
	GuildFeatureCommerce                      GuildFeature = "COMMERCE"
	GuildFeatureCommunity                     GuildFeature = "COMMUNITY"
	GuildFeatureDiscoverable                  GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable                    GuildFeature = "FEATURABLE"
	GuildFeatureInviteSplash                  GuildFeature = "INVITE_SPLASH"
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeatureMonetizationEnabled           GuildFeature = "MONETIZATION_ENABLED"
	GuildFeatureMoreStickers                  GuildFeature = "MORE_STICKERS"
	GuildFeatureNews                          GuildFeature = "NEWS"
	GuildFeaturePartnered                     GuildFeature = "PARTNERED"
	GuildFeaturePreviewEnabled                GuildFeature = "PREVIEW_ENABLED"
	GuildFeaturePrivateThreads                GuildFeature = "PRIVATE_THREADS"
	GuildFeatureRoleIcons                     GuildFeature = "ROLE_ICONS"
	GuildFeatureSevenDayThreadArchive         GuildFeature = "SEVEN_DAY_THREAD_ARCHIVE"
	GuildFeatureThreeDayThreadArchive         GuildFeature = "THREE_DAY_THREAD_ARCHIVE"
	GuildFeatureTicketedEventsEnabled         GuildFeature = "TICKETED_EVENTS_ENABLED"
	GuildFeatureVanityURL                     GuildFeature = "VANITY_URL"
	GuildFeatureVerified                      GuildFeature = "VERIFIED"
	GuildFeatureVipRegions                    GuildFeature = "VIP_REGIONS"
	GuildFeatureWelcomeScreenEnabled          GuildFeature = "WELCOME_SCREEN_ENABLED"
)

// Guild represents a discord Guild
type Guild struct {
	ID                          snowflake.Snowflake        `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        *string                    `json:"icon"`
	Splash                      *string                    `json:"splash"`
	DiscoverySplash             *string                    `json:"discovery_splash"`
	OwnerID                     snowflake.Snowflake        `json:"owner_id"`
	AfkChannelID                *snowflake.Snowflake       `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             snowflake.Snowflake        `json:"widget_channel_id"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	DefaultMessageNotifications MessageNotificationsLevel  `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Roles                       []Role                     `json:"roles"`
	Emojis                      []Emoji                    `json:"emojis"`
	Features                    []GuildFeature             `json:"features"`
	MFALevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               *snowflake.Snowflake       `json:"application_id"`
	SystemChannelID             *snowflake.Snowflake       `json:"system_channel_id"`
	SystemChannelFlags          SystemChannelFlags         `json:"system_channel_flags"`
	RulesChannelID              *snowflake.Snowflake       `json:"rules_channel_id"`
	MemberCount                 int                        `json:"member_count"`
	MaxPresences                *int                       `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	VanityURLCode               *string                    `json:"vanity_url_code"`
	Description                 *string                    `json:"description"`
	Banner                      *string                    `json:"banner"`
	PremiumTier                 PremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      *snowflake.Snowflake       `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	WelcomeScreen               WelcomeScreen              `json:"welcome_screen"`
	NSFWLevel                   NSFWLevel                  `json:"nsfw_level"`
	BoostProgressBarEnabled     bool                       `json:"premium_progress_bar_enabled"`
	Stickers                    []Sticker                  `json:"stickers"`
	JoinedAt                    Time                       `json:"joined_at"`

	// only over GET /guilds/{guild.id}
	ApproximateMemberCount   int `json:"approximate_member_count"`
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

type GatewayGuild struct {
	Guild
	Large                bool                  `json:"large"`
	Unavailable          bool                  `json:"unavailable"`
	VoiceStates          []VoiceState          `json:"voice_states"`
	Members              []Member              `json:"members"`
	Channels             []GuildChannel        `json:"channels"`
	Threads              []GuildThread         `json:"threads"`
	Presences            []Presence            `json:"presences"`
	StageInstances       []StageInstance       `json:"stage_instances"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
}

func (g *GatewayGuild) UnmarshalJSON(data []byte) error {
	type gatewayGuild GatewayGuild
	var v struct {
		Channels []UnmarshalChannel `json:"channels"`
		Threads  []UnmarshalChannel `json:"threads"`
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

	g.Threads = make([]GuildThread, len(v.Threads))
	for i := range v.Threads {
		g.Threads[i] = v.Threads[i].Channel.(GuildThread)
	}

	return nil
}

type UnavailableGuild struct {
	ID          snowflake.Snowflake `json:"id"`
	Unavailable bool                `json:"unavailable"`
}

// OAuth2Guild is returned on the route.GetGuilds route
type OAuth2Guild struct {
	ID          snowflake.Snowflake `json:"id"`
	Name        string              `json:"name"`
	Icon        *string             `json:"icon"`
	Owner       bool                `json:"owner"`
	Permissions Permissions         `json:"permissions"`
	Features    []GuildFeature      `json:"features"`
}

// WelcomeScreen is the Welcome Screen of a Guild
type WelcomeScreen struct {
	Description     *string               `json:"description,omitempty"`
	WelcomeChannels []GuildWelcomeChannel `json:"welcome_channels"`
}

// GuildWelcomeChannel is one of the channels in a WelcomeScreen
type GuildWelcomeChannel struct {
	ChannelID   snowflake.Snowflake  `json:"channel_id"`
	Description string               `json:"description"`
	EmojiID     *snowflake.Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string              `json:"emoji_name,omitempty"`
}

// GuildPreview is used for previewing public Guild(s) before joining them
type GuildPreview struct {
	ID                       snowflake.Snowflake `json:"id"`
	Name                     string              `json:"name"`
	Icon                     *string             `json:"icon"`
	DiscoverySplash          *string             `json:"discovery_splash"`
	Splash                   *string             `json:"splash"`
	Features                 []GuildFeature      `json:"features"`
	Description              *string             `json:"description"`
	ApproximateMemberCount   *int                `json:"approximate_member_count"`
	ApproximatePresenceCount *int                `json:"approximate_presence_count"`
	Emojis                   []Emoji             `json:"emojis"`
}

// GuildCreate is the payload used to create a Guild
type GuildCreate struct {
	Name                            string                     `json:"name"`
	Icon                            *Icon                      `json:"icon,omitempty"`
	VerificationLevel               VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotificationLevel MessageNotificationsLevel  `json:"default_message_notification_level,omitempty"`
	ExplicitContentFilterLevel      ExplicitContentFilterLevel `json:"explicit_content_filter_level,omitempty"`
	Roles                           []GuildCreateRole          `json:"roles,omitempty"`
	Channels                        []GuildCreateChannel       `json:"channels,omitempty"`
	AFKChannelID                    snowflake.Snowflake        `json:"afk_channel_id,omitempty"`
	AFKTimeout                      int                        `json:"afk_timeout,omitempty"`
	SystemChannelID                 snowflake.Snowflake        `json:"system_channel_id,omitempty"`
	SystemChannelFlags              SystemChannelFlags         `json:"system_channel_flags,omitempty"`
}

// GuildUpdate is the payload used to update a Guild
type GuildUpdate struct {
	Name                            string                      `json:"name,omitempty"`
	VerificationLevel               *VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotificationLevel *MessageNotificationsLevel  `json:"default_message_notification_level,omitempty"`
	ExplicitContentFilterLevel      *ExplicitContentFilterLevel `json:"explicit_content_filter_level,omitempty"`
	AFKChannelID                    *snowflake.Snowflake        `json:"afk_channel_id,omitempty"`
	AFKTimeout                      *int                        `json:"afk_timeout,omitempty"`
	Icon                            *json.Nullable[Icon]        `json:"icon,omitempty"`
	OwnerID                         *snowflake.Snowflake        `json:"owner_id,omitempty"`
	Splash                          *json.Nullable[Icon]        `json:"splash,omitempty"`
	DiscoverySplash                 *json.Nullable[Icon]        `json:"discovery_splash,omitempty"`
	Banner                          *json.Nullable[Icon]        `json:"banner,omitempty"`
	SystemChannelID                 *snowflake.Snowflake        `json:"system_channel_id,omitempty"`
	SystemChannelFlags              *SystemChannelFlags         `json:"system_channel_flags,omitempty"`
	RulesChannelID                  *snowflake.Snowflake        `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID          *snowflake.Snowflake        `json:"public_updates_channel_id,omitempty"`
	PreferredLocale                 *string                     `json:"preferred_locale,omitempty"`
	Features                        []GuildFeature              `json:"features,omitempty"`
	Description                     *string                     `json:"description,omitempty"`
	BoostProgressBarEnabled         *bool                       `json:"premium_progress_bar_enabled,omitempty"`
}

type NSFWLevel int

//goland:noinspection GoUnusedConst
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
