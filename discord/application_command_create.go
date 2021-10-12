package discord

import "github.com/DisgoOrg/disgo/json"

type ApplicationCommandCreate interface {
	json.Marshaler
	Type() ApplicationCommandType
}

type SlashCommandCreate struct {
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
}

func (c SlashCommandCreate) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommandCreate
	}{
		Type:                     c.Type(),
		ApplicationCommandCreate: c,
	}
	return json.Marshal(v)
}

func (_ SlashCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

type UserCommandCreate struct {
	Name              string `json:"name"`
	DefaultPermission bool   `json:"default_permission,omitempty"`
}

func (c UserCommandCreate) MarshalJSON() ([]byte, error) {
	v := struct {
		Type        ApplicationCommandType `json:"type"`
		Description string                 `json:"description"`
		ApplicationCommandCreate
	}{
		Type:                     c.Type(),
		ApplicationCommandCreate: c,
	}
	return json.Marshal(v)
}

func (_ UserCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

type MessageCommandCreate struct {
	Name              string `json:"name"`
	DefaultPermission bool   `json:"default_permission,omitempty"`
}

func (c MessageCommandCreate) MarshalJSON() ([]byte, error) {
	v := struct {
		Type        ApplicationCommandType `json:"type"`
		Description string                 `json:"description"`
		ApplicationCommandCreate
	}{
		Type:                     c.Type(),
		ApplicationCommandCreate: c,
	}
	return json.Marshal(v)
}

func (_ MessageCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}
