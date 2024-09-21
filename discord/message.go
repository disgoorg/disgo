package discord

import (
	"fmt"
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
	MessageTypePurchaseNotification MessageType = iota + 11
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
	Components           []ContainerComponent  `json:"components,omitempty"`
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

// ActionRows returns all ActionRowComponent(s) from this Message
func (m Message) ActionRows() []ActionRowComponent {
	var actionRows []ActionRowComponent
	for i := range m.Components {
		if actionRow, ok := m.Components[i].(ActionRowComponent); ok {
			actionRows = append(actionRows, actionRow)
		}
	}
	return actionRows
}

// InteractiveComponents returns the InteractiveComponent(s) from this Message
func (m Message) InteractiveComponents() []InteractiveComponent {
	var interactiveComponents []InteractiveComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			interactiveComponents = append(interactiveComponents, m.Components[i].Components()[ii])
		}
	}
	return interactiveComponents
}

// ComponentByID returns the Component with the specific CustomID
func (m Message) ComponentByID(customID string) InteractiveComponent {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if m.Components[i].Components()[ii].ID() == customID {
				return m.Components[i].Components()[ii]
			}
		}
	}
	return nil
}

// Buttons returns all ButtonComponent(s) from this Message
func (m Message) Buttons() []ButtonComponent {
	var buttons []ButtonComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if button, ok := m.Components[i].Components()[ii].(ButtonComponent); ok {
				buttons = append(buttons, button)
			}
		}
	}
	return buttons
}

// ButtonByID returns a ButtonComponent with the specific customID from this Message
func (m Message) ButtonByID(customID string) (ButtonComponent, bool) {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if button, ok := m.Components[i].Components()[ii].(ButtonComponent); ok && button.ID() == customID {
				return button, true
			}
		}
	}
	return ButtonComponent{}, false
}

// SelectMenus returns all SelectMenuComponent(s) from this Message
func (m Message) SelectMenus() []SelectMenuComponent {
	var selectMenus []SelectMenuComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if selectMenu, ok := m.Components[i].Components()[ii].(SelectMenuComponent); ok {
				selectMenus = append(selectMenus, selectMenu)
			}
		}
	}
	return selectMenus
}

// SelectMenuByID returns a SelectMenuComponent with the specific customID from this Message
func (m Message) SelectMenuByID(customID string) (SelectMenuComponent, bool) {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if selectMenu, ok := m.Components[i].Components()[ii].(SelectMenuComponent); ok && selectMenu.ID() == customID {
				return selectMenu, true
			}
		}
	}
	return nil, false
}

// UserSelectMenus returns all UserSelectMenuComponent(s) from this Message
func (m Message) UserSelectMenus() []UserSelectMenuComponent {
	var userSelectMenus []UserSelectMenuComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if userSelectMenu, ok := m.Components[i].Components()[ii].(UserSelectMenuComponent); ok {
				userSelectMenus = append(userSelectMenus, userSelectMenu)
			}
		}
	}
	return userSelectMenus
}

// UserSelectMenuByID returns a UserSelectMenuComponent with the specific customID from this Message
func (m Message) UserSelectMenuByID(customID string) (UserSelectMenuComponent, bool) {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if userSelectMenu, ok := m.Components[i].Components()[ii].(UserSelectMenuComponent); ok && userSelectMenu.ID() == customID {
				return userSelectMenu, true
			}
		}
	}
	return UserSelectMenuComponent{}, false
}

// RoleSelectMenus returns all RoleSelectMenuComponent(s) from this Message
func (m Message) RoleSelectMenus() []RoleSelectMenuComponent {
	var roleSelectMenus []RoleSelectMenuComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if roleSelectMenu, ok := m.Components[i].Components()[ii].(RoleSelectMenuComponent); ok {
				roleSelectMenus = append(roleSelectMenus, roleSelectMenu)
			}
		}
	}
	return roleSelectMenus
}

