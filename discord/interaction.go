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
	InteractionTypeModalSubmit
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

	case InteractionTypeModalSubmit:
		v := ModalSubmitInteraction{}
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
	ID() snowflake.Snowflake
	Name() string
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
func (d SlashCommandInteractionData) ID() snowflake.Snowflake {
	return d.CommandID
}
func (d SlashCommandInteractionData) Name() string {
	return d.CommandName
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
	Users    map[snowflake.Snowflake]User            `json:"users,omitempty"`
	Members  map[snowflake.Snowflake]ResolvedMember  `json:"members,omitempty"`
	Roles    map[snowflake.Snowflake]Role            `json:"roles,omitempty"`
	Channels map[snowflake.Snowflake]ResolvedChannel `json:"channels,omitempty"`
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
func (d UserCommandInteractionData) ID() snowflake.Snowflake {
	return d.CommandID
}
func (d UserCommandInteractionData) Name() string {
	return d.CommandName
}

type UserCommandResolved struct {
	Users   map[snowflake.Snowflake]User           `json:"users,omitempty"`
	Members map[snowflake.Snowflake]ResolvedMember `json:"members,omitempty"`
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
func (d MessageCommandInteractionData) ID() snowflake.Snowflake {
	return d.CommandID
}
func (d MessageCommandInteractionData) Name() string {
	return d.CommandName
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
	ID() CustomID
}

type ButtonInteractionData struct {
	CustomID CustomID `json:"custom_id"`
}

func (ButtonInteractionData) componentInteractionData() {}
func (ButtonInteractionData) Type() ComponentType {
	return ComponentTypeButton
}
func (d ButtonInteractionData) ID() CustomID {
	return d.CustomID
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
func (d SelectMenuInteractionData) ID() CustomID {
	return d.CustomID
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
		Options []UnmarshalAutocompleteOption `json:"options"`
		autocompleteInteractionData
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	*d = AutocompleteInteractionData(iData.autocompleteInteractionData)
	if len(iData.Options) > 0 {
		d.Options = make([]AutocompleteOption, len(iData.Options))
		for i := range iData.Options {
			d.Options[i] = iData.Options[i].AutocompleteOption
		}
	}

	return nil
}

var (
	_ Interaction = (*ModalSubmitInteraction)(nil)
)

type ModalSubmitInteraction struct {
	BaseInteraction
	Data ModalSubmitInteractionData `json:"data"`
}

func (ModalSubmitInteraction) interaction() {}

func (ModalSubmitInteraction) Type() InteractionType {
	return InteractionTypeModalSubmit
}

type ModalSubmitInteractionData struct {
	CustomID   CustomID             `json:"custom_id"`
	Components []ContainerComponent `json:"components"`
}

func (d *ModalSubmitInteractionData) UnmarshalJSON(data []byte) error {
	type modalSubmitInteractionData ModalSubmitInteractionData
	var iData struct {
		Components []UnmarshalComponent `json:"components"`
		modalSubmitInteractionData
	}

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}

	*d = ModalSubmitInteractionData(iData.modalSubmitInteractionData)

	if len(iData.Components) > 0 {
		d.Components = make([]ContainerComponent, len(iData.Components))
		for i := range iData.Components {
			d.Components[i] = iData.Components[i].Component.(ContainerComponent)
		}
	}

	return nil
}

// to consider using them in Resolved

type ResolvedMember struct {
	Member
	Permissions Permissions `json:"permissions,omitempty"`
}

type ResolvedChannel struct {
	ID             snowflake.Snowflake `json:"id"`
	Name           string              `json:"name"`
	Type           ChannelType         `json:"type"`
	Permissions    Permissions         `json:"permissions"`
	ThreadMetadata ThreadMetadata      `json:"thread_metadata"`
	ParentID       snowflake.Snowflake `json:"parent_id"`
}
