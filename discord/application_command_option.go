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

type ApplicationCommandOption interface {
	json.Marshaler
	Type() ApplicationCommandOptionType
	applicationCommandOption()
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

	var (
		applicationCommandOption ApplicationCommandOption
		err                      error
	)

	switch oType.Type {
	case ApplicationCommandOptionTypeSubCommand:
		var v ApplicationCommandOptionSubCommand
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeSubCommandGroup:
		var v ApplicationCommandOptionSubCommandGroup
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeString:
		var v ApplicationCommandOptionString
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeInt:
		var v ApplicationCommandOptionInt
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeBool:
		var v ApplicationCommandOptionBool
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeUser:
		var v ApplicationCommandOptionUser
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeChannel:
		var v ApplicationCommandOptionChannel
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeRole:
		var v ApplicationCommandOptionRole
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeMentionable:
		var v ApplicationCommandOptionMentionable
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	case ApplicationCommandOptionTypeFloat:
		var v ApplicationCommandOptionFloat
		err = json.Unmarshal(data, &v)
		applicationCommandOption = v

	default:
		err = fmt.Errorf("unkown application command option with type %d received", oType.Type)
	}

	if err != nil {
		return err
	}

	u.ApplicationCommandOption = applicationCommandOption
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
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionSubCommand
	}{
		Type:                               o.Type(),
		applicationCommandOptionSubCommand: applicationCommandOptionSubCommand(o),
	})
}

func (o *ApplicationCommandOptionSubCommand) UnmarshalJSON(data []byte) error {
	type applicationCommandOptionSubCommand ApplicationCommandOptionSubCommand
	var v struct {
		Options []UnmarshalApplicationCommandOption `json:"options"`
		applicationCommandOptionSubCommand
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*o = ApplicationCommandOptionSubCommand(v.applicationCommandOptionSubCommand)

	if len(v.Options) > 0 {
		o.Options = make([]ApplicationCommandOption, len(v.Options))
		for i := range v.Options {
			o.Options[i] = v.Options[i].ApplicationCommandOption
		}
	}

	return nil
}

func (_ ApplicationCommandOptionSubCommand) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

func (_ ApplicationCommandOptionSubCommand) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionSubCommandGroup)(nil)

type ApplicationCommandOptionSubCommandGroup struct {
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Options     []ApplicationCommandOptionSubCommand `json:"options,omitempty"`
}

func (o ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionSubCommandGroup ApplicationCommandOptionSubCommandGroup
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionSubCommandGroup
	}{
		Type:                                    o.Type(),
		applicationCommandOptionSubCommandGroup: applicationCommandOptionSubCommandGroup(o),
	})
}

func (_ ApplicationCommandOptionSubCommandGroup) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

func (_ ApplicationCommandOptionSubCommandGroup) applicationCommandOption() {}

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
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionString
	}{
		Type:                           o.Type(),
		applicationCommandOptionString: applicationCommandOptionString(o),
	})
}

func (_ ApplicationCommandOptionString) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

func (_ ApplicationCommandOptionString) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionInt)(nil)

type ApplicationCommandOptionInt struct {
	Name         string                              `json:"name"`
	Description  string                              `json:"description"`
	Required     bool                                `json:"required,omitempty"`
	Choices      []ApplicationCommandOptionChoiceInt `json:"choices,omitempty"`
	Autocomplete bool                                `json:"autocomplete,omitempty"`
	MinValue     *json.NullInt                       `json:"min_value,omitempty"`
	MaxValue     *json.NullInt                       `json:"max_value,omitempty"`
}

func (o ApplicationCommandOptionInt) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionInt ApplicationCommandOptionInt
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionInt
	}{
		Type:                        o.Type(),
		applicationCommandOptionInt: applicationCommandOptionInt(o),
	})
}

func (_ ApplicationCommandOptionInt) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInt
}

func (_ ApplicationCommandOptionInt) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionBool)(nil)

