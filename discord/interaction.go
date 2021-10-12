package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

// InteractionType is the type of Interaction
type InteractionType int

// Supported InteractionType(s)
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
	InteractionTypeComponent
	InteractionTypeAutocomplete
)

// Interaction is used for easier unmarshalling of different Interaction(s)
type Interaction interface {
	Type() InteractionType
}

type UnmarshalInteraction struct {
	Interaction
}

func (i *UnmarshalInteraction) UnmarshalJSON(data []byte) error {
	var iType struct {
		Type InteractionType `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(data, &iType); err != nil {
		return err
	}

	var (
		interaction Interaction
		err         error
	)

	switch iType.Type {
	case InteractionTypePing:
		v := PingInteraction{}
		err = json.Unmarshal(data, &v)
		interaction = v

	case InteractionTypeApplicationCommand:
		var cType struct {
			Type ApplicationCommandType `json:"type"`
		}

		if err = json.Unmarshal(iType.Data, &cType); err != nil {
			return err
		}

		switch cType.Type {
		case ApplicationCommandTypeSlash:
			v := SlashCommandInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		case ApplicationCommandTypeUser:
			v := UserCommandInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		case ApplicationCommandTypeMessage:
			v := MessageCommandInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		default:
			return fmt.Errorf("unkown application command interaction with type %d received", iType.Type)
		}

	case InteractionTypeComponent:
		var cType struct {
			Type ComponentType `json:"component_type"`
		}

		if err = json.Unmarshal(iType.Data, &cType); err != nil {
			return err
		}

		switch cType.Type {
		case ComponentTypeButton:
			v := ButtonInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		case ComponentTypeSelectMenu:
			v := SelectMenuInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		default:
			return fmt.Errorf("unkown component interaction with type %d received", iType.Type)
		}

	case InteractionTypeAutocomplete:
		v := AutocompleteInteraction{}
		err = json.Unmarshal(data, &v)
		interaction = v

	default:
		return fmt.Errorf("unkown interaction with type %d received", iType.Type)
	}
	if err != nil {
		return err
	}

	i.Interaction = interaction
	return nil
}

var _ Interaction = (*PingInteraction)(nil)

type PingInteraction struct {
	ID            Snowflake `json:"id"`
	ApplicationID Snowflake `json:"application_id"`
	Token         string    `json:"token"`
	Version       int       `json:"version"`
}

func (_ PingInteraction) Type() InteractionType {
	return InteractionTypePing
}

type ApplicationCommandInteraction interface {
	Interaction
	ApplicationCommandType() ApplicationCommandType
}

var (
	_ Interaction                   = (*SlashCommandInteraction)(nil)
	_ ApplicationCommandInteraction = (*SlashCommandInteraction)(nil)
)

type SlashCommandInteraction struct {
	ID            Snowflake                   `json:"id"`
	ApplicationID Snowflake                   `json:"application_id"`
	Data          SlashCommandInteractionData `json:"data"`
	GuildID       *Snowflake                  `json:"guild_id,omitempty"`
	ChannelID     Snowflake                   `json:"channel_id"`
	Member        *Member                     `json:"member,omitempty"`
	User          *User                       `json:"user,omitempty"`
	Token         string                      `json:"token"`
	Version       int                         `json:"version"`
}

type SlashCommandInteractionData struct {
	ID          Snowflake            `json:"id"`
	CommandName string               `json:"name"`
	Resolved    SlashCommandResolved `json:"resolved"`
	Options     []SlashCommandOption `json:"options"`
}

type SlashCommandResolved struct {
	Users    map[Snowflake]User    `json:"users,omitempty"`
	Members  map[Snowflake]Member  `json:"members,omitempty"`
	Roles    map[Snowflake]Role    `json:"roles,omitempty"`
	Channels map[Snowflake]Channel `json:"channels,omitempty"`
}

func (d *SlashCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData struct {
		SlashCommandInteractionData
		Options []unmarshalSlashCommandOption `json:"options"`
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	if len(iData.Options) > 0 {
		d.Options = make([]SlashCommandOption, len(iData.Options))
		for i, option := range iData.Options {
			d.Options[i] = option.SlashCommandOption
		}
	}

	return nil
}

func (_ SlashCommandInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (_ SlashCommandInteraction) ApplicationCommandType() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

var (
	_ Interaction                   = (*UserCommandInteraction)(nil)
	_ ApplicationCommandInteraction = (*UserCommandInteraction)(nil)
)

type UserCommandInteraction struct {
	ID            Snowflake                  `json:"id"`
	ApplicationID Snowflake                  `json:"application_id"`
	Data          UserCommandInteractionData `json:"data"`
	GuildID       *Snowflake                 `json:"guild_id,omitempty"`
	ChannelID     Snowflake                  `json:"channel_id"`
	Member        *Member                    `json:"member,omitempty"`
	User          *User                      `json:"user,omitempty"`
	Token         string                     `json:"token"`
	Version       int                        `json:"version"`
}

type UserCommandInteractionData struct {
	ID          Snowflake           `json:"id"`
	CommandName string              `json:"name"`
	Resolved    UserCommandResolved `json:"resolved"`
	TargetID    Snowflake           `json:"target_id"`
}

type UserCommandResolved struct {
	Users map[Snowflake]User `json:"users,omitempty"`
}

func (_ UserCommandInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (_ UserCommandInteraction) ApplicationCommandType() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

var (
	_ Interaction                   = (*MessageCommandInteraction)(nil)
	_ ApplicationCommandInteraction = (*MessageCommandInteraction)(nil)
)

type MessageCommandInteraction struct {
	ID            Snowflake                     `json:"id"`
	ApplicationID Snowflake                     `json:"application_id"`
	Data          MessageCommandInteractionData `json:"data"`
	GuildID       *Snowflake                    `json:"guild_id,omitempty"`
	ChannelID     Snowflake                     `json:"channel_id"`
	Member        *Member                       `json:"member,omitempty"`
	User          *User                         `json:"user,omitempty"`
	Token         string                        `json:"token"`
	Version       int                           `json:"version"`
}

type MessageCommandInteractionData struct {
	ID          Snowflake              `json:"id"`
	CommandName string                 `json:"name"`
	Resolved    MessageCommandResolved `json:"resolved"`
	TargetID    Snowflake              `json:"target_id"`
}

type MessageCommandResolved struct {
	Messages map[Snowflake]Message `json:"messages,omitempty"`
}

func (_ MessageCommandInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (_ MessageCommandInteraction) ApplicationCommandType() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

type ComponentInteraction interface {
	Interaction
	ComponentType() ComponentType
}

var (
	_ Interaction          = (*ButtonInteraction)(nil)
	_ ComponentInteraction = (*ButtonInteraction)(nil)
)

type ButtonInteraction struct {
	ID            Snowflake             `json:"id"`
	ApplicationID Snowflake             `json:"application_id"`
	Data          ButtonInteractionData `json:"data"`
	GuildID       *Snowflake            `json:"guild_id,omitempty"`
	ChannelID     Snowflake             `json:"channel_id"`
	Member        *Member               `json:"member,omitempty"`
	User          *User                 `json:"user,omitempty"`
	Token         string                `json:"token"`
	Version       int                   `json:"version"`
	Message       Message               `json:"message"`
}

type ButtonInteractionData struct {
	CustomID string `json:"custom_id"`
}

func (_ ButtonInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (_ ButtonInteraction) ComponentType() ComponentType {
	return ComponentTypeButton
}

var (
	_ Interaction          = (*SelectMenuInteraction)(nil)
	_ ComponentInteraction = (*SelectMenuInteraction)(nil)
)

type SelectMenuInteraction struct {
	ID            Snowflake                 `json:"id"`
	ApplicationID Snowflake                 `json:"application_id"`
	Data          SelectMenuInteractionData `json:"data"`
	GuildID       *Snowflake                `json:"guild_id,omitempty"`
	ChannelID     Snowflake                 `json:"channel_id"`
	Member        *Member                   `json:"member,omitempty"`
	User          *User                     `json:"user,omitempty"`
	Token         string                    `json:"token"`
	Version       int                       `json:"version"`
	Message       Message                   `json:"message"`
}

type SelectMenuInteractionData struct {
	CustomID string `json:"custom_id"`
	Values   string `json:"values"`
}

func (_ SelectMenuInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (_ SelectMenuInteraction) ComponentType() ComponentType {
	return ComponentTypeSelectMenu
}

var (
	_ Interaction = (*AutocompleteInteraction)(nil)
)

type AutocompleteInteraction struct {
	ID            Snowflake                   `json:"id"`
	ApplicationID Snowflake                   `json:"application_id"`
	Data          AutocompleteInteractionData `json:"data"`
	GuildID       *Snowflake                  `json:"guild_id,omitempty"`
	ChannelID     Snowflake                   `json:"channel_id"`
	Member        *Member                     `json:"member,omitempty"`
	User          *User                       `json:"user,omitempty"`
	Token         string                      `json:"token"`
	Version       int                         `json:"version"`
}

func (_ AutocompleteInteraction) Type() InteractionType {
	return InteractionTypeAutocomplete
}

type AutocompleteInteractionData struct {
	ID      Snowflake            `json:"id"`
	Name    string               `json:"name"`
	Options []AutocompleteOption `json:"options"`
}

func (d *AutocompleteInteractionData) UnmarshalJSON(data []byte) error {
	var iData struct {
		AutocompleteInteractionData
		Options []unmarshalAutocompleteOption `json:"options"`
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	if len(iData.Options) > 0 {
		d.Options = make([]AutocompleteOption, len(iData.Options))
		for i, option := range iData.Options {
			d.Options[i] = option.AutocompleteOption
		}
	}

	return nil
}

// to consider using them in Resolved
/*
type ResolvedMember struct {
	GuildID      Snowflake   `json:"guild_id"`
	User         User        `json:"user"`
	Nick         *string     `json:"nick"`
	RoleIDs      []Snowflake `json:"roles,omitempty"`
	JoinedAt     Time        `json:"joined_at"`
	PremiumSince *Time       `json:"premium_since,omitempty"`
	Permissions  Permissions `json:"permissions,omitempty"`
}

type ResolvedChannel struct {
	ID          Snowflake   `json:"id"`
	Name        string      `json:"name"`
	Type        ChannelType `json:"type"`
	Permissions Permissions `json:"permissions"`
}*/
