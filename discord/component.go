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
	ComponentTypeActionRow ComponentType = iota + 1
	ComponentTypeButton
	ComponentTypeStringSelectMenu
	ComponentTypeTextInput
	ComponentTypeUserSelectMenu
	ComponentTypeRoleSelectMenu
	ComponentTypeMentionableSelectMenu
	ComponentTypeChannelSelectMenu
	ComponentTypeSection
	ComponentTypeTextDisplay
	ComponentTypeThumbnail
	ComponentTypeMediaGallery
	ComponentTypeFile
	ComponentTypeSeparator
	ComponentTypeContainer
)

type Component interface {
	json.Marshaler
	Type() ComponentType
	GetID() int
	component()
}

type InteractiveComponent interface {
	Component
	GetCustomID() string
	interactiveComponent()
}

// LayoutComponent is an interface for all components that can be present as a top level component in a [Message].
// [ActionRowComponent]
// [SectionComponent]
// [TextDisplayComponent]
// [MediaGalleryComponent]
// [FileComponent]
// [SeparatorComponent]
// [ContainerComponent]
type LayoutComponent interface {
	Component
	layoutComponent()
}

type MessageComponent interface {
	Component
	messageComponent()
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
		var v ActionRowComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeButton:
		var v ButtonComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeStringSelectMenu:
		var v StringSelectMenuComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeTextInput:
		var v TextInputComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeUserSelectMenu:
		var v UserSelectMenuComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeRoleSelectMenu:
		var v RoleSelectMenuComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeMentionableSelectMenu:
		var v MentionableSelectMenuComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeChannelSelectMenu:
		var v ChannelSelectMenuComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeSection:
		var v SectionComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeTextDisplay:
		var v TextDisplayComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeThumbnail:
		var v ThumbnailComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeMediaGallery:
		var v MediaGalleryComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeFile:
		var v FileComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeSeparator:
		var v SeparatorComponent
		err = json.Unmarshal(data, &v)
		component = v

	case ComponentTypeContainer:
		var v ContainerComponent
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
	_ Component = (*ActionRowComponent)(nil)
)

func NewActionRow(components ...Component) ActionRowComponent {
	return ActionRowComponent{
		Components: components,
	}
}

type ActionRowComponent struct {
	ID         int         `json:"id,omitempty"`
	Components []Component `json:"components"`
}

func (c ActionRowComponent) MarshalJSON() ([]byte, error) {
	type actionRowComponent ActionRowComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		actionRowComponent
	}{
		Type:               c.Type(),
		actionRowComponent: actionRowComponent(c),
	})
}

func (ActionRowComponent) Type() ComponentType {
	return ComponentTypeActionRow
}

func (c ActionRowComponent) GetID() int {
	return c.ID
}

func (c ActionRowComponent) GetComponents() []Component {
	return c.Components
}

func (ActionRowComponent) component()       {}
func (ActionRowComponent) layoutComponent() {}

// Buttons returns all ButtonComponent(s) in the ActionRowComponent
func (c ActionRowComponent) Buttons() []ButtonComponent {
	var buttons []ButtonComponent
	for i := range c.Components {
		if button, ok := c.Components[i].(ButtonComponent); ok {
			buttons = append(buttons, button)
		}
	}
	return buttons
}

// SelectMenus returns all SelectMenuComponent(s) in the ActionRowComponent
func (c ActionRowComponent) SelectMenus() []SelectMenuComponent {
	var selectMenus []SelectMenuComponent
	for i := range c.Components {
		if selectMenu, ok := c.Components[i].(SelectMenuComponent); ok {
			selectMenus = append(selectMenus, selectMenu)
		}
	}
	return selectMenus
}

// TextInputs returns all TextInputComponent(s) in the ActionRowComponent
func (c ActionRowComponent) TextInputs() []TextInputComponent {
	var textInputs []TextInputComponent
	for i := range c.Components {
		if textInput, ok := c.Components[i].(TextInputComponent); ok {
			textInputs = append(textInputs, textInput)
		}
	}
	return textInputs
}

