package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
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
	applicationCommand()
}

type UnmarshalApplicationCommand struct {
	ApplicationCommand
}

func (u *UnmarshalApplicationCommand) UnmarshalJSON(data []byte) error {
	var aType struct {
		Type ApplicationCommandType `json:"type"`
	}

	if err := json.Unmarshal(data, &aType); err != nil {
		return err
	}

	var (
		applicationCommand ApplicationCommand
		err                error
	)

	switch aType.Type {
	case ApplicationCommandTypeSlash:
		v := SlashCommand{}
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypeUser:
		v := UserCommand{}
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	case ApplicationCommandTypeMessage:
		v := MessageCommand{}
		err = json.Unmarshal(data, &v)
		applicationCommand = v

	default:
		return fmt.Errorf("unkown application command with type %d received", aType.Type)
	}
	if err != nil {
		return err
	}

	u.ApplicationCommand = applicationCommand
	return nil
}

type SlashCommand struct {
	ID                Snowflake                  `json:"id"`
	ApplicationID     Snowflake                  `json:"application_id"`
	GuildID           *Snowflake                 `json:"guild_id,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description,omitempty"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
	Version           Snowflake                  `json:"version"`
}

func (c SlashCommand) MarshalJSON() ([]byte, error) {
	type slashCommand SlashCommand
	v := struct {
		Type ApplicationCommandType `json:"type"`
		slashCommand
	}{
		Type:         c.Type(),
		slashCommand: slashCommand(c),
	}
	return json.Marshal(v)
}

func (c *SlashCommand) UnmarshalJSON(data []byte) error {
	type slashCommand SlashCommand
	var sc struct {
		slashCommand
		Options []UnmarshalApplicationCommandOption `json:"options,omitempty"`
	}

	if err := json.Unmarshal(data, &sc); err != nil {
		return err
	}

	if len(sc.Options) > 0 {
		c.Options = make([]ApplicationCommandOption, len(sc.Options))
		for i := range sc.Options {
			c.Options[i] = sc.Options[i]
		}
	}

	*c = SlashCommand(sc.slashCommand)

	return nil
}

func (_ SlashCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (_ SlashCommand) applicationCommand() {}

type UserCommand struct {
	ID                Snowflake  `json:"id"`
	ApplicationID     Snowflake  `json:"application_id"`
	GuildID           *Snowflake `json:"guild_id,omitempty"`
	Name              string     `json:"name"`
	DefaultPermission bool       `json:"default_permission,omitempty"`
	Version           Snowflake  `json:"version"`
}

func (c UserCommand) MarshalJSON() ([]byte, error) {
	type userCommand UserCommand
	v := struct {
		Type ApplicationCommandType `json:"type"`
		userCommand
	}{
		Type:        c.Type(),
		userCommand: userCommand(c),
	}
	return json.Marshal(v)
}

func (c UserCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (_ UserCommand) applicationCommand() {}

type MessageCommand struct {
	ID                Snowflake  `json:"id"`
	ApplicationID     Snowflake  `json:"application_id"`
	GuildID           *Snowflake `json:"guild_id,omitempty"`
	Name              string     `json:"name"`
	DefaultPermission bool       `json:"default_permission,omitempty"`
	Version           Snowflake  `json:"version"`
}

func (c MessageCommand) MarshalJSON() ([]byte, error) {
	type messageCommand MessageCommand
	v := struct {
		Type ApplicationCommandType `json:"type"`
		messageCommand
	}{
		Type:           c.Type(),
		messageCommand: messageCommand(c),
	}
	return json.Marshal(v)
}

func (_ MessageCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (_ MessageCommand) applicationCommand() {}
