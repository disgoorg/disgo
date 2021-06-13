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
	System        *bool     `json:"system"`
	MfaEnabled    *bool     `json:"mfa_enabled"`
	Locale        *string   `json:"locale"`
	Verified      *bool     `json:"verified"`
	Email         *string   `json:"email"`
	Flags         *int      `json:"flags"`
	PremiumType   *int      `json:"premium_type"`
	PublicFlags   *int      `json:"public_flags"`
}

// AvatarURL returns the Avatar URL of the User
func (u *User) AvatarURL(size int) string {
	if u.Avatar == nil {
		discrim, _ := strconv.Atoi(u.Discriminator)
		route, err := restclient.DefaultUserAvatar.Compile(nil, restclient.PNG, size, discrim%5)
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

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u *User) OpenDMChannel() (*DMChannel, error) {
	return u.Disgo.RestClient().OpenDMChannel(u.ID)
}