// UpdateComponent returns a new ActionRowComponent with the Component which has the customID replaced
// TODO: fix this
func (c ActionRowComponent) UpdateComponent(customID string, component Component) ActionRowComponent {
	//for i, cc := range c.Components {
	//	if cc.ID() == customID {
	//		c[i] = component
	//		return c
	//	}
	//}
	return c
}

// AddComponents returns a new ActionRowComponent with the provided Component(s) added
func (c ActionRowComponent) AddComponents(components ...Component) ActionRowComponent {
	c.Components = append(c.Components, components...)
	return c
}

// RemoveComponent returns a new ActionRowComponent with the provided Component at the index removed
func (c ActionRowComponent) RemoveComponent(index int) ActionRowComponent {
	if len(c.Components) > index {
		c.Components = append(c.Components[:index], c.Components[index+1:]...)
	}
	return c
}

// ButtonStyle defines how the ButtonComponent looks like (https://discord.com/assets/7bb017ce52cfd6575e21c058feb3883b.png)
type ButtonStyle int

// Supported ButtonStyle(s)
const (
	ButtonStylePrimary ButtonStyle = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
	ButtonStylePremium
)

// NewButton creates a new ButtonComponent with the provided parameters. Link ButtonComponent(s) need a URL and other ButtonComponent(s) need a customID
func NewButton(style ButtonStyle, label string, customID string, url string, skuID snowflake.ID) ButtonComponent {
	return ButtonComponent{
		Style:    style,
		CustomID: customID,
		URL:      url,
		Label:    label,
		SkuID:    skuID,
	}
}

// NewPrimaryButton creates a new ButtonComponent with ButtonStylePrimary & the provided parameters
func NewPrimaryButton(label string, customID string) ButtonComponent {
	return ButtonComponent{
		Style:    ButtonStylePrimary,
		Label:    label,
		CustomID: customID,
	}
}

// NewSecondaryButton creates a new ButtonComponent with ButtonStyleSecondary & the provided parameters
func NewSecondaryButton(label string, customID string) ButtonComponent {
	return ButtonComponent{
		Style:    ButtonStyleSecondary,
		Label:    label,
		CustomID: customID,
	}
}

// NewSuccessButton creates a new ButtonComponent with ButtonStyleSuccess & the provided parameters
func NewSuccessButton(label string, customID string) ButtonComponent {
	return ButtonComponent{
		Style:    ButtonStyleSuccess,
		Label:    label,
		CustomID: customID,
	}
}

// NewDangerButton creates a new ButtonComponent with ButtonStyleDanger & the provided parameters
func NewDangerButton(label string, customID string) ButtonComponent {
	return ButtonComponent{
		Style:    ButtonStyleDanger,
		Label:    label,
		CustomID: customID,
	}
}

// NewLinkButton creates a new link ButtonComponent with ButtonStyleLink & the provided parameters
func NewLinkButton(label string, url string) ButtonComponent {
	return ButtonComponent{
		Style: ButtonStyleLink,
		Label: label,
		URL:   url,
	}
}

// NewPremiumButton creates a new ButtonComponent with ButtonStylePremium & the provided parameters
func NewPremiumButton(skuID snowflake.ID) ButtonComponent {
	return ButtonComponent{
		Style: ButtonStylePremium,
		SkuID: skuID,
	}
}

var (
	_ Component = (*ButtonComponent)(nil)
)

