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
