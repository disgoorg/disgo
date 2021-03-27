package api

// DefaultInteractionAllowedMentions gives you the default AllowedMentions for an Interaction
var DefaultInteractionAllowedMentions = AllowedMentions{
	Parse:       []AllowedMentionType{AllowedMentionTypeUser},
	Roles:       []Snowflake{},
	Users:       []Snowflake{},
	RepliedUser: false,
}

// DefaultMessageAllowedMentions gives you the default AllowedMentions for a Message
var DefaultMessageAllowedMentions = AllowedMentions{
	Parse:       []AllowedMentionType{AllowedMentionTypeUser, AllowedMentionTypeRole, AllowedMentionTypeEveryone},
	Roles:       []Snowflake{},
	Users:       []Snowflake{},
	RepliedUser: true,
}

// AllowedMentions are used for avoiding mentioning users in Message and Interaction
type AllowedMentions struct {
	Parse       []AllowedMentionType `json:"parse"`
	Roles       []Snowflake          `json:"roles"`
	Users       []Snowflake          `json:"users"`
	RepliedUser bool                 `json:"replied_user"`
}

// AllowedMentionType ?
type AllowedMentionType string

// All AllowedMentionType(s)
const (
	AllowedMentionTypeRole     AllowedMentionType = "roles"
	AllowedMentionTypeUser     AllowedMentionType = "user"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)
