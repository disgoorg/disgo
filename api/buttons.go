package api

type Buttons struct {
	Component
	Buttons []*Button `json:"components,omitempty"`
}

type Button struct {
	Component
	Style    *Style `json:"style,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	Label    string `json:"label,omitempty"`
	URL      string `json:"url,omitempty"`
	Emoji    *Emoji `json:"emoji,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

type Emoji struct {
	Name string    `json:"name,omitempty"`
	ID   Snowflake `json:"id,omitempty"`
}
