package discord

import "github.com/DisgoOrg/disgo/json"

var _ Component = (*ActionRowComponent)(nil)

func NewActionRow(components ...Component) ActionRowComponent {
	return components
}

type ActionRowComponent []Component

func (r ActionRowComponent) MarshalJSON() ([]byte, error) {
	v := struct {
		Type       ComponentType `json:"type"`
		Components []Component   `json:"components"`
	}{
		Type:       r.Type(),
		Components: r,
	}
	return json.Marshal(v)
}

func (r *ActionRowComponent) UnmarshalJSON(data []byte) error {
	var actionRow struct {
		Components []unmarshalComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &actionRow); err != nil {
		return err
	}

	if len(actionRow.Components) > 0 {
		*r = make([]Component, len(actionRow.Components))
		for i, component := range actionRow.Components {
			(*r)[i] = component.Component
		}
	}

	return nil
}

func (_ ActionRowComponent) Type() ComponentType {
	return ComponentTypeActionRow
}

// SetComponents returns a new ActionRowComponent with the provided Component(s)
func (r *ActionRowComponent) SetComponents(components ...Component) ActionRowComponent {
	*r = components
	return *r
}

// SetComponent returns a new ActionRowComponent with the Component which has the customID replaced
func (r *ActionRowComponent) SetComponent(customID string, component Component) ActionRowComponent {
	for i, c := range *r {
		switch com := c.(type) {
		case ButtonComponent:
			if com.CustomID == customID {
				(*r)[i] = component
				break
			}
		case SelectMenuComponent:
			if com.CustomID == customID {
				(*r)[i] = component
				break
			}
		default:
			continue
		}
	}
	return *r
}

// AddComponents returns a new ActionRowComponent with the provided Component(s) added
func (r *ActionRowComponent) AddComponents(components ...Component) ActionRowComponent {
	*r = append(*r, components...)
	return *r
}

// RemoveComponent returns a new ActionRowComponent with the provided Component at the index removed
func (r *ActionRowComponent) RemoveComponent(index int) ActionRowComponent {
	if len(*r) > index {
		*r = append((*r)[:index], (*r)[index+1:]...)
	}
	return *r
}
