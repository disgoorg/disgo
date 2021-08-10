package discord

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

// SystemChannelFlag contains the settings for the Guild(s) system channel
type SystemChannelFlag int

// Constants for SystemChannelFlag
//goland:noinspection GoUnusedConst
const (
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlag = 1 << iota
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

// MessageNotifications indicates whether users receive @ mentions on a new message
type MessageNotifications int

// Constants for MessageNotifications
//goland:noinspection GoUnusedConst
const (
	MessageNotificationsAllMessages MessageNotifications = iota
	MessageNotificationsOnlyMentions
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
	GuildFeatureInviteSplash                  GuildFeature = "INVITE_SPLASH"
	GuildFeatureVipRegions                    GuildFeature = "VIP_REGIONS"
	GuildFeatureVanityURL                     GuildFeature = "VANITY_URL"
	GuildFeatureVERIFIED                      GuildFeature = "VERIFIED"
	GuildFeaturePARTNERED                     GuildFeature = "PARTNERED"
	GuildFeatureCOMMUNITY                     GuildFeature = "COMMUNITY"
	GuildFeatureCOMMERCE                      GuildFeature = "COMMERCE"
	GuildFeatureNews                          GuildFeature = "NEWS"
	GuildFeatureDiscoverable                  GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable                    GuildFeature = "FEATURABLE"
	GuildFeatureAnimatedIcon                  GuildFeature = "ANIMATED_ICON"
	GuildFeatureBANNER                        GuildFeature = "BANNER"
	GuildFeatureWelcomeScreenEnabled          GuildFeature = "WELCOME_SCREEN_ENABLED"
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeaturePreviewEnabled                GuildFeature = "PREVIEW_ENABLED"
)

// Guild represents a discord Guild
type Guild struct {
	ID                          Snowflake                  `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        *string                    `json:"icon"`
	Region                      string                     `json:"region"`
	OwnerID                     Snowflake                  `json:"owner_id"`
	JoinedAt                    *Time                      `json:"joined_at"`
	DiscoverySplash             *string                    `json:"discovery_splash"`
	Splash                      *string                    `json:"splash"`
	AfkChannelID                *Snowflake                 `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	MemberCount                 *int                       `json:"member_count"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	Large                       *bool                      `json:"large"`
	DefaultMessageNotifications MessageNotifications       `json:"default_message_notifications"`
	MaxPresences                *int                       `json:"max_presences"`
	MaxMembers                  *int                       `json:"max_members"`
	Unavailable                 bool                       `json:"unavailable"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Features                    []GuildFeature             `json:"features"`
	MfaLevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               Snowflake                  `json:"application_id"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             Snowflake                  `json:"widget_channel_id"`
	SystemChannelID             Snowflake                  `json:"system_channel_id"`
	SystemChannelFlags          SystemChannelFlag          `json:"system_channel_flags"`
	RulesChannelID              Snowflake                  `json:"rules_channel_id"`
	VanityURLCode               *string                    `json:"vanity_url_code"`
	Description                 *string                    `json:"description"`
	Banner                      *string                    `json:"banner"`
	PremiumTier                 PremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    *int                       `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      *Snowflake                 `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        *int                       `json:"max_video_channel_users"`
	ApproximateMemberCount      *int                       `json:"approximate_member_count"`
	ApproximatePresenceCount    *int                       `json:"approximate_presence_count"`
	WelcomeScreen               *GuildWelcomeScreen        `json:"welcome_screen"`
	Roles                       []Role                     `json:"roles"`
	Emojis                      []Emoji                    `json:"emojis"`
	Members                     []Member                   `json:"members"`
	Channels                    []Channel                  `json:"channels"`
	VoiceStates                 []VoiceState               `json:"voice_states"`
	Presences                   []Presence                 `json:"presences"`
}

// PartialGuild is returned on the restclient.GetGuilds route
type PartialGuild struct {
	ID          Snowflake      `json:"id"`
	Name        string         `json:"name"`
	Icon        string         `json:"icon"`
	Owner       bool           `json:"owner"`
	Permissions Permissions    `json:"permissions"`
	Features    []GuildFeature `json:"features"`
}

// GuildWelcomeScreen is the Welcome Screen of a Guild
type GuildWelcomeScreen struct {
	Description     *string               `json:"description,omitempty"`
	WelcomeChannels []GuildWelcomeChannel `json:"welcome_channels"`
}

// GuildWelcomeChannel is one of the channels in a GuildWelcomeScreen
type GuildWelcomeChannel struct {
	ChannelID   Snowflake  `json:"channel_id"`
	Description string     `json:"description"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}

// GuildPreview is used for previewing public Guild(s) before joining them
type GuildPreview struct {
	ID                       Snowflake      `json:"id"`
	Name                     string         `json:"name"`
	Icon                     *string        `json:"icon"`
	DiscoverySplash          *string        `json:"discovery_splash"`
	Splash                   *string        `json:"splash"`
	Features                 []GuildFeature `json:"features"`
	Description              *string        `json:"description"`
	ApproximateMemberCount   *int           `json:"approximate_member_count"`
	ApproximatePresenceCount *int           `json:"approximate_presence_count"`
	Emojis                   []Emoji        `json:"emojis"`
}

// GuildCreate is the payload used to create a Guild
type GuildCreate struct {
	Name                            string                     `json:"name"`
	Region                          string                     `json:"region,omitempty"`
	Icon                            string                     `json:"icon,omitempty"`
	VerificationLevel               VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotificationLevel MessageNotifications       `json:"default_message_notification_level"`
	ExplicitContentFilterLevel      ExplicitContentFilterLevel `json:"explicit_content_filter_level"`
	Roles                           []RoleCreate               `json:"roles,omitempty"`
	Channels                        []ChannelCreate            `json:"channels,omitempty"`
	AFKChannelID                    Snowflake                  `json:"afk_channel_id,omitempty"`
	AFKTimeout                      int                        `json:"afk_timeout,omitempty"`
	SystemChannelID                 Snowflake                  `json:"system_channel_id,omitempty"`
	SystemChannelFlags              SystemChannelFlag          `json:"system_channel_flags,omitempty"`
}

// GuildUpdate is the payload used to update a Guild
type GuildUpdate struct {
	Name                            *string                     `json:"name,omitempty"`
	Region                          *string                     `json:"region,omitempty"`
	VerificationLevel               *VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotificationLevel *MessageNotifications       `json:"default_message_notification_level,omitempty"`
	ExplicitContentFilterLevel      *ExplicitContentFilterLevel `json:"explicit_content_filter_level,omitempty"`
	AFKChannelID                    *Snowflake                  `json:"afk_channel_id,omitempty"`
	AFKTimeout                      *int                        `json:"afk_timeout,omitempty"`
	Icon                            *string                     `json:"icon,omitempty"`
	OwnerID                         *Snowflake                  `json:"owner_id,omitempty"`
	Splash                          Icon                        `json:"splash,omitempty"`
	DiscoverySplash                 Icon                        `json:"discovery_splash,omitempty"`
	Banner                          Icon                        `json:"banner,omitempty"`
	SystemChannelID                 *Snowflake                  `json:"system_channel_id,omitempty"`
	SystemChannelFlags              *SystemChannelFlag          `json:"system_channel_flags,omitempty"`
	RulesChannelID                  *Snowflake                  `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID          *Snowflake                  `json:"public_updates_channel_id,omitempty"`
	PreferredLocale                 *string                     `json:"preferred_locale,omitempty"`
	Features                        *[]GuildFeature             `json:"features,omitempty"`
	Description                     *string                     `json:"description,omitempty"`
}
