package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

type AutocompleteOption interface {
	Type() ApplicationCommandOptionType
	Name() string
	Focused() bool
	autocompleteOption()
}

type UnmarshalAutocompleteOption struct {
	AutocompleteOption
}

func (o *UnmarshalAutocompleteOption) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	var (
		autocompleteOption AutocompleteOption
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

	case ApplicationCommandOptionTypeString:
		var v AutocompleteOptionString
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeInt:
		var v AutocompleteOptionInt
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeBool:
		var v AutocompleteOptionBool
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeUser:
		var v AutocompleteOptionUser
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeChannel:
		var v AutocompleteOptionChannel
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeRole:
		var v AutocompleteOptionRole
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeMentionable:
		var v AutocompleteOptionMentionable
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	case ApplicationCommandOptionTypeFloat:
		var v AutocompleteOptionFloat
		err = json.Unmarshal(data, &v)
		autocompleteOption = v

	default:
		err = fmt.Errorf("unkown autocomplete option with type %d received", oType.Type)
	}

	if err != nil {
		return err
	}

	o.AutocompleteOption = autocompleteOption
	return nil
}

var _ AutocompleteOption = (*AutocompleteOptionSubCommand)(nil)

type AutocompleteOptionSubCommand struct {
	CommandName string               `json:"name"`
	Description string               `json:"description"`
	Options     []AutocompleteOption `json:"options,omitempty"`
}

func (AutocompleteOptionSubCommand) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

func (o AutocompleteOptionSubCommand) Name() string {
	return o.CommandName
}

func (o AutocompleteOptionSubCommand) Focused() bool {
	return false
}

func (AutocompleteOptionSubCommand) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionSubCommandGroup)(nil)

type AutocompleteOptionSubCommandGroup struct {
	GroupName   string                         `json:"name"`
	Description string                         `json:"description"`
	Options     []AutocompleteOptionSubCommand `json:"options,omitempty"`
}

func (AutocompleteOptionSubCommandGroup) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

func (o AutocompleteOptionSubCommandGroup) Name() string {
	return o.GroupName
}

func (o AutocompleteOptionSubCommandGroup) Focused() bool {
	return false
}

func (AutocompleteOptionSubCommandGroup) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionString)(nil)

type AutocompleteOptionString struct {
	OptionName    string `json:"name"`
	Value         string `json:"value"`
	OptionFocused bool   `json:"focused"`
}

func (AutocompleteOptionString) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

func (o AutocompleteOptionString) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionString) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionString) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionInt)(nil)

type AutocompleteOptionInt struct {
	OptionName    string `json:"name"`
	Value         int    `json:"value"`
	OptionFocused bool   `json:"focused"`
}

func (AutocompleteOptionInt) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInt
}

func (o AutocompleteOptionInt) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionInt) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionInt) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionBool)(nil)

type AutocompleteOptionBool struct {
	OptionName    string `json:"name"`
	Value         bool   `json:"value"`
	OptionFocused bool   `json:"focused"`
}

func (AutocompleteOptionBool) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBool
}

func (o AutocompleteOptionBool) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionBool) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionBool) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionUser)(nil)

type AutocompleteOptionUser struct {
	OptionName    string              `json:"name"`
	Value         snowflake.Snowflake `json:"value"`
	OptionFocused bool                `json:"focused"`
}

func (AutocompleteOptionUser) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

func (o AutocompleteOptionUser) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionUser) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionUser) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionChannel)(nil)

type AutocompleteOptionChannel struct {
	OptionName    string              `json:"name"`
	Value         snowflake.Snowflake `json:"value"`
	OptionFocused bool                `json:"focused"`
}

func (AutocompleteOptionChannel) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

func (o AutocompleteOptionChannel) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionChannel) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionChannel) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionRole)(nil)

type AutocompleteOptionRole struct {
	OptionName    string              `json:"name"`
	Value         snowflake.Snowflake `json:"value"`
	OptionFocused bool                `json:"focused"`
}

func (AutocompleteOptionRole) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

func (o AutocompleteOptionRole) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionRole) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionRole) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionMentionable)(nil)

type AutocompleteOptionMentionable struct {
	OptionName    string              `json:"name"`
	Value         snowflake.Snowflake `json:"value"`
	OptionFocused bool                `json:"focused"`
}

func (AutocompleteOptionMentionable) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

func (o AutocompleteOptionMentionable) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionMentionable) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionMentionable) autocompleteOption() {}

var _ AutocompleteOption = (*AutocompleteOptionFloat)(nil)

type AutocompleteOptionFloat struct {
	OptionName    string  `json:"name"`
	Value         float64 `json:"value"`
	OptionFocused bool    `json:"focused"`
}

func (AutocompleteOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}

func (o AutocompleteOptionFloat) Name() string {
	return o.OptionName
}

func (o AutocompleteOptionFloat) Focused() bool {
	return o.OptionFocused
}

func (AutocompleteOptionFloat) autocompleteOption() {}
