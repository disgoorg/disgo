package discord

import "github.com/disgoorg/disgo/json"

type ApplicationCommandUpdate interface {
	json.Marshaler
	Type() ApplicationCommandType
	Name() *string
	applicationCommandUpdate()
}

type SlashCommandUpdate struct {
	CommandName              *string                     `json:"name,omitempty"`
	CommandNameLocalizations *map[Locale]string          `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations *map[Locale]string          `json:"description_localizations,omitempty"`
	Options                  *[]ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *Permissions                `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
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

func (SlashCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (c SlashCommandUpdate) Name() *string {
	return c.CommandName
}

func (SlashCommandUpdate) applicationCommandUpdate() {}

type UserCommandUpdate struct {
	CommandName              *string            `json:"name"`
	CommandNameLocalizations *map[Locale]string `json:"name_localizations,omitempty"`
	DefaultMemberPermissions *Permissions       `json:"default_member_permissions,omitempty"`
	DMPermission             *bool              `json:"dm_permission,omitempty"`
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

func (UserCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (c UserCommandUpdate) Name() *string {
	return c.CommandName
}

func (UserCommandUpdate) applicationCommandUpdate() {}

type MessageCommandUpdate struct {
	CommandName              *string            `json:"name"`
	CommandNameLocalizations *map[Locale]string `json:"name_localizations,omitempty"`
	DefaultMemberPermissions *Permissions       `json:"default_member_permissions,omitempty"`
	DMPermission             *bool              `json:"dm_permission,omitempty"`
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

func (MessageCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (c MessageCommandUpdate) Name() *string {
	return c.CommandName
}

func (MessageCommandUpdate) applicationCommandUpdate() {}
