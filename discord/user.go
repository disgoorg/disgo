package discord

import (
	"strconv"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
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
	UserFlagNone UserFlags = 0
)

var _ Mentionable = (*User)(nil)

// User is a struct for interacting with discord's users
type User struct {
	ID            snowflake.Snowflake `json:"id"`
	Username      string              `json:"username"`
	Discriminator string              `json:"discriminator"`
	Avatar        *string             `json:"avatar"`
	Banner        *string             `json:"banner"`
	AccentColor   *int                `json:"accent_color"`
	BotUser       bool                `json:"bot"`
	System        bool                `json:"system"`
	PublicFlags   UserFlags           `json:"public_flags"`
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
	return formatAssetURL(route.UserAvatar, opts, u.ID, u.Avatar)
}

func (u User) DefaultAvatarURL(opts ...CDNOpt) string {
	discrim, err := strconv.Atoi(u.Discriminator)
	if err != nil {
		return ""
	}
	if avatar := formatAssetURL(route.DefaultUserAvatar, opts, discrim%5); avatar != nil {
		return *avatar
	}
	return ""
}

func (u User) BannerURL(opts ...CDNOpt) *string {
	return formatAssetURL(route.UserBanner, opts, u.ID, u.Avatar)
}

// OAuth2User represents a full User returned by the oauth2 endpoints
type OAuth2User struct {
	User
	// Requires ApplicationScopeIdentify
	MfaEnabled  bool        `json:"mfa_enabled"`
	Locale      string      `json:"locale"`
	Flags       UserFlags   `json:"flags"`
	PremiumType PremiumType `json:"premium_type"`

	// Requires ApplicationScopeEmail
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
)

// SelfUserUpdate is the payload used to update the OAuth2User
type SelfUserUpdate struct {
	Username string               `json:"username"`
	Avatar   *json.Nullable[Icon] `json:"avatar"`
}
