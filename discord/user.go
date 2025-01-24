package discord

import (
	"strconv"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/internal/flags"
)

// UserFlags defines certain flags/badges a user can have (https://discord.com/developers/docs/resources/user#user-object-user-flags)
type UserFlags int

// All UserFlags
const (
	UserFlagDiscordEmployee UserFlags = 1 << iota
	UserFlagPartneredServerOwner
	UserFlagHypeSquadEvents
	UserFlagBugHunterLevel1
	_
	_
	UserFlagHouseBravery
	UserFlagHouseBrilliance
	UserFlagHouseBalance
	UserFlagEarlySupporter
	UserFlagTeamUser
	_
	_
	_
	UserFlagBugHunterLevel2
	_
	UserFlagVerifiedBot
	UserFlagEarlyVerifiedBotDeveloper
	UserFlagDiscordCertifiedModerator
	UserFlagBotHTTPInteractions
	_
	_
	UserFlagActiveDeveloper
	UserFlagsNone UserFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f UserFlags) Add(bits ...UserFlags) UserFlags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f UserFlags) Remove(bits ...UserFlags) UserFlags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f UserFlags) Has(bits ...UserFlags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f UserFlags) Missing(bits ...UserFlags) bool {
	return flags.Missing(f, bits...)
}

var _ Mentionable = (*User)(nil)

// User is a struct for interacting with discord's users
type User struct {
	ID                   snowflake.ID          `json:"id"`
	Username             string                `json:"username"`
	Discriminator        string                `json:"discriminator"`
	GlobalName           *string               `json:"global_name"`
	Avatar               *string               `json:"avatar"`
	Banner               *string               `json:"banner"`
	AccentColor          *int                  `json:"accent_color"`
	Bot                  bool                  `json:"bot"`
	System               bool                  `json:"system"`
	PublicFlags          UserFlags             `json:"public_flags"`
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration_data"`
}

// String returns a mention of the user
func (u User) String() string {
	return UserMention(u.ID)
}

// Mention returns a mention of the user
func (u User) Mention() string {
	return u.String()
}

// Tag returns a formatted string of "Username#Discriminator", falling back to the username if discriminator is "0"
func (u User) Tag() string {
	return UserTag(u.Username, u.Discriminator)
}

// EffectiveName returns the global (display) name of the user if set, falling back to the username
func (u User) EffectiveName() string {
	if u.GlobalName != nil {
		return *u.GlobalName
	}
	return u.Username
}

// EffectiveAvatarURL returns the avatar URL of the user if set, falling back to the default avatar URL
func (u User) EffectiveAvatarURL(opts ...CDNOpt) string {
	if u.Avatar == nil {
		return u.DefaultAvatarURL(opts...)
	}
	return formatAssetURL(UserAvatar, opts, u.ID, *u.Avatar)
}

// AvatarURL returns the avatar URL of the user if set or nil
func (u User) AvatarURL(opts ...CDNOpt) *string {
	if u.Avatar == nil {
		return nil
	}
	url := formatAssetURL(UserAvatar, opts, u.ID, *u.Avatar)
	return &url
}

// DefaultAvatarURL calculates and returns the default avatar URL
func (u User) DefaultAvatarURL(opts ...CDNOpt) string {
	discriminator, err := strconv.Atoi(u.Discriminator)
	if err != nil {
		return ""
	}
	index := discriminator % 5
	if index == 0 { // new username system
		index = int((u.ID >> 22) % 6)
	}
	return formatAssetURL(DefaultUserAvatar, opts, index)
}

// BannerURL returns the banner URL if set or nil
func (u User) BannerURL(opts ...CDNOpt) *string {
	if u.Banner == nil {
		return nil
	}
	url := formatAssetURL(UserBanner, opts, u.ID, *u.Banner)
	return &url
}

// AvatarDecorationURL returns the avatar decoration URL if set or nil
func (u User) AvatarDecorationURL(opts ...CDNOpt) *string {
	if u.AvatarDecorationData == nil {
		return nil
	}
	url := formatAssetURL(AvatarDecoration, opts, u.AvatarDecorationData.Asset)
	return &url
}

func (u User) CreatedAt() time.Time {
	return u.ID.Time()
}

// OAuth2User represents a full User returned by the oauth2 endpoints
type OAuth2User struct {
	User
	// Requires OAuth2ScopeIdentify
	MfaEnabled  bool        `json:"mfa_enabled"`
	Locale      string      `json:"locale"`
	Flags       UserFlags   `json:"flags"`
	PremiumType PremiumType `json:"premium_type"`

	// Requires OAuth2ScopeEmail
	Verified bool   `json:"verified"`
	Email    string `json:"email"`
}

// PremiumType defines the different discord nitro tiers a user can have (https://discord.com/developers/docs/resources/user#user-object-premium-types)
type PremiumType int

// All PremiumType(s)
const (
	PremiumTypeNone PremiumType = iota
	PremiumTypeNitroClassic
	PremiumTypeNitro
	PremiumTypeNitroBasic
)

// UserUpdate is the payload used to update the OAuth2User
type UserUpdate struct {
	Username string               `json:"username,omitempty"`
	Avatar   *json.Nullable[Icon] `json:"avatar,omitempty"`
	Banner   *json.Nullable[Icon] `json:"banner,omitempty"`
}

type ApplicationRoleConnection struct {
	PlatformName     *string           `json:"platform_name"`
	PlatformUsername *string           `json:"platform_username"`
	Metadata         map[string]string `json:"metadata"`
}

type ApplicationRoleConnectionUpdate struct {
	PlatformName     *string            `json:"platform_name,omitempty"`
	PlatformUsername *string            `json:"platform_username,omitempty"`
	Metadata         *map[string]string `json:"metadata,omitempty"`
}

type AvatarDecorationData struct {
	Asset string       `json:"asset"`
	SkuID snowflake.ID `json:"sku_id"`
}
