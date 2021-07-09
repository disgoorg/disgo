package api

import (
	"strings"
	"time"

	"github.com/DisgoOrg/restclient"
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

// SystemChannelFlag contains the settings for the Guild(s) system channel
type SystemChannelFlag int

// Constants for SystemChannelFlag
const (
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlag = 1 << iota
	SystemChannelFlagSuppressPremiumSubscriptions
)

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

// MessageNotificationLevel indicates whether users receive @ mentions on a new message
type MessageNotificationLevel int

// Constants for MessageNotifications
const (
	MessageNotificationLevelAllMessages MessageNotificationLevel = iota
	MessageNotificationLevelOnlyMentions
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

// GuildWelcomeScreen is the Welcome Screen of a Guild
type GuildWelcomeScreen struct {
	Description     *string                `json:"description,omitempty"`
	WelcomeChannels []*GuildWelcomeChannel `json:"welcome_channels"`
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
	Disgo                    Disgo
	ID                       Snowflake      `json:"id"`
	Name                     string         `json:"name"`
	Icon                     *string        `json:"icon"`
	DiscoverySplash          *string        `json:"discovery_splash"`
	Splash                   *string        `json:"splash"`
	Features                 []GuildFeature `json:"features"`
	Description              *string        `json:"description"`
	ApproximateMemberCount   *int           `json:"approximate_member_count"`
	ApproximatePresenceCount *int           `json:"approximate_presence_count"`
	Emojis                   []*Emoji       `json:"emojis"`
}

// FullGuild represents a Guild objects sent by discord with the GatewayEventGuildCreate
type FullGuild struct {
	*Guild
	Roles       []*Role        `json:"roles"`
	Emojis      []*Emoji       `json:"emojis"`
	Members     []*Member      `json:"members"`
	Channels    []*ChannelImpl `json:"channels"`
	Threads     []*ChannelImpl `json:"threads"`
	VoiceStates []*VoiceState  `json:"voice_states"`
	//Presences   []*Presence     `json:"presences"`
}

// Guild represents a discord Guild
type Guild struct {
	Disgo                       Disgo
	Ready                       bool
	ID                          Snowflake                  `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        *string                    `json:"icon"`
	Region                      string                     `json:"region"`
	OwnerID                     Snowflake                  `json:"owner_id"`
	JoinedAt                    *time.Time                 `json:"joined_at"`
	DiscoverySplash             *string                    `json:"discovery_splash"`
	Splash                      *string                    `json:"splash"`
	AfkChannelID                *Snowflake                 `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	MemberCount                 *int                       `json:"member_count"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	Large                       *bool                      `json:"large"`
	DefaultMessageNotifications MessageNotificationLevel   `json:"default_message_notifications"`
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
}

// PublicRole returns the @everyone Role
func (g *Guild) PublicRole() *Role {
	return g.Disgo.Cache().Role(g.ID)
}

// Roles return all Role(s) in this Guild
func (g *Guild) Roles() []*Role {
	return g.Disgo.Cache().Roles(g.ID)
}

// SelfMember returns the Member for the current logged-in User for this Guild
func (g *Guild) SelfMember() *Member {
	return g.Disgo.Cache().Member(g.ID, g.Disgo.ClientID())
}

// Disconnect sends a api.GatewayCommand to disconnect from this Guild
func (g *Guild) Disconnect() error {
	return g.Disgo.AudioController().Disconnect(g.ID)
}

// Update updates the current Guild
func (g *Guild) Update(guildUpdate GuildUpdate) (*Guild, restclient.RestError) {
	return g.Disgo.RestClient().UpdateGuild(g.ID, guildUpdate)
}

// Delete deletes the current Guild
func (g *Guild) Delete() restclient.RestError {
	return g.Disgo.RestClient().DeleteGuild(g.ID)
}

// CreateRole allows you to create a new Role
func (g *Guild) CreateRole(roleCreate RoleCreate) (*Role, restclient.RestError) {
	return g.Disgo.RestClient().CreateRole(g.ID, roleCreate)
}

// UpdateRole allows you to update a Role
func (g *Guild) UpdateRole(roleID Snowflake, roleUpdate RoleUpdate) (*Role, restclient.RestError) {
	return g.Disgo.RestClient().UpdateRole(g.ID, roleID, roleUpdate)
}

// DeleteRole allows you to delete a Role
func (g *Guild) DeleteRole(roleID Snowflake) restclient.RestError {
	return g.Disgo.RestClient().DeleteRole(g.ID, roleID)
}

// Leave leaves the Guild
func (g *Guild) Leave() restclient.RestError {
	return g.Disgo.RestClient().LeaveGuild(g.ID)
}

// GetMember returns the Member for the current logged in User for this Guild
func (g *Guild) GetMember(userID Snowflake) *Member {
	return g.Disgo.Cache().Member(g.ID, userID)
}

// AddMember adds a member to the Guild with the oauth2 access token
func (g *Guild) AddMember(userID Snowflake, memberAdd MemberAdd) (*Member, restclient.RestError) {
	return g.Disgo.RestClient().AddMember(g.ID, userID, memberAdd)
}

// UpdateMember updates an existing member of the Guild
func (g *Guild) UpdateMember(userID Snowflake, memebrUpdate MemberUpdate) (*Member, restclient.RestError) {
	return g.Disgo.RestClient().UpdateMember(g.ID, userID, memebrUpdate)
}

// KickMember kicks an existing member from the Guild
func (g *Guild) KickMember(userID Snowflake, reason string) restclient.RestError {
	return g.Disgo.RestClient().RemoveMember(g.ID, userID, reason)
}

// IconURL returns the Icon of a Guild
func (g *Guild) IconURL(size int) *string {
	if g.Icon == nil {
		return nil
	}
	animated := strings.HasPrefix(*g.Icon, "a_")
	format := restclient.PNG
	if animated {
		format = restclient.GIF
	}
	route, err := restclient.GuildIcon.Compile(nil, format, size, g.ID.String(), *g.Icon)
	if err != nil {
		return nil
	}
	u := route.Route()
	return &u
}

// GetCommand fetches a specific Guild Command
func (g *Guild) GetCommand(commandID Snowflake) (*Command, restclient.RestError) {
	return g.Disgo.GetGuildCommand(g.ID, commandID)
}

// GetCommands fetches all Guild Command(s)
func (g *Guild) GetCommands() ([]*Command, restclient.RestError) {
	return g.Disgo.GetGuildCommands(g.ID)
}

// CreateCommand creates a new Command for this Guild
func (g *Guild) CreateCommand(command CommandCreate) (*Command, restclient.RestError) {
	return g.Disgo.CreateGuildCommand(g.ID, command)
}

// EditCommand edits a specific Guild Command
func (g *Guild) EditCommand(commandID Snowflake, command CommandUpdate) (*Command, restclient.RestError) {
	return g.Disgo.EditGuildCommand(g.ID, commandID, command)
}

// DeleteCommand creates a new Command for this Guild
func (g *Guild) DeleteCommand(commandID Snowflake) restclient.RestError {
	return g.Disgo.DeleteGuildCommand(g.ID, commandID)
}

// SetCommands overrides all Command(s) for this Guild
func (g *Guild) SetCommands(commands ...CommandCreate) ([]*Command, restclient.RestError) {
	return g.Disgo.SetGuildCommands(g.ID, commands...)
}

// GetCommandsPermissions returns the GuildCommandPermissions for a all Command(s) in a Guild
func (g *Guild) GetCommandsPermissions() ([]*GuildCommandPermissions, restclient.RestError) {
	return g.Disgo.GetGuildCommandsPermissions(g.ID)
}

// GetCommandPermissions returns the GuildCommandPermissions for a specific Command in a Guild
func (g *Guild) GetCommandPermissions(commandID Snowflake) (*GuildCommandPermissions, restclient.RestError) {
	return g.Disgo.GetGuildCommandPermissions(g.ID, commandID)
}

// SetCommandsPermissions sets the GuildCommandPermissions for a all Command(s)
func (g *Guild) SetCommandsPermissions(commandPermissions ...SetGuildCommandPermissions) ([]*GuildCommandPermissions, error) {
	return g.Disgo.SetGuildCommandsPermissions(g.ID, commandPermissions...)
}

// SetCommandPermissions sets the GuildCommandPermissions for a specific Command
func (g *Guild) SetCommandPermissions(commandID Snowflake, permissions SetGuildCommandPermissions) (*GuildCommandPermissions, error) {
	return g.Disgo.SetGuildCommandPermissions(g.ID, commandID, permissions)
}

// GetTemplates gets a specific GuildTemplate
func (g *Guild) GetTemplates() ([]*GuildTemplate, restclient.RestError) {
	return g.Disgo.RestClient().GetGuildTemplates(g.ID)
}

// CreateTemplate creates a new GuildTemplate
func (g *Guild) CreateTemplate(createGuildTemplate CreateGuildTemplate) (*GuildTemplate, restclient.RestError) {
	return g.Disgo.RestClient().CreateGuildTemplate(g.ID, createGuildTemplate)
}

// SyncTemplate syncs the current Guild status to an existing GuildTemplate
func (g *Guild) SyncTemplate(code string) (*GuildTemplate, restclient.RestError) {
	return g.Disgo.RestClient().SyncGuildTemplate(g.ID, code)
}

// UpdateTemplate updates a specific GuildTemplate
func (g *Guild) UpdateTemplate(code string, updateGuildTemplate UpdateGuildTemplate) (*GuildTemplate, restclient.RestError) {
	return g.Disgo.RestClient().UpdateGuildTemplate(g.ID, code, updateGuildTemplate)
}

// DeleteTemplate deletes a specific GuildTemplate
func (g *Guild) DeleteTemplate(code string) (*GuildTemplate, restclient.RestError) {
	return g.Disgo.RestClient().DeleteGuildTemplate(g.ID, code)
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

// GuildCreate is the payload used to create a Guild
type GuildCreate struct {
	Name                            string                     `json:"name"`
	Region                          string                     `json:"region,omitempty"`
	Icon                            string                     `json:"icon,omitempty"`
	VerificationLevel               VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotificationLevel MessageNotificationLevel   `json:"default_message_notification_level"`
	ExplicitContentFilterLevel      ExplicitContentFilterLevel `json:"explicit_content_filter_level"`
	Roles                           []RoleCreate               `json:"roles,omitempty"`
	Channels                        []GuildChannelCreate       `json:"channels,omitempty"`
	AFKChannelID                    Snowflake                  `json:"afk_channel_id,omitempty"`
	AFKTimeout                      int                        `json:"afk_timeout,omitempty"`
	SystemChannelID                 Snowflake                  `json:"system_channel_id,omitempty"`
	SystemChannelFlags              SystemChannelFlag          `json:"system_channel_flags,omitempty"`
}

type GuildChannelCreate struct {
	ID       Snowflake   `json:"id,omitempty"`
	Name     string      `json:"name"`
	Type     ChannelType `json:"type"`
	ParentID Snowflake   `json:"parent_id,omitempty"`
}

type GuildUpdate struct {
	Name                            *string                      `json:"name"`
	Region                          *string                     `json:"region"`
	VerificationLevel               *VerificationLevel          `json:"verification_level"`
	DefaultMessageNotificationLevel *MessageNotificationLevel   `json:"default_message_notification_level"`
	ExplicitContentFilterLevel      *ExplicitContentFilterLevel `json:"explicit_content_filter_level"`
	AFKChannelID                    *Snowflake                  `json:"afk_channel_id"`
	AFKTimeout                      *int                         `json:"afk_timeout"`
	Icon                            *string                     `json:"icon"`
	OwnerID                         Snowflake                   `json:"owner_id"`
	Splash                          []byte                 `json:"splash"`
	DiscoverySplash                 []byte                 `json:"discovery_splash"`
	Banner                          []byte                 `json:"banner"`
	SystemChannelID                 *Snowflake                  `json:"system_channel_id"`
	SystemChannelFlags              *SystemChannelFlag          `json:"system_channel_flags"`
	RulesChannelID                  *Snowflake                  `json:"rules_channel_id"`
	PublicUpdatesChannelID          *Snowflake                  `json:"public_updates_channel_id"`
	PreferredLocale                 *string                     `json:"preferred_locale"`
	Features                        []GuildFeature              `json:"features"`
	Description                     *string                     `json:"description"`
}

type PruneCount struct {
	Pruned int `json:"pruned"`
}
