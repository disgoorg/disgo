package api

import (
	"strings"
	"time"

	"github.com/DiscoOrg/disgo/api/endpoints"
)

type PremiumTier int

const (
	PremiumTierNone PremiumTier = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

type SystemChannelFlag int

const (
	SystemChannelFlagSuppressJoinNotifications = 1 << iota
	SystemChannelFlagSuppressPremiumSubscriptions
)

type VerificationLevel int

const (
	VerificationLevelNone = iota
	VerificationLevelLow
	VerificationLevelMedium
	VerificationLevelHigh
	VerificationLevelVeryHigh
)

type MessageNotifications int

const (
	MessageNotificationsAllMessages = iota
	MessageNotificationsOnlyMentions
)

type ExplicitContentFilterLevel int

const (
	explicitContentFilterLevelDisabled = iota
	ExplicitContentFilterLevelMembersWithoutRoles
	ExplicitContentFilterLevelAllMembers
)

type MfaLevel int

const (
	mfaLevelNone = iota
	mfaLevelElevated
)

type GuildFeature string

const (
	GuildFeatureInviteSplash                  = "INVITE_SPLASH"
	GuildFeatureVipRegions                    = "VIP_REGIONS"
	GuildFeatureVanityUrl                     = "VANITY_URL"
	GuildFeatureVERIFIED                      = "VERIFIED"
	GuildFeaturePARTNERED                     = "PARTNERED"
	GuildFeatureCOMMUNITY                     = "COMMUNITY"
	GuildFeatureCOMMERCE                      = "COMMERCE"
	guildFeatureNews                          = "NEWS"
	guildFeatureDiscoverable                  = "DISCOVERABLE"
	guildFeatureFeaturable                    = "FEATURABLE"
	GuildFeatureAnimatedIcon                  = "ANIMATED_ICON"
	GuildFeatureBANNER                        = "BANNER"
	GuildFeatureWelcomeScreenEnabled          = "WELCOME_SCREEN_ENABLED"
	GuildFeatureMemberVerificationGateEnabled = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeaturePreviewEnabled                = "PREVIEW_ENABLED"
)

// Guild represents a discord guild_events
type Guild struct {
	Disgo                       Disgo
	ID                          Snowflake                  `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        *string                    `json:"icon"`
	Region                      string                     `json:"region"`
	OwnerID                     Snowflake                  `json:"owner_id"`
	Owner                       bool                       `json:"owner"`
	JoinedAt                    time.Time                  `json:"joined_at"`
	DiscoverySplash             string                     `json:"discovery_splash"`
	Splash                      string                     `json:"splash"`
	AfkChannelID                *Snowflake                 `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	MemberCount                 int                        `json:"member_count"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	Large                       bool                       `json:"large"`
	DefaultMessageNotifications MessageNotifications       `json:"default_message_notifications"`
	Roles                       []*Role                    `json:"roles"`
	Emojis                      []*Emote                   `json:"emojis"`
	Members                     []*Member                  `json:"members"`
	MaxPresences                int                        `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	Channels                    []*GuildChannel            `json:"channels"`
	VoiceStates                 []*VoiceState              `json:"voice_states"`
	Unavailable                 bool                       `json:"unavailable"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Features                    []GuildFeature             `json:"features"`
	MfaLevel                    MfaLevel                   `json:"mfa_level"`
	ApplicationID               Snowflake                  `json:"application_id"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             Snowflake                  `json:"widget_channel_id"`
	SystemChannelID             Snowflake                  `json:"system_channel_id"`
	SystemChannelFlags          SystemChannelFlag          `json:"system_channel_flags"`
	RulesChannelID              Snowflake                  `json:"rules_channel_id"`
	VanityURLCode               string                     `json:"vanity_url_code"`
	Description                 string                     `json:"description"`
	Banner                      string                     `json:"banner"`
	PremiumTier                 PremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      Snowflake                  `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	ApproximateMemberCount      int                        `json:"approximate_member_count"`
	ApproximatePresenceCount    int                        `json:"approximate_presence_count"`
	Permissions                 int64                      `json:"permissions"`
	//Presences                   []*Presence                `json:"presences"`
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
	u := endpoints.GuildIcon.Compile(format, g.ID.String(), *g.Icon)
	return &u
}

func (g Guild) CreateCommand(name string, description string) GuildCommandBuilder {
	return NewGuildCommandBuilder(g.Disgo, g.ID, name, description)
}

// UnavailableGuild represents a unavailable discord guild_events
type UnavailableGuild struct {
	ID          Snowflake
	Unavailable bool
}
