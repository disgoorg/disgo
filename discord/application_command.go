package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

type ApplicationCommandType int

//goland:noinspection GoUnusedConst
const (
	ApplicationCommandTypeSlash = iota + 1
	ApplicationCommandTypeUser
	ApplicationCommandTypeMessage
)

type ApplicationCommand interface {
	json.Marshaler
	Type() ApplicationCommandType
	ID() snowflake.Snowflake
	Name() string
	applicationCommand()
}

type UnmarshalApplicationCommand struct {
	ApplicationCommand
}

func (u *UnmarshalApplicationCommand) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ApplicationCommandType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	var (
		applicationCommand ApplicationCommand
		err                error
	)

	switch cType.Type {
	case ApplicationCommandTypeSlash:
		var v SlashCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypeUser:
		var v UserCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypeMessage:
		var v MessageCommand
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	default:
		err = fmt.Errorf("unkown application command with type %d received", cType.Type)
	}

	if err != nil {
		return err
	}

	u.ApplicationCommand = applicationCommand
	return nil
}

var _ ApplicationCommand = (*SlashCommand)(nil)

type SlashCommand struct {
	CommandID         snowflake.Snowflake        `json:"id"`
	ApplicationID     snowflake.Snowflake        `json:"application_id"`
	GuildID           *snowflake.Snowflake       `json:"guild_id,omitempty"`
	CommandName       string                     `json:"name"`
	Description       string                     `json:"description,omitempty"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
	Version           snowflake.Snowflake        `json:"version"`
}

func (c SlashCommand) MarshalJSON() ([]byte, error) {
	type slashCommand SlashCommand
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		slashCommand
	}{
		Type:         c.Type(),
		slashCommand: slashCommand(c),
	})
}

func (c *SlashCommand) UnmarshalJSON(data []byte) error {
	type slashCommand SlashCommand
	var sc struct {
		Options []UnmarshalApplicationCommandOption `json:"options,omitempty"`
		slashCommand
	}

	if err := json.Unmarshal(data, &sc); err != nil {
		return err
	}

	*c = SlashCommand(sc.slashCommand)

	if len(sc.Options) > 0 {
		c.Options = make([]ApplicationCommandOption, len(sc.Options))
		for i := range sc.Options {
			c.Options[i] = sc.Options[i].ApplicationCommandOption
		}
	}
	return nil
}

func (c SlashCommand) ID() snowflake.Snowflake {
	return c.CommandID
}

func (SlashCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (c SlashCommand) Name() string {
	return c.CommandName
}

func (SlashCommand) applicationCommand() {}

var _ ApplicationCommand = (*UserCommand)(nil)

type UserCommand struct {
	CommandID         snowflake.Snowflake  `json:"id"`
	ApplicationID     snowflake.Snowflake  `json:"application_id"`
	GuildID           *snowflake.Snowflake `json:"guild_id,omitempty"`
	CommandName       string               `json:"name"`
	DefaultPermission bool                 `json:"default_permission,omitempty"`
	Version           snowflake.Snowflake  `json:"version"`
}

func (c UserCommand) MarshalJSON() ([]byte, error) {
	type userCommand UserCommand
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		userCommand
	}{
		Type:        c.Type(),
		userCommand: userCommand(c),
	})
}

func (c UserCommand) ID() snowflake.Snowflake {
	return c.CommandID
}

func (c UserCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (c UserCommand) Name() string {
	return c.CommandName
}

func (UserCommand) applicationCommand() {}

var _ ApplicationCommand = (*MessageCommand)(nil)

type MessageCommand struct {
	CommandID         snowflake.Snowflake  `json:"id"`
	ApplicationID     snowflake.Snowflake  `json:"application_id"`
	GuildID           *snowflake.Snowflake `json:"guild_id,omitempty"`
	CommandName       string               `json:"name"`
	DefaultPermission bool                 `json:"default_permission,omitempty"`
	Version           snowflake.Snowflake  `json:"version"`
}

func (c MessageCommand) MarshalJSON() ([]byte, error) {
	type messageCommand MessageCommand
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		messageCommand
	}{
		Type:           c.Type(),
		messageCommand: messageCommand(c),
	})
}

func (c MessageCommand) ID() snowflake.Snowflake {
	return c.CommandID
}

func (MessageCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (c MessageCommand) Name() string {
	return c.CommandName
}

func (MessageCommand) applicationCommand() {}
