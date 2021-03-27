package api

import (
	"bytes"
	"encoding/json"
)

// SlashCommands is a slice of SlashCommand
type SlashCommands []SlashCommand

// MarshalJSON is used for marshalling multiple commands into a []byte
func (commands SlashCommands) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	for i, command := range commands {
		commandBytes, err := json.Marshal(command)
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

// SlashCommand is the base "command" model that belongs to an application.
type SlashCommand struct {
	ID            Snowflake        `json:"id,omitempty"`
	ApplicationID Snowflake        `json:"application_id,omitempty"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Options       []*CommandOption `json:"options,omitempty"`
}

// SlashCommandOptionType specifies the type of the arguments used in SlashCommand.Options
type SlashCommandOptionType int

// Constants for each slash command option type
const (
	OptionTypeSubCommand SlashCommandOptionType = iota + 1
	OptionTypeSubCommandGroup
	OptionTypeString
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeUser
	OptionTypeChannel
	OptionTypeRole
)

// CommandOption are the arguments used in SlashCommand.Options
type CommandOption struct {
	Type        SlashCommandOptionType `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Required    bool                   `json:"required,omitempty"`
	Choices     []OptionChoice         `json:"choices,omitempty"`
	Options     []CommandOption        `json:"options,omitempty"`
}

// OptionChoice contains the data for a user using your command
type OptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
