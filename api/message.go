package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/DisgoOrg/restclient"
)

// The MessageType indicates the Message type
type MessageType int

// Constants for the MessageType
const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	ChannelPinnedMessage
	MessageTypeGuildMemberJoin
	MessageTypeUserPremiumGuildSubscription
	MessageTypeUserPremiumGuildSubscriptionTier1
	MMessageTypeUserPremiumGuildSubscriptionTier2
	MessageTypeUserPremiumGuildSubscriptionTier3
	MessageTypeChannelFollowAdd
	_
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
	MessageTypeGuildDiscoveryInitialWarning
	MessageTypeGuildDiscoveryFinalWarning
	_
	MessageTypeReply
	MessageTypeCommand
)

// The MessageFlags of a Message
type MessageFlags int64

// Constants for MessageFlags
const (
	MessageFlagCrossposted MessageFlags = 1 << iota
	MessageFlagIsCrosspost
	MessageFlagSuppressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	_
	MessageFlagEphemeral
	MessageFlagLoading              // Message is an interaction of type 5, awaiting further response
	MessageFlagNone    MessageFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f MessageFlags) Add(bits ...MessageFlags) MessageFlags {
	total := MessageFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	f |= total
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f MessageFlags) Remove(bits ...MessageFlags) MessageFlags {
	total := MessageFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	f &^= total
	return f
}

