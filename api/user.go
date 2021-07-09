package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DisgoOrg/restclient"
)

// User is a struct for interacting with discord's users
type User struct {
	Disgo         Disgo
	ID            Snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Avatar        *string   `json:"avatar"`
	IsBot         bool      `json:"bot"`
	System        bool      `json:"system"`
	PublicFlags   UserFlags `json:"public_flags"`
}

// AvatarURL returns the Avatar URL of the User
func (u *User) AvatarURL(size int) string {
	if u.Avatar == nil {
		discriminator, _ := strconv.Atoi(u.Discriminator)
		route, err := restclient.DefaultUserAvatar.Compile(nil, restclient.PNG, size, discriminator%5)
		if err != nil {
			return ""
		}
		return route.Route()
	}
	format := restclient.PNG
	if strings.HasPrefix(*u.Avatar, "a_") {
		format = restclient.GIF
	}
	route, err := restclient.UserAvatar.Compile(nil, format, size, u.ID.String(), *u.Avatar)
	if err != nil {
		return ""
	}
	return route.Route()
}

// Mention returns the user as a mention
func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

// Tag returns the user's Username and Discriminator
func (u *User) Tag() string {
	return fmt.Sprintf("%s#%s", u.Username, u.Discriminator)
}

func (u *User) String() string {
	return u.Mention()
}

// CreateDMChannel creates a DMChannel between the user and the Disgo client
func (u *User) CreateDMChannel() (DMChannel, restclient.RestError) {
	return u.Disgo.RestClient().CreateDMChannel(u.ID)
}

// UserFlags defines certain flags/badges a user can have (https://discord.com/developers/docs/resources/user#user-object-user-flags)
type UserFlags int

// All UserFlags
const (
	UserFlagDiscordEmployee UserFlags = 1 << iota
	UserFlagPartneredServerOwner
	UserFlagHypeSquadEvents
	UserFlagBugHunterLevel1
	UserFlagHouseBravery
	UserFlagHouseBrilliance
	UserFlagHouseBalance
	UserFlagEarlySupporter
	UserFlagTeamUser
	UserFlagBugHunterLevel2
	UserFlagVerifiedBot
	UserFlagEarlyVerifiedBotDeveloper
	UserFlagDiscordCertifiedModerator
	UserFlagNone UserFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f UserFlags) Add(bits ...UserFlags) UserFlags {
	total := UserFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	f |= total
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f UserFlags) Remove(bits ...UserFlags) UserFlags {
	total := UserFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	f &^= total
	return f
}

// HasAll will ensure that the bit includes all of the bits entered
func (f UserFlags) HasAll(bits ...UserFlags) bool {
	for _, bit := range bits {
		if !f.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (f UserFlags) Has(bit UserFlags) bool {
	return (f & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (f UserFlags) MissingAny(bits ...UserFlags) bool {
	for _, bit := range bits {
		if !f.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of UserFlags.Has
func (f UserFlags) Missing(bit UserFlags) bool {
	return !f.Has(bit)
}