package discord

import (
	"bytes"
	"fmt"
	"iter"
	"strconv"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/internal/flags"
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
	MessageTypeChannelPinnedMessage
	MessageTypeUserJoin
	MessageTypeGuildBoost
	MessageTypeGuildBoostTier1
	MessageTypeGuildBoostTier2
	MessageTypeGuildBoostTier3
	MessageTypeChannelFollowAdd
	_
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
	MessageTypeGuildDiscoveryGracePeriodInitialWarning
	MessageTypeGuildDiscoveryGracePeriodFinalWarning
	MessageTypeThreadCreated
	MessageTypeReply
	MessageTypeSlashCommand
	MessageTypeThreadStarterMessage
	MessageTypeGuildInviteReminder
	MessageTypeContextMenuCommand
	MessageTypeAutoModerationAction
	MessageTypeRoleSubscriptionPurchase
	MessageTypeInteractionPremiumUpsell
	MessageTypeStageStart
	MessageTypeStageEnd
	MessageTypeStageSpeaker
	_
	MessageTypeStageTopic
	MessageTypeGuildApplicationPremiumSubscription
	_
	_
	_
	MessageTypeGuildIncidentAlertModeEnabled
	MessageTypeGuildIncidentAlertModeDisabled
	MessageTypeGuildIncidentReportRaid
	MessageTypeGuildIncidentReportFalseAlarm
	_
	_
	_
	_
	MessageTypePurchaseNotification
	_
	MessageTypePollResult
)

func (t MessageType) System() bool {
	switch t {
	case MessageTypeDefault, MessageTypeReply, MessageTypeSlashCommand, MessageTypeThreadStarterMessage, MessageTypeContextMenuCommand:
		return false

	default:
		return true
	}
}

func (t MessageType) Deleteable() bool {
	switch t {
	case MessageTypeRecipientAdd, MessageTypeRecipientRemove, MessageTypeCall,
		MessageTypeChannelNameChange, MessageTypeChannelIconChange, MessageTypeThreadStarterMessage:
		return false

	default:
		return true
	}
}

const MessageURLFmt = "https://discord.com/channels/%s/%d/%d"

func MessageURL(guildID snowflake.ID, channelID snowflake.ID, messageID snowflake.ID) string {
	return fmt.Sprintf(MessageURLFmt, guildID, channelID, messageID)
}

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	ID                   snowflake.ID          `json:"id"`
	GuildID              *snowflake.ID         `json:"guild_id"`
	Reactions            []MessageReaction     `json:"reactions"`
	Attachments          []Attachment          `json:"attachments"`
	TTS                  bool                  `json:"tts"`
	Embeds               []Embed               `json:"embeds,omitempty"`
	Components           []LayoutComponent     `json:"components,omitempty"`
	CreatedAt            time.Time             `json:"timestamp"`
	Mentions             []User                `json:"mentions"`
	MentionEveryone      bool                  `json:"mention_everyone"`
	MentionRoles         []snowflake.ID        `json:"mention_roles"`
	MentionChannels      []MentionChannel      `json:"mention_channels"`
	Pinned               bool                  `json:"pinned"`
	EditedTimestamp      *time.Time            `json:"edited_timestamp"`
	Author               User                  `json:"author"`
	Member               *Member               `json:"member"`
	Content              string                `json:"content,omitempty"`
	ChannelID            snowflake.ID          `json:"channel_id"`
	Type                 MessageType           `json:"type"`
	Flags                MessageFlags          `json:"flags"`
	MessageReference     *MessageReference     `json:"message_reference,omitempty"`
	MessageSnapshots     []MessageSnapshot     `json:"message_snapshots,omitempty"`
	Interaction          *MessageInteraction   `json:"interaction,omitempty"`
	WebhookID            *snowflake.ID         `json:"webhook_id,omitempty"`
	Activity             *MessageActivity      `json:"activity,omitempty"`
	Application          *MessageApplication   `json:"application,omitempty"`
	ApplicationID        *snowflake.ID         `json:"application_id,omitempty"`
	StickerItems         []MessageSticker      `json:"sticker_items,omitempty"`
	ReferencedMessage    *Message              `json:"referenced_message,omitempty"`
	LastUpdated          *time.Time            `json:"last_updated,omitempty"`
	Thread               *MessageThread        `json:"thread,omitempty"`
	Position             *int                  `json:"position,omitempty"`
	RoleSubscriptionData *RoleSubscriptionData `json:"role_subscription_data,omitempty"`
	InteractionMetadata  *InteractionMetadata  `json:"interaction_metadata,omitempty"`
	Resolved             *ResolvedData         `json:"resolved,omitempty"`
	Poll                 *Poll                 `json:"poll,omitempty"`
	Call                 *MessageCall          `json:"call,omitempty"`
	Nonce                Nonce                 `json:"nonce,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v struct {
		Components []UnmarshalComponent `json:"components"`
		message
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = Message(v.message)

	if len(v.Components) > 0 {
		m.Components = unmarshalComponents(v.Components)
	}

	if m.Member != nil && m.GuildID != nil {
		m.Member.GuildID = *m.GuildID
	}

	return nil
}

// AllComponents returns an [iter.Seq] of all components in the message.
func (m Message) AllComponents() iter.Seq[Component] {
	return componentIter(m.Components)
}

// JumpURL returns the URL which can be used to jump to the message in the discord client.
func (m Message) JumpURL() string {
	guildID := "@me"
	if m.GuildID != nil {
		guildID = m.GuildID.String()
	}
	return fmt.Sprintf(MessageURLFmt, guildID, m.ChannelID, m.ID) // duplicate code, but there isn't a better way without sacrificing user convenience
}

type MentionChannel struct {
	ID      snowflake.ID `json:"id"`
	GuildID snowflake.ID `json:"guild_id"`
	Type    ChannelType  `json:"type"`
	Name    string       `json:"name"`
}

type MessageThread struct {
	GuildThread
	Member ThreadMember `json:"member"`
}

type MessageSticker struct {
	ID         snowflake.ID      `json:"id"`
	Name       string            `json:"name"`
	FormatType StickerFormatType `json:"format_type"`
}

// MessageReaction contains information about the reactions of a message
type MessageReaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details"`
	Me           bool                 `json:"me"`
	MeBurst      bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []string             `json:"burst_colors"`
}

type ReactionCountDetails struct {
	Burst  int `json:"burst"`
	Normal int `json:"normal"`
}

type MessageReactionType int

const (
	MessageReactionTypeNormal MessageReactionType = iota
	MessageReactionTypeBurst
)

// MessageActivityType is the type of MessageActivity https://com/developers/docs/resources/channel#message-object-message-activity-types
type MessageActivityType int

// Constants for MessageActivityType
const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	_
	MessageActivityTypeJoinRequest
)

// MessageActivity is used for rich presence-related chat embeds in a Message
type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID *string             `json:"party_id,omitempty"`
}

// MessageApplication is used for rich presence-related chat embeds in a Message
type MessageApplication struct {
	ID          snowflake.ID `json:"id"`
	CoverImage  *string      `json:"cover_image,omitempty"`
	Description string       `json:"description"`
	Icon        *string      `json:"icon,omitempty"`
	Name        string       `json:"name"`
}

// MessageReference is a reference to another message
type MessageReference struct {
	Type            MessageReferenceType `json:"type,omitempty"`
	MessageID       *snowflake.ID        `json:"message_id"`
	ChannelID       *snowflake.ID        `json:"channel_id,omitempty"`
	GuildID         *snowflake.ID        `json:"guild_id,omitempty"`
	FailIfNotExists bool                 `json:"fail_if_not_exists,omitempty"`
}

type MessageReferenceType int

const (
	MessageReferenceTypeDefault MessageReferenceType = iota
	MessageReferenceTypeForward
)

type MessageSnapshot struct {
	Message PartialMessage `json:"message"`
}

type PartialMessage struct {
	Type            MessageType       `json:"type"`
	Content         string            `json:"content,omitempty"`
	Embeds          []Embed           `json:"embeds,omitempty"`
	Attachments     []Attachment      `json:"attachments"`
	CreatedAt       time.Time         `json:"timestamp"`
	EditedTimestamp *time.Time        `json:"edited_timestamp"`
	Flags           MessageFlags      `json:"flags"`
	Mentions        []User            `json:"mentions"`
	MentionRoles    []snowflake.ID    `json:"mention_roles"`
	Stickers        []Sticker         `json:"stickers"`
	StickerItems    []MessageSticker  `json:"sticker_items,omitempty"`
	Components      []LayoutComponent `json:"components,omitempty"`
}

func (m *PartialMessage) UnmarshalJSON(data []byte) error {
	type partialMessage PartialMessage
	var v struct {
		Components []UnmarshalComponent `json:"components"`
		partialMessage
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = PartialMessage(v.partialMessage)

	if len(v.Components) > 0 {
		m.Components = unmarshalComponents(v.Components)
	}

	return nil
}

func (m PartialMessage) AllComponents() iter.Seq[Component] {
	return componentIter(m.Components)
}

// MessageInteraction is sent on the Message object when the message is a response to an interaction
type MessageInteraction struct {
	ID   snowflake.ID    `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User User            `json:"user"`
}

type MessageBulkDelete struct {
	Messages []snowflake.ID `json:"messages"`
}

// The MessageFlags of a Message
type MessageFlags int

// Constants for MessageFlags
const (
	MessageFlagCrossposted MessageFlags = 1 << iota
	MessageFlagIsCrosspost
	MessageFlagSuppressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	MessageFlagHasThread
	MessageFlagEphemeral
	MessageFlagLoading // Message is an interaction of type 5, awaiting further response
	MessageFlagFailedToMentionSomeRolesInThread
	_
	_
	_
	MessageFlagSuppressNotifications
	MessageFlagIsVoiceMessage
	MessageFlagHasSnapshot
	// MessageFlagIsComponentsV2 should be set when you want to send v2 components.
	// After setting this, you will not be allowed to send message content and embeds anymore.
	MessageFlagIsComponentsV2
	MessageFlagsNone MessageFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f MessageFlags) Add(bits ...MessageFlags) MessageFlags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f MessageFlags) Remove(bits ...MessageFlags) MessageFlags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f MessageFlags) Has(bits ...MessageFlags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f MessageFlags) Missing(bits ...MessageFlags) bool {
	return flags.Missing(f, bits...)
}

type RoleSubscriptionData struct {
	RoleSubscriptionListingID snowflake.ID `json:"role_subscription_listing_id"`
	TierName                  string       `json:"tier_name"`
	TotalMonthsSubscribed     int          `json:"total_months_subscribed"`
	IsRenewal                 bool         `json:"is_renewal"`
}

type InteractionMetadata struct {
	ID                           snowflake.ID                                `json:"id"`
	Type                         InteractionType                             `json:"type"`
	User                         User                                        `json:"user"`
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]snowflake.ID `json:"authorizing_integration_owners"`
	OriginalResponseMessageID    *snowflake.ID                               `json:"original_response_message_id"`
	// This field will only be present for application command interactions of ApplicationCommandTypeUser.
	// See https://discord.com/developers/docs/resources/message#message-interaction-metadata-object-application-command-interaction-metadata-structure
	TargetUser *User `json:"target_user"`
	// This field will only be present for application command interactions of ApplicationCommandTypeMessage.
	// See https://discord.com/developers/docs/resources/message#message-interaction-metadata-object-application-command-interaction-metadata-structure
	TargetMessageID *snowflake.ID `json:"target_message_id"`
	// This field will only be present for InteractionTypeComponent interactions.
	// See https://discord.com/developers/docs/resources/message#message-interaction-metadata-object-message-component-interaction-metadata-structure
	InteractedMessageID *snowflake.ID `json:"interacted_message_id"`
	// This field will only be present for InteractionTypeModalSubmit interactions.
	// See https://discord.com/developers/docs/resources/message#message-interaction-metadata-object-modal-submit-interaction-metadata-structure
	TriggeringInteractionMetadata *InteractionMetadata `json:"triggering_interaction_metadata"`
}

type MessageCall struct {
	Participants   []snowflake.ID `json:"participants"`
	EndedTimestamp *time.Time     `json:"ended_timestamp"`
}

func unmarshalComponents(components []UnmarshalComponent) []LayoutComponent {
	c := make([]LayoutComponent, len(components))
	for i := range components {
		c[i] = components[i].Component.(LayoutComponent)
	}
	return c
}

// Nonce is a string or int used when sending a message to discord.
type Nonce string

// UnmarshalJSON unmarshals the Nonce from a string or int.
func (n *Nonce) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		i, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		*n = Nonce(strconv.FormatInt(i, 10))
	} else {
		*n = Nonce(unquoted)
	}

	return nil
}
