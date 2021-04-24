package api

import (
	"errors"
	"time"
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
	MessageTypeApplicationCommand
)

// The MessageFlags of a Message
type MessageFlags int64

// Add allows you to add multiple bits together, producing a new bit
func (f MessageFlags) Add(bits ...MessageFlags) *MessageFlags {
	total := MessageFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	f |= total
	return &f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f MessageFlags) Remove(bits ...MessageFlags) *MessageFlags {
	total := MessageFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	f &^= total
	return &f
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

// Constants for MessageFlags
const (
	MessageFlagNone        MessageFlags = 0
	MessageFlagCrossposted MessageFlags = 1 << (iota - 1)
	MessageFlagIsCrosspost
	MessageFlagSuppressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	_
	MessageFlagEphemeral
	MessageFlagLoading // Message is an interaction of type 5, awaiting further response
)

//MessageAttachment is used for files sent in a Message
type MessageAttachment struct {
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
	ID                Snowflake            `json:"id"`
	GuildID           *Snowflake           `json:"guild_id"`
	Reactions         []*MessageReaction   `json:"reactions"`
	Attachments       []*MessageAttachment `json:"attachments"`
	TTS               bool                 `json:"tts"`
	Embeds            []*Embed             `json:"embeds,omitempty"`
	CreatedAt         time.Time            `json:"timestamp"`
	Mentions          []interface{}        `json:"mentions"`
	MentionEveryone   bool                 `json:"mention_everyone"`
	MentionRoles      []*Role              `json:"mention_roles"`
	MentionChannels   []*Channel           `json:"mention_channels"`
	Pinned            bool                 `json:"pinned"`
	EditedTimestamp   *time.Time           `json:"edited_timestamp"`
	Author            *User                `json:"author"`
	Member            *Member              `json:"member"`
	Content           *string              `json:"content,omitempty"`
	ChannelID         Snowflake            `json:"channel_id"`
	Type              MessageType          `json:"type"`
	Flags             *MessageFlags        `json:"flags"`
	MessageReference  *MessageReference    `json:"message_reference,omitempty"`
	Interaction       *MessageInteraction  `json:"message_interaction,omitempty"`
	WebhookID         *Snowflake           `json:"webhook_id,omitempty"`
	Activity          *MessageActivity     `json:"activity,omitempty"`
	Application       *MessageApplication  `json:"application,omitempty"`
	Stickers          []*MessageSticker    `json:"stickers,omitempty"`
	ReferencedMessage *Message             `json:"referenced_message,omitempty"`
	LastUpdated       *time.Time
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

// AddReactionByEmote allows you to add an Emote to a message_events via reaction
func (m *Message) AddReactionByEmote(emote Emote) error {
	return m.AddReaction(emote.Reaction())
}

// AddReaction allows you to add a reaction to a message_events from a string, for example a custom emoji ID, or a native emoji
func (m *Message) AddReaction(emoji string) error {
	return m.Disgo.RestClient().AddReaction(m.ChannelID, m.ID, emoji)
}

// Edit allows you to edit an existing Message sent by you
func (m *Message) Edit(message *MessageUpdate) (*Message, error) {
	return m.Disgo.RestClient().EditMessage(m.ChannelID, m.ID, message)
}

// Delete allows you to edit an existing Message sent by you
func (m *Message) Delete() error {
	return m.Disgo.RestClient().DeleteMessage(m.ChannelID, m.ID)
}

// Crosspost crossposts an existing message
func (m *Message) Crosspost() (*Message, error) {
	channel := m.Channel()
	if channel != nil && channel.Type != ChannelTypeNews {
		return nil, errors.New("channel type is not NEWS")
	}
	return m.Disgo.RestClient().CrosspostMessage(m.ChannelID, m.ID)
}

// Reply allows you to reply to an existing Message
func (m *Message) Reply(message *MessageCreate) (*Message, error) {
	message.MessageReference = &MessageReference{
		MessageID: &m.ID,
	}
	return m.Disgo.RestClient().SendMessage(m.ChannelID, message)
}

// MessageReaction contains information about the reactions of a message_events
type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emote `json:"emoji"`
}

// MessageUpdate is used to edit a Message
type MessageUpdate struct {
	Content         *string          `json:"content,omitempty"`
	Embed           *Embed           `json:"embed,omitempty"`
	Flags           *MessageFlags    `json:"flags,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
}

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Content          *string           `json:"content,omitempty"`
	TTS              *bool             `json:"tts,omitempty"`
	Embed            *Embed            `json:"embed,omitempty"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
}

// MessageBulkDelete is used to bulk delete Message(s)
type MessageBulkDelete struct {
	Messages []Snowflake `json:"messages"`
}
