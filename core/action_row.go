package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

var _ Component = (*ActionRow)(nil)

// NewActionRow creates a new ActionRow holding the provided Component(s)
func NewActionRow(components ...Component) ActionRow {
	return ActionRow{
		Component: discord.Component{
			Type: discord.ComponentTypeActionRow,
		},
		Components: components,
	}
}

type ActionRow struct {
	discord.Component
	Components []Component `json:"components"`
}

// Type returns the discord.ComponentType of this Component
func (r ActionRow) Type() discord.ComponentType {
	return r.Component.Type
}

// SetComponents replaces the Components in this ActionRow, returning the updated ActionRow object
func (r ActionRow) SetComponents(components ...Component) ActionRow {
	r.Components = components
	return r
}

// SetComponent replaces the Component bound to the provided customID, returning the updated ActionRow object
func (r ActionRow) SetComponent(customID string, component Component) ActionRow {
	for i, c := range r.Components {
		switch com := c.(type) {
		case Button:
			if com.CustomID == customID {
				r.Components[i] = component
				break
			}
		case SelectMenu:
			if com.CustomID == customID {
				r.Components[i] = component
				break
			}
		default:
			continue
		}
	}
	return r
}

// AddComponents adds the provided Components to this ActionRow, returning the updated ActionRow object
func (r ActionRow) AddComponents(components ...Component) ActionRow {
	r.Components = append(r.Components, components...)
	return r
}

// RemoveComponent removes the Component at the provided index, returning the updated ActionRow object
func (r ActionRow) RemoveComponent(index int) ActionRow {
	if len(r.Components) > index {
		r.Components = append(r.Components[:index], r.Components[index+1:]...)
	}
	return r
}
