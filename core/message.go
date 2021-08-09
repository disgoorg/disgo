package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Message struct {
	discord.Message
	Disgo      Disgo
	Member     *Member
	Author     *User
	Components []Component `json:"components"`
}

// Guild gets the guild_events the message_events was sent in
func (m *Message) Guild() *Guild {
	if m.GuildID == nil {
		return nil
	}
	return m.Disgo.Cache().GuildCache().Get(*m.GuildID)
}

// Channel gets the channel the message_events was sent in
func (m *Message) Channel() MessageChannel {
	return m.Disgo.Cache().ChannelCache().MessageChannel(m.ChannelID)
}

// AddReactionByEmote allows you to add an Emoji to a message_events via reaction
func (m *Message) AddReactionByEmote(emote Emoji) rest.Error {
	return m.AddReaction(emote.Reaction())
}

// AddReaction allows you to add a reaction to a message_events from a string, for example a custom emoji ID, or a native emoji
func (m *Message) AddReaction(emoji string) rest.Error {
	return m.Disgo.RestServices().ChannelsService().AddReaction(m.ChannelID, m.ID, emoji)
}

// Update allows you to edit an existing Message sent by you
func (m *Message) Update(messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := m.Disgo.RestServices().ChannelsService().UpdateMessage(m.ChannelID, m.ID, messageUpdate)
	if err != nil {
		return nil, err
	}
	return m.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// Delete allows you to edit an existing Message sent by you
func (m *Message) Delete() rest.Error {
	return m.Disgo.RestServices().ChannelsService().DeleteMessage(m.ChannelID, m.ID)
}

// Crosspost crossposts an existing message
func (m *Message) Crosspost() (*Message, rest.Error) {
	channel := m.Channel()
	if channel != nil && channel.IsNewsChannel() {
		return nil, rest.NewError(nil, discord.ErrChannelNotTypeNews)
	}
	message, err := m.Disgo.RestServices().ChannelsService().CrosspostMessage(m.ChannelID, m.ID)
	if err != nil {
		return nil, err
	}
	return m.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// Reply allows you to reply to an existing Message
func (m *Message) Reply(messageCreate discord.MessageCreate) (*Message, rest.Error) {
	messageCreate.MessageReference = &discord.MessageReference{MessageID: &m.ID}
	message, err := m.Disgo.RestServices().ChannelsService().CreateMessage(m.ChannelID, messageCreate)
	if err != nil {
		return nil, err
	}
	return m.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
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
	return m.Flags.Has(discord.MessageFlagEphemeral)
}
