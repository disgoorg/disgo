package discord

import "github.com/DisgoOrg/disgo/json"

// The MessageType indicates the Message type
type MessageType int

// Constants for the MessageType
//goland:noinspection GoUnusedConst
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
	MessageTypeThreadCreated
	MessageTypeReply
	MessageTypeSlashCommand
	MessageTypeThreadStarterMessage
	MessageTypeGuildInviteReminder
	MessageTypeContextMenuCommand
)

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	ID                Snowflake           `json:"id"`
	GuildID           *Snowflake          `json:"guild_id"`
	Reactions         []MessageReaction   `json:"reactions"`
	Attachments       []Attachment        `json:"attachments"`
	TTS               bool                `json:"tts"`
	Embeds            []Embed             `json:"embeds,omitempty"`
	Components        []Component         `json:"components,omitempty"`
	CreatedAt         Time                `json:"timestamp"`
	Mentions          []interface{}       `json:"mentions"`
	MentionEveryone   bool                `json:"mention_everyone"`
	MentionRoles      []Role              `json:"mention_roles"`
	MentionChannels   []Channel           `json:"mention_channels"`
	Pinned            bool                `json:"pinned"`
	EditedTimestamp   *Time               `json:"edited_timestamp"`
	Author            User                `json:"author"`
	Member            *Member             `json:"member"`
	Content           string              `json:"content,omitempty"`
	ChannelID         Snowflake           `json:"channel_id"`
	Type              MessageType         `json:"type"`
	Flags             MessageFlags        `json:"flags"`
	MessageReference  *MessageReference   `json:"message_reference,omitempty"`
	Interaction       *MessageInteraction `json:"interaction,omitempty"`
	WebhookID         *Snowflake          `json:"webhook_id,omitempty"`
	Activity          *MessageActivity    `json:"activity,omitempty"`
	Application       *MessageApplication `json:"application,omitempty"`
	Stickers          []MessageSticker    `json:"sticker_items,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message,omitempty"`
	LastUpdated       *Time               `json:"last_updated,omitempty"`
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var message struct {
		*Message
		Components []unmarshalComponent `json:"components"`
	}

	if err := json.Unmarshal(b, &message); err != nil {
		return err
	}

	if len(message.Components) > 0 {
		m.Components = make([]Component, len(message.Components))
		for i := range message.Components {
			m.Components[i] = message.Components[i].Component
		}
	}

	return nil
}

type MessageSticker struct {
	ID         Snowflake         `json:"id"`
	Name       string            `json:"name"`
	FormatType StickerFormatType `json:"format_type"`
}

// MessageReaction contains information about the reactions of a message_events
type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

// MessageActivityType is the type of MessageActivity https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
type MessageActivityType int

//Constants for MessageActivityType
//goland:noinspection GoUnusedConst
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

// MessageReference is a reference to another message
type MessageReference struct {
	MessageID       *Snowflake `json:"message_id"`
	ChannelID       *Snowflake `json:"channel_id,omitempty"`
	GuildID         *Snowflake `json:"guild_id,omitempty"`
	FailIfNotExists bool       `json:"fail_if_not_exists,omitempty"`
}

// MessageInteraction is sent on the Message object when the message_events is a response to an interaction
type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User User            `json:"user"`
}

type MessageBulkDelete struct {
	Messages []Snowflake `json:"message s"`
}

// The MessageFlags of a Message
type MessageFlags int64

// Constants for MessageFlags
//goland:noinspection GoUnusedConst
const (
	MessageFlagCrossposted MessageFlags = 1 << iota
	MessageFlagIsCrosspost
	MessageFlagSuppressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	MessageFlagHasThread
	MessageFlagEphemeral
	MessageFlagLoading              // Message is an interaction of type 5, awaiting further response
	MessageFlagNone    MessageFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f MessageFlags) Add(bits ...MessageFlags) MessageFlags {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f MessageFlags) Remove(bits ...MessageFlags) MessageFlags {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func (f MessageFlags) Has(bits ...MessageFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (f MessageFlags) Missing(bits ...MessageFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}
