package discord

import "github.com/DisgoOrg/disgo/json"

var (
	_ Component          = (*ActionRowComponent)(nil)
	_ ContainerComponent = (*ActionRowComponent)(nil)
)

//goland:noinspection GoUnusedExportedFunction
func NewActionRow(components ...InteractiveComponent) ActionRowComponent {
	return components
}

type ActionRowComponent []InteractiveComponent

func (c ActionRowComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       ComponentType          `json:"type"`
		Components []InteractiveComponent `json:"components"`
	}{
		Type:       c.Type(),
		Components: c,
	})
}

func (c *ActionRowComponent) UnmarshalJSON(data []byte) error {
	var actionRow struct {
		Components []UnmarshalComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &actionRow); err != nil {
		return err
	}

	if len(actionRow.Components) > 0 {
		*c = make([]InteractiveComponent, len(actionRow.Components))
		for i, component := range actionRow.Components {
			(*c)[i] = component.Component.(InteractiveComponent)
		}
	}

	return nil
}

func (c ActionRowComponent) Type() ComponentType {
	return ComponentTypeActionRow
}

func (c ActionRowComponent) component()          {}
func (c ActionRowComponent) containerComponent() {}

func (c ActionRowComponent) Components() []InteractiveComponent {
	return c
}

// Buttons returns all ButtonComponent(s) in the ActionRowComponent
func (c ActionRowComponent) Buttons() []ButtonComponent {
	var buttons []ButtonComponent
	for i := range c {
		if button, ok := c[i].(ButtonComponent); ok {
			buttons = append(buttons, button)
		}
	}
	return buttons
}

// SelectMenus returns all SelectMenuComponent(s) in the ActionRowComponent
func (c ActionRowComponent) SelectMenus() []SelectMenuComponent {
	var selectMenus []SelectMenuComponent
	for i := range c {
		if selectMenu, ok := c[i].(SelectMenuComponent); ok {
			selectMenus = append(selectMenus, selectMenu)
		}
	}
	return selectMenus
}

// UpdateComponent returns a new ActionRowComponent with the Component which has the customID replaced
func (c ActionRowComponent) UpdateComponent(customID CustomID, component InteractiveComponent) ActionRowComponent {
	for i, cc := range c {
		if cc.ID() == customID {
			c[i] = component
			return c
		}
	}
	return c
}

// AddComponents returns a new ActionRowComponent with the provided Component(s) added
func (c ActionRowComponent) AddComponents(components ...InteractiveComponent) ActionRowComponent {
	return append(c, components...)
}

// RemoveComponent returns a new ActionRowComponent with the provided Component at the index removed
func (c ActionRowComponent) RemoveComponent(index int) ActionRowComponent {
	if len(c) > index {
		return append(c[:index], c[index+1:]...)
	}
	return c
}
