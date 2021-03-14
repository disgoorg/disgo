package api

import (
	"time"

	"github.com/chebyrash/promise"
)

type MessageType int

const (
	MessageTypeDefault = iota
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
	_
	_
	_
	MessageTypeReply
	MessageTypeApplicationCommand
)

// Message is a struct for messages sent in discord text-based channels
type Message struct {
	Disgo           Disgo
	ID              Snowflake     `json:"id"`
	GuildID         *Snowflake     `json:"guild_id"`
	Reactions       []Reactions   `json:"reactions"`
	Attachments     []interface{} `json:"attachments"`
	Tts             bool          `json:"tts"`
	Embeds          []Embed       `json:"embeds,omitempty"`
	CreatedAt       time.Time     `json:"timestamp"`
	MentionEveryone bool          `json:"mention_everyone"`
	Pinned          bool          `json:"pinned"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Author          User          `json:"author"`
	MentionRoles    []interface{} `json:"mention_roles"`
	Content         string        `json:"content"`
	ChannelID       Snowflake     `json:"channel_id"`
	Mentions        []interface{} `json:"mentions"`
	MessageType     MessageType   `json:"type"`
	LastUpdated     *time.Time
}

// missing Member, mention channels, nonce, webhook id, type, activity, application, message_reference, flags, stickers
// referenced_message, interaction
// https://discord.com/developers/docs/resources/channel#message-object

// Guild gets the guild_events the message_events was sent in
func (m Message) Guild() *Guild {
	if m.GuildID == nil {
		return nil
	}
	return m.Disgo.Cache().Guild(*m.GuildID)
}

// Channel gets the channel the message_events was sent in
func (m Message) Channel() *MessageChannel {
	return nil //m.Disgo.Cache().MessageChannel(m.ChannelID)
}

// AddReactionByEmote allows you to add an Emote to a message_events via reaction
func (m Message) AddReactionByEmote(emote Emote) *promise.Promise {
	return m.AddReaction(emote.Reaction())
}

// AddReaction allows you to add a reaction to a message_events from a string, for example a custom emoji ID, or a native
// emoji
func (m Message) AddReaction(emoji string) *promise.Promise {
	return m.Disgo.RestClient().AddReaction(m.ChannelID, m.ID, emoji)
}

// Reactions contains information about the reactions of a message_events
type Reactions struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emote `json:"emoji"`
}

// Embed allows you to send embeds to discord
type Embed struct {
}
