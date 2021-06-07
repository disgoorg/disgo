package api

// ComponentType defines different Component(s)
type ComponentType int

// Supported ComponentType(s)
const (
	ComponentTypeActionRow = iota + 1
	ComponentTypeButton
)

// Component is a general interface each Component needs to implement
type Component interface {
	Type() ComponentType
}

func newComponentImpl(componentType ComponentType) ComponentImpl {
	return ComponentImpl{ComponentType: componentType}
}

// ComponentImpl is used to embed in each different ComponentType
type ComponentImpl struct {
	ComponentType ComponentType `json:"type"`
}

// Type returns the ComponentType of this Component
func (t ComponentImpl) Type() ComponentType {
	return t.ComponentType
}

// UnmarshalComponent is used for easier unmarshalling of different Component(s)
type UnmarshalComponent struct {
	ComponentType ComponentType        `json:"type"`
	Style         ButtonStyle          `json:"style"`
	Label         *string              `json:"label"`
	Emoji         *Emoji               `json:"emoji"`
	CustomID      string               `json:"custom_id"`
	URL           string               `json:"url"`
	Disabled      bool                 `json:"disabled"`
	Components    []UnmarshalComponent `json:"components"`
}
