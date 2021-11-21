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
	return interactiveComponents
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
