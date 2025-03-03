package discord

import (
	"fmt"
	"iter"

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

// Component is an interface for all components.
// [ActionRowComponent]
// [ButtonComponent]
// [StringSelectMenuComponent]
// [TextInputComponent]
// [UserSelectMenuComponent]
// [RoleSelectMenuComponent]
// [MentionableSelectMenuComponent]
// [ChannelSelectMenuComponent]
// [SectionComponent]
// [TextDisplayComponent]
// [ThumbnailComponent]
// [MediaGalleryComponent]
// [FileComponent]
// [SeparatorComponent]
// [ContainerComponent]
// [UnknownComponent]
type Component interface {
	json.Marshaler
	// Type returns the ComponentType of the Component.
	Type() ComponentType
	// GetID returns the id of the Component. This is used to uniquely identify a Component in a [Message] and needs to be unique.
	GetID() int
	// component is a marker to simulate unions.
	component()
}

// ComponentIter is an optional interface a Component can implement to return an iterator over its children.
type ComponentIter interface {
	Children() iter.Seq[Component]
}

// InteractiveComponent is an interface for all components that can be present in an [ActionRowComponent].
// [ButtonComponent]
// [StringSelectMenuComponent]
// [TextInputComponent] (currently only supported in modals)
// [UserSelectMenuComponent]
// [RoleSelectMenuComponent]
// [MentionableSelectMenuComponent]
// [ChannelSelectMenuComponent]
// [ButtonComponent]
// [SelectMenuComponent]
// [UnknownComponent]
type InteractiveComponent interface {
	Component
	// GetCustomID returns the customID of the Component. This can be used to identify or transport data with the Component.
	GetCustomID() string
	// interactiveComponent is a marker to simulate unions.
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
// [UnknownComponent]
type LayoutComponent interface {
	Component
	layoutComponent()
}

// SectionSubComponent is an interface for all components that can be present in a [SectionComponent].
// [TextDisplayComponent]
// [UnknownComponent]
type SectionSubComponent interface {
	Component
	// sectionSubComponent is a marker to simulate unions.
	sectionSubComponent()
}

// SectionAccessoryComponent is an interface for all components that can be present as an accessory in [SectionComponent.Accessory].
// [ThumbnailComponent]
// [UnknownComponent]
type SectionAccessoryComponent interface {
	Component
	// sectionAccessoryComponent is a marker to simulate unions.
	sectionAccessoryComponent()
}

// ContainerSubComponent is an interface for all components that can be present in a [ContainerComponent].
// [ActionRowComponent]
// [SectionComponent]
// [TextDisplayComponent]
// [MediaGalleryComponent]
// [FileComponent]
// [SeparatorComponent]
// [UnknownComponent]
type ContainerSubComponent interface {
	Component
	// containerSubComponent is a marker to simulate unions.
	containerSubComponent()
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
	_ Component             = (*ActionRowComponent)(nil)
	_ LayoutComponent       = (*ActionRowComponent)(nil)
	_ ContainerSubComponent = (*ActionRowComponent)(nil)
	_ ComponentIter         = (*ActionRowComponent)(nil)
)

func NewActionRow(components ...InteractiveComponent) ActionRowComponent {
	return ActionRowComponent{
		Components: components,
	}
}

type ActionRowComponent struct {
	ID         int                    `json:"id,omitempty"`
	Components []InteractiveComponent `json:"components"`
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

func (c *ActionRowComponent) UnmarshalJSON(data []byte) error {
	var actionRowComponent struct {
		ID         int                  `json:"id,omitempty"`
		Components []UnmarshalComponent `json:"components"`
	}
	if err := json.Unmarshal(data, &actionRowComponent); err != nil {
		return err
	}

	c.ID = actionRowComponent.ID
	components := make([]InteractiveComponent, 0, len(actionRowComponent.Components))
	for _, component := range actionRowComponent.Components {
		components = append(components, component.Component.(InteractiveComponent))
	}
	c.Components = components
	return nil
}

func (ActionRowComponent) Type() ComponentType {
	return ComponentTypeActionRow
}

func (c ActionRowComponent) GetID() int {
	return c.ID
}

func (ActionRowComponent) component()             {}
func (ActionRowComponent) layoutComponent()       {}
func (ActionRowComponent) containerSubComponent() {}

func (c ActionRowComponent) Children() iter.Seq[Component] {
	return func(yield func(Component) bool) {
		for _, cc := range c.Components {
			if !yield(cc) {
				return
			}
		}
	}
}

func (c ActionRowComponent) WithID(id int) ActionRowComponent {
	c.ID = id
	return c
}

// UpdateComponent returns a new ActionRowComponent with the Component which has the id replaced with the provided Component.
func (c ActionRowComponent) UpdateComponent(id int, component InteractiveComponent) ActionRowComponent {
	for i, cc := range c.Components {
		if cc.GetID() == id {
			c.Components[i] = component
			return c
		}
	}
	return c
}

// AddComponents returns a new ActionRowComponent with the provided Component(s) added
func (c ActionRowComponent) AddComponents(components ...InteractiveComponent) ActionRowComponent {
	c.Components = append(c.Components, components...)
	return c
}

// RemoveComponent returns a new ActionRowComponent with the provided Component which has the provided id removed.
func (c ActionRowComponent) RemoveComponent(id int) ActionRowComponent {
	for i, cc := range c.Components {
		if cc.GetID() == id {
			c.Components = append(c.Components[:i], c.Components[i+1:]...)
			return c
		}
	}
	return c
}

// ButtonStyle defines how the ButtonComponent looks like. [Discord Docs]
//
// [Discord Docs]: https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
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

// NewButton creates a new [ButtonComponent] with the provided parameters. Link ButtonComponent(s) need a URL and other ButtonComponent(s) need a customID
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
	_ Component                 = (*ButtonComponent)(nil)
	_ InteractiveComponent      = (*ButtonComponent)(nil)
	_ SectionAccessoryComponent = (*ButtonComponent)(nil)
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

func (ButtonComponent) component()                 {}
func (ButtonComponent) interactiveComponent()      {}
func (ButtonComponent) sectionAccessoryComponent() {}

// WithID returns a new ButtonComponent with the provided id
func (c ButtonComponent) WithID(id int) ButtonComponent {
	c.ID = id
	return c
}

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

// NewTextInput creates a new [TextInputComponent] with the provided parameters.
func NewTextInput(customID string, style TextInputStyle, label string) TextInputComponent {
	return TextInputComponent{
		CustomID: customID,
		Style:    style,
		Label:    label,
	}
}

// NewShortTextInput creates a new [TextInputComponent] with [TextInputStyleShort] & the provided parameters
func NewShortTextInput(customID string, label string) TextInputComponent {
	return NewTextInput(customID, TextInputStyleShort, label)
}

// NewParagraphTextInput creates a new [TextInputComponent] with [TextInputStyleParagraph] & the provided parameters
func NewParagraphTextInput(customID string, label string) TextInputComponent {
	return NewTextInput(customID, TextInputStyleParagraph, label)
}

// TextInputComponent is a component that allows users to input text. [Discord Docs]
//
// [Discord Docs]: https://discord.com/developers/docs/interactions/message-components#text-inputs
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

// WithID returns a new TextInputComponent with the provided id
func (c TextInputComponent) WithID(id int) TextInputComponent {
	c.ID = id
	return c
}

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
	_ Component             = (*SectionComponent)(nil)
	_ LayoutComponent       = (*SectionComponent)(nil)
	_ ContainerSubComponent = (*SectionComponent)(nil)
)

func NewSection(components ...SectionSubComponent) SectionComponent {
	return SectionComponent{
		Components: components,
	}
}

type SectionComponent struct {
	ID         int                       `json:"id,omitempty"`
	Components []SectionSubComponent     `json:"components"`
	Accessory  SectionAccessoryComponent `json:"accessory"`
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

func (c *SectionComponent) UnmarshalJSON(data []byte) error {
	var sectionComponent struct {
		ID         int                  `json:"id,omitempty"`
		Components []UnmarshalComponent `json:"components"`
		Accessory  UnmarshalComponent   `json:"accessory"`
	}
	if err := json.Unmarshal(data, &sectionComponent); err != nil {
		return err
	}

	c.ID = sectionComponent.ID

	components := make([]SectionSubComponent, 0, len(sectionComponent.Components))
	for _, component := range sectionComponent.Components {
		components = append(components, component.Component.(SectionSubComponent))
	}
	c.Components = components

	c.Accessory = sectionComponent.Accessory.Component.(SectionAccessoryComponent)
	return nil
}

func (SectionComponent) Type() ComponentType {
	return ComponentTypeSection
}

func (c SectionComponent) GetID() int {
	return c.ID
}

func (SectionComponent) component()             {}
func (SectionComponent) layoutComponent()       {}
func (SectionComponent) containerSubComponent() {}

func (c SectionComponent) Children() iter.Seq[Component] {
	return func(yield func(Component) bool) {
		for _, cc := range c.Components {
			if !yield(cc) {
				return
			}
		}

		if c.Accessory != nil {
			if !yield(c.Accessory) {
				return
			}
		}
	}
}

func (c SectionComponent) WithID(id int) SectionComponent {
	c.ID = id
	return c
}

var (
	_ Component             = (*TextDisplayComponent)(nil)
	_ ContainerSubComponent = (*TextDisplayComponent)(nil)
	_ SectionSubComponent   = (*TextDisplayComponent)(nil)
)

func NewTextDisplay(content string) TextDisplayComponent {
	return TextDisplayComponent{
		Content: content,
	}
}

// TextDisplayComponent is a component that displays text.
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

func (TextDisplayComponent) component()             {}
func (TextDisplayComponent) sectionSubComponent()   {}
func (TextDisplayComponent) containerSubComponent() {}

func (c TextDisplayComponent) WithID(i int) SectionSubComponent {
	c.ID = i
	return c
}

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

func (c ThumbnailComponent) WithID(id int) ThumbnailComponent {
	c.ID = id
	return c
}

type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description string            `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler,omitempty"`
}

var (
	_ Component             = (*MediaGalleryComponent)(nil)
	_ LayoutComponent       = (*MediaGalleryComponent)(nil)
	_ ContainerSubComponent = (*MediaGalleryComponent)(nil)
)

func NewMediaGallery(items ...MediaGalleryItem) MediaGalleryComponent {
	return MediaGalleryComponent{
		Items: items,
	}
}

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

func (MediaGalleryComponent) component()             {}
func (MediaGalleryComponent) layoutComponent()       {}
func (MediaGalleryComponent) containerSubComponent() {}

func (c MediaGalleryComponent) WithID(id int) MediaGalleryComponent {
	c.ID = id
	return c
}

type SeparatorSpacingSize int

const (
	SeparatorSpacingSizeNone SeparatorSpacingSize = iota
	SeparatorSpacingSizeSmall
	SeparatorSpacingSizeLarge
)

var (
	_ Component             = (*SeparatorComponent)(nil)
	_ LayoutComponent       = (*MediaGalleryComponent)(nil)
	_ ContainerSubComponent = (*MediaGalleryComponent)(nil)
)

func NewSeparator(spacing SeparatorSpacingSize) SeparatorComponent {
	return SeparatorComponent{
		Spacing: spacing,
	}
}

func NewSmallSeparator() SeparatorComponent {
	return NewSeparator(SeparatorSpacingSizeSmall)
}

func NewLargeSeparator() SeparatorComponent {
	return NewSeparator(SeparatorSpacingSizeLarge)
}

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

func (SeparatorComponent) component()             {}
func (SeparatorComponent) layoutComponent()       {}
func (SeparatorComponent) containerSubComponent() {}

func (c SeparatorComponent) WithID(i int) LayoutComponent {
	c.ID = i
	return c
}

var (
	_ Component = (*FileComponent)(nil)
)

func NewFileComponent(url string, spoiler bool) FileComponent {
	return FileComponent{
		File: UnfurledMediaItem{
			URL: url,
		},
		Spoiler: spoiler,
	}
}

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

func (c FileComponent) WithID(id int) FileComponent {
	c.ID = id
	return c
}

func (c FileComponent) WithSpoiler(spoiler bool) FileComponent {
	c.Spoiler = spoiler
	return c
}

var (
	_ Component       = (*ContainerComponent)(nil)
	_ LayoutComponent = (*ContainerComponent)(nil)
)

func NewContainer(components ...ContainerSubComponent) ContainerComponent {
	return ContainerComponent{
		Components: components,
	}
}

type ContainerComponent struct {
	ID          int                     `json:"id,omitempty"`
	AccentColor *int                    `json:"accent_color,omitempty"`
	Spoiler     bool                    `json:"spoiler,omitempty"`
	Components  []ContainerSubComponent `json:"components"`
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

func (c *ContainerComponent) UnmarshalJSON(data []byte) error {
	var containerComponent struct {
		ID          int                  `json:"id,omitempty"`
		AccentColor *int                 `json:"accent_color,omitempty"`
		Spoiler     bool                 `json:"spoiler,omitempty"`
		Components  []UnmarshalComponent `json:"components"`
	}
	if err := json.Unmarshal(data, &containerComponent); err != nil {
		return err
	}

	c.ID = containerComponent.ID
	c.AccentColor = containerComponent.AccentColor
	c.Spoiler = containerComponent.Spoiler

	components := make([]ContainerSubComponent, 0, len(containerComponent.Components))
	for _, component := range containerComponent.Components {
		components = append(components, component.Component.(ContainerSubComponent))
	}
	c.Components = components
	return nil
}

func (ContainerComponent) Type() ComponentType {
	return ComponentTypeContainer
}

func (c ContainerComponent) GetID() int {
	return c.ID
}

func (ContainerComponent) component()       {}
func (ContainerComponent) layoutComponent() {}

func (c ContainerComponent) WithID(id int) ContainerComponent {
	c.ID = id
	return c
}

func (c ContainerComponent) Children() iter.Seq[Component] {
	return func(yield func(Component) bool) {
		for _, cc := range c.Components {
			if !yield(cc) {
				return
			}
			if ic, ok := cc.(ComponentIter); ok {
				for cc := range ic.Children() {
					if !yield(cc) {
						return
					}
				}
			}
		}
	}
}

var (
	_ Component                 = (*UnknownComponent)(nil)
	_ InteractiveComponent      = (*UnknownComponent)(nil)
	_ LayoutComponent           = (*UnknownComponent)(nil)
	_ SelectMenuComponent       = (*UnknownComponent)(nil)
	_ SectionSubComponent       = (*UnknownComponent)(nil)
	_ SectionAccessoryComponent = (*UnknownComponent)(nil)
	_ ContainerSubComponent     = (*UnknownComponent)(nil)
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

func (UnknownComponent) component()                 {}
func (UnknownComponent) interactiveComponent()      {}
func (UnknownComponent) layoutComponent()           {}
func (UnknownComponent) selectMenuComponent()       {}
func (UnknownComponent) containerSubComponent()     {}
func (UnknownComponent) sectionSubComponent()       {}
func (UnknownComponent) sectionAccessoryComponent() {}