type ButtonComponent struct {
	ID       int             `json:"id,omitempty"`
	Style    ButtonStyle     `json:"style"`
	Label    string          `json:"label,omitempty"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`
	CustomID string          `json:"custom_id,omitempty"`
	SkuID    snowflake.ID    `json:"sku_id,omitempty"`
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

func (c ButtonComponent) GetID() int {
	return c.ID
}

func (c ButtonComponent) GetCustomID() string {
	return c.CustomID
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

// WithSkuID returns a new ButtonComponent with the provided skuID
func (c ButtonComponent) WithSkuID(skuID snowflake.ID) ButtonComponent {
	c.SkuID = skuID
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
	_ Component = (*TextInputComponent)(nil)
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
	ID          int            `json:"id,omitempty"`
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

func (c TextInputComponent) GetID() int {
	return c.ID
}

func (c TextInputComponent) GetCustomID() string {
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
	TextInputStyleShort TextInputStyle = iota + 1
	TextInputStyleParagraph
)

type UnfurledMediaItem struct {
	// URL supports arbitrary urls and attachment://<filename> references
	URL string `json:"url"`
}

var (
	_ Component       = (*SectionComponent)(nil)
	_ LayoutComponent = (*SectionComponent)(nil)
)

type SectionComponent struct {
	ID         int         `json:"id,omitempty"`
	Components []Component `json:"components"`
	Accessory  Component   `json:"accessory"`
}

func (c SectionComponent) MarshalJSON() ([]byte, error) {
	type sectionComponent SectionComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		sectionComponent
	}{
		Type:             c.Type(),
		sectionComponent: sectionComponent(c),
	})
}

func (SectionComponent) Type() ComponentType {
	return ComponentTypeSection
}

func (c SectionComponent) GetID() int {
	return c.ID
}

func (c SectionComponent) GetComponents() []Component {
	return c.Components
}

func (SectionComponent) component()       {}
func (SectionComponent) layoutComponent() {}

var (
	_ Component = (*TextDisplayComponent)(nil)
)

type TextDisplayComponent struct {
	ID      int    `json:"id,omitempty"`
	Content string `json:"content"`
}

func (c TextDisplayComponent) MarshalJSON() ([]byte, error) {
	type textDisplayComponent TextDisplayComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		textDisplayComponent
	}{
		Type:                 c.Type(),
		textDisplayComponent: textDisplayComponent(c),
	})
}

func (TextDisplayComponent) Type() ComponentType {
	return ComponentTypeTextDisplay
}

func (c TextDisplayComponent) GetID() int {
	return c.ID
}

func (TextDisplayComponent) component() {}

var (
	_ Component = (*ThumbnailComponent)(nil)
)

type ThumbnailComponent struct {
	ID          int               `json:"id,omitempty"`
	Media       UnfurledMediaItem `json:"media"`
	Description string            `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

func (c ThumbnailComponent) MarshalJSON() ([]byte, error) {
	type thumbnailComponent ThumbnailComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		thumbnailComponent
	}{
		Type:               c.Type(),
		thumbnailComponent: thumbnailComponent(c),
	})
}

func (ThumbnailComponent) Type() ComponentType {
	return ComponentTypeThumbnail
}

func (c ThumbnailComponent) GetID() int {
	return c.ID
}

func (ThumbnailComponent) component() {}

type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description string            `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

var (
	_ Component = (*MediaGalleryComponent)(nil)
)

type MediaGalleryComponent struct {
	ID    int                `json:"id,omitempty"`
	Items []MediaGalleryItem `json:"items"`
}

func (c MediaGalleryComponent) MarshalJSON() ([]byte, error) {
	type mediaGalleryComponent MediaGalleryComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		mediaGalleryComponent
	}{
		Type:                  c.Type(),
		mediaGalleryComponent: mediaGalleryComponent(c),
	})
}

func (MediaGalleryComponent) Type() ComponentType {
	return ComponentTypeMediaGallery
}

func (c MediaGalleryComponent) GetID() int {
	return c.ID
}

func (MediaGalleryComponent) component() {}

type SeparatorSpacingSize int

const (
	SeparatorSpacingSizeNone SeparatorSpacingSize = iota
	SeparatorSpacingSizeSmall
	SeparatorSpacingSizeLarge
)

var (
	_ Component = (*SeparatorComponent)(nil)
)

type SeparatorComponent struct {
	ID      int                  `json:"id,omitempty"`
	Divider bool                 `json:"divider,omitempty"`
	Spacing SeparatorSpacingSize `json:"spacing,omitempty"`
}

