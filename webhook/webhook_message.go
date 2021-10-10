package webhook

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Message struct {
	discord.Message
	WebhookClient *Client
	Components    []core.Component `json:"components"`
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
func (m *Message) ActionRows() []core.ActionRow {
	var actionRows []core.ActionRow
	for _, component := range m.Components {
		if actionRow, ok := component.(core.ActionRow); ok {
			actionRows = append(actionRows, actionRow)
		}
	}
	return actionRows
}

// ComponentByID returns the first Component with the specific customID
func (m *Message) ComponentByID(customID string) core.Component {
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow.Components {
			switch c := component.(type) {
			case core.Button:
				if c.CustomID == customID {
					return c
				}
			case core.SelectMenu:
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
func (m *Message) Buttons() []core.Button {
	var buttons []core.Button
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow.Components {
			if button, ok := component.(core.Button); ok {
				buttons = append(buttons, button)
			}
		}
	}
	return buttons
}

// ButtonByID returns a Button with the specific customID from this Message
func (m *Message) ButtonByID(customID string) *core.Button {
	for _, button := range m.Buttons() {
		if button.CustomID == customID {
			return &button
		}
	}
	return nil
}

// SelectMenus returns all SelectMenu(s) from this Message
func (m *Message) SelectMenus() []core.SelectMenu {
	var selectMenus []core.SelectMenu
	for _, actionRow := range m.ActionRows() {
		for _, component := range actionRow.Components {
			if selectMenu, ok := component.(core.SelectMenu); ok {
				selectMenus = append(selectMenus, selectMenu)
			}
		}
	}
	return selectMenus
}

// SelectMenuByID returns a SelectMenu with the specific customID from this Message
func (m *Message) SelectMenuByID(customID string) *core.SelectMenu {
	for _, selectMenu := range m.SelectMenus() {
		if selectMenu.CustomID == customID {
			return &selectMenu
		}
	}
	return nil
}
