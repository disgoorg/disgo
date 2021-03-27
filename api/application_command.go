package api

import (
	"bytes"
	"encoding/json"
)

type ApplicationCommands []ApplicationCommand

func (commands ApplicationCommands) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	for i, command := range commands {
		commandBytes, err:= json.Marshal(command)
		if err != nil {
			return nil, err
		}
		buffer.Write(commandBytes)
		if i < len(commands) {
			buffer.WriteString(",")
		}
	}
	buffer = bytes.NewBufferString("]")
	return buffer.Bytes(), nil
}

// ApplicationCommand is the base "command" model that belongs to an application.
type ApplicationCommand struct {
	ID            Snowflake                  `json:"id,omitempty"`
	ApplicationID Snowflake                  `json:"application_id,omitempty"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Options       []ApplicationCommandOption `json:"options,omitempty"`
}

// ApplicationCommandOptionType specifies the type of the arguments used in ApplicationCommand.Options
type ApplicationCommandOptionType int

// Constants for each slash command option type
const (
	ApplicationCommandOptionTypeSubCommand ApplicationCommandOptionType = iota + 1
	ApplicationCommandOptionTypeSubCommandGroup
	ApplicationCommandOptionTypeString
	ApplicationCommandOptionTypeInteger
	ApplicationCommandOptionTypeBoolean
	ApplicationCommandOptionTypeUser
	ApplicationCommandOptionTypeChannel
	ApplicationCommandOptionTypeRole
)

// ApplicationCommandOption are the arguments used in ApplicationCommand.Options
type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType     `json:"type"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Required    bool                             `json:"required,omitempty"`
	Choices     []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options     []ApplicationCommandOption       `json:"options,omitempty"`
}

// ApplicationCommandOptionChoice contains the data for a user using your command
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
