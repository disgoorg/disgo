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
	type slashCommandUpdate SlashCommandUpdate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		slashCommandUpdate
	}{
		Type:               c.Type(),
		slashCommandUpdate: slashCommandUpdate(c),
	})
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
	type userCommandUpdate UserCommandUpdate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		userCommandUpdate
	}{
		Type:              c.Type(),
		userCommandUpdate: userCommandUpdate(c),
	})
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
	type messageCommandUpdate MessageCommandUpdate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		messageCommandUpdate
	}{
		Type:                 c.Type(),
		messageCommandUpdate: messageCommandUpdate(c),
	})
}

func (_ MessageCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (_ MessageCommandUpdate) applicationCommandUpdate() {}
