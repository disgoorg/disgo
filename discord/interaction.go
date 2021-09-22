package discord

// InteractionType is the type of Interaction
type InteractionType int

// Supported InteractionType(s)
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeCommand
	InteractionTypeComponent
	InteractionTypeApplicationCommandAutoComplete
)

// InteractionCallbackType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionCallbackType int

// Constants for the InteractionCallbackType(s)
const (
	InteractionCallbackTypePong InteractionCallbackType = iota + 1
	_
	_
	InteractionCallbackTypeChannelMessageWithSource
	InteractionCallbackTypeDeferredChannelMessageWithSource
	InteractionCallbackTypeDeferredUpdateMessage
	InteractionCallbackTypeUpdateMessage
	InteractionCallbackTypeApplicationCommandAutoCompleteResult
)

// Interaction is used for easier unmarshalling of different Interaction(s)
type Interaction struct {
	ID            Snowflake        `json:"id"`
	ApplicationID Snowflake        `json:"application_id"`
	Type          InteractionType  `json:"type"`
	Data          *InteractionData `json:"data,omitempty"`
	GuildID       *Snowflake       `json:"guild_id,omitempty"`
	ChannelID     *Snowflake       `json:"channel_id,omitempty"`
	Member        *Member          `json:"member,omitempty"`
	User          *User            `json:"user,omitempty"`
	Token         string           `json:"token"`
	Version       int              `json:"version"`
	Message       Message          `json:"message,omitempty"`
}

type InteractionData struct {
	// Application Command Interactions
	ID          Snowflake              `json:"id"`
	CommandType ApplicationCommandType `json:"type"`
	Name        string                 `json:"name"`
	Resolved    Resolved               `json:"resolved"`

	// Slash Command Interactions
	Options []ReceivedApplicationCommandOption `json:"options"`

	// Context Command Interactions
	TargetID Snowflake `json:"target_id"`

	// Component Interactions
	ComponentType ComponentType `json:"component_type"`
	CustomID      string        `json:"custom_id"`
	Values        []string      `json:"values"`
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[Snowflake]User    `json:"users,omitempty"`
	Members  map[Snowflake]Member  `json:"members,omitempty"`
	Roles    map[Snowflake]Role    `json:"roles,omitempty"`
	Channels map[Snowflake]Channel `json:"channels,omitempty"`
	Messages map[Snowflake]Message `json:"messages,omitempty"`
}

// to consider using them in Resolved
/*
type ResolvedMember struct {
	GuildID      Snowflake   `json:"guild_id"`
	User         User        `json:"user"`
	Nick         *string     `json:"nick"`
	RoleIDs      []Snowflake `json:"roles,omitempty"`
	JoinedAt     Time        `json:"joined_at"`
	PremiumSince *Time       `json:"premium_since,omitempty"`
	Permissions  Permissions `json:"permissions,omitempty"`
}

type ResolvedChannel struct {
	ID          Snowflake   `json:"id"`
	Name        string      `json:"name"`
	Type        ChannelType `json:"type"`
	Permissions Permissions `json:"permissions"`
}*/

type ReceivedApplicationCommandOption struct {
	Name    string                             `json:"name"`
	Type    ApplicationCommandOptionType       `json:"type"`
	Value   interface{}                        `json:"value,omitempty"`
	Options []ReceivedApplicationCommandOption `json:"options,omitempty"`
	Focused bool                               `json:"focused,omitempty"`
}

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionCallbackType `json:"type"`
	Data interface{}             `json:"data,omitempty"`
}

// ToBody returns the InteractionResponse ready for body
func (r *InteractionResponse) ToBody() (interface{}, error) {
	if r.Data == nil {
		return r, nil
	}
	switch v := r.Data.(type) {
	case MessageCreate:
		if len(v.Files) > 0 {
			return PayloadWithFiles(r, v.Files...)
		}
	case MessageUpdate:
		if len(v.Files) > 0 {
			return PayloadWithFiles(r, v.Files...)
		}
	}
	return r, nil
}

type ApplicationCommandAutoCompleteResult struct {
	Choices []ApplicationCommandOptionChoice `json:"choices"`
}
