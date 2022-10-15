package discord

import "github.com/disgoorg/disgo/json"

type SelectMenuComponent interface {
	InteractiveComponent
	// Placeholder() string
	// MinValues() *int
	// MaxValues() int
	// Disabled() bool
	selectMenu()
}

var (
	_ Component            = (*StringSelectMenuComponent)(nil)
	_ InteractiveComponent = (*StringSelectMenuComponent)(nil)
	_ SelectMenuComponent  = (*StringSelectMenuComponent)(nil)
)

// NewStringSelectMenu builds a new SelectMenuComponent from the provided values
func NewStringSelectMenu(customID string, placeholder string, options ...StringSelectMenuOption) StringSelectMenuComponent {
	return StringSelectMenuComponent{
		CustomID:    customID,
		Placeholder: placeholder,
		Options:     options,
	}
}

type StringSelectMenuComponent struct {
	CustomID    string                   `json:"custom_id"`
	Placeholder string                   `json:"placeholder,omitempty"`
	MinValues   *int                     `json:"min_values,omitempty"`
	MaxValues   int                      `json:"max_values,omitempty"`
	Disabled    bool                     `json:"disabled,omitempty"`
	Options     []StringSelectMenuOption `json:"options,omitempty"`
}

func (c StringSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type component StringSelectMenuComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		component
	}{
		Type:      c.Type(),
		component: component(c),
	})
}

func (StringSelectMenuComponent) Type() ComponentType {
	return ComponentTypeStringSelectMenu
}

func (c StringSelectMenuComponent) ID() string {
	return c.CustomID
}

func (StringSelectMenuComponent) component()            {}
func (StringSelectMenuComponent) interactiveComponent() {}
func (StringSelectMenuComponent) selectMenu()           {}

// NewSelectMenuOption builds a new SelectMenuOption
func NewSelectMenuOption(label string, value string) StringSelectMenuOption {
	return StringSelectMenuOption{
		Label: label,
		Value: value,
	}
}

// StringSelectMenuOption represents an option in a StringSelectMenuComponent
type StringSelectMenuOption struct {
	Label       string          `json:"label"`
	Value       string          `json:"value"`
	Description string          `json:"description,omitempty"`
	Emoji       *ComponentEmoji `json:"emoji,omitempty"`
	Default     bool            `json:"default,omitempty"`
}

var (
	_ Component            = (*UserSelectMenuComponent)(nil)
	_ InteractiveComponent = (*UserSelectMenuComponent)(nil)
	_ SelectMenuComponent  = (*UserSelectMenuComponent)(nil)
)

// NewUserSelectMenu builds a new SelectMenuComponent from the provided values
func NewUserSelectMenu(customID string, placeholder string) UserSelectMenuComponent {
	return UserSelectMenuComponent{
		CustomID:    customID,
		Placeholder: placeholder,
	}
}

type UserSelectMenuComponent struct {
	CustomID    string `json:"custom_id"`
	Placeholder string `json:"placeholder,omitempty"`
	MinValues   *int   `json:"min_values,omitempty"`
	MaxValues   int    `json:"max_values,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

func (c UserSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type component UserSelectMenuComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		component
	}{
		Type:      c.Type(),
		component: component(c),
	})
}

func (UserSelectMenuComponent) Type() ComponentType {
	return ComponentTypeUserSelectMenu
}

func (c UserSelectMenuComponent) ID() string {
	return c.CustomID
}

func (UserSelectMenuComponent) component()            {}
func (UserSelectMenuComponent) interactiveComponent() {}
func (UserSelectMenuComponent) selectMenu()           {}

var (
	_ Component            = (*UserSelectMenuComponent)(nil)
	_ InteractiveComponent = (*UserSelectMenuComponent)(nil)
	_ SelectMenuComponent  = (*UserSelectMenuComponent)(nil)
)

// NewRoleSelectMenu builds a new SelectMenuComponent from the provided values
func NewRoleSelectMenu(customID string, placeholder string) RoleSelectMenuComponent {
	return RoleSelectMenuComponent{
		CustomID:    customID,
		Placeholder: placeholder,
	}
}

type RoleSelectMenuComponent struct {
	CustomID    string `json:"custom_id"`
	Placeholder string `json:"placeholder,omitempty"`
	MinValues   *int   `json:"min_values,omitempty"`
	MaxValues   int    `json:"max_values,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

func (c RoleSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type component RoleSelectMenuComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		component
	}{
		Type:      c.Type(),
		component: component(c),
	})
}

func (RoleSelectMenuComponent) Type() ComponentType {
	return ComponentTypeRoleSelectMenu
}

func (c RoleSelectMenuComponent) ID() string {
	return c.CustomID
}

func (RoleSelectMenuComponent) component()            {}
func (RoleSelectMenuComponent) interactiveComponent() {}
func (RoleSelectMenuComponent) selectMenu()           {}

var (
	_ Component            = (*MentionableSelectMenuComponent)(nil)
	_ InteractiveComponent = (*MentionableSelectMenuComponent)(nil)
	_ SelectMenuComponent  = (*MentionableSelectMenuComponent)(nil)
)

// NewMentionableSelectMenu builds a new SelectMenuComponent from the provided values
func NewMentionableSelectMenu(customID string, placeholder string) MentionableSelectMenuComponent {
	return MentionableSelectMenuComponent{
		CustomID:    customID,
		Placeholder: placeholder,
	}
}

type MentionableSelectMenuComponent struct {
	CustomID    string `json:"custom_id"`
	Placeholder string `json:"placeholder,omitempty"`
	MinValues   *int   `json:"min_values,omitempty"`
	MaxValues   int    `json:"max_values,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

func (c MentionableSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type component MentionableSelectMenuComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		component
	}{
		Type:      c.Type(),
		component: component(c),
	})
}

func (MentionableSelectMenuComponent) Type() ComponentType {
	return ComponentTypeMentionableSelectMenu
}

func (c MentionableSelectMenuComponent) ID() string {
	return c.CustomID
}

func (MentionableSelectMenuComponent) component()            {}
func (MentionableSelectMenuComponent) interactiveComponent() {}
func (MentionableSelectMenuComponent) selectMenu()           {}

var (
	_ Component            = (*ChannelSelectMenuComponent)(nil)
	_ InteractiveComponent = (*ChannelSelectMenuComponent)(nil)
	_ SelectMenuComponent  = (*ChannelSelectMenuComponent)(nil)
)

// NewChannelSelectMenu builds a new SelectMenuComponent from the provided values
func NewChannelSelectMenu(customID string, placeholder string) ChannelSelectMenuComponent {
	return ChannelSelectMenuComponent{
		CustomID:    customID,
		Placeholder: placeholder,
	}
}

type ChannelSelectMenuComponent struct {
	CustomID     string          `json:"custom_id"`
	Placeholder  string          `json:"placeholder,omitempty"`
	MinValues    *int            `json:"min_values,omitempty"`
	MaxValues    int             `json:"max_values,omitempty"`
	Disabled     bool            `json:"disabled,omitempty"`
	ChannelTypes []ComponentType `json:"channel_types,omitempty"`
}

func (c ChannelSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type component ChannelSelectMenuComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		component
	}{
		Type:      c.Type(),
		component: component(c),
	})
}

func (ChannelSelectMenuComponent) Type() ComponentType {
	return ComponentTypeChannelSelectMenu
}

func (c ChannelSelectMenuComponent) ID() string {
	return c.CustomID
}

func (ChannelSelectMenuComponent) component()            {}
func (ChannelSelectMenuComponent) interactiveComponent() {}
func (ChannelSelectMenuComponent) selectMenu()           {}
