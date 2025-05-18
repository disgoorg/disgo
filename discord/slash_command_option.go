package discord

import (
	"fmt"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

type internalSlashCommandOption interface {
	name() string
	slashCommandOption()
}

type UnmarshalSlashCommandOption struct {
	internalSlashCommandOption
}

func (o *UnmarshalSlashCommandOption) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type ApplicationCommandOptionType `json:"type"`
	}
	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	var (
		slashCommandOption internalSlashCommandOption
		err                error
	)

	switch oType.Type {
	case ApplicationCommandOptionTypeSubCommand:
		var v SlashCommandOptionSubCommand
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeSubCommandGroup:
		var v SlashCommandOptionSubCommandGroup
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	default:
		var v SlashCommandOption
		err = json.Unmarshal(data, &v)
		slashCommandOption = v
	}
	if err != nil {
		return err
	}

	o.internalSlashCommandOption = slashCommandOption
	return nil
}

var _ internalSlashCommandOption = (*SlashCommandOptionSubCommand)(nil)

type SlashCommandOptionSubCommand struct {
	Name    string                       `json:"name"`
	Type    ApplicationCommandOptionType `json:"type"`
	Options []SlashCommandOption         `json:"options,omitempty"`
}

func (o SlashCommandOptionSubCommand) name() string {
	return o.Name
}
func (SlashCommandOptionSubCommand) slashCommandOption() {}

var _ internalSlashCommandOption = (*SlashCommandOptionSubCommandGroup)(nil)

type SlashCommandOptionSubCommandGroup struct {
	Name    string                         `json:"name"`
	Type    ApplicationCommandOptionType   `json:"type"`
	Options []SlashCommandOptionSubCommand `json:"options,omitempty"`
}

func (o SlashCommandOptionSubCommandGroup) name() string {
	return o.Name
}
func (SlashCommandOptionSubCommandGroup) slashCommandOption() {}

var _ internalSlashCommandOption = (*SlashCommandOption)(nil)

type SlashCommandOption struct {
	Name  string                       `json:"name"`
	Type  ApplicationCommandOptionType `json:"type"`
	Value json.RawMessage              `json:"value"`
}

func (o SlashCommandOption) name() string {
	return o.Name
}
func (SlashCommandOption) slashCommandOption() {}

// String returns the string value of the option.
// If the type is not ApplicationCommandOptionTypeString, it panics.
func (o SlashCommandOption) String() string {
	var v string
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Int returns the int value of the option.
// If the type is not ApplicationCommandOptionTypeInt, it panics.
func (o SlashCommandOption) Int() int {
	var v int
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Float returns the float value of the option.
// If the type is not ApplicationCommandOptionTypeFloat, it panics.
func (o SlashCommandOption) Float() float64 {
	var v float64
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Bool returns the bool value of the option.
// If the type is not ApplicationCommandOptionTypeBool, it panics.
func (o SlashCommandOption) Bool() bool {
	var v bool
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Snowflake returns the snowflake value of the option.
// If the type is not ApplicationCommandOptionTypeUser, ApplicationCommandOptionTypeChannel, ApplicationCommandOptionTypeRole or ApplicationCommandOptionTypeMentionable, it panics.
func (o SlashCommandOption) Snowflake() snowflake.ID {
	var v snowflake.ID
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}
