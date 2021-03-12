package api

type InteractionType int

const (
	InteractionTypePing = iota + 1
	InteractionTypeApplicationCommand
)

type Interaction struct {
	ID        Snowflake       `json:"id"`
	Type      InteractionType `json:"type"`
	Data      InteractionData `json:"data,omitempty"`
	GuildID   Snowflake       `json:"guild_id,omitempty"`
	ChannelID Snowflake       `json:"channel_id,omitempty"`
	Member    *Member         `json:"member,omitempty"`
	User      *User           `json:"User,omitempty"`
	Token     string          `json:"token"`
	Version   int             `json:"version"`
}

type InteractionData struct {
	ID      Snowflake    `json:"id"`
	Name    string       `json:"name"`
	Options []OptionData `json:"options,omitempty"`
}

type OptionData struct {
	Name    string       `json:"name"`
	Value   interface{}  `json:"value,omitempty"`
	Options []OptionData `json:"options,omitempty"`
}
