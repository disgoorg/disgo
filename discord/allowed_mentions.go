package discord

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
	AllowedMentionTypeRoles    AllowedMentionType = "roles"
	AllowedMentionTypeUsers    AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)
