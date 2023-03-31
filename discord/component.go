package discord

import (
	"fmt"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

// ComponentType defines different Component(s)
type ComponentType int

// Supported ComponentType(s)
const (
	ComponentTypeActionRow = iota + 1
	ComponentTypeButton
	ComponentTypeStringSelectMenu
	ComponentTypeTextInput
	ComponentTypeUserSelectMenu
	ComponentTypeRoleSelectMenu
	ComponentTypeMentionableSelectMenu
	ComponentTypeChannelSelectMenu
)

type Component interface {
	json.Marshaler
	Type() ComponentType
	component()
}

type ContainerComponent interface {
	Component
	Components() []InteractiveComponent
	containerComponent()
}

type InteractiveComponent interface {
	Component
	ID() string
	interactiveComponent()
}

type UnmarshalComponent struct {
	Component
}

func (u *UnmarshalComponent) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ComponentType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	var (
		component Component
		err       error
	)

	switch cType.Type {
	case ComponentTypeActionRow:
		v := ActionRowComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeButton:
		v := ButtonComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeStringSelectMenu:
		v := StringSelectMenuComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeTextInput:
		v := TextInputComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeUserSelectMenu:
		v := UserSelectMenuComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeRoleSelectMenu:
		v := RoleSelectMenuComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeMentionableSelectMenu:
		v := MentionableSelectMenuComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeChannelSelectMenu:
		v := ChannelSelectMenuComponent{}
		err = json.Unmarshal(data, &v)
		component = v

	default:
		err = fmt.Errorf("unknown component with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	u.Component = component
	return nil
}

type ComponentEmoji struct {
	ID       snowflake.ID `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"`
	Animated bool         `json:"animated,omitempty"`
}

var (
	_ Component          = (*ActionRowComponent)(nil)
	_ ContainerComponent = (*ActionRowComponent)(nil)
)

func NewActionRow(components ...InteractiveComponent) ActionRowComponent {
	return components
}

type ActionRowComponent []InteractiveComponent

func (c ActionRowComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       ComponentType          `json:"type"`
		Components []InteractiveComponent `json:"components"`
	}{
		Type:       c.Type(),
		Components: c,
	})
}

func (c *ActionRowComponent) UnmarshalJSON(data []byte) error {
	var actionRow struct {
		Components []UnmarshalComponent `json:"components"`
	}

	if err := json.Unmarshal(data, &actionRow); err != nil {
		return err
	}

	if len(actionRow.Components) > 0 {
		*c = make([]InteractiveComponent, len(actionRow.Components))
		for i, component := range actionRow.Components {
			(*c)[i] = component.Component.(InteractiveComponent)
		}
	}

	return nil
}

func (ActionRowComponent) Type() ComponentType {
	return ComponentTypeActionRow
}

func (ActionRowComponent) component()          {}
func (ActionRowComponent) containerComponent() {}

func (c ActionRowComponent) Components() []InteractiveComponent {
	return c
}

// Buttons returns all ButtonComponent(s) in the ActionRowComponent
func (c ActionRowComponent) Buttons() []ButtonComponent {
	var buttons []ButtonComponent
	for i := range c {
		if button, ok := c[i].(ButtonComponent); ok {
			buttons = append(buttons, button)
		}
	}
	return buttons
}

// SelectMenus returns all SelectMenuComponent(s) in the ActionRowComponent
func (c ActionRowComponent) SelectMenus() []SelectMenuComponent {
	var selectMenus []SelectMenuComponent
	for i := range c {
		if selectMenu, ok := c[i].(SelectMenuComponent); ok {
			selectMenus = append(selectMenus, selectMenu)
		}
	}
	return selectMenus
}

// TextInputs returns all TextInputComponent(s) in the ActionRowComponent
func (c ActionRowComponent) TextInputs() []TextInputComponent {
	var textInputs []TextInputComponent
	for i := range c {
		if textInput, ok := c[i].(TextInputComponent); ok {
			textInputs = append(textInputs, textInput)
		}
	}
	return textInputs
}

// UpdateComponent returns a new ActionRowComponent with the Component which has the customID replaced
func (c ActionRowComponent) UpdateComponent(customID string, component InteractiveComponent) ActionRowComponent {
	for i, cc := range c {
		if cc.ID() == customID {
			c[i] = component
			return c
		}
	}
	return c
}

// AddComponents returns a new ActionRowComponent with the provided Component(s) added
func (c ActionRowComponent) AddComponents(components ...InteractiveComponent) ActionRowComponent {
	return append(c, components...)
}

// RemoveComponent returns a new ActionRowComponent with the provided Component at the index removed
func (c ActionRowComponent) RemoveComponent(index int) ActionRowComponent {
	if len(c) > index {
		return append(c[:index], c[index+1:]...)
	}
	return c
}

// ButtonStyle defines how the ButtonComponent looks like (https://discord.com/assets/7bb017ce52cfd6575e21c058feb3883b.png)
type ButtonStyle int

// Supported ButtonStyle(s)
const (
	ButtonStylePrimary = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

// NewButton creates a new ButtonComponent with the provided parameters. Link ButtonComponent(s) need a URL and other ButtonComponent(s) need a customID
func NewButton(style ButtonStyle, label string, customID string, url string) ButtonComponent {
	return ButtonComponent{
		Style:    style,
		CustomID: customID,
		URL:      url,
		Label:    label,
	}
}

// NewPrimaryButton creates a new ButtonComponent with ButtonStylePrimary & the provided parameters
func NewPrimaryButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStylePrimary, label, customID, "")
}

// NewSecondaryButton creates a new ButtonComponent with ButtonStyleSecondary & the provided parameters
func NewSecondaryButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStyleSecondary, label, customID, "")
}

