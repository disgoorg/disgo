package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// MessageFilter used to filter Message(s) in a collectors.MessageCollector
type MessageFilter func(message *Message) bool

type Message struct {
	discord.Message
	Bot      *Bot
	Member   *Member
	Author   *User
	Stickers []*MessageSticker
}

func (m *Message) CreateThread(threadCreateWithMessage discord.ThreadCreateWithMessage, opts ...rest.RequestOpt) (GuildThread, error) {
	channel, err := m.Bot.RestServices.ThreadService().CreateThreadWithMessage(m.ChannelID, m.ID, threadCreateWithMessage, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateChannel(channel.(discord.Channel), CacheStrategyNo).(GuildThread), nil
}

// Guild gets the guild_events the message_events was sent in
func (m *Message) Guild() *Guild {
	if m.GuildID == nil {
		return nil
	}
	return m.Bot.Caches.GuildCache().Get(*m.GuildID)
}

// Channel gets the MessageChannel the Message was sent in
func (m *Message) Channel() MessageChannel {
	return m.Bot.Caches.ChannelCache().Get(m.ChannelID).(MessageChannel)
}

// AddReactionByEmote allows you to add an Emoji to a message_events via reaction
func (m *Message) AddReactionByEmote(emote Emoji, opts ...rest.RequestOpt) error {
	return m.AddReaction(emote.Reaction(), opts...)
}

// AddReaction allows you to add a reaction to a message_events from a string, for _examples a custom emoji CommandID, or a native emoji
func (m *Message) AddReaction(emoji string, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.ChannelService().AddReaction(m.ChannelID, m.ID, emoji, opts...)
}

// Update allows you to edit an existing Message sent by you
func (m *Message) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := m.Bot.RestServices.ChannelService().UpdateMessage(m.ChannelID, m.ID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// Delete allows you to edit an existing Message sent by you
func (m *Message) Delete(opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.ChannelService().DeleteMessage(m.ChannelID, m.ID, opts...)
}

// Crosspost crossposts an existing message
func (m *Message) Crosspost(opts ...rest.RequestOpt) (*Message, error) {
	channel := m.Channel()
	if channel != nil && channel.Type() == discord.ChannelTypeGuildNews {
		return nil, discord.ErrChannelNotTypeNews
	}
	message, err := m.Bot.RestServices.ChannelService().CrosspostMessage(m.ChannelID, m.ID, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// Reply allows you to reply to an existing Message
func (m *Message) Reply(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	messageCreate.MessageReference = &discord.MessageReference{MessageID: &m.ID}
	message, err := m.Bot.RestServices.ChannelService().CreateMessage(m.ChannelID, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// ActionRows returns all discord.ActionRowComponent(s) from this Message
func (m *Message) ActionRows() []discord.ActionRowComponent {
	var actionRows []discord.ActionRowComponent
	for i := range m.Components {
		if actionRow, ok := m.Components[i].(discord.ActionRowComponent); ok {
			actionRows = append(actionRows, actionRow)
		}
	}
	return actionRows
}

// InteractiveComponents returns the discord.InteractiveComponent(s) from this Message
func (m *Message) InteractiveComponents() []discord.InteractiveComponent {
	var interactiveComponents []discord.InteractiveComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			interactiveComponents = append(interactiveComponents, m.Components[i].Components()[ii])
		}
	}
	return nil
}

// ComponentByID returns the discord.Component with the specific discord.CustomID
func (m *Message) ComponentByID(customID discord.CustomID) discord.InteractiveComponent {
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
func (m *Message) Buttons() []discord.ButtonComponent {
	var buttons []discord.ButtonComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if button, ok := m.Components[i].Components()[ii].(discord.ButtonComponent); ok {
				buttons = append(buttons, button)
			}
		}
	}
	return buttons
}

// ButtonByID returns a ButtonComponent with the specific customID from this Message
func (m *Message) ButtonByID(customID discord.CustomID) *discord.ButtonComponent {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if button, ok := m.Components[i].Components()[ii].(*discord.ButtonComponent); ok && button.ID() == customID {
				return button
			}
		}
	}
	return nil
}

// SelectMenus returns all SelectMenuComponent(s) from this Message
func (m *Message) SelectMenus() []discord.SelectMenuComponent {
	var selectMenus []discord.SelectMenuComponent
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if button, ok := m.Components[i].Components()[ii].(discord.SelectMenuComponent); ok {
				selectMenus = append(selectMenus, button)
			}
		}
	}
	return selectMenus
}

// SelectMenuByID returns a SelectMenuComponent with the specific customID from this Message
func (m *Message) SelectMenuByID(customID discord.CustomID) *discord.SelectMenuComponent {
	for i := range m.Components {
		for ii := range m.Components[i].Components() {
			if button, ok := m.Components[i].Components()[ii].(*discord.SelectMenuComponent); ok && button.ID() == customID {
				return button
			}
		}
	}
	return nil
}

// IsEphemeral returns true if the Message has MessageFlagEphemeral
func (m *Message) IsEphemeral() bool {
	return m.Flags.Has(discord.MessageFlagEphemeral)
}

// IsWebhookMessage returns true if the Message was sent by a Webhook
func (m *Message) IsWebhookMessage() bool {
	return m.WebhookID != nil
}

// MessageReactionAddFilter used to filter MessageReactionAddEvent in a collectors.MessageReactionAddCollector
type MessageReactionAddFilter func(e *MessageReactionAdd) bool

type MessageReactionAdd struct {
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   *discord.Snowflake
	Member    *Member
	Emoji     discord.ReactionEmoji
}

// MessageReactionRemoveFilter used to filter MessageReactionRemoveEvent in a collectors.MessageReactionRemoveCollector
type MessageReactionRemoveFilter func(e *MessageReactionRemove) bool

type MessageReactionRemove struct {
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	MessageID discord.Snowflake
	GuildID   *discord.Snowflake
	Emoji     discord.ReactionEmoji
}
