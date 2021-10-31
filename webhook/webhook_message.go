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

// ActionRows returns all ActionRowComponent(s) from this Message
func (m *Message) ActionRows() []discord.ActionRowComponent {
	var actionRows []discord.ActionRowComponent
	for _, component := range m.Components {
		if actionRow, ok := component.(discord.ActionRowComponent); ok {
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
			case discord.ButtonComponent:
				if c.CustomID == customID {
					return c
				}
			case discord.SelectMenuComponent:
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

// Buttons returns all ButtonComponent(s) from this Message
func (m *Message) Buttons() []discord.ButtonComponent {
	var buttons []discord.ButtonComponent
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow {
			if button, ok := component.(discord.ButtonComponent); ok {
				buttons = append(buttons, button)
			}
		}
	}
	return buttons
}

// ButtonByID returns a ButtonComponent with the specific customID from this Message
func (m *Message) ButtonByID(customID string) *discord.ButtonComponent {
	for _, button := range m.Buttons() {
		if button.CustomID == customID {
			return &button
		}
	}
	return nil
}

// SelectMenus returns all SelectMenuComponent(s) from this Message
func (m *Message) SelectMenus() []discord.SelectMenuComponent {
	var selectMenus []discord.SelectMenuComponent
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow {
			if selectMenu, ok := component.(discord.SelectMenuComponent); ok {
				selectMenus = append(selectMenus, selectMenu)
			}
		}
	}
	return selectMenus
}

// SelectMenuByID returns a SelectMenuComponent with the specific customID from this Message
func (m *Message) SelectMenuByID(customID string) *discord.SelectMenuComponent {
	for _, selectMenu := range m.SelectMenus() {
		if selectMenu.CustomID == customID {
			return &selectMenu
		}
	}
	return nil
}
