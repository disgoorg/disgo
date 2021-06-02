package api

// NewActionRow creates a new ActionRow holding th provided Component(s)
func NewActionRow(components ...Component) ActionRow {
	return ActionRow{
		ComponentImpl: newComponentImpl(ComponentTypeActionRow),
		Components:    components,
	}
}

// ActionRow holds up to 5 Component(s) in a row
type ActionRow struct {
	ComponentImpl
	Components []Component `json:"components"`
}
