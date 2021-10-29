package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Message struct {
	discord.Message
	WebhookClient *Client
}

// Update allows you to edit an existing Message sent by you
func (m *Message) Update(messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return m.WebhookClient.UpdateMessage(m.ID, messageUpdate, opts...)
}

// Delete allows you to edit an existing Message sent by you
func (m *Message) Delete(opts ...rest.RequestOpt) error {
	return m.WebhookClient.DeleteMessage(m.ID, opts...)
}

// ActionRows returns all ActionRow(s) from this Message
func (m *Message) ActionRows() []discord.ActionRow {
	var actionRows []discord.ActionRow
	for _, component := range m.Components {
		if actionRow, ok := component.(discord.ActionRow); ok {
			actionRows = append(actionRows, actionRow)
		}
	}
	return actionRows
}

// ComponentByID returns the first Component with the specific customID
func (m *Message) ComponentByID(customID string) discord.Component {
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow {
			switch c := component.(type) {
			case discord.Button:
				if c.CustomID == customID {
					return c
				}
			case discord.SelectMenu:
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
func (m *Message) Buttons() []discord.Button {
	var buttons []discord.Button
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow {
			if button, ok := component.(discord.Button); ok {
				buttons = append(buttons, button)
			}
		}
	}
	return buttons
}

// ButtonByID returns a Button with the specific customID from this Message
func (m *Message) ButtonByID(customID string) *discord.Button {
	for _, button := range m.Buttons() {
		if button.CustomID == customID {
			return &button
		}
	}
	return nil
}

// SelectMenus returns all SelectMenu(s) from this Message
func (m *Message) SelectMenus() []discord.SelectMenu {
	var selectMenus []discord.SelectMenu
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow {
			if selectMenu, ok := component.(discord.SelectMenu); ok {
				selectMenus = append(selectMenus, selectMenu)
			}
		}
	}
	return selectMenus
}

// SelectMenuByID returns a SelectMenu with the specific customID from this Message
func (m *Message) SelectMenuByID(customID string) *discord.SelectMenu {
	for _, selectMenu := range m.SelectMenus() {
		if selectMenu.CustomID == customID {
			return &selectMenu
		}
	}
	return nil
}
