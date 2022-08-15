package discord

import (
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/validate"
)

type ApplicationCommandCreate interface {
	json.Marshaler
	Type() ApplicationCommandType
	CommandName() string
	applicationCommandCreate()
}

type SlashCommandCreate struct {
	Name                     string                     `json:"name"`
	CommandNameLocalizations map[Locale]string          `json:"name_localizations,omitempty"`
	Description              string                     `json:"description"`
	DescriptionLocalizations map[Locale]string          `json:"description_localizations,omitempty"`
	Options                  []ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions Permissions                `json:"default_member_permissions,omitempty"`
	DMPermission             bool                       `json:"dm_permission"`
}

func (c SlashCommandCreate) Validate() (err error) {
	err = validate.Validate(
		validate.New(c.Name, validate.Required[string], validate.StringRange(1, ApplicationCommandNameMaxLength)),
		validate.New(c.Description, validate.Required[string], validate.StringRange(1, ApplicationCommandDescriptionMaxLength)))
	if err != nil {
		return
	}
	options := c.Options
	err = validate.Validate(validate.New(options,
		validate.SliceNoneNil[ApplicationCommandOption],
		validate.SliceMaxLen[ApplicationCommandOption](ApplicationCommandMaxOptions)))
	if err != nil {
		return
	}
	for _, option := range options {
		if err = option.Validate(); err != nil {
			return
		}
	}
	return nil
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
	Name                     string            `json:"name"`
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

func (c UserCommandCreate) CommandName() string {
	return c.Name
}

func (UserCommandCreate) applicationCommandCreate() {}

type MessageCommandCreate struct {
	Name                     string            `json:"name"`
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

func (c MessageCommandCreate) CommandName() string {
	return c.Name
}

func (MessageCommandCreate) applicationCommandCreate() {}

const (
	ApplicationCommandNameMaxLength        = 32
	ApplicationCommandDescriptionMaxLength = 100

	ApplicationCommandMaxOptions = 25
)
