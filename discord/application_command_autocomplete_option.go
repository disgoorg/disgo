package discord

import (
	"fmt"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
)

type internalAutocompleteOption interface {
	name() string
	autocompleteOption()
}

type UnmarshalAutocompleteOption struct {
	internalAutocompleteOption
}

func (o *UnmarshalAutocompleteOption) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	var (
		autocompleteOption internalAutocompleteOption
		err                error
	)

	switch oType.Type {
	case ApplicationCommandOptionTypeSubCommand:
		var v AutocompleteOptionSubCommand
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeSubCommandGroup:
		var v AutocompleteOptionSubCommandGroup
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	default:
		var v AutocompleteOption
		err = json.Unmarshal(data, &v)
		autocompleteOption = v
	}

	if err != nil {
		return err
	}

	o.internalAutocompleteOption = autocompleteOption
	return nil
}

var _ internalAutocompleteOption = (*AutocompleteOptionSubCommand)(nil)

type AutocompleteOptionSubCommand struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Type        ApplicationCommandOptionType `json:"type"`
	Options     []AutocompleteOption         `json:"options,omitempty"`
}

func (o AutocompleteOptionSubCommand) name() string {
	return o.Name
}
func (AutocompleteOptionSubCommand) autocompleteOption() {}

var _ internalAutocompleteOption = (*AutocompleteOptionSubCommandGroup)(nil)

type AutocompleteOptionSubCommandGroup struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Type        ApplicationCommandOptionType   `json:"type"`
	Options     []AutocompleteOptionSubCommand `json:"options,omitempty"`
}

func (o AutocompleteOptionSubCommandGroup) name() string {
	return o.Name
}
func (AutocompleteOptionSubCommandGroup) autocompleteOption() {}

var _ internalAutocompleteOption = (*AutocompleteOption)(nil)

type AutocompleteOption struct {
	Name    string                       `json:"name"`
	Type    ApplicationCommandOptionType `json:"type"`
	Value   json.RawMessage              `json:"value"`
	Focused bool                         `json:"focused"`
}

func (o AutocompleteOption) name() string {
	return o.Name
}
func (AutocompleteOption) autocompleteOption() {}

// String returns the string value of the option.
// If the type is not ApplicationCommandOptionTypeString, it panics.
func (o AutocompleteOption) String() string {
	var v string
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Int returns the int value of the option.
// If the type is not ApplicationCommandOptionTypeInt, it panics.
func (o AutocompleteOption) Int() int {
	var v int
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Float returns the float value of the option.
// If the type is not ApplicationCommandOptionTypeFloat, it panics.
func (o AutocompleteOption) Float() float64 {
	var v float64
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Bool returns the bool value of the option.
// If the type is not ApplicationCommandOptionTypeBool, it panics.
func (o AutocompleteOption) Bool() bool {
	var v bool
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}

// Snowflake returns the snowflake value of the option.
// If the type is not ApplicationCommandOptionTypeUser, ApplicationCommandOptionTypeChannel, ApplicationCommandOptionTypeRole or ApplicationCommandOptionTypeMentionable, it panics.
func (o AutocompleteOption) Snowflake() snowflake.ID {
	var v snowflake.ID
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal value of option %s: %v", o.Name, err))
	}
	return v
}
