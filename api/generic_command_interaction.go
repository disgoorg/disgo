package api

type GenericCommandInteraction struct {
	*Interaction
	Data *GenericCommandInteractionData `json:"data,omitempty"`
}

// CommandID returns the ID of the api.Command which got used
func (i *GenericCommandInteraction) CommandID() Snowflake {
	return i.Data.ID
}

// CommandName the name of the api.Command which got used
func (i *GenericCommandInteraction) CommandName() string {
	return i.Data.CommandName
}

// GenericCommandInteractionData is the command data payload
type GenericCommandInteractionData struct {
	ID          Snowflake   `json:"id"`
	Type        CommandType `json:"type"`
	CommandName string      `json:"name"`
	Resolved    *Resolved   `json:"resolved,omitempty"`
	TargetID    *Snowflake  `json:"target_id"` // TODO: remove this once discord sends CommandType
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[Snowflake]*User    `json:"users,omitempty"`
	Members  map[Snowflake]*Member  `json:"members,omitempty"`
	Roles    map[Snowflake]*Role    `json:"roles,omitempty"`
	Channels map[Snowflake]*Channel `json:"channels,omitempty"`
	Messages map[Snowflake]*Message `json:"messages,omitempty"`
}
