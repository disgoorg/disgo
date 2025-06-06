package discord

import (
	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/omit"
)

type ApplicationCommandUpdate interface {
	json.Marshaler
	Type() ApplicationCommandType
	CommandName() *string
	applicationCommandUpdate()
}

type SlashCommandUpdate struct {
	Name                     *string                       `json:"name,omitempty"`
	NameLocalizations        *map[Locale]string            `json:"name_localizations,omitempty"`
	Description              *string                       `json:"description,omitempty"`
	DescriptionLocalizations *map[Locale]string            `json:"description_localizations,omitempty"`
	Options                  *[]ApplicationCommandOption   `json:"options,omitempty"`
	DefaultMemberPermissions omit.Omit[*Permissions]       `json:"default_member_permissions,omitzero"`
	IntegrationTypes         *[]ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 *[]InteractionContextType     `json:"contexts,omitempty"`
	NSFW                     *bool                         `json:"nsfw,omitempty"`
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

func (c SlashCommandUpdate) CommandName() *string {
	return c.Name
}

func (SlashCommandUpdate) applicationCommandUpdate() {}

type UserCommandUpdate struct {
	Name                     *string                       `json:"name,omitempty"`
	NameLocalizations        *map[Locale]string            `json:"name_localizations,omitempty"`
	DefaultMemberPermissions omit.Omit[*Permissions]       `json:"default_member_permissions,omitzero"`
	IntegrationTypes         *[]ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 *[]InteractionContextType     `json:"contexts,omitempty"`
	NSFW                     *bool                         `json:"nsfw,omitempty"`
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

func (c UserCommandUpdate) CommandName() *string {
	return c.Name
}

func (UserCommandUpdate) applicationCommandUpdate() {}

type MessageCommandUpdate struct {
	Name                     *string                       `json:"name,omitempty"`
	NameLocalizations        *map[Locale]string            `json:"name_localizations,omitempty"`
	DefaultMemberPermissions omit.Omit[*Permissions]       `json:"default_member_permissions,omitzero"`
	IntegrationTypes         *[]ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 *[]InteractionContextType     `json:"contexts,omitempty"`
	NSFW                     *bool                         `json:"nsfw,omitempty"`
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

func (c MessageCommandUpdate) CommandName() *string {
	return c.Name
}

func (MessageCommandUpdate) applicationCommandUpdate() {}

type EntryPointCommandUpdate struct {
	Name                     *string                       `json:"name,omitempty"`
	NameLocalizations        *map[Locale]string            `json:"name_localizations,omitempty"`
	DefaultMemberPermissions omit.Omit[*Permissions]       `json:"default_member_permissions,omitzero"`
	IntegrationTypes         *[]ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 *[]InteractionContextType     `json:"contexts,omitempty"`
	NSFW                     *bool                         `json:"nsfw,omitempty"`
	Handler                  *EntryPointCommandHandlerType `json:"handler,omitempty"`
}

func (c EntryPointCommandUpdate) MarshalJSON() ([]byte, error) {
	type entryPointCommandUpdate EntryPointCommandUpdate
	return json.Marshal(struct {
		Type ApplicationCommandType `json:"type"`
		entryPointCommandUpdate
	}{
		Type:                    c.Type(),
		entryPointCommandUpdate: entryPointCommandUpdate(c),
	})
}

func (EntryPointCommandUpdate) Type() ApplicationCommandType {
	return ApplicationCommandTypePrimaryEntryPoint
}

func (c EntryPointCommandUpdate) CommandName() *string {
	return c.Name
}

func (EntryPointCommandUpdate) applicationCommandUpdate() {}
