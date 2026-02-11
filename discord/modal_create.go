package discord

import "slices"

var _ InteractionResponseData = (*ModalCreate)(nil)

func NewModalCreate(customID string, title string, components []LayoutComponent) ModalCreate {
	return ModalCreate{
		CustomID:   customID,
		Title:      title,
		Components: components,
	}
}

type ModalCreate struct {
	CustomID   string            `json:"custom_id"`
	Title      string            `json:"title"`
	Components []LayoutComponent `json:"components"`
}

func (ModalCreate) interactionCallbackData() {}

// WithCustomID returns a new ModalCreate with the provided custom ID.
func (m ModalCreate) WithCustomID(customID string) ModalCreate {
	m.CustomID = customID
	return m
}

// WithTitle returns a new ModalCreate with the provided title.
func (m ModalCreate) WithTitle(title string) ModalCreate {
	m.Title = title
	return m
}

// WithComponents returns a new ModalCreate with the provided LayoutComponent(s).
func (m ModalCreate) WithComponents(components ...LayoutComponent) ModalCreate {
	m.Components = components
	return m
}

// UpdateComponent returns a new ModalCreate with the provided LayoutComponent at the index.
func (m ModalCreate) UpdateComponent(i int, container LayoutComponent) ModalCreate {
	if len(m.Components) > i {
		m.Components = slices.Clone(m.Components)
		m.Components[i] = container
	}
	return m
}

// AddLabel returns a new ModalCreate with a new LabelComponent containing the provided label and component added.
func (m ModalCreate) AddLabel(label string, component LabelSubComponent) ModalCreate {
	m.Components = append(m.Components, NewLabel(label, component))
	return m
}

// AddComponents returns a new ModalCreate with the provided LayoutComponent(s) added.
func (m ModalCreate) AddComponents(containers ...LayoutComponent) ModalCreate {
	m.Components = append(m.Components, containers...)
	return m
}

// RemoveComponent returns a new ModalCreate with the LayoutComponent at the index removed.
func (m ModalCreate) RemoveComponent(i int) ModalCreate {
	if len(m.Components) > i {
		m.Components = slices.Delete(slices.Clone(m.Components), i, i+1)
	}
	return m
}

// ClearComponents returns a new ModalCreate with no LayoutComponent(s).
func (m ModalCreate) ClearComponents() ModalCreate {
	m.Components = []LayoutComponent{}
	return m
}
