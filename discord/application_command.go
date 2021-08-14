package discord

type ApplicationCommandType int

//goland:noinspection GoUnusedConst
const (
	ApplicationCommandTypeSlash = iota + 1
	ApplicationCommandTypeUser
	ApplicationCommandTypeMessage
)

// ApplicationCommandOptionType specifies the type of the arguments used in ApplicationCommand.Options
type ApplicationCommandOptionType int

// Constants for each slash command option type
//goland:noinspection GoUnusedConst
const (
	CommandOptionTypeSubCommand ApplicationCommandOptionType = iota + 1
	CommandOptionTypeSubCommandGroup
	CommandOptionTypeString
	CommandOptionTypeInteger
	CommandOptionTypeBoolean
	CommandOptionTypeUser
	CommandOptionTypeChannel
	CommandOptionTypeRole
	CommandOptionTypeMentionable
	CommandOptionTypeNumber
)

// ApplicationCommand is the base "command" model that belongs to an application.
type ApplicationCommand struct {
	ID                Snowflake                  `json:"id"`
	Type              ApplicationCommandType     `json:"type"`
	ApplicationID     Snowflake                  `json:"application_id"`
	GuildID           *Snowflake                 `json:"guild_id,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description,omitempty"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
}

// ApplicationCommandOption are the arguments used in ApplicationCommand.Options
type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType     `json:"type"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Required    bool                             `json:"required,omitempty"`
	Choices     []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options     []ApplicationCommandOption       `json:"options,omitempty"`
}

// ApplicationCommandOptionChoice contains the data for a user using your command. Value can either be a string, int, float or boolean
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// ApplicationCommandCreate is used to create a ApplicationCommand. all fields are optional
type ApplicationCommandCreate struct {
	Type              ApplicationCommandType     `json:"type"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission"`
}

// ApplicationCommandUpdate is used to update an existing ApplicationCommand. all fields are optional
type ApplicationCommandUpdate struct {
	Name              *string                    `json:"name,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission *bool                      `json:"default_permission,omitempty"`
}
