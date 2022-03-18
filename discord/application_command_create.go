package discord

import "github.com/DisgoOrg/disgo/json"

type ApplicationCommandCreate interface {
	json.Marshaler
	Type() ApplicationCommandType
	Name() string
	applicationCommandCreate()
}

type SlashCommandCreate struct {
	CommandName       string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                       `json:"default_permission"`
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

func (SlashCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (c SlashCommandCreate) Name() string {
	return c.CommandName
}

func (SlashCommandCreate) applicationCommandCreate() {}

type UserCommandCreate struct {
	CommandName       string `json:"name"`
	DefaultPermission bool   `json:"default_permission"`
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

func (UserCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (c UserCommandCreate) Name() string {
	return c.CommandName
}

func (UserCommandCreate) applicationCommandCreate() {}

type MessageCommandCreate struct {
	CommandName       string `json:"name"`
	DefaultPermission bool   `json:"default_permission"`
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

func (MessageCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (c MessageCommandCreate) Name() string {
	return c.CommandName
}

func (MessageCommandCreate) applicationCommandCreate() {}
