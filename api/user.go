package api

import (
	"github.com/chebyrash/promise"
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

// Mention returns the user as a mention
func (u User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

// Tag returns the user's Username and Discriminator
func (u User) Tag() string {
	return u.Username + "#" + u.Discriminator
}

func (u User) String() string {
	return u.Mention()
}

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u User) OpenDMChannel() *promise.Promise {
	return u.Disgo.RestClient().OpenDMChannel(u.ID).Then(func(channel promise.Any) promise.Any {
		channel.(*DMChannel).Disgo = u.Disgo
		return channel
	})
}
