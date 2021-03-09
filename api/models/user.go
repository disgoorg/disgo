package models

import (
	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/api"
)

type User struct {
	Disgo         api.Disgo
	ID            Snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Avatar        *string   `json:"avatar"`
	Bot           *bool     `json:"bot"`
	System        *bool     `json:"system"`
	MfaEnabled    *bool     `json:"mfa_enabled"`
	Locale        *string   `json:"locale"`
	Verified      *bool     `json:"verified"`
	Email         *string   `json:"email"`
	Flags         *int      `json:"flags"`
	PremiumType   *int             `json:"premium_type"`
	PublicFlags   *int             `json:"public_flags"`
}

func (u User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

func (u User) OpenDMChannel() *promise.Promise {
	return u.Disgo.RestClient().OpenDMChannel(u.ID).Then(func(channel promise.Any) promise.Any {
		channel.(*DMChannel).Disgo = u.Disgo
		return channel
	})
}

type CreateDMChannel struct {
	RecipientID Snowflake `json:"recipient_id"`
}
