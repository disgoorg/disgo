package discord

import "slices"

// BaseThemeType is the type of BaseTheme
type BaseThemeType int

// Constants for BaseThemeType
const (
	BaseThemeTypeUnset = iota
	BaseThemeTypeDark
	BaseThemeTypeLight
	BaseThemeTypeDarker
	BaseThemeTypeMidnight
)

// NewSharedClientTheme returns a new SharedClientTheme with no fields set
func NewSharedClientTheme() SharedClientTheme {
	return SharedClientTheme{}
}

// SharedClientTheme is a custom client-side theme that can be shared via a message
type SharedClientTheme struct {
	Colors        []string       `json:"colors"`
	GradientAngle int            `json:"gradient_angle"`
	BaseMix       int            `json:"base_mix"`
	BaseTheme     *BaseThemeType `json:"base_theme,omitempty"`
}

// WithColors returns a new SharedClientTheme with the provided colors set (max 5 hex-encoded strings, e.g. "5865F2")
func (t SharedClientTheme) WithColors(colors ...string) SharedClientTheme {
	t.Colors = colors
	return t
}

// WithColor returns a new SharedClientTheme with the color at the provided index set
func (t SharedClientTheme) WithColor(i int, color string) SharedClientTheme {
	if len(t.Colors) > i {
		t.Colors = slices.Insert(t.Colors, i, color)
	}
	return t
}

// AddColors returns a new SharedClientTheme with the provided colors added
func (t SharedClientTheme) AddColors(colors ...string) SharedClientTheme {
	t.Colors = append(t.Colors, colors...)
	return t
}

// RemoveColor returns a new SharedClientTheme with the color at the provided index removed
func (t SharedClientTheme) RemoveColor(i int) SharedClientTheme {
	if len(t.Colors) > i {
		t.Colors = slices.Delete(slices.Clone(t.Colors), i, i+1)
	}
	return t
}

// ClearColors returns a new SharedClientTheme with no colors
func (t SharedClientTheme) ClearColors() SharedClientTheme {
	t.Colors = []string{}
	return t
}

// WithGradientAngle returns a new SharedClientTheme with the gradient angle set (0–360)
func (t SharedClientTheme) WithGradientAngle(angle int) SharedClientTheme {
	t.GradientAngle = angle
	return t
}

// WithBaseMix returns a new SharedClientTheme with the base mix set (0–100)
func (t SharedClientTheme) WithBaseMix(mix int) SharedClientTheme {
	t.BaseMix = mix
	return t
}

// WithBaseTheme returns a new SharedClientTheme with the base theme set
func (t SharedClientTheme) WithBaseTheme(baseTheme BaseThemeType) SharedClientTheme {
	t.BaseTheme = &baseTheme
	return t
}

// ClearBaseTheme returns a new SharedClientTheme with no base theme
func (t SharedClientTheme) ClearBaseTheme() SharedClientTheme {
	t.BaseTheme = nil
	return t
}