package discord

import "github.com/DisgoOrg/disgo/json"

type ApplicationCommandUpdate interface {
	json.Marshaler
	Type() ApplicationCommandType
	applicationCommandUpdate()
}

type SlashCommandUpdate struct {
	Name              *string                     `json:"name,omitempty"`
	Description       *string                     `json:"description,omitempty"`
	Options           *[]ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission *bool                       `json:"default_permission,omitempty"`
}

func (c SlashCommandUpdate) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommandUpdate
	}{
		Type:                     c.Type(),
		ApplicationCommandUpdate: c,
	}
	return json.Marshal(v)
}

func (_ SlashCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (_ SlashCommandUpdate) applicationCommandUpdate() {}

type UserCommandUpdate struct {
	Name              *string `json:"name"`
	DefaultPermission *bool   `json:"default_permission,omitempty"`
}

func (c UserCommandUpdate) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommandUpdate
	}{
		Type:                     c.Type(),
		ApplicationCommandUpdate: c,
	}
	return json.Marshal(v)
}

func (_ UserCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (_ UserCommandUpdate) applicationCommandUpdate() {}

type MessageCommandUpdate struct {
	Name              *string `json:"name"`
	DefaultPermission *bool   `json:"default_permission,omitempty"`
}

func (c MessageCommandUpdate) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ApplicationCommandType `json:"type"`
		ApplicationCommandUpdate
	}{
		Type:                     c.Type(),
		ApplicationCommandUpdate: c,
	}
	return json.Marshal(v)
}

func (_ MessageCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (_ MessageCommandUpdate) applicationCommandUpdate() {}
