package discord

import (
	"fmt"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
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
	ID            snowflake.ID    `json:"id"`
	Type          InteractionType `json:"type"`
	ApplicationID snowflake.ID    `json:"application_id"`
	Token         string          `json:"token"`
	Version       int             `json:"version"`
	GuildID       *snowflake.ID   `json:"guild_id,omitempty"`
	// Deprecated: Use Channel instead
	ChannelID      snowflake.ID       `json:"channel_id,omitempty"`
	Channel        InteractionChannel `json:"channel,omitempty"`
	Locale         Locale             `json:"locale,omitempty"`
	GuildLocale    *Locale            `json:"guild_locale,omitempty"`
	Member         *ResolvedMember    `json:"member,omitempty"`
	User           *User              `json:"user,omitempty"`
	AppPermissions *Permissions       `json:"app_permissions,omitempty"`
}

// Interaction is used for easier unmarshalling of different Interaction(s)
type Interaction interface {
	Type() InteractionType
	ID() snowflake.ID
	ApplicationID() snowflake.ID
	Token() string
	Version() int
	GuildID() *snowflake.ID
	// Deprecated: Use Interaction.Channel instead
	ChannelID() snowflake.ID
	Channel() InteractionChannel
	Locale() Locale
	GuildLocale() *Locale
	Member() *ResolvedMember
	User() User
	AppPermissions() *Permissions
	CreatedAt() time.Time

	interaction()
}

func UnmarshalInteraction(data []byte) (Interaction, error) {
	var iType struct {
		Type InteractionType `json:"type"`
	}

	if err := json.Unmarshal(data, &iType); err != nil {
		return nil, err
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
		err = fmt.Errorf("unknown rawInteraction with type %d received", iType.Type)
	}
	if err != nil {
		return nil, err
	}

	return interaction, nil
}

type ResolvedMember struct {
	Member
	Permissions Permissions `json:"permissions,omitempty"`
}

type ResolvedChannel struct {
	ID             snowflake.ID   `json:"id"`
	Name           string         `json:"name"`
	Type           ChannelType    `json:"type"`
	Permissions    Permissions    `json:"permissions"`
	ThreadMetadata ThreadMetadata `json:"thread_metadata"`
	ParentID       snowflake.ID   `json:"parent_id"`
}

type InteractionChannel struct {
	MessageChannel
	Permissions Permissions `json:"permissions"`
}

func (c *InteractionChannel) UnmarshalJSON(data []byte) error {
	var v struct {
		Permissions Permissions `json:"permissions"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var vc UnmarshalChannel
	if err := json.Unmarshal(data, &vc); err != nil {
		return err
	}
	msgChannel, ok := vc.Channel.(MessageChannel)
	if !ok {
		return fmt.Errorf("unknown channel type: %T", vc.Channel)
	}
	c.MessageChannel = msgChannel
	c.Permissions = v.Permissions

	return nil
}

func (c InteractionChannel) MarshalJSON() ([]byte, error) {
	mData, err := json.Marshal(c.MessageChannel)
	if err != nil {
		return nil, err
	}

	pData, err := json.Marshal(struct {
		Permissions Permissions `json:"permissions"`
	}{
		Permissions: c.Permissions,
	})
	if err != nil {
		return nil, err
	}

	return json.Merge(mData, pData)
}
