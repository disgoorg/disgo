package api

func NewActionRow(components ...Component) *ActionRow {
	return &ActionRow{
		ComponentImpl: newComponentImpl(ComponentTypeActionRow),
		Components:    components,
	}
}

type ActionRow struct {
	ComponentImpl
	Components []Component `json:"components"`
}
