package discord

import "github.com/DisgoOrg/disgo/json"

var _ Component = (*ActionRow)(nil)

func NewActionRow(components ...Component) ActionRow {
	return components
}

type ActionRow []Component

func (r ActionRow) MarshalJSON() ([]byte, error) {
	v := struct {
		Type       ComponentType `json:"type"`
		Components []Component   `json:"components"`
	}{
		Type:       r.Type(),
		Components: r,
	}
	return json.Marshal(v)
}

func (r *ActionRow) UnmarshalJSON(data []byte) error {
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

func (_ ActionRow) Type() ComponentType {
	return ComponentTypeActionRow
}

// SetComponents returns a new ActionRow with the provided Component(s)
func (r *ActionRow) SetComponents(components ...Component) ActionRow {
	*r = components
	return *r
}

// SetComponent returns a new ActionRow with the Component which has the customID replaced
func (r *ActionRow) SetComponent(customID string, component Component) ActionRow {
	for i, c := range *r {
		switch com := c.(type) {
		case Button:
			if com.CustomID == customID {
				(*r)[i] = component
				break
			}
		case SelectMenu:
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

// AddComponents returns a new ActionRow with the provided Component(s) added
func (r *ActionRow) AddComponents(components ...Component) ActionRow {
	*r = append(*r, components...)
	return *r
}

// RemoveComponent returns a new ActionRow with the provided Component at the index removed
func (r *ActionRow) RemoveComponent(index int) ActionRow {
	if len(*r) > index {
		*r = append((*r)[:index], (*r)[index+1:]...)
	}
	return *r
}
