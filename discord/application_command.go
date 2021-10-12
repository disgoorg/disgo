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
}

type unmarshalApplicationCommand struct {
	ApplicationCommand
}

func (u *unmarshalApplicationCommand) UnmarshalJSON(data []byte) error {
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
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommand
	}{
		Type:               c.Type(),
		ApplicationCommand: c,
	}
	return json.Marshal(v)
}

func (_ SlashCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

type UserCommand struct {
	ID                Snowflake  `json:"id"`
	ApplicationID     Snowflake  `json:"application_id"`
	GuildID           *Snowflake `json:"guild_id,omitempty"`
	Name              string     `json:"name"`
	DefaultPermission bool       `json:"default_permission,omitempty"`
	Version           Snowflake  `json:"version"`
}

func (c UserCommand) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommand
	}{
		Type:               c.Type(),
		ApplicationCommand: c,
	}
	return json.Marshal(v)
}

func (c UserCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

type MessageCommand struct {
	ID                Snowflake  `json:"id"`
	ApplicationID     Snowflake  `json:"application_id"`
	GuildID           *Snowflake `json:"guild_id,omitempty"`
	Name              string     `json:"name"`
	DefaultPermission bool       `json:"default_permission,omitempty"`
	Version           Snowflake  `json:"version"`
}

func (c MessageCommand) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommand
	}{
		Type:               c.Type(),
		ApplicationCommand: c,
	}
	return json.Marshal(v)
}

func (_ MessageCommand) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}
