package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

// ApplicationCommandOptionType specifies the type of the arguments used in ApplicationCommand.Options
type ApplicationCommandOptionType int

// Constants for each slash command option type
//goland:noinspection GoUnusedConst
const (
	ApplicationCommandOptionTypeSubCommand ApplicationCommandOptionType = iota + 1
	ApplicationCommandOptionTypeSubCommandGroup
	ApplicationCommandOptionTypeString
	ApplicationCommandOptionTypeInt
	ApplicationCommandOptionTypeBool
	ApplicationCommandOptionTypeUser
	ApplicationCommandOptionTypeChannel
	ApplicationCommandOptionTypeRole
	ApplicationCommandOptionTypeMentionable
	ApplicationCommandOptionTypeFloat
)

var applicationCommandOptions = map[ApplicationCommandOptionType]func() ApplicationCommandOption{
	ApplicationCommandOptionTypeSubCommand: func() ApplicationCommandOption {
		return &ApplicationCommandOptionSubCommand{}
	},
	ApplicationCommandOptionTypeSubCommandGroup: func() ApplicationCommandOption {
		return &ApplicationCommandOptionSubCommandGroup{}
	},
	ApplicationCommandOptionTypeString: func() ApplicationCommandOption {
		return &ApplicationCommandOptionString{}
	},
	ApplicationCommandOptionTypeInt: func() ApplicationCommandOption {
		return &ApplicationCommandOptionInt{}
	},
	ApplicationCommandOptionTypeBool: func() ApplicationCommandOption {
		return &ApplicationCommandOptionBool{}
	},
	ApplicationCommandOptionTypeUser: func() ApplicationCommandOption {
		return &ApplicationCommandOptionUser{}
	},
	ApplicationCommandOptionTypeChannel: func() ApplicationCommandOption {
		return &ApplicationCommandOptionChannel{}
	},
	ApplicationCommandOptionTypeRole: func() ApplicationCommandOption {
		return &ApplicationCommandOptionRole{}
	},
	ApplicationCommandOptionTypeMentionable: func() ApplicationCommandOption {
		return &ApplicationCommandOptionMentionable{}
	},
	ApplicationCommandOptionTypeFloat: func() ApplicationCommandOption {
		return &ApplicationCommandOptionFloat{}
	},
}

type ApplicationCommandOption interface {
	json.Marshaler
	Type() ApplicationCommandOptionType
}

type UnmarshalApplicationCommandOption struct {
	ApplicationCommandOption
}

func (u *UnmarshalApplicationCommandOption) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	fn, ok := applicationCommandOptions[oType.Type]
	if !ok {
		return fmt.Errorf("unkown application command option with type %d received", oType.Type)
	}

	v := fn()

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	u.ApplicationCommandOption = v
	return nil
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionSubCommand)(nil)

type ApplicationCommandOptionSubCommand struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Options     []ApplicationCommandOption `json:"options,omitempty"`
}

func (o ApplicationCommandOptionSubCommand) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionSubCommand ApplicationCommandOptionSubCommand
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionSubCommand
	}{
		Type:                               o.Type(),
		applicationCommandOptionSubCommand: applicationCommandOptionSubCommand(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionSubCommand) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionSubCommandGroup)(nil)

type ApplicationCommandOptionSubCommandGroup struct {
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Options     []ApplicationCommandOptionSubCommand `json:"options,omitempty"`
}

func (o ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionSubCommandGroup ApplicationCommandOptionSubCommandGroup
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionSubCommandGroup
	}{
		Type:                                    o.Type(),
		applicationCommandOptionSubCommandGroup: applicationCommandOptionSubCommandGroup(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionSubCommandGroup) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionString)(nil)

type ApplicationCommandOptionString struct {
	Name         string                                 `json:"name"`
	Description  string                                 `json:"description"`
	Required     bool                                   `json:"required,omitempty"`
	Choices      []ApplicationCommandOptionChoiceString `json:"choices,omitempty"`
	Autocomplete bool                                   `json:"autocomplete,omitempty"`
}

func (o ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionString ApplicationCommandOptionString
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionString
	}{
		Type:                           o.Type(),
		applicationCommandOptionString: applicationCommandOptionString(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionString) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionInt)(nil)

type ApplicationCommandOptionInt struct {
	Name         string                              `json:"name"`
	Description  string                              `json:"description"`
	Required     bool                                `json:"required,omitempty"`
	Choices      []ApplicationCommandOptionChoiceInt `json:"choices,omitempty"`
	Autocomplete bool                                `json:"autocomplete,omitempty"`
	MinValue     OptionalInt                         `json:"min_value,omitempty"`
	MaxValue     OptionalInt                         `json:"max_value,omitempty"`
}

func (o ApplicationCommandOptionInt) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionInt ApplicationCommandOptionInt
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionInt
	}{
		Type:                        o.Type(),
		applicationCommandOptionInt: applicationCommandOptionInt(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionInt) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInt
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionBool)(nil)

type ApplicationCommandOptionBool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionBool) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionBool ApplicationCommandOptionBool
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionBool
	}{
		Type:                         o.Type(),
		applicationCommandOptionBool: applicationCommandOptionBool(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionBool) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBool
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionUser)(nil)

type ApplicationCommandOptionUser struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionUser ApplicationCommandOptionUser
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionUser
	}{
		Type:                         o.Type(),
		applicationCommandOptionUser: applicationCommandOptionUser(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionUser) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionChannel)(nil)

type ApplicationCommandOptionChannel struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Required     bool          `json:"required,omitempty"`
	ChannelTypes []ChannelType `json:"channel_types"`
}

func (o ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionChannel ApplicationCommandOptionChannel
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionChannel
	}{
		Type:                            o.Type(),
		applicationCommandOptionChannel: applicationCommandOptionChannel(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionChannel) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionRole)(nil)

type ApplicationCommandOptionRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionRole) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionRole ApplicationCommandOptionRole
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionRole
	}{
		Type:                         o.Type(),
		applicationCommandOptionRole: applicationCommandOptionRole(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionRole) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionMentionable)(nil)

type ApplicationCommandOptionMentionable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionMentionable) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionMentionable ApplicationCommandOptionMentionable
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionMentionable
	}{
		Type:                                o.Type(),
		applicationCommandOptionMentionable: applicationCommandOptionMentionable(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionMentionable) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionFloat)(nil)

type ApplicationCommandOptionFloat struct {
	Name         string                                `json:"name"`
	Description  string                                `json:"description"`
	Required     bool                                  `json:"required,omitempty"`
	Choices      []ApplicationCommandOptionChoiceFloat `json:"choices,omitempty"`
	Autocomplete bool                                  `json:"autocomplete,omitempty"`
	MinValue     OptionalFloat                         `json:"min_value,omitempty"`
	MaxValue     OptionalFloat                         `json:"max_value,omitempty"`
}

func (o ApplicationCommandOptionFloat) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionFloat ApplicationCommandOptionFloat
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionFloat
	}{
		Type:                          o.Type(),
		applicationCommandOptionFloat: applicationCommandOptionFloat(o),
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}

type ApplicationCommandOptionChoice interface {
	choice()
}

type ApplicationCommandOptionChoiceInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (_ ApplicationCommandOptionChoiceInt) choice() {}

type ApplicationCommandOptionChoiceString struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (_ ApplicationCommandOptionChoiceString) choice() {}

type ApplicationCommandOptionChoiceFloat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (_ ApplicationCommandOptionChoiceFloat) choice() {}
