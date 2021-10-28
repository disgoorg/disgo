package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

var autocompleteOptions = map[ApplicationCommandOptionType]func() AutocompleteOption{
	ApplicationCommandOptionTypeSubCommand: func() AutocompleteOption {
		return &AutocompleteOptionSubCommand{}
	},
	ApplicationCommandOptionTypeSubCommandGroup: func() AutocompleteOption {
		return &AutocompleteOptionSubCommandGroup{}
	},
	ApplicationCommandOptionTypeString: func() AutocompleteOption {
		return &AutocompleteOptionString{}
	},
	ApplicationCommandOptionTypeInt: func() AutocompleteOption {
		return &AutocompleteOptionInt{}
	},
	ApplicationCommandOptionTypeBool: func() AutocompleteOption {
		return &AutocompleteOptionBool{}
	},
	ApplicationCommandOptionTypeUser: func() AutocompleteOption {
		return &AutocompleteOptionUser{}
	},
	ApplicationCommandOptionTypeChannel: func() AutocompleteOption {
		return &AutocompleteOptionChannel{}
	},
	ApplicationCommandOptionTypeRole: func() AutocompleteOption {
		return &AutocompleteOptionRole{}
	},
	ApplicationCommandOptionTypeMentionable: func() AutocompleteOption {
		return &AutocompleteOptionMentionable{}
	},
	ApplicationCommandOptionTypeFloat: func() AutocompleteOption {
		return &AutocompleteOptionFloat{}
	},
}

type AutocompleteOption interface {
	Type() ApplicationCommandOptionType
}

type unmarshalAutocompleteOption struct {
	AutocompleteOption
}

func (o unmarshalAutocompleteOption) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	fn, ok := autocompleteOptions[oType.Type]
	if !ok {
		return fmt.Errorf("unkown application command autocomplete option with type %d received", oType.Type)
	}

	v := fn()

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	o.AutocompleteOption = v
	return nil
}

var _ AutocompleteOption = (*AutocompleteOptionSubCommand)(nil)

type AutocompleteOptionSubCommand struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Options     []AutocompleteOption `json:"options,omitempty"`
}

func (_ AutocompleteOptionSubCommand) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

var _ AutocompleteOption = (*AutocompleteOptionSubCommandGroup)(nil)

type AutocompleteOptionSubCommandGroup struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Options     []AutocompleteOptionSubCommand `json:"options,omitempty"`
}

func (_ AutocompleteOptionSubCommandGroup) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

var _ AutocompleteOption = (*AutocompleteOptionString)(nil)

type AutocompleteOptionString struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Focused bool   `json:"focused"`
}

func (_ AutocompleteOptionString) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

var _ AutocompleteOption = (*AutocompleteOptionInt)(nil)

type AutocompleteOptionInt struct {
	Name    string `json:"name"`
	Value   int    `json:"value"`
	Focused bool   `json:"focused"`
}

func (_ AutocompleteOptionInt) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInt
}

var _ AutocompleteOption = (*AutocompleteOptionBool)(nil)

type AutocompleteOptionBool struct {
	Name    string `json:"name"`
	Value   bool   `json:"value"`
	Focused bool   `json:"focused"`
}

func (_ AutocompleteOptionBool) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBool
}

var _ AutocompleteOption = (*AutocompleteOptionUser)(nil)

type AutocompleteOptionUser struct {
	Name    string    `json:"name"`
	Value   Snowflake `json:"value"`
	Focused bool      `json:"focused"`
}

func (_ AutocompleteOptionUser) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

var _ AutocompleteOption = (*AutocompleteOptionChannel)(nil)

type AutocompleteOptionChannel struct {
	Name    string    `json:"name"`
	Value   Snowflake `json:"value"`
	Focused bool      `json:"focused"`
}

func (_ AutocompleteOptionChannel) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

var _ AutocompleteOption = (*AutocompleteOptionRole)(nil)

type AutocompleteOptionRole struct {
	Name    string    `json:"name"`
	Value   Snowflake `json:"value"`
	Focused bool      `json:"focused"`
}

func (_ AutocompleteOptionRole) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

var _ AutocompleteOption = (*AutocompleteOptionMentionable)(nil)

type AutocompleteOptionMentionable struct {
	Name    string    `json:"name"`
	Value   Snowflake `json:"value"`
	Focused bool      `json:"focused"`
}

func (_ AutocompleteOptionMentionable) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

var _ AutocompleteOption = (*AutocompleteOptionFloat)(nil)

type AutocompleteOptionFloat struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Focused bool    `json:"focused"`
}

func (_ AutocompleteOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}
