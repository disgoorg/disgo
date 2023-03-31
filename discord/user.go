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
	ID            snowflake.ID `json:"id"`
	Username      string       `json:"username"`
	Discriminator string       `json:"discriminator"`
	Avatar        *string      `json:"avatar"`
	Banner        *string      `json:"banner"`
	AccentColor   *int         `json:"accent_color"`
	Bot           bool         `json:"bot"`
	System        bool         `json:"system"`
	PublicFlags   UserFlags    `json:"public_flags"`
}

func (u User) String() string {
	return UserMention(u.ID)
}

func (u User) Mention() string {
	return u.String()
}

func (u User) Tag() string {
	return UserTag(u.Username, u.Discriminator)
}

func (u User) EffectiveAvatarURL(opts ...CDNOpt) string {
	if u.Avatar == nil {
		return u.DefaultAvatarURL(opts...)
	}
	if avatar := u.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (u User) AvatarURL(opts ...CDNOpt) *string {
	if u.Avatar == nil {
		return nil
	}
	url := formatAssetURL(UserAvatar, opts, u.ID, *u.Avatar)
	return &url
}

func (u User) DefaultAvatarURL(opts ...CDNOpt) string {
	discriminator, err := strconv.Atoi(u.Discriminator)
	if err != nil {
		return ""
	}
	return formatAssetURL(DefaultUserAvatar, opts, discriminator%5)
}

func (u User) BannerURL(opts ...CDNOpt) *string {
	if u.Banner == nil {
		return nil
	}
	url := formatAssetURL(UserBanner, opts, u.ID, *u.Banner)
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

// SelfUserUpdate is the payload used to update the OAuth2User
type SelfUserUpdate struct {
	Username string               `json:"username,omitempty"`
	Avatar   *json.Nullable[Icon] `json:"avatar,omitempty"`
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
