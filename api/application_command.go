package api

import (
	"bytes"
	"encoding/json"
)

type Commands []Command

func (commands Commands) MarshalJSON() ([]byte, error) {
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
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

// Command is the base "command" model that belongs to an application.
type Command struct {
	ID            Snowflake       `json:"id,omitempty"`
	ApplicationID Snowflake       `json:"application_id,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Options       []*CommandOption `json:"options,omitempty"`
}

// CommandOptionType specifies the type of the arguments used in Command.Options
type CommandOptionType int

// Constants for each slash command option type
const (
	CommandOptionTypeSubCommand CommandOptionType = iota + 1
	CommandOptionTypeSubCommandGroup
	CommandOptionTypeString
	CommandOptionTypeInteger
	CommandOptionTypeBoolean
	CommandOptionTypeUser
	CommandOptionTypeChannel
	CommandOptionTypeRole
)

// CommandOption are the arguments used in Command.Options
type CommandOption struct {
	Type        CommandOptionType `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required,omitempty"`
	Choices     []OptionChoice    `json:"choices,omitempty"`
	Options     []CommandOption   `json:"options,omitempty"`
}

// OptionChoice contains the data for a user using your command
type OptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