func (c SeparatorComponent) MarshalJSON() ([]byte, error) {
	type separatorComponent SeparatorComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		separatorComponent
	}{
		Type:               c.Type(),
		separatorComponent: separatorComponent(c),
	})
}

func (SeparatorComponent) Type() ComponentType {
	return ComponentTypeSeparator
}

func (c SeparatorComponent) GetID() int {
	return c.ID
}

func (SeparatorComponent) component() {}

var (
	_ Component = (*FileComponent)(nil)
)

type FileComponent struct {
	ID int `json:"id,omitempty"`
	// File only supports attachment://<filename> references
	File    UnfurledMediaItem `json:"file"`
	Spoiler bool              `json:"spoiler,omitempty"`
}

func (c FileComponent) MarshalJSON() ([]byte, error) {
	type fileComponent FileComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		fileComponent
	}{
		Type:          c.Type(),
		fileComponent: fileComponent(c),
	})
}

func (FileComponent) Type() ComponentType {
	return ComponentTypeFile
}

func (c FileComponent) GetID() int {
	return c.ID
}

func (FileComponent) component() {}

var (
	_ Component       = (*ContainerComponent)(nil)
	_ LayoutComponent = (*ContainerComponent)(nil)
)

type ContainerComponent struct {
	ID          int         `json:"id,omitempty"`
	AccentColor *int        `json:"accent_color,omitempty"`
	Spoiler     bool        `json:"spoiler,omitempty"`
	Components  []Component `json:"components"`
}

func (c ContainerComponent) MarshalJSON() ([]byte, error) {
	type containerComponent ContainerComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		containerComponent
	}{
		Type:               c.Type(),
		containerComponent: containerComponent(c),
	})
}

func (ContainerComponent) Type() ComponentType {
	return ComponentTypeContainer
}

func (c ContainerComponent) GetID() int {
	return c.ID
}

func (c ContainerComponent) GetComponents() []Component {
	return c.Components
}

func (ContainerComponent) component()       {}
func (ContainerComponent) layoutComponent() {}

var (
	_ Component            = (*UnknownComponent)(nil)
	_ InteractiveComponent = (*UnknownComponent)(nil)
	_ LayoutComponent      = (*UnknownComponent)(nil)
	_ SelectMenuComponent  = (*UnknownComponent)(nil)
)

type UnknownComponent struct {
	ComponentType ComponentType
	ID            int
	Data          json.RawMessage
}

func (c UnknownComponent) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		Type ComponentType `json:"type"`
		ID   int           `json:"id,omitempty"`
	}{
		Type: c.ComponentType,
		ID:   c.ID,
	})
	if err != nil {
		return nil, err
	}

	return json.Merge(c.Data, data)
}

func (c *UnknownComponent) UnmarshalJSON(data []byte) error {
	var unknownComponent struct {
		Type ComponentType `json:"type"`
		ID   int           `json:"id,omitempty"`
	}
	if err := json.Unmarshal(data, &unknownComponent); err != nil {
		return err
	}

	c.ComponentType = unknownComponent.Type
	c.ID = unknownComponent.ID
	c.Data = data
	return nil
}

func (c UnknownComponent) Type() ComponentType {
	return c.ComponentType
}

func (c UnknownComponent) GetID() int {
	return c.ID
}

func (c UnknownComponent) GetCustomID() string {
	var data struct {
		CustomID string `json:"custom_id"`
	}
	if err := json.Unmarshal(c.Data, &data); err != nil {
		return ""
	}

	return data.CustomID
}

func (c UnknownComponent) GetComponents() []Component {
	var data struct {
		Components []UnmarshalComponent `json:"components"`
	}
	if err := json.Unmarshal(c.Data, &data); err != nil {
		return nil
	}

	var components []Component
	for _, component := range data.Components {
		components = append(components, component.Component)
	}
	return components
}

func (UnknownComponent) component()            {}
func (UnknownComponent) interactiveComponent() {}
func (UnknownComponent) layoutComponent()      {}
func (UnknownComponent) selectMenuComponent()  {}
