package discord

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
	MessageTypeCommand
)

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

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	ID                Snowflake            `json:"id"`
	GuildID           *Snowflake           `json:"guild_id"`
	Reactions         []MessageReaction    `json:"reactions"`
	Attachments       []Attachment         `json:"attachments"`
	TTS               bool                 `json:"tts"`
	Embeds            []Embed              `json:"embeds,omitempty"`
	Components        []UnmarshalComponent `json:"-"`
	CreatedAt         Time                 `json:"timestamp"`
	Mentions          []interface{}        `json:"mentions"`
	MentionEveryone   bool                 `json:"mention_everyone"`
	MentionRoles      []Role               `json:"mention_roles"`
	MentionChannels   []Channel            `json:"mention_channels"`
	Pinned            bool                 `json:"pinned"`
	EditedTimestamp   *Time                `json:"edited_timestamp"`
	Author            User                 `json:"author"`
	Member            *Member              `json:"member"`
	Content           *string              `json:"content,omitempty"`
	ChannelID         Snowflake            `json:"channel_id"`
	Type              MessageType          `json:"type"`
	Flags             MessageFlags         `json:"flags"`
	MessageReference  *MessageReference    `json:"message_reference,omitempty"`
	Interaction       *MessageInteraction  `json:"interaction,omitempty"`
	WebhookID         *Snowflake           `json:"webhook_id,omitempty"`
	Activity          *MessageActivity     `json:"activity,omitempty"`
	Application       *MessageApplication  `json:"application,omitempty"`
	Stickers          []Sticker            `json:"stickers,omitempty"`
	ReferencedMessage *Message             `json:"referenced_message,omitempty"`
	LastUpdated       *Time                `json:"last_updated,omitempty"`
}

// MessageReaction contains information about the reactions of a message_events
type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

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
