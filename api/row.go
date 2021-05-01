package api

func NewRow(components ...Component) *Row {
	return &Row{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeButtons,
		},
		Components: components,
	}
}

type Row struct {
	ComponentImpl
	Components []Component `json:"components,omitempty"`
}
