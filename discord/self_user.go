package discord

// SelfUser represents the current logged in User
type SelfUser struct {
	User
	MfaEnabled  *bool   `json:"mfa_enabled"`
	Locale      *string `json:"locale"`
	Verified    *bool   `json:"verified"`
	Email       *string `json:"email"`
	Flags       *int    `json:"flags"`
	PremiumType *int    `json:"premium_type"`
}

// SelfUserUpdate is the payload used to update the SelfUser
type SelfUserUpdate struct {
	Username string `json:"username"`
	Avatar   *Icon  `json:"avatar"`
}
