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
	InteractionType() InteractionType
	interaction()
}

type InteractionFields struct {
	ID            Snowflake  `json:"id"`
	ApplicationID Snowflake  `json:"application_id"`
	Token         string     `json:"token"`
	Version       int        `json:"version"`
	GuildID       *Snowflake `json:"guild_id,omitempty"`
	ChannelID     Snowflake  `json:"channel_id"`
	Member        *Member    `json:"member,omitempty"`
	User          *User      `json:"user,omitempty"`
}

type UnmarshalInteraction struct {
	Interaction
}

func (i *UnmarshalInteraction) UnmarshalJSON(data []byte) error {
	var iType struct {
		InteractionType InteractionType `json:"type"`
		Data            struct {
			ApplicationCommandType ApplicationCommandType `json:"type"`
			ComponentType          ComponentType          `json:"component_type"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &iType); err != nil {
		return err
	}

	var (
		interaction Interaction
		err         error
	)

	switch iType.InteractionType {
	case InteractionTypePing:
		v := PingInteraction{}
		err = json.Unmarshal(data, &v)
		interaction = v

	case InteractionTypeApplicationCommand:
		switch iType.Data.ApplicationCommandType {
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
			return fmt.Errorf("unkown application command interaction with type %d received", iType.Data.ApplicationCommandType)
		}

	case InteractionTypeComponent:
		switch iType.Data.ComponentType {
		case ComponentTypeButton:
			v := ButtonInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		case ComponentTypeSelectMenu:
			v := SelectMenuInteraction{}
			err = json.Unmarshal(data, &v)
			interaction = v

		default:
			return fmt.Errorf("unkown component interaction with type %d received", iType.Data.ComponentType)
		}

	case InteractionTypeAutocomplete:
		v := AutocompleteInteraction{}
		err = json.Unmarshal(data, &v)
		interaction = v

	default:
		return fmt.Errorf("unkown interaction with type %d received", iType.InteractionType)
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

func (_ PingInteraction) interaction() {}

func (_ PingInteraction) InteractionType() InteractionType {
	return InteractionTypePing
}

type ApplicationCommandInteraction interface {
	Interaction
	ApplicationCommandType() ApplicationCommandType
	applicationCommandInteraction()
}

var (
	_ Interaction                   = (*SlashCommandInteraction)(nil)
	_ ApplicationCommandInteraction = (*SlashCommandInteraction)(nil)
)

type SlashCommandInteraction struct {
	InteractionFields
	Data SlashCommandInteractionData `json:"data"`
}

func (_ SlashCommandInteraction) interaction() {}
func (_ SlashCommandInteraction) applicationCommandInteraction() {}

type SlashCommandInteractionData struct {
	CommandID   Snowflake            `json:"id"`
	CommandName string               `json:"name"`
	Resolved    SlashCommandResolved `json:"resolved"`
	Options     []SlashCommandOption `json:"options"`
}

func (d *SlashCommandInteractionData) UnmarshalJSON(data []byte) error {
	type slashCommandInteractionData SlashCommandInteractionData
	var iData struct {
		Options []UnmarshalSlashCommandOption `json:"options"`
		slashCommandInteractionData
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	*d = SlashCommandInteractionData(iData.slashCommandInteractionData)

	if len(iData.Options) > 0 {
		d.Options = make([]SlashCommandOption, len(iData.Options))
		for i := range iData.Options {
			d.Options[i] = iData.Options[i].SlashCommandOption
		}
	}

	return nil
}

type SlashCommandResolved struct {
	Users    map[Snowflake]User    `json:"users,omitempty"`
	Members  map[Snowflake]Member  `json:"members,omitempty"`
	Roles    map[Snowflake]Role    `json:"roles,omitempty"`
	Channels map[Snowflake]Channel `json:"channels,omitempty"`
}

func (_ SlashCommandInteraction) InteractionType() InteractionType {
	return InteractionTypeComponent
}

func (_ SlashCommandInteraction) ApplicationCommandType() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

var (
	_ Interaction                   = (*UserCommandInteraction)(nil)
	_ ApplicationCommandInteraction = (*UserCommandInteraction)(nil)
)

func (_ UserCommandInteraction) interaction() {}
func (_ UserCommandInteraction) applicationCommandInteraction() {}

type UserCommandInteraction struct {
	InteractionFields
	Data UserCommandInteractionData `json:"data"`
}

type UserCommandInteractionData struct {
	CommandID   Snowflake           `json:"id"`
	CommandName string              `json:"name"`
	Resolved    UserCommandResolved `json:"resolved"`
	TargetID    Snowflake           `json:"target_id"`
}

type UserCommandResolved struct {
	Users   map[Snowflake]User   `json:"users,omitempty"`
	Members map[Snowflake]Member `json:"members,omitempty"`
}

func (_ UserCommandInteraction) InteractionType() InteractionType {
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
	InteractionFields
	Data MessageCommandInteractionData `json:"data"`
}

func (_ MessageCommandInteraction) interaction() {}
func (_ MessageCommandInteraction) applicationCommandInteraction() {}

type MessageCommandInteractionData struct {
	CommandID   Snowflake              `json:"id"`
	CommandName string                 `json:"name"`
	Resolved    MessageCommandResolved `json:"resolved"`
	TargetID    Snowflake              `json:"target_id"`
}

type MessageCommandResolved struct {
	Messages map[Snowflake]Message `json:"messages,omitempty"`
}

func (_ MessageCommandInteraction) InteractionType() InteractionType {
	return InteractionTypeComponent
}

func (_ MessageCommandInteraction) ApplicationCommandType() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

type ComponentInteraction interface {
	Interaction
	ComponentType() ComponentType
	componentInteraction()
}

var (
	_ Interaction          = (*ButtonInteraction)(nil)
	_ ComponentInteraction = (*ButtonInteraction)(nil)
)

type ButtonInteraction struct {
	InteractionFields
	Data    ButtonInteractionData `json:"data"`
	Message Message               `json:"message"`
}

type ButtonInteractionData struct {
	CustomID CustomID `json:"custom_id"`
}

func (_ ButtonInteraction) interaction() {}
func (_ ButtonInteraction) componentInteraction() {}

func (_ ButtonInteraction) InteractionType() InteractionType {
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
	InteractionFields
	Data    SelectMenuInteractionData `json:"data"`
	Message Message                   `json:"message"`
}

type SelectMenuInteractionData struct {
	CustomID CustomID `json:"custom_id"`
	Values   []string `json:"values"`
}

func (_ SelectMenuInteraction) interaction() {}
func (_ SelectMenuInteraction) componentInteraction() {}

func (_ SelectMenuInteraction) InteractionType() InteractionType {
	return InteractionTypeComponent
}

func (_ SelectMenuInteraction) ComponentType() ComponentType {
	return ComponentTypeSelectMenu
}

var (
	_ Interaction = (*AutocompleteInteraction)(nil)
)

type AutocompleteInteraction struct {
	InteractionFields
	Data AutocompleteInteractionData `json:"data"`
}

func (_ AutocompleteInteraction) interaction() {}

func (_ AutocompleteInteraction) InteractionType() InteractionType {
	return InteractionTypeAutocomplete
}

type AutocompleteInteractionData struct {
	CommandID   Snowflake            `json:"id"`
	CommandName string               `json:"name"`
	Options     []AutocompleteOption `json:"options"`
}

func (d *AutocompleteInteractionData) UnmarshalJSON(data []byte) error {
	type autocompleteInteractionData AutocompleteInteractionData
	var iData struct {
		autocompleteInteractionData
		Options []UnmarshalAutocompleteOption `json:"options"`
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

	*d = AutocompleteInteractionData(iData.autocompleteInteractionData)

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
	CommandID          Snowflake   `json:"id"`
	CommandName        string      `json:"name"`
	InteractionType        ChannelType `json:"type"`
	Permissions Permissions `json:"permissions"`
}*/
