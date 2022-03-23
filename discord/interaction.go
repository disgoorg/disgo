package discord

import (
	"fmt"

	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/json"
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

type rawInteraction struct {
	ID            snowflake.Snowflake  `json:"id"`
	Type          InteractionType      `json:"type"`
	ApplicationID snowflake.Snowflake  `json:"application_id"`
	Token         string               `json:"token"`
	Version       int                  `json:"version"`
	GuildID       *snowflake.Snowflake `json:"guild_id,omitempty"`
	ChannelID     snowflake.Snowflake  `json:"channel_id,omitempty"`
	Locale        Locale               `json:"locale,omitempty"`
	GuildLocale   *Locale              `json:"guild_locale,omitempty"`
	Member        *ResolvedMember      `json:"member,omitempty"`
	User          *User                `json:"user,omitempty"`
}

// Interaction is used for easier unmarshalling of different Interaction(s)
type Interaction interface {
	Type() InteractionType
	BaseInteraction

	interaction()
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
		vInteraction Interaction
		err          error
	)

	switch iType.Type {
	case InteractionTypePing:
		v := PingInteraction{}
		err = json.Unmarshal(data, &v)
		vInteraction = v

	case InteractionTypeApplicationCommand:
		v := ApplicationCommandInteraction{}
		err = json.Unmarshal(data, &v)
		vInteraction = v

	case InteractionTypeComponent:
		v := ComponentInteraction{}
		err = json.Unmarshal(data, &v)
		vInteraction = v

	case InteractionTypeAutocomplete:
		v := AutocompleteInteraction{}
		err = json.Unmarshal(data, &v)
		vInteraction = v

	case InteractionTypeModalSubmit:
		v := ModalSubmitInteraction{}
		err = json.Unmarshal(data, &v)
		vInteraction = v

	default:
		return fmt.Errorf("unkown rawInteraction with type %d received", iType.Type)
	}
	if err != nil {
		return err
	}

	i.Interaction = vInteraction
	return nil
}

type (
	ResolvedMember struct {
		Member
		Permissions Permissions `json:"permissions,omitempty"`
	}
	ResolvedChannel struct {
		ID             snowflake.Snowflake `json:"id"`
		Name           string              `json:"name"`
		Type           ChannelType         `json:"type"`
		Permissions    Permissions         `json:"permissions"`
		ThreadMetadata ThreadMetadata      `json:"thread_metadata"`
		ParentID       snowflake.Snowflake `json:"parent_id"`
	}
)
