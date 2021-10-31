package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

type SlashCommandOption interface {
	Type() ApplicationCommandOptionType
}

type unmarshalSlashCommandOption struct {
	SlashCommandOption
}

func (o unmarshalSlashCommandOption) UnmarshalJSON(data []byte) error {
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
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Options     []SlashCommandOption `json:"options,omitempty"`
}

func (_ SlashCommandOptionSubCommand) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

var _ SlashCommandOption = (*SlashCommandOptionSubCommandGroup)(nil)

type SlashCommandOptionSubCommandGroup struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Options     []SlashCommandOptionSubCommand `json:"options,omitempty"`
}

func (_ SlashCommandOptionSubCommandGroup) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

var _ SlashCommandOption = (*SlashCommandOptionString)(nil)

type SlashCommandOptionString struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (_ SlashCommandOptionString) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

var _ SlashCommandOption = (*SlashCommandOptionInt)(nil)

type SlashCommandOptionInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (_ SlashCommandOptionInt) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInt
}

var _ SlashCommandOption = (*SlashCommandOptionBool)(nil)

type SlashCommandOptionBool struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

func (_ SlashCommandOptionBool) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBool
}

var _ SlashCommandOption = (*SlashCommandOptionUser)(nil)

type SlashCommandOptionUser struct {
	Name  string    `json:"name"`
	Value Snowflake `json:"value"`
}

func (_ SlashCommandOptionUser) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

var _ SlashCommandOption = (*SlashCommandOptionChannel)(nil)

type SlashCommandOptionChannel struct {
	Name  string    `json:"name"`
	Value Snowflake `json:"value"`
}

func (_ SlashCommandOptionChannel) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

var _ SlashCommandOption = (*SlashCommandOptionRole)(nil)

type SlashCommandOptionRole struct {
	Name  string    `json:"name"`
	Value Snowflake `json:"value"`
}

func (_ SlashCommandOptionRole) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

var _ SlashCommandOption = (*SlashCommandOptionMentionable)(nil)

type SlashCommandOptionMentionable struct {
	Name  string    `json:"name"`
	Value Snowflake `json:"value"`
}

func (_ SlashCommandOptionMentionable) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

var _ SlashCommandOption = (*SlashCommandOptionFloat)(nil)

type SlashCommandOptionFloat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (_ SlashCommandOptionFloat) Type() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeFloat
}