// HasAll will ensure that the bit includes all of the bits entered
func (f MessageFlags) HasAll(bits ...MessageFlags) bool {
	for _, bit := range bits {
		if !f.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (f MessageFlags) Has(bit MessageFlags) bool {
	return (f & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (f MessageFlags) MissingAny(bits ...MessageFlags) bool {
	for _, bit := range bits {
		if !f.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (f MessageFlags) Missing(bit MessageFlags) bool {
	return !f.Has(bit)
}

//Attachment is used for files sent in a Message
type Attachment struct {
	ID       Snowflake `json:"id,omitempty"`
	Filename string    `json:"filename"`
	Size     int       `json:"size"`
	URL      string    `json:"url"`
	ProxyURL string    `json:"proxy_url"`
	Height   *int      `json:"height"`
	Width    *int      `json:"width"`
}

// MessageActivityType is the type of MessageActivity
type MessageActivityType int

//Constants for MessageActivityType
const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	_
	MessageActivityTypeJoinRequest
)

//MessageActivity is used for rich presence-related chat embeds in a Message
type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID *string             `json:"party_id,omitempty"`
}

//MessageApplication is used for rich presence-related chat embeds in a Message
type MessageApplication struct {
	ID          Snowflake `json:"id"`
	CoverImage  *string   `json:"cover_image,omitempty"`
	Description string    `json:"description"`
	Icon        *string   `json:"icon,omitempty"`
	Name        string    `json:"name"`
}

// MessageStickerFormatType is the Format Type of a MessageSticker
type MessageStickerFormatType int

// Constants for MessageStickerFormatType
const (
	MessageStickerFormatPNG MessageStickerFormatType = iota + 1
	MessageStickerFormatAPNG
	MessageStickerFormatLottie
)

// MessageSticker is a sticker sent with a Message
type MessageSticker struct {
	ID          Snowflake                `json:"id"`
	PackID      Snowflake                `json:"pack_id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Tags        *string                  `json:"tags"`
	FormatType  MessageStickerFormatType `json:"format_type"`
}

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	Disgo             Disgo
	ID                Snowflake           `json:"id"`
	GuildID           *Snowflake          `json:"guild_id"`
	Reactions         []MessageReaction   `json:"reactions"`
	Attachments       []Attachment        `json:"attachments"`
	TTS               bool                `json:"tts"`
	Embeds            []Embed             `json:"embeds,omitempty"`
	Components        []Component         `json:"-"`
	CreatedAt         time.Time           `json:"timestamp"`
	Mentions          []interface{}       `json:"mentions"`
	MentionEveryone   bool                `json:"mention_everyone"`
	MentionRoles      []*Role             `json:"mention_roles"`
	MentionChannels   []*Channel          `json:"mention_channels"`
	Pinned            bool                `json:"pinned"`
	EditedTimestamp   *time.Time          `json:"edited_timestamp"`
	Author            *User               `json:"author"`
	Member            *Member             `json:"member"`
	Content           *string             `json:"content,omitempty"`
	ChannelID         Snowflake           `json:"channel_id"`
	Type              MessageType         `json:"type"`
	Flags             MessageFlags        `json:"flags"`
	MessageReference  *MessageReference   `json:"message_reference,omitempty"`
	Interaction       *MessageInteraction `json:"message_interaction,omitempty"`
	WebhookID         *Snowflake          `json:"webhook_id,omitempty"`
	Activity          *MessageActivity    `json:"activity,omitempty"`
	Application       *MessageApplication `json:"application,omitempty"`
	Stickers          []*MessageSticker   `json:"stickers,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message,omitempty"`
	LastUpdated       *time.Time          `json:"last_updated,omitempty"`
}

// Unmarshal is used to unmarshal a Message we received from discord
func (m *Message) Unmarshal(data []byte) error {
	var fullM struct {
		*Message
		Components []UnmarshalComponent `json:"components,omitempty"`
	}
	err := json.Unmarshal(data, &fullM)
	if err != nil {
		return err
	}
	*m = *fullM.Message
	for _, component := range fullM.Components {
		m.Components = append(m.Components, createComponent(component))
	}
	return nil
}

// Marshal is used to marshal a Message we send to discord
func (m *Message) Marshal() ([]byte, error) {
	fullM := struct {
		*Message
		Components []Component `json:"components,omitempty"`
	}{
		Message:    m,
		Components: m.Components,
	}
	fullM.Message = m
	data, err := json.Marshal(fullM)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func createComponent(unmarshalComponent UnmarshalComponent) Component {
	switch unmarshalComponent.ComponentType {
	case ComponentTypeActionRow:
		components := make([]Component, len(unmarshalComponent.Components))
		for i, unmarshalC := range unmarshalComponent.Components {
			components[i] = createComponent(unmarshalC)
		}
		return ActionRow{
			ComponentImpl: ComponentImpl{
				ComponentType: ComponentTypeActionRow,
			},
			Components: components,
		}

	case ComponentTypeButton:
		return Button{
			ComponentImpl: ComponentImpl{
				ComponentType: ComponentTypeButton,
			},
			Style:    unmarshalComponent.Style,
			Label:    unmarshalComponent.Label,
			Emoji:    unmarshalComponent.Emoji,
			CustomID: unmarshalComponent.CustomID,
			URL:      unmarshalComponent.URL,
			Disabled: unmarshalComponent.Disabled,
		}

	default:
		return nil
	}
}

// MessageReference is a reference to another message
type MessageReference struct {
	MessageID       *Snowflake `json:"message_id"`
	ChannelID       *Snowflake `json:"channel_id,omitempty"`
	GuildID         *Snowflake `json:"guild_id,omitempty"`
	FailIfNotExists *bool      `json:"fail_if_not_exists,omitempty"`
}

// MessageInteraction is sent on the Message object when the message_events is a response to an interaction
type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User User            `json:"user"`
}

// Guild gets the guild_events the message_events was sent in
func (m *Message) Guild() *Guild {
	if m.GuildID == nil {
		return nil
	}
	return m.Disgo.Cache().Guild(*m.GuildID)
}

// Channel gets the channel the message_events was sent in
func (m *Message) Channel() *MessageChannel {
	return m.Disgo.Cache().MessageChannel(m.ChannelID)
}

// AddReactionByEmote allows you to add an Emoji to a message_events via reaction
func (m *Message) AddReactionByEmote(emote Emoji) restclient.RestError {
	return m.AddReaction(emote.Reaction())
}

// AddReaction allows you to add a reaction to a message_events from a string, for example a custom emoji ID, or a native emoji
func (m *Message) AddReaction(emoji string) restclient.RestError {
	return m.Disgo.RestClient().AddReaction(m.ChannelID, m.ID, emoji)
}

// Update allows you to edit an existing Message sent by you
func (m *Message) Update(message MessageUpdate) (*Message, restclient.RestError) {
	return m.Disgo.RestClient().UpdateMessage(m.ChannelID, m.ID, message)
}

// Delete allows you to edit an existing Message sent by you
func (m *Message) Delete() restclient.RestError {
	return m.Disgo.RestClient().DeleteMessage(m.ChannelID, m.ID)
}

// Crosspost crossposts an existing message
func (m *Message) Crosspost() (*Message, restclient.RestError) {
	channel := m.Channel()
	if channel != nil && channel.Type != ChannelTypeNews {
		return nil, restclient.NewError(nil, errors.New("channel type is not NEWS"))
	}
	return m.Disgo.RestClient().CrosspostMessage(m.ChannelID, m.ID)
}

// Reply allows you to reply to an existing Message
func (m *Message) Reply(message MessageCreate) (*Message, restclient.RestError) {
	message.MessageReference = &MessageReference{
		MessageID: &m.ID,
	}
	return m.Disgo.RestClient().CreateMessage(m.ChannelID, message)
}

// ActionRows returns all ActionRow(s) from this Message
func (m *Message) ActionRows() []ActionRow {
	if m.IsEphemeral() {
		return nil
	}
	var actionRows []ActionRow
	for _, component := range m.Components {
		if actionRow, ok := component.(ActionRow); ok {
			actionRows = append(actionRows, actionRow)
		}
	}
	return actionRows
}

// ComponentByID returns the first Component with the specific customID
func (m *Message) ComponentByID(customID string) Component {
	if m.IsEphemeral() {
		return nil
	}
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow.Components {
			switch c := component.(type) {
			case Button:
				if c.CustomID == customID {
					return c
				}
			case SelectMenu:
				if c.CustomID == customID {
					return c
				}
			default:
				continue
			}
		}
	}
	return nil
}

// Buttons returns all Button(s) from this Message
func (m *Message) Buttons() []Button {
	if m.IsEphemeral() {
		return nil
	}
	var buttons []Button
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow.Components {
			if button, ok := component.(Button); ok {
				buttons = append(buttons, button)
			}
		}
	}
	return buttons
}

// ButtonByID returns a Button with the specific customID from this Message
func (m *Message) ButtonByID(customID string) *Button {
	if m.IsEphemeral() {
		return nil
	}
	for _, button := range m.Buttons() {
		if button.CustomID == customID {
			return &button
		}
	}
	return nil
}

// SelectMenus returns all SelectMenu(s) from this Message
func (m *Message) SelectMenus() []SelectMenu {
	if m.IsEphemeral() {
		return nil
	}
	var selectMenus []SelectMenu
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow.Components {
			if selectMenu, ok := component.(SelectMenu); ok {
				selectMenus = append(selectMenus, selectMenu)
			}
		}
	}
	return selectMenus
}

// SelectMenuByID returns a SelectMenu with the specific customID from this Message
func (m *Message) SelectMenuByID(customID string) *SelectMenu {
	if m.IsEphemeral() {
		return nil
	}
	for _, selectMenu := range m.SelectMenus() {
		if selectMenu.CustomID == customID {
			return &selectMenu
		}
	}
	return nil
}

// IsEphemeral returns true if the Message has MessageFlagEphemeral
func (m *Message) IsEphemeral() bool {
	return m.Flags.Has(MessageFlagEphemeral)
}

// MessageReaction contains information about the reactions of a message_events
type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

// MessageBulkDelete is used to bulk delete Message(s)
type MessageBulkDelete struct {
	Messages []Snowflake `json:"messages"`
}
