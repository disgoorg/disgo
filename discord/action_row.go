package discord

// ActionRow holds up to 5 Component(s) in a row
type ActionRow struct {
	Type       ComponentType `json:"type"`
	Components []Component   `json:"components"`
}
