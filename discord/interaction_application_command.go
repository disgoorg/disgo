package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Interaction = (*ApplicationCommandInteraction)(nil)
)

type ApplicationCommandInteraction struct {
	baseInteractionImpl
	Data ApplicationCommandInteractionData `json:"data"`
}

func (i *ApplicationCommandInteraction) UnmarshalJSON(data []byte) error {
	type applicationCommandInteraction ApplicationCommandInteraction
	var vInteraction struct {
		Data json.RawMessage `json:"data"`
		applicationCommandInteraction
	}

	if err := json.Unmarshal(data, &vInteraction); err != nil {
		return err
	}
	var cType struct {
		Type ApplicationCommandType `json:"type"`
	}
	if err := json.Unmarshal(vInteraction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ApplicationCommandInteractionData
		err             error
	)

	switch cType.Type {
	case ApplicationCommandTypeSlash:
		v := SlashCommandInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeUser:
		v := UserCommandInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeMessage:
		v := MessageCommandInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	default:
		return fmt.Errorf("unkown application rawInteraction data with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	*i = ApplicationCommandInteraction(vInteraction.applicationCommandInteraction)

	i.Data = interactionData
	return nil
}

func (ApplicationCommandInteraction) Type() InteractionType {
	return InteractionTypeApplicationCommand
}

func (ApplicationCommandInteraction) interaction() {}

type ApplicationCommandInteractionData interface {
	Type() ApplicationCommandType
	ID() snowflake.Snowflake
	Name() string

	applicationCommandInteractionData()
}

type rawSlashCommandInteractionData struct {
	ID       snowflake.Snowflake  `json:"id"`
	Name     string               `json:"name"`
	Resolved SlashCommandResolved `json:"resolved"`
	Options  []SlashCommandOption `json:"options"`
}

func (d *rawSlashCommandInteractionData) UnmarshalJSON(data []byte) error {
	type alias rawSlashCommandInteractionData
	var v struct {
		Options []UnmarshalSlashCommandOption `json:"options"`
		alias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*d = rawSlashCommandInteractionData(v.alias)
	d.Options = make([]SlashCommandOption, len(v.Options))
	for i := range v.Options {
		d.Options[i] = v.Options[i].SlashCommandOption
	}

	return nil
}

type SlashCommandInteractionData struct {
	id                  snowflake.Snowflake
	name                string
	SubCommandName      *string
	SubCommandGroupName *string
	Resolved            SlashCommandResolved
	Options             map[string]SlashCommandOption
}

func (d *SlashCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawSlashCommandInteractionData

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.Resolved = iData.Resolved

	d.Options = make(map[string]SlashCommandOption)
	if len(iData.Options) > 0 {
		flattenedOptions := iData.Options

		unmarshalOption := flattenedOptions[0]
		if option, ok := unmarshalOption.(SlashCommandOptionSubCommandGroup); ok {
			d.SubCommandGroupName = &option.OptionName
			flattenedOptions = make([]SlashCommandOption, len(option.Options))
			for ii := range option.Options {
				flattenedOptions[ii] = option.Options[ii]
			}
			unmarshalOption = option.Options[0]
		}
		if option, ok := unmarshalOption.(SlashCommandOptionSubCommand); ok {
			d.SubCommandName = &option.OptionName
			flattenedOptions = option.Options
		}

		for _, option := range flattenedOptions {
			d.Options[option.Name()] = option
		}
	}
	return nil
}

func (d *SlashCommandInteractionData) MarshalJSON() ([]byte, error) {
	options := make([]SlashCommandOption, len(d.Options))
	i := 0
	for _, option := range d.Options {
		options[i] = option
		i++
	}

	if d.SubCommandName != nil {
		options = []SlashCommandOption{
			SlashCommandOptionSubCommand{
				OptionName: *d.SubCommandName,
				Options:    options,
			},
		}
	}
	if d.SubCommandGroupName != nil {
		subCommandOptions := make([]SlashCommandOptionSubCommand, len(options))
		for ii := range options {
			subCommandOptions[ii] = options[ii].(SlashCommandOptionSubCommand)
		}
		options = []SlashCommandOption{
			SlashCommandOptionSubCommandGroup{
				OptionName: *d.SubCommandGroupName,
				Options:    subCommandOptions,
			},
		}
	}

	return json.Marshal(rawSlashCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Resolved: d.Resolved,
		Options:  options,
	})
}

func (SlashCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (d SlashCommandInteractionData) ID() snowflake.Snowflake {
	return d.id
}

func (d SlashCommandInteractionData) Name() string {
	return d.name
}

func (SlashCommandInteractionData) applicationCommandInteractionData() {}

type SlashCommandResolved struct {
	Users    map[snowflake.Snowflake]User            `json:"users,omitempty"`
	Members  map[snowflake.Snowflake]ResolvedMember  `json:"members,omitempty"`
	Roles    map[snowflake.Snowflake]Role            `json:"roles,omitempty"`
	Channels map[snowflake.Snowflake]ResolvedChannel `json:"channels,omitempty"`
}

type ContextCommandInteractionData interface {
	ApplicationCommandInteractionData
	TargetID() snowflake.Snowflake

	contextCommandInteractionData()
}

var (
	_ ApplicationCommandInteractionData = (*UserCommandInteractionData)(nil)
	_ ContextCommandInteractionData     = (*UserCommandInteractionData)(nil)
)

type rawUserCommandInteractionData struct {
	ID       snowflake.Snowflake `json:"id"`
	Name     string              `json:"name"`
	Resolved UserCommandResolved `json:"resolved"`
	TargetID snowflake.Snowflake `json:"target_id"`
}

type UserCommandInteractionData struct {
	id       snowflake.Snowflake
	name     string
	Resolved UserCommandResolved `json:"resolved"`
	targetID snowflake.Snowflake
}

func (d *UserCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawUserCommandInteractionData
	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.Resolved = iData.Resolved
	d.targetID = iData.TargetID
	return nil
}

func (d *UserCommandInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawUserCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Resolved: d.Resolved,
		TargetID: d.targetID,
	})
}

func (UserCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}
func (d UserCommandInteractionData) ID() snowflake.Snowflake {
	return d.id
}
func (d UserCommandInteractionData) Name() string {
	return d.name
}
func (d UserCommandInteractionData) TargetID() snowflake.Snowflake {
	return d.targetID
}
func (d UserCommandInteractionData) TargetUser() User {
	return d.Resolved.Users[d.targetID]
}
func (d UserCommandInteractionData) TargetMember() ResolvedMember {
	return d.Resolved.Members[d.targetID]
}

func (UserCommandInteractionData) applicationCommandInteractionData() {}
func (UserCommandInteractionData) contextCommandInteractionData()     {}

type UserCommandResolved struct {
	Users   map[snowflake.Snowflake]User           `json:"users,omitempty"`
	Members map[snowflake.Snowflake]ResolvedMember `json:"members,omitempty"`
}

var (
	_ ApplicationCommandInteractionData = (*MessageCommandInteractionData)(nil)
	_ ContextCommandInteractionData     = (*MessageCommandInteractionData)(nil)
)

type rawMessageCommandInteractionData struct {
	ID       snowflake.Snowflake    `json:"id"`
	Name     string                 `json:"name"`
	Resolved MessageCommandResolved `json:"resolved"`
	TargetID snowflake.Snowflake    `json:"target_id"`
}

type MessageCommandInteractionData struct {
	id       snowflake.Snowflake
	name     string
	Resolved MessageCommandResolved `json:"resolved"`
	targetID snowflake.Snowflake
}

func (d *MessageCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawMessageCommandInteractionData
	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.Resolved = iData.Resolved
	d.targetID = iData.TargetID
	return nil
}

func (d *MessageCommandInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawMessageCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Resolved: d.Resolved,
		TargetID: d.targetID,
	})
}

func (MessageCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}
func (d MessageCommandInteractionData) ID() snowflake.Snowflake {
	return d.id
}
func (d MessageCommandInteractionData) Name() string {
	return d.name
}
func (d MessageCommandInteractionData) TargetID() snowflake.Snowflake {
	return d.targetID
}
func (d MessageCommandInteractionData) TargetMessage() Message {
	return d.Resolved.Messages[d.targetID]
}

func (MessageCommandInteractionData) applicationCommandInteractionData() {}
func (MessageCommandInteractionData) contextCommandInteractionData()     {}

type MessageCommandResolved struct {
	Messages map[snowflake.Snowflake]Message `json:"messages,omitempty"`
}