// RoleSelectMenuByID returns a RoleSelectMenuComponent with the specific customID from this Message
func (m Message) RoleSelectMenuByID(customID string) (RoleSelectMenuComponent, bool) {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if roleSelectMenu, ok := m.Components[i].Components()[ii].(RoleSelectMenuComponent); ok && roleSelectMenu.ID() == customID {
				return roleSelectMenu, true
			}
		}
	}
	return RoleSelectMenuComponent{}, false
}

// MentionableSelectMenus returns all MentionableSelectMenuComponent(s) from this Message
func (m Message) MentionableSelectMenus() []MentionableSelectMenuComponent {
	var mentionableSelectMenus []MentionableSelectMenuComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if mentionableSelectMenu, ok := m.Components[i].Components()[ii].(MentionableSelectMenuComponent); ok {
				mentionableSelectMenus = append(mentionableSelectMenus, mentionableSelectMenu)
			}
		}
	}
	return mentionableSelectMenus
}

// MentionableSelectMenuByID returns a MentionableSelectMenuComponent with the specific customID from this Message
func (m Message) MentionableSelectMenuByID(customID string) (MentionableSelectMenuComponent, bool) {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if mentionableSelectMenu, ok := m.Components[i].Components()[ii].(MentionableSelectMenuComponent); ok && mentionableSelectMenu.ID() == customID {
				return mentionableSelectMenu, true
			}
		}
	}
	return MentionableSelectMenuComponent{}, false
}

// ChannelSelectMenus returns all ChannelSelectMenuComponent(s) from this Message
func (m Message) ChannelSelectMenus() []ChannelSelectMenuComponent {
	var channelSelectMenus []ChannelSelectMenuComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if channelSelectMenu, ok := m.Components[i].Components()[ii].(ChannelSelectMenuComponent); ok {
				channelSelectMenus = append(channelSelectMenus, channelSelectMenu)
			}
		}
	}
	return channelSelectMenus
}

// ChannelSelectMenuByID returns a ChannelSelectMenuComponent with the specific customID from this Message
func (m Message) ChannelSelectMenuByID(customID string) (ChannelSelectMenuComponent, bool) {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if channelSelectMenu, ok := m.Components[i].Components()[ii].(ChannelSelectMenuComponent); ok && channelSelectMenu.ID() == customID {
				return channelSelectMenu, true
			}
		}
	}
	return ChannelSelectMenuComponent{}, false
}

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
	Type            MessageType          `json:"type"`
	Content         string               `json:"content,omitempty"`
	Embeds          []Embed              `json:"embeds,omitempty"`
	Attachments     []Attachment         `json:"attachments"`
	CreatedAt       time.Time            `json:"timestamp"`
	EditedTimestamp *time.Time           `json:"edited_timestamp"`
	Flags           MessageFlags         `json:"flags"`
	Mentions        []User               `json:"mentions"`
	MentionRoles    []snowflake.ID       `json:"mention_roles"`
	Stickers        []Sticker            `json:"stickers"`
	StickerItems    []MessageSticker     `json:"sticker_items,omitempty"`
	Components      []ContainerComponent `json:"components,omitempty"`
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
	ID                            snowflake.ID                                `json:"id"`
	Type                          InteractionType                             `json:"type"`
	User                          User                                        `json:"user"`
	AuthorizingIntegrationOwners  map[ApplicationIntegrationType]snowflake.ID `json:"authorizing_integration_owners"`
	OriginalResponseMessageID     *snowflake.ID                               `json:"original_response_message_id"`
	Name                          *string                                     `json:"name"`
	InteractedMessageID           *snowflake.ID                               `json:"interacted_message_id"`
	TriggeringInteractionMetadata *InteractionMetadata                        `json:"triggering_interaction_metadata"`
}

type MessageCall struct {
	Participants   []snowflake.ID `json:"participants"`
	EndedTimestamp *time.Time     `json:"ended_timestamp"`
}

func unmarshalComponents(components []UnmarshalComponent) []ContainerComponent {
	containerComponents := make([]ContainerComponent, len(components))
	for i := range components {
		containerComponents[i] = components[i].Component.(ContainerComponent)
	}
	return containerComponents
}
