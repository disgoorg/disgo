package discord

// InteractionType is the type of Interaction
type InteractionType int

// Supported InteractionType(s)
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeCommand
	InteractionTypeComponent
)

// InteractionResponseType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionResponseType int

// Constants for the InteractionResponseType(s)
const (
	InteractionResponseTypePong InteractionResponseType = iota + 1
	_
	_
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
	InteractionResponseTypeDeferredUpdateMessage
	InteractionResponseTypeUpdateMessage
)

// UnmarshalInteraction is used for easier unmarshalling of different Interaction(s)
type UnmarshalInteraction struct {
	ID            Snowflake                `json:"id"`
	ApplicationID Snowflake                `json:"application_id"`
	Type          InteractionType          `json:"type"`
	Data          UnmarshalInteractionData `json:"data,omitempty"`
	GuildID       *Snowflake               `json:"guild_id,omitempty"`
	ChannelID     *Snowflake               `json:"channel_id,omitempty"`
	Member        *Member                  `json:"member,omitempty"`
	User          User                     `json:"User,omitempty"`
	Token         string                   `json:"token"`
	Version       int                      `json:"version"`
	Message       Message                  `json:"message,omitempty"`
}

type UnmarshalInteractionData struct {
	ID            Snowflake         `json:"id"`
	Name          string            `json:"name"`
	Resolved      Resolved          `json:"resolved"`
	Options       []UnmarshalOption `json:"options"`
	ComponentType ComponentType     `json:"component_type"`
	CustomID      string            `json:"custom_id"`
	Values        []string          `json:"values"`
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[Snowflake]User    `json:"users,omitempty"`
	Members  map[Snowflake]Member  `json:"members,omitempty"`
	Roles    map[Snowflake]Role    `json:"roles,omitempty"`
	Channels map[Snowflake]Channel `json:"channels,omitempty"`
}

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
}

type UnmarshalOption struct {
	Name    string            `json:"name"`
	Type    CommandOptionType `json:"type"`
	Value   interface{}       `json:"value"`
	Options []UnmarshalOption `json:"options"`
}

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionResponseType `json:"type"`
	Data interface{}             `json:"data,omitempty"`
}

// ToBody returns the InteractionResponse ready for body
func (r InteractionResponse) ToBody() (interface{}, error) {
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
