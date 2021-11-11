package discord

import "github.com/DisgoOrg/disgo/json"

var _ Component = (*ActionRowComponent)(nil)

func NewActionRow(components ...Component) ActionRowComponent {
	return components
}

type ActionRowComponent []Component

func (c ActionRowComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       ComponentType `json:"type"`
		Components []Component   `json:"components"`
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
		*c = make([]Component, len(actionRow.Components))
		for i, component := range actionRow.Components {
			(*c)[i] = component.Component
		}
	}

	return nil
}

func (c ActionRowComponent) Type() ComponentType {
	return ComponentTypeActionRow
}

func (c ActionRowComponent) component() {}

// SetComponents returns a new ActionRowComponent with the provided Component(s)
func (c ActionRowComponent) SetComponents(components ...Component) ActionRowComponent {
	return components
}

// SetComponent returns a new ActionRowComponent with the Component which has the customID replaced
func (c ActionRowComponent) SetComponent(customID string, component Component) ActionRowComponent {
	for i, cc := range c {
		switch com := cc.(type) {
		case ButtonComponent:
			if com.CustomID == customID {
				c[i] = component
				break
			}

		case SelectMenuComponent:
			if com.CustomID == customID {
				c[i] = component
				break
			}

		default:
			continue
		}
	}
	return c
}

// AddComponents returns a new ActionRowComponent with the provided Component(s) added
func (c ActionRowComponent) AddComponents(components ...Component) ActionRowComponent {
	return append(c, components...)
}

// RemoveComponent returns a new ActionRowComponent with the provided Component at the index removed
func (c ActionRowComponent) RemoveComponent(index int) ActionRowComponent {
	if len(c) > index {
		return append(c[:index], c[index+1:]...)
	}
	return c
}
