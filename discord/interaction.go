package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
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
	interaction()
}

type BaseInteraction struct {
	ID            snowflake.Snowflake  `json:"id"`
	ApplicationID snowflake.Snowflake  `json:"application_id"`
	Token         string               `json:"token"`
	Version       int                  `json:"version"`
	GuildID       *snowflake.Snowflake `json:"guild_id,omitempty"`
	ChannelID     snowflake.Snowflake  `json:"channel_id"`
	Locale        Locale               `json:"locale"`
	GuildLocale   *Locale              `json:"guild_locale,omitempty"`
	Member        *Member              `json:"member,omitempty"`
	User          *User                `json:"user,omitempty"`
}

type UnmarshalInteraction struct {
	Interaction
}

func (i *UnmarshalInteraction) UnmarshalJSON(data []byte) error {
	var iType struct {
		Type InteractionType `json:"type"`
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
		v := ApplicationCommandInteraction{}
		err = json.Unmarshal(data, &v)
		interaction = v

	case InteractionTypeComponent:
		v := ComponentInteraction{}
		err = json.Unmarshal(data, &v)
		interaction = v

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
	ID            snowflake.Snowflake `json:"id"`
	ApplicationID snowflake.Snowflake `json:"application_id"`
	Token         string              `json:"token"`
	Version       int                 `json:"version"`
}

func (PingInteraction) interaction() {}
func (PingInteraction) Type() InteractionType {
	return InteractionTypePing
}

var (
	_ Interaction = (*ApplicationCommandInteraction)(nil)
)

type ApplicationCommandInteraction struct {
	BaseInteraction
	Data ApplicationCommandInteractionData `json:"data"`
}

func (ApplicationCommandInteraction) interaction() {}
func (ApplicationCommandInteraction) Type() InteractionType {
	return InteractionTypeApplicationCommand
}

func (i *ApplicationCommandInteraction) UnmarshalJSON(data []byte) error {
	type applicationCommandInteraction ApplicationCommandInteraction
	var interaction struct {
		Data json.RawMessage `json:"data"`
		applicationCommandInteraction
	}

	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}
	var cType struct {
		Type ApplicationCommandType `json:"type"`
	}
	if err := json.Unmarshal(interaction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ApplicationCommandInteractionData
		err             error
	)

	switch cType.Type {
	case ApplicationCommandTypeSlash:
		v := SlashCommandInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeUser:
		v := UserCommandInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeMessage:
		v := MessageCommandInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	default:
		return fmt.Errorf("unkown application interaction data with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	*i = ApplicationCommandInteraction(interaction.applicationCommandInteraction)

	i.Data = interactionData
	return nil
}

type ApplicationCommandInteractionData interface {
	applicationCommandInteractionData()
	Type() ApplicationCommandType
}

type SlashCommandInteractionData struct {
	CommandID   snowflake.Snowflake  `json:"id"`
	CommandName string               `json:"name"`
	Resolved    SlashCommandResolved `json:"resolved"`
	Options     []SlashCommandOption `json:"options"`
}

func (SlashCommandInteractionData) applicationCommandInteractionData() {}
func (SlashCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
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
	Users    map[snowflake.Snowflake]User    `json:"users,omitempty"`
	Members  map[snowflake.Snowflake]Member  `json:"members,omitempty"`
	Roles    map[snowflake.Snowflake]Role    `json:"roles,omitempty"`
	Channels map[snowflake.Snowflake]Channel `json:"channels,omitempty"`
}

var (
	_ ApplicationCommandInteractionData = (*UserCommandInteractionData)(nil)
)

type UserCommandInteractionData struct {
	CommandID   snowflake.Snowflake `json:"id"`
	CommandName string              `json:"name"`
	Resolved    UserCommandResolved `json:"resolved"`
	TargetID    snowflake.Snowflake `json:"target_id"`
}

func (UserCommandInteractionData) applicationCommandInteractionData() {}
func (UserCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

type UserCommandResolved struct {
	Users   map[snowflake.Snowflake]User   `json:"users,omitempty"`
	Members map[snowflake.Snowflake]Member `json:"members,omitempty"`
}

var (
	_ ApplicationCommandInteractionData = (*MessageCommandInteractionData)(nil)
)

type MessageCommandInteractionData struct {
	CommandID   snowflake.Snowflake    `json:"id"`
	CommandName string                 `json:"name"`
	Resolved    MessageCommandResolved `json:"resolved"`
	TargetID    snowflake.Snowflake    `json:"target_id"`
}

func (MessageCommandInteractionData) applicationCommandInteractionData() {}
func (MessageCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

type MessageCommandResolved struct {
	Messages map[snowflake.Snowflake]Message `json:"messages,omitempty"`
}

var (
	_ Interaction = (*ComponentInteraction)(nil)
)

type ComponentInteraction struct {
	BaseInteraction
	Data    ComponentInteractionData `json:"data"`
	Message Message                  `json:"message"`
}

func (ComponentInteraction) interaction() {}
func (ComponentInteraction) Type() InteractionType {
	return InteractionTypeComponent
}

func (i *ComponentInteraction) UnmarshalJSON(data []byte) error {
	type componentInteraction ComponentInteraction
	var interaction struct {
		Data json.RawMessage `json:"data"`
		componentInteraction
	}

	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}

	var cType struct {
		Type ComponentType `json:"component_type"`
	}

	if err := json.Unmarshal(interaction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ComponentInteractionData
		err             error
	)
	switch cType.Type {
	case ComponentTypeButton:
		v := ButtonInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ComponentTypeSelectMenu:
		v := SelectMenuInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	default:
		return fmt.Errorf("unkown component interaction data with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	*i = ComponentInteraction(interaction.componentInteraction)

	i.Data = interactionData
	return nil
}

type ComponentInteractionData interface {
	componentInteractionData()
	Type() ComponentType
}

type ButtonInteractionData struct {
	CustomID CustomID `json:"custom_id"`
}

func (ButtonInteractionData) componentInteractionData() {}
func (ButtonInteractionData) Type() ComponentType {
	return ComponentTypeButton
}

var (
	_ ComponentInteractionData = (*SelectMenuInteractionData)(nil)
)

type SelectMenuInteractionData struct {
	CustomID CustomID `json:"custom_id"`
	Values   []string `json:"values"`
}

func (SelectMenuInteractionData) componentInteractionData() {}
func (SelectMenuInteractionData) Type() ComponentType {
	return ComponentTypeSelectMenu
}

var (
	_ Interaction = (*AutocompleteInteraction)(nil)
)

type AutocompleteInteraction struct {
	BaseInteraction
	Data AutocompleteInteractionData `json:"data"`
}

func (AutocompleteInteraction) interaction() {}
func (AutocompleteInteraction) Type() InteractionType {
	return InteractionTypeAutocomplete
}

type AutocompleteInteractionData struct {
	CommandID   snowflake.Snowflake  `json:"id"`
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
	GuildID      snowflake.Snowflake   `json:"guild_id"`
	User         User        `json:"user"`
	Nick         *string     `json:"nick"`
	RoleIDs      []snowflake.Snowflake `json:"roles,omitempty"`
	JoinedAt     Time        `json:"joined_at"`
	PremiumSince *Time       `json:"premium_since,omitempty"`
	Permissions  Permissions `json:"permissions,omitempty"`
}

type ResolvedChannel struct {
	CommandID          snowflake.Snowflake   `json:"id"`
	CommandName        string      `json:"name"`
	InteractionType        ChannelType `json:"type"`
	Permissions Permissions `json:"permissions"`
}*/
