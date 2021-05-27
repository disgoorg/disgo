package api

type ComponentType int

const (
	ComponentTypeActionRow = iota + 1
	ComponentTypeButton
)

type Component interface {
	Type() ComponentType
}

func newComponentImpl(componentType ComponentType) ComponentImpl {
	return ComponentImpl{ComponentType: componentType}
}

type ComponentImpl struct {
	ComponentType ComponentType `json:"type"`
}

func (t ComponentImpl) Type() ComponentType {
	return t.ComponentType
}

type UnmarshalComponent struct {
	ComponentType ComponentType         `json:"type"`
	Style         ButtonStyle           `json:"style"`
	Label         *string               `json:"label"`
	Emote         *Emote                `json:"emoji"`
	CustomID      string                `json:"custom_id"`
	URL           string                `json:"url"`
	Disabled      bool                  `json:"disabled"`
	Components    []*UnmarshalComponent `json:"components"`
}
