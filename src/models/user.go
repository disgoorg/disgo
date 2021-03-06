package models

type User struct {
	ID            Snowflake  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	Avatar        *string  `json:"avatar"`
	Bot           *bool   `json:"bot"`
	System        *bool   `json:"system"`
	MfaEnabled    *bool   `json:"mfa_enabled"`
	Locale        *string `json:"locale"`
	Verified      *bool    `json:"verified"`
	Email         *string  `json:"email"`
	Flags         *int     `json:"flags"`
	PremiumType   *int     `json:"premium_type"`
	PublicFlags   *int     `json:"public_flags"`
}