// NewSuccessButton creates a new ButtonComponent with ButtonStyleSuccess & the provided parameters
func NewSuccessButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStyleSuccess, label, customID, "")
}

// NewDangerButton creates a new ButtonComponent with ButtonStyleDanger & the provided parameters
func NewDangerButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStyleDanger, label, customID, "")
}

// NewLinkButton creates a new link ButtonComponent with ButtonStyleLink & the provided parameters
func NewLinkButton(label string, url string) ButtonComponent {
	return NewButton(ButtonStyleLink, label, "", url)
}

var (
	_ Component            = (*ButtonComponent)(nil)
	_ InteractiveComponent = (*ButtonComponent)(nil)
)

type ButtonComponent struct {
	Style    ButtonStyle     `json:"style"`
	Label    string          `json:"label,omitempty"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`
	CustomID string          `json:"custom_id,omitempty"`
	URL      string          `json:"url,omitempty"`
	Disabled bool            `json:"disabled,omitempty"`
}

func (c ButtonComponent) MarshalJSON() ([]byte, error) {
	type buttonComponent ButtonComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		buttonComponent
	}{
		Type:            c.Type(),
		buttonComponent: buttonComponent(c),
	})
}

func (ButtonComponent) Type() ComponentType {
	return ComponentTypeButton
}

func (c ButtonComponent) ID() string {
	return c.CustomID
}

func (c ButtonComponent) SetID(id string) InteractiveComponent {
	c.CustomID = id
	return c
}

func (ButtonComponent) component()            {}
func (ButtonComponent) interactiveComponent() {}

// WithStyle returns a new ButtonComponent with the provided style
func (c ButtonComponent) WithStyle(style ButtonStyle) ButtonComponent {
	c.Style = style
	return c
}

// WithLabel returns a new ButtonComponent with the provided label
func (c ButtonComponent) WithLabel(label string) ButtonComponent {
	c.Label = label
	return c
}

// WithEmoji returns a new ButtonComponent with the provided Emoji
func (c ButtonComponent) WithEmoji(emoji ComponentEmoji) ButtonComponent {
	c.Emoji = &emoji
	return c
}

// WithCustomID returns a new ButtonComponent with the provided custom id
func (c ButtonComponent) WithCustomID(customID string) ButtonComponent {
	c.CustomID = customID
	return c
}

// WithURL returns a new ButtonComponent with the provided URL
func (c ButtonComponent) WithURL(url string) ButtonComponent {
	c.URL = url
	return c
}

// AsEnabled returns a new ButtonComponent but enabled
func (c ButtonComponent) AsEnabled() ButtonComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new ButtonComponent but disabled
func (c ButtonComponent) AsDisabled() ButtonComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new ButtonComponent but disabled/enabled
func (c ButtonComponent) WithDisabled(disabled bool) ButtonComponent {
	c.Disabled = disabled
	return c
}

var (
	_ Component            = (*TextInputComponent)(nil)
	_ InteractiveComponent = (*TextInputComponent)(nil)
)

func NewTextInput(customID string, style TextInputStyle, label string) TextInputComponent {
	return TextInputComponent{
		CustomID: customID,
		Style:    style,
		Label:    label,
	}
}

func NewShortTextInput(customID string, label string) TextInputComponent {
	return NewTextInput(customID, TextInputStyleShort, label)
}

func NewParagraphTextInput(customID string, label string) TextInputComponent {
	return NewTextInput(customID, TextInputStyleParagraph, label)
}

type TextInputComponent struct {
	CustomID    string         `json:"custom_id"`
	Style       TextInputStyle `json:"style"`
	Label       string         `json:"label"`
	MinLength   *int           `json:"min_length,omitempty"`
	MaxLength   int            `json:"max_length,omitempty"`
	Required    bool           `json:"required"`
	Placeholder string         `json:"placeholder,omitempty"`
	Value       string         `json:"value,omitempty"`
}

func (c TextInputComponent) MarshalJSON() ([]byte, error) {
	type textInputComponent TextInputComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		textInputComponent
	}{
		Type:               c.Type(),
		textInputComponent: textInputComponent(c),
	})
}

func (TextInputComponent) Type() ComponentType {
	return ComponentTypeTextInput
}

func (c TextInputComponent) ID() string {
	return c.CustomID
}

func (TextInputComponent) component()            {}
func (TextInputComponent) interactiveComponent() {}

// WithCustomID returns a new SelectMenuComponent with the provided customID
func (c TextInputComponent) WithCustomID(customID string) TextInputComponent {
	c.CustomID = customID
	return c
}

// WithStyle returns a new SelectMenuComponent with the provided TextInputStyle
func (c TextInputComponent) WithStyle(style TextInputStyle) TextInputComponent {
	c.Style = style
	return c
}

// WithMinLength returns a new TextInputComponent with the provided minLength
func (c TextInputComponent) WithMinLength(minLength int) TextInputComponent {
	c.MinLength = &minLength
	return c
}

// WithMaxLength returns a new TextInputComponent with the provided maxLength
func (c TextInputComponent) WithMaxLength(maxLength int) TextInputComponent {
	c.MaxLength = maxLength
	return c
}

// WithRequired returns a new TextInputComponent with the provided required
func (c TextInputComponent) WithRequired(required bool) TextInputComponent {
	c.Required = required
	return c
}

// WithPlaceholder returns a new TextInputComponent with the provided placeholder
func (c TextInputComponent) WithPlaceholder(placeholder string) TextInputComponent {
	c.Placeholder = placeholder
	return c
}

// WithValue returns a new TextInputComponent with the provided value
func (c TextInputComponent) WithValue(value string) TextInputComponent {
	c.Value = value
	return c
}

type TextInputStyle int

const (
	TextInputStyleShort = iota + 1
	TextInputStyleParagraph
)