type ApplicationCommandOptionBool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionBool) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionBool ApplicationCommandOptionBool
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionBool
	}{
		Type:                         o.Type(),
		applicationCommandOptionBool: applicationCommandOptionBool(o),
	})
}

func (_ ApplicationCommandOptionBool) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBool
}

func (_ ApplicationCommandOptionBool) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionUser)(nil)

type ApplicationCommandOptionUser struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionUser ApplicationCommandOptionUser
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionUser
	}{
		Type:                         o.Type(),
		applicationCommandOptionUser: applicationCommandOptionUser(o),
	})
}

func (_ ApplicationCommandOptionUser) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

func (_ ApplicationCommandOptionUser) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionChannel)(nil)

type ApplicationCommandOptionChannel struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Required     bool          `json:"required,omitempty"`
	ChannelTypes []ChannelType `json:"channel_types,omitempty"`
}

func (o ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionChannel ApplicationCommandOptionChannel
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionChannel
	}{
		Type:                            o.Type(),
		applicationCommandOptionChannel: applicationCommandOptionChannel(o),
	})
}

func (_ ApplicationCommandOptionChannel) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

func (_ ApplicationCommandOptionChannel) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionRole)(nil)

type ApplicationCommandOptionRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionRole) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionRole ApplicationCommandOptionRole
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionRole
	}{
		Type:                         o.Type(),
		applicationCommandOptionRole: applicationCommandOptionRole(o),
	})
}

func (_ ApplicationCommandOptionRole) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

func (_ ApplicationCommandOptionRole) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionMentionable)(nil)

type ApplicationCommandOptionMentionable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

func (o ApplicationCommandOptionMentionable) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionMentionable ApplicationCommandOptionMentionable
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionMentionable
	}{
		Type:                                o.Type(),
		applicationCommandOptionMentionable: applicationCommandOptionMentionable(o),
	})
}

func (_ ApplicationCommandOptionMentionable) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

func (_ ApplicationCommandOptionMentionable) applicationCommandOption() {}

var _ ApplicationCommandOption = (*ApplicationCommandOptionFloat)(nil)

type ApplicationCommandOptionFloat struct {
	Name         string                                `json:"name"`
	Description  string                                `json:"description"`
	Required     bool                                  `json:"required,omitempty"`
	Choices      []ApplicationCommandOptionChoiceFloat `json:"choices,omitempty"`
	Autocomplete bool                                  `json:"autocomplete,omitempty"`
	MinValue     *json.NullFloat                       `json:"min_value,omitempty"`
	MaxValue     *json.NullFloat                       `json:"max_value,omitempty"`
}

func (o ApplicationCommandOptionFloat) MarshalJSON() ([]byte, error) {
	type applicationCommandOptionFloat ApplicationCommandOptionFloat
	return json.Marshal(struct {
		Type ApplicationCommandOptionType `json:"type"`
		applicationCommandOptionFloat
	}{
		Type:                          o.Type(),
		applicationCommandOptionFloat: applicationCommandOptionFloat(o),
	})
}

func (_ ApplicationCommandOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}

func (_ ApplicationCommandOptionFloat) applicationCommandOption() {}

type ApplicationCommandOptionChoice interface {
	applicationCommandOptionChoice()
}

var _ ApplicationCommandOptionChoice = (*ApplicationCommandOptionChoiceInt)(nil)

type ApplicationCommandOptionChoiceInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (_ ApplicationCommandOptionChoiceInt) applicationCommandOptionChoice() {}

var _ ApplicationCommandOptionChoice = (*ApplicationCommandOptionChoiceString)(nil)

type ApplicationCommandOptionChoiceString struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (_ ApplicationCommandOptionChoiceString) applicationCommandOptionChoice() {}

var _ ApplicationCommandOptionChoice = (*ApplicationCommandOptionChoiceInt)(nil)

type ApplicationCommandOptionChoiceFloat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (_ ApplicationCommandOptionChoiceFloat) applicationCommandOptionChoice() {}
