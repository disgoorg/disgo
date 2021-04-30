package api

type ComponentType int

const (
	ComponentTypeButtons = iota + 1
	ComponentTypeButton
)

type Style int

const (
	StyleBlurple = iota + 1
	StyleGrey
	StyleGreen
	StyleRed
	StyleHyperlink
)

type Component struct {
	Type ComponentType `json:"type"`
}
