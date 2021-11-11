package discord

import "github.com/DisgoOrg/disgo/json"

type ApplicationCommandCreate interface {
	json.Marshaler
	Type() ApplicationCommandType
	applicationCommandCreate()
}

type SlashCommandCreate struct {
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
}

func (c SlashCommandCreate) MarshalJSON() ([]byte, error) {
	type slashCommandCreate SlashCommandCreate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		slashCommandCreate
	}{
		Type:               c.Type(),
		slashCommandCreate: slashCommandCreate(c),
	})
}

func (_ SlashCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (_ SlashCommandCreate) applicationCommandCreate() {}

type UserCommandCreate struct {
	Name              string `json:"name"`
	DefaultPermission bool   `json:"default_permission,omitempty"`
}

func (c UserCommandCreate) MarshalJSON() ([]byte, error) {
	type userCommandCreate UserCommandCreate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		userCommandCreate
	}{
		Type:              c.Type(),
		userCommandCreate: userCommandCreate(c),
	})
}

func (_ UserCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (_ UserCommandCreate) applicationCommandCreate() {}

type MessageCommandCreate struct {
	Name              string `json:"name"`
	DefaultPermission bool   `json:"default_permission,omitempty"`
}

func (c MessageCommandCreate) MarshalJSON() ([]byte, error) {
	type messageCommandCreate MessageCommandCreate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		messageCommandCreate
	}{
		Type:                 c.Type(),
		messageCommandCreate: messageCommandCreate(c),
	})
}

func (_ MessageCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (_ MessageCommandCreate) applicationCommandCreate() {}
