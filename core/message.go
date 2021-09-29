package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Message struct {
	discord.Message
	Bot        *Bot
	Member     *Member
	Author     *User
	Components []Component
	Stickers   []*MessageSticker
}

// Guild returns the Guild this Message was sent in
func (m *Message) Guild() *Guild {
	if m.GuildID == nil {
		return nil
	}
	return m.Bot.Caches.GuildCache().Get(*m.GuildID)
}

// Channel returns the Channel this Message was sent in
func (m *Message) Channel() *Channel {
	return m.Bot.Caches.ChannelCache().Get(m.ChannelID)
}

// AddReactionByEmote adds a reaction to the Message with the specified Emoji
func (m *Message) AddReactionByEmote(emote Emoji, opts ...rest.RequestOpt) rest.Error {
	return m.AddReaction(emote.Reaction(), opts...)
}

// AddReaction adds a reaction to the Message with the specified string containing a custom emoji ID or a native emoji unicode
func (m *Message) AddReaction(emoji string, opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.ChannelService().AddReaction(m.ChannelID, m.ID, emoji, opts...)
}

// Update edits the Message with the content provided in discord.MessageUpdate
func (m *Message) Update(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := m.Bot.RestServices.ChannelService().UpdateMessage(m.ChannelID, m.ID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// Delete deletes this Message
func (m *Message) Delete(opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.ChannelService().DeleteMessage(m.ChannelID, m.ID, opts...)
}

// Crosspost crossposts this Message
func (m *Message) Crosspost(opts ...rest.RequestOpt) (*Message, rest.Error) {
	channel := m.Channel()
	if channel != nil && channel.IsNewsChannel() {
		return nil, rest.NewError(nil, discord.ErrChannelNotTypeNews)
	}
	message, err := m.Bot.RestServices.ChannelService().CrosspostMessage(m.ChannelID, m.ID, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// Reply replies to this Message
func (m *Message) Reply(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	messageCreate.MessageReference = &discord.MessageReference{MessageID: &m.ID}
	message, err := m.Bot.RestServices.ChannelService().CreateMessage(m.ChannelID, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// ActionRows returns all ActionRow(s) from this Message
func (m *Message) ActionRows() []ActionRow {
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
	for _, button := range m.Buttons() {
		if button.CustomID == customID {
			return &button
		}
	}
	return nil
}

// SelectMenus returns all SelectMenu(s) from this Message
func (m *Message) SelectMenus() []SelectMenu {
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
	for _, selectMenu := range m.SelectMenus() {
		if selectMenu.CustomID == customID {
			return &selectMenu
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
