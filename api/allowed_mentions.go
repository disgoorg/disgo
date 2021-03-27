package api

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
