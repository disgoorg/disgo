package api

type ComponentType int

const (
	ComponentTypeButtons = iota + 1
	ComponentTypeButton
)

type ComponentStyle int

const (
	ComponentStyleBlurple = iota + 1
	ComponentStyleGrey
	ComponentStyleGreen
	ComponentStyleRed
	ComponentStyleHyperlink
)

type Component struct {
	ComponentType ComponentType   `json:"type"`
	Style         *ComponentStyle `json:"style,omitempty"`
	CustomID      string          `json:"custom_id,omitempty"`
	Label         string          `json:"label,omitempty"`
	URL           string          `json:"url,omitempty"`
	Emoji         *Emoji          `json:"emoji,omitempty"`
	Disabled      bool            `json:"disabled,omitempty"`
	Components    []*Component    `json:"components,omitempty"`
}

type Emoji struct {
	Name string    `json:"name,omitempty"`
	ID   Snowflake `json:"id,omitempty"`
}
