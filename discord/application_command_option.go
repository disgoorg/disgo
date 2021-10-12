package discord

import "github.com/DisgoOrg/disgo/json"

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
}

var _ ApplicationCommandOption = (*ApplicationCommandOptionSubCommand)(nil)

type ApplicationCommandOptionSubCommand struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Options     []ApplicationCommandOption `json:"options,omitempty"`
}

func (c ApplicationCommandOptionSubCommand) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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
}

func (c ApplicationCommandOptionInt) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionBool) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionRole) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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

func (c ApplicationCommandOptionMentionable) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
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
}

func (c ApplicationCommandOptionFloat) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandOptionType `json:"type"`
		ApplicationCommandOption
	}{
		Type:                     c.Type(),
		ApplicationCommandOption: c,
	}
	return json.Marshal(v)
}

func (_ ApplicationCommandOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}

type choiceType int

const (
	choiceTypeInt = iota
	choiceTypeString
	choiceTypeFloat
)

type ApplicationCommandOptionChoice interface {
	choiceType() choiceType
}

type ApplicationCommandOptionChoiceInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (_ ApplicationCommandOptionChoiceInt) choiceType() choiceType {
	return choiceTypeInt
}

type ApplicationCommandOptionChoiceString struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (_ ApplicationCommandOptionChoiceString) choiceType() choiceType {
	return choiceTypeString
}

type ApplicationCommandOptionChoiceFloat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (_ ApplicationCommandOptionChoiceFloat) choiceType() choiceType {
	return choiceTypeFloat
}
