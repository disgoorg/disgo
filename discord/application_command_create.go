package discord

import "github.com/disgoorg/disgo/json"

type ApplicationCommandCreate interface {
	json.Marshaler
	Type() ApplicationCommandType
	Name() string
	applicationCommandCreate()
}

type SlashCommandCreate struct {
	CommandName              string                     `json:"name"`
	CommandNameLocalizations map[Locale]string          `json:"name_localizations,omitempty"`
	Description              string                     `json:"description"`
	DescriptionLocalizations map[Locale]string          `json:"description_localizations,omitempty"`
	Options                  []ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions Permissions                `json:"default_member_permissions,omitempty"`
	DMPermission             bool                       `json:"dm_permission"`
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
	CommandName              string            `json:"name"`
	CommandNameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	DefaultMemberPermissions Permissions       `json:"default_member_permissions"`
	DMPermission             bool              `json:"dm_permission"`
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
	CommandName              string            `json:"name"`
	CommandNameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	DefaultMemberPermissions Permissions       `json:"default_member_permissions"`
	DMPermission             bool              `json:"dm_permission"`
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
