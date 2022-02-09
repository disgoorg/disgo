package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

type SlashCommandOption interface {
	Type() ApplicationCommandOptionType
	Name() string
	slashCommandOption()
}

type UnmarshalSlashCommandOption struct {
	SlashCommandOption
}

func (o *UnmarshalSlashCommandOption) UnmarshalJSON(data []byte) error {
	var oType struct {
		Type ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &oType); err != nil {
		return err
	}

	var (
		slashCommandOption SlashCommandOption
		err                error
	)

	switch oType.Type {
	case ApplicationCommandOptionTypeSubCommand:
		v := SlashCommandOptionSubCommand{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeSubCommandGroup:
		v := SlashCommandOptionSubCommandGroup{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeString:
		v := SlashCommandOptionString{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeInt:
		v := SlashCommandOptionInt{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeBool:
		v := SlashCommandOptionBool{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeUser:
		v := SlashCommandOptionUser{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeChannel:
		v := SlashCommandOptionChannel{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeRole:
		v := SlashCommandOptionRole{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeMentionable:
		v := SlashCommandOptionMentionable{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeFloat:
		v := SlashCommandOptionFloat{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	case ApplicationCommandOptionTypeAttachment:
		v := SlashCommandOptionAttachment{}
		err = json.Unmarshal(data, &v)
		slashCommandOption = v

	default:
		return fmt.Errorf("unkown application command option with type %d received", oType.Type)
	}
	if err != nil {
		return err
	}

	o.SlashCommandOption = slashCommandOption
	return nil
}

var _ SlashCommandOption = (*SlashCommandOptionSubCommand)(nil)

type SlashCommandOptionSubCommand struct {
	OptionName  string               `json:"name"`
	Description string               `json:"description"`
	Options     []SlashCommandOption `json:"options,omitempty"`
}

func (o *SlashCommandOptionSubCommand) UnmarshalJSON(data []byte) error {
	type slashCommandOptionSubCommand SlashCommandOptionSubCommand
	var v struct {
		Options []UnmarshalSlashCommandOption `json:"options"`
		slashCommandOptionSubCommand
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*o = SlashCommandOptionSubCommand(v.slashCommandOptionSubCommand)

	if len(v.Options) > 0 {
		o.Options = make([]SlashCommandOption, len(v.Options))
		for i := range v.Options {
			o.Options[i] = v.Options[i].SlashCommandOption
		}
	}

	return nil
}

func (SlashCommandOptionSubCommand) slashCommandOption() {}
func (SlashCommandOptionSubCommand) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

func (o SlashCommandOptionSubCommand) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionSubCommandGroup)(nil)

type SlashCommandOptionSubCommandGroup struct {
	OptionName  string                         `json:"name"`
	Description string                         `json:"description"`
	Options     []SlashCommandOptionSubCommand `json:"options,omitempty"`
}

func (SlashCommandOptionSubCommandGroup) slashCommandOption() {}
func (SlashCommandOptionSubCommandGroup) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

func (o SlashCommandOptionSubCommandGroup) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionString)(nil)

type SlashCommandOptionString struct {
	OptionName string `json:"name"`
	Value      string `json:"value"`
}

func (SlashCommandOptionString) slashCommandOption() {}
func (SlashCommandOptionString) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

func (o SlashCommandOptionString) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionInt)(nil)

type SlashCommandOptionInt struct {
	OptionName string `json:"name"`
	Value      int    `json:"value"`
}

func (SlashCommandOptionInt) slashCommandOption() {}
func (SlashCommandOptionInt) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInt
}

func (o SlashCommandOptionInt) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionBool)(nil)

type SlashCommandOptionBool struct {
	OptionName string `json:"name"`
	Value      bool   `json:"value"`
}

func (SlashCommandOptionBool) slashCommandOption() {}
func (SlashCommandOptionBool) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBool
}

func (o SlashCommandOptionBool) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionUser)(nil)

type SlashCommandOptionUser struct {
	OptionName string              `json:"name"`
	Value      snowflake.Snowflake `json:"value"`
}

func (SlashCommandOptionUser) slashCommandOption() {}
func (SlashCommandOptionUser) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

func (o SlashCommandOptionUser) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionChannel)(nil)

type SlashCommandOptionChannel struct {
	OptionName string              `json:"name"`
	Value      snowflake.Snowflake `json:"value"`
}

func (SlashCommandOptionChannel) slashCommandOption() {}
func (SlashCommandOptionChannel) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

func (o SlashCommandOptionChannel) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionRole)(nil)

type SlashCommandOptionRole struct {
	OptionName string              `json:"name"`
	Value      snowflake.Snowflake `json:"value"`
}

func (SlashCommandOptionRole) slashCommandOption() {}
func (SlashCommandOptionRole) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

func (o SlashCommandOptionRole) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionMentionable)(nil)

type SlashCommandOptionMentionable struct {
	OptionName string              `json:"name"`
	Value      snowflake.Snowflake `json:"value"`
}

func (SlashCommandOptionMentionable) slashCommandOption() {}
func (SlashCommandOptionMentionable) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

func (o SlashCommandOptionMentionable) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionFloat)(nil)

type SlashCommandOptionFloat struct {
	OptionName string  `json:"name"`
	Value      float64 `json:"value"`
}

func (SlashCommandOptionFloat) slashCommandOption() {}
func (SlashCommandOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}

func (o SlashCommandOptionFloat) Name() string {
	return o.OptionName
}

var _ SlashCommandOption = (*SlashCommandOptionAttachment)(nil)

type SlashCommandOptionAttachment struct {
	OptionName string              `json:"name"`
	Value      snowflake.Snowflake `json:"value"`
}

func (SlashCommandOptionAttachment) slashCommandOption() {}
func (SlashCommandOptionAttachment) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeAttachment
}

func (o SlashCommandOptionAttachment) Name() string {
	return o.OptionName
}
