package api

// InteractionType is the type of Interaction
type InteractionType int

// Constants for InteractionType
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
)

// An Interaction is the slash command object you receive when a user uses one of your commands
type Interaction struct {
	ID        Snowflake        `json:"id"`
	Type      InteractionType  `json:"type"`
	Data      *InteractionData `json:"data,omitempty"`
	GuildID   *Snowflake       `json:"guild_id,omitempty"`
	ChannelID *Snowflake       `json:"channel_id,omitempty"`
	Member    *Member          `json:"member,omitempty"`
	User      *User            `json:"User,omitempty"`
	Token     string           `json:"token"`
	Version   int              `json:"version"`
}

// InteractionData is the command data payload
type InteractionData struct {
	ID      Snowflake    `json:"id"`
	Name    string       `json:"name"`
	Options []OptionData `json:"options,omitempty"`
}

// OptionData is used for options or subcommands in your slash commands
type OptionData struct {
	Name    string                       `json:"name"`
	Type    ApplicationCommandOptionType `json:"type"`
	Value   interface{}                  `json:"value,omitempty"`
	Options []OptionData                 `json:"options,omitempty"`
}


func (v OptionData) String() string {
	return v.Value.(string)
}

func (v OptionData) Bool() bool {
	return v.Value.(bool)
}

func (v OptionData) Snowflake() Snowflake {
	return Snowflake(v.String())
}
