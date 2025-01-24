package discord

import "github.com/disgoorg/json"

type ApplicationCommandCreate interface {
	json.Marshaler
	Type() ApplicationCommandType
	CommandName() string
	applicationCommandCreate()
}

type SlashCommandCreate struct {
	Name                     string                      `json:"name"`
	NameLocalizations        map[Locale]string           `json:"name_localizations,omitempty"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[Locale]string           `json:"description_localizations,omitempty"`
	Options                  []ApplicationCommandOption  `json:"options,omitempty"`
	DefaultMemberPermissions *json.Nullable[Permissions] `json:"default_member_permissions,omitempty"` // different behavior for 0 and null, optional
	// Deprecated: Use Contexts instead
	DMPermission     *bool                        `json:"dm_permission,omitempty"`
	IntegrationTypes []ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts         []InteractionContextType     `json:"contexts,omitempty"`
	NSFW             *bool                        `json:"nsfw,omitempty"`
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

func (c SlashCommandCreate) CommandName() string {
	return c.Name
}

func (SlashCommandCreate) applicationCommandCreate() {}

type UserCommandCreate struct {
	Name                     string                      `json:"name"`
	NameLocalizations        map[Locale]string           `json:"name_localizations,omitempty"`
	DefaultMemberPermissions *json.Nullable[Permissions] `json:"default_member_permissions,omitempty"`
	// Deprecated: Use Contexts instead
	DMPermission     *bool                        `json:"dm_permission,omitempty"`
	IntegrationTypes []ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts         []InteractionContextType     `json:"contexts,omitempty"`
	NSFW             *bool                        `json:"nsfw,omitempty"`
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

func (c UserCommandCreate) CommandName() string {
	return c.Name
}

func (UserCommandCreate) applicationCommandCreate() {}

type MessageCommandCreate struct {
	Name                     string                      `json:"name"`
	NameLocalizations        map[Locale]string           `json:"name_localizations,omitempty"`
	DefaultMemberPermissions *json.Nullable[Permissions] `json:"default_member_permissions,omitempty"`
	// Deprecated: Use Contexts instead
	DMPermission     *bool                        `json:"dm_permission,omitempty"`
	IntegrationTypes []ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts         []InteractionContextType     `json:"contexts,omitempty"`
	NSFW             *bool                        `json:"nsfw,omitempty"`
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

func (c MessageCommandCreate) CommandName() string {
	return c.Name
}

func (MessageCommandCreate) applicationCommandCreate() {}

type EntryPointCommandCreate struct {
	Name                     string                      `json:"name"`
	NameLocalizations        map[Locale]string           `json:"name_localizations,omitempty"`
	DefaultMemberPermissions *json.Nullable[Permissions] `json:"default_member_permissions,omitempty"`
	// Deprecated: Use Contexts instead
	DMPermission     *bool                        `json:"dm_permission,omitempty"`
	IntegrationTypes []ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts         []InteractionContextType     `json:"contexts,omitempty"`
	NSFW             *bool                        `json:"nsfw,omitempty"`
	Handler          EntryPointCommandHandlerType `json:"handler,omitempty"`
}

func (c EntryPointCommandCreate) MarshalJSON() ([]byte, error) {
	type entryPointCommandCreate EntryPointCommandCreate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		entryPointCommandCreate
	}{
		Type:                    c.Type(),
		entryPointCommandCreate: entryPointCommandCreate(c),
	})
}

func (EntryPointCommandCreate) Type() ApplicationCommandType {
	return ApplicationCommandTypePrimaryEntryPoint
}

func (c EntryPointCommandCreate) CommandName() string {
	return c.Name
}

func (EntryPointCommandCreate) applicationCommandCreate() {}
