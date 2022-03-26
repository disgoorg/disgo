package discord

import "github.com/disgoorg/snowflake"

// DefaultAllowedMentions gives you the default AllowedMentions for a Message
var DefaultAllowedMentions = AllowedMentions{
	Parse:       []AllowedMentionType{AllowedMentionTypeUsers, AllowedMentionTypeRoles, AllowedMentionTypeEveryone},
	Roles:       []snowflake.Snowflake{},
	Users:       []snowflake.Snowflake{},
	RepliedUser: true,
}

// AllowedMentions are used for avoiding mentioning users in Message and Interaction
type AllowedMentions struct {
	Parse       []AllowedMentionType  `json:"parse"`
	Roles       []snowflake.Snowflake `json:"roles"`
	Users       []snowflake.Snowflake `json:"users"`
	RepliedUser bool                  `json:"replied_user"`
}

// AllowedMentionType ?
type AllowedMentionType string

// All AllowedMentionType(s)
const (
	AllowedMentionTypeRoles    AllowedMentionType = "roles"
	AllowedMentionTypeUsers    AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)
