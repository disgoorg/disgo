package api

import (
	"strings"
	"time"

	"github.com/DiscoOrg/disgo/api/endpoints"
)

// PremiumTier tells you the boost level of a guild
type PremiumTier int

// Constants for PremiumTier
const (
	PremiumTierNone PremiumTier = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

// SystemChannelFlag contains the settings for the guilds system channel
type SystemChannelFlag int

// Constants for SystemChannelFlag
const (
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlag = 1 << iota
	SystemChannelFlagSuppressPremiumSubscriptions
)

// The VerificationLevel of a guild that members must be to send messages
type VerificationLevel int

// Constants for VerificationLevel
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
const (
	MessageNotificationsAllMessages MessageNotifications = iota
	MessageNotificationsOnlyMentions
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

// The GuildFeature (s) that a guild contains
type GuildFeature string

// Constants for GuildFeature
const (
	GuildFeatureInviteSplash                  GuildFeature = "INVITE_SPLASH"
	GuildFeatureVipRegions                    GuildFeature = "VIP_REGIONS"
	GuildFeatureVanityURL                     GuildFeature = "VANITY_URL"
	GuildFeatureVERIFIED                      GuildFeature = "VERIFIED"
	GuildFeaturePARTNERED                     GuildFeature = "PARTNERED"
	GuildFeatureCOMMUNITY                     GuildFeature = "COMMUNITY"
	GuildFeatureCOMMERCE                      GuildFeature = "COMMERCE"
	guildFeatureNews                          GuildFeature = "NEWS"
	guildFeatureDiscoverable                  GuildFeature = "DISCOVERABLE"
	guildFeatureFeaturable                    GuildFeature = "FEATURABLE"
	GuildFeatureAnimatedIcon                  GuildFeature = "ANIMATED_ICON"
	GuildFeatureBANNER                        GuildFeature = "BANNER"
	GuildFeatureWelcomeScreenEnabled          GuildFeature = "WELCOME_SCREEN_ENABLED"
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeaturePreviewEnabled                GuildFeature = "PREVIEW_ENABLED"
)

// Guild represents a discord guild_events
type Guild struct {
	Disgo   Disgo
	ID      Snowflake `json:"id"`
	Name    string    `json:"name"`
	Icon    *string   `json:"icon"`
	Region  string    `json:"region"`
	OwnerID Snowflake `json:"owner_id"`
	// only used in GET USER GUILDS
	//Owner                       *bool                       `json:"owner"`
	//Permissions                 *Permissions                `json:"permissions"`
	JoinedAt                    *time.Time                 `json:"joined_at"`
	DiscoverySplash             *string                    `json:"discovery_splash"`
	Splash                      *string                    `json:"splash"`
	AfkChannelID                *Snowflake                 `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	MemberCount                 *int                       `json:"member_count"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	Large                       *bool                      `json:"large"`
	DefaultMessageNotifications MessageNotifications       `json:"default_message_notifications"`
	Roles                       []*Role                    `json:"roles"`
	Emojis                      []*Emote                   `json:"emojis"`
	Members                     []*Member                  `json:"members"`
	MaxPresences                *int                       `json:"max_presences"`
	MaxMembers                  *int                       `json:"max_members"`
	Channels                    []*GuildChannel            `json:"channels"`
	VoiceStates                 []*VoiceState              `json:"voice_states"`
	Unavailable                 bool                      `json:"unavailable"`
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
	//Presences                   []*Presence                `json:"presences"`
}

func (g Guild) CreateRole(role UpdateRole) (*Role, error) {
	return g.Disgo.RestClient().CreateRole(g.ID, role)
}

// AddMember adds a member to the guild with the oauth2 access token
func (g Guild) AddMember(userID Snowflake, addGuildMemberData AddGuildMemberData) (*Member, error) {
	return g.Disgo.RestClient().AddMember(g.ID, userID, addGuildMemberData)
}

// IconURL returns the Icon of a guild_events
func (g Guild) IconURL() *string {
	if g.Icon == nil {
		return nil
	}
	animated := strings.HasPrefix(*g.Icon, "a_")
	format := endpoints.PNG
	if animated {
		format = endpoints.GIF
	}
	u := endpoints.GuildIcon.Compile(format, g.ID.String(), *g.Icon).Route()
	return &u
}

// CreateCommand returns a GuildCommandBuilder for that guild
func (g Guild) CreateCommand(name string, description string) GuildCommandBuilder {
	return NewGuildCommandBuilder(g.Disgo, g.ID, name, description)
}
