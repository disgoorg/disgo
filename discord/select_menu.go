package discord

import "github.com/disgoorg/json"

type SelectMenuComponent interface {
	InteractiveComponent
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

// WithCustomID returns a new StringSelectMenuComponent with the provided customID
func (c StringSelectMenuComponent) WithCustomID(customID string) StringSelectMenuComponent {
	c.CustomID = customID
	return c
}

// WithPlaceholder returns a new StringSelectMenuComponent with the provided placeholder
func (c StringSelectMenuComponent) WithPlaceholder(placeholder string) StringSelectMenuComponent {
	c.Placeholder = placeholder
	return c
}

// WithMinValues returns a new StringSelectMenuComponent with the provided minValue
func (c StringSelectMenuComponent) WithMinValues(minValue int) StringSelectMenuComponent {
	c.MinValues = &minValue
	return c
}

// WithMaxValues returns a new StringSelectMenuComponent with the provided maxValue
func (c StringSelectMenuComponent) WithMaxValues(maxValue int) StringSelectMenuComponent {
	c.MaxValues = maxValue
	return c
}

// AsEnabled returns a new StringSelectMenuComponent but enabled
func (c StringSelectMenuComponent) AsEnabled() StringSelectMenuComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new StringSelectMenuComponent but disabled
func (c StringSelectMenuComponent) AsDisabled() StringSelectMenuComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new StringSelectMenuComponent with the provided disabled
func (c StringSelectMenuComponent) WithDisabled(disabled bool) StringSelectMenuComponent {
	c.Disabled = disabled
	return c
}

// SetOptions returns a new StringSelectMenuComponent with the provided StringSelectMenuOption(s)
func (c StringSelectMenuComponent) SetOptions(options ...StringSelectMenuOption) StringSelectMenuComponent {
	c.Options = options
	return c
}

// SetOption returns a new StringSelectMenuComponent with the StringSelectMenuOption which has the value replaced
func (c StringSelectMenuComponent) SetOption(value string, option StringSelectMenuOption) StringSelectMenuComponent {
	for i, o := range c.Options {
		if o.Value == value {
			c.Options[i] = option
			break
		}
	}
	return c
}

// AddOptions returns a new StringSelectMenuComponent with the provided StringSelectMenuOption(s) added
func (c StringSelectMenuComponent) AddOptions(options ...StringSelectMenuOption) StringSelectMenuComponent {
	c.Options = append(c.Options, options...)
	return c
}

// RemoveOption returns a new StringSelectMenuComponent with the provided StringSelectMenuOption at the index removed
func (c StringSelectMenuComponent) RemoveOption(index int) StringSelectMenuComponent {
	if len(c.Options) > index {
		c.Options = append(c.Options[:index], c.Options[index+1:]...)
	}
	return c
}

// NewStringSelectMenuOption builds a new StringSelectMenuOption
func NewStringSelectMenuOption(label string, value string) StringSelectMenuOption {
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

// WithLabel returns a new StringSelectMenuOption with the provided label
func (o StringSelectMenuOption) WithLabel(label string) StringSelectMenuOption {
	o.Label = label
	return o
}

// WithValue returns a new StringSelectMenuOption with the provided value
func (o StringSelectMenuOption) WithValue(value string) StringSelectMenuOption {
	o.Value = value
	return o
}

// WithDescription returns a new StringSelectMenuOption with the provided description
func (o StringSelectMenuOption) WithDescription(description string) StringSelectMenuOption {
	o.Description = description
	return o
}

// WithEmoji returns a new StringSelectMenuOption with the provided Emoji
func (o StringSelectMenuOption) WithEmoji(emoji ComponentEmoji) StringSelectMenuOption {
	o.Emoji = &emoji
	return o
}

// WithDefault returns a new StringSelectMenuOption as default/non-default
func (o StringSelectMenuOption) WithDefault(defaultOption bool) StringSelectMenuOption {
	o.Default = defaultOption
	return o
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

// WithCustomID returns a new UserSelectMenuComponent with the provided customID
func (c UserSelectMenuComponent) WithCustomID(customID string) UserSelectMenuComponent {
	c.CustomID = customID
	return c
}

// WithPlaceholder returns a new UserSelectMenuComponent with the provided placeholder
func (c UserSelectMenuComponent) WithPlaceholder(placeholder string) UserSelectMenuComponent {
	c.Placeholder = placeholder
	return c
}

// WithMinValues returns a new UserSelectMenuComponent with the provided minValue
func (c UserSelectMenuComponent) WithMinValues(minValue int) UserSelectMenuComponent {
	c.MinValues = &minValue
	return c
}

// WithMaxValues returns a new UserSelectMenuComponent with the provided maxValue
func (c UserSelectMenuComponent) WithMaxValues(maxValue int) UserSelectMenuComponent {
	c.MaxValues = maxValue
	return c
}

// AsEnabled returns a new UserSelectMenuComponent but enabled
func (c UserSelectMenuComponent) AsEnabled() UserSelectMenuComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new UserSelectMenuComponent but disabled
func (c UserSelectMenuComponent) AsDisabled() UserSelectMenuComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new UserSelectMenuComponent with the provided disabled
func (c UserSelectMenuComponent) WithDisabled(disabled bool) UserSelectMenuComponent {
	c.Disabled = disabled
	return c
}

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

// WithCustomID returns a new RoleSelectMenuComponent with the provided customID
func (c RoleSelectMenuComponent) WithCustomID(customID string) RoleSelectMenuComponent {
	c.CustomID = customID
	return c
}

// WithPlaceholder returns a new RoleSelectMenuComponent with the provided placeholder
func (c RoleSelectMenuComponent) WithPlaceholder(placeholder string) RoleSelectMenuComponent {
	c.Placeholder = placeholder
	return c
}

// WithMinValues returns a new RoleSelectMenuComponent with the provided minValue
func (c RoleSelectMenuComponent) WithMinValues(minValue int) RoleSelectMenuComponent {
	c.MinValues = &minValue
	return c
}

// WithMaxValues returns a new RoleSelectMenuComponent with the provided maxValue
func (c RoleSelectMenuComponent) WithMaxValues(maxValue int) RoleSelectMenuComponent {
	c.MaxValues = maxValue
	return c
}

// AsEnabled returns a new RoleSelectMenuComponent but enabled
func (c RoleSelectMenuComponent) AsEnabled() RoleSelectMenuComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new RoleSelectMenuComponent but disabled
func (c RoleSelectMenuComponent) AsDisabled() RoleSelectMenuComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new RoleSelectMenuComponent with the provided disabled
func (c RoleSelectMenuComponent) WithDisabled(disabled bool) RoleSelectMenuComponent {
	c.Disabled = disabled
	return c
}

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

// WithCustomID returns a new MentionableSelectMenuComponent with the provided customID
func (c MentionableSelectMenuComponent) WithCustomID(customID string) MentionableSelectMenuComponent {
	c.CustomID = customID
	return c
}

// WithPlaceholder returns a new MentionableSelectMenuComponent with the provided placeholder
func (c MentionableSelectMenuComponent) WithPlaceholder(placeholder string) MentionableSelectMenuComponent {
	c.Placeholder = placeholder
	return c
}

// WithMinValues returns a new MentionableSelectMenuComponent with the provided minValue
func (c MentionableSelectMenuComponent) WithMinValues(minValue int) MentionableSelectMenuComponent {
	c.MinValues = &minValue
	return c
}

// WithMaxValues returns a new MentionableSelectMenuComponent with the provided maxValue
func (c MentionableSelectMenuComponent) WithMaxValues(maxValue int) MentionableSelectMenuComponent {
	c.MaxValues = maxValue
	return c
}

// AsEnabled returns a new MentionableSelectMenuComponent but enabled
func (c MentionableSelectMenuComponent) AsEnabled() MentionableSelectMenuComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new MentionableSelectMenuComponent but disabled
func (c MentionableSelectMenuComponent) AsDisabled() MentionableSelectMenuComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new MentionableSelectMenuComponent with the provided disabled
func (c MentionableSelectMenuComponent) WithDisabled(disabled bool) MentionableSelectMenuComponent {
	c.Disabled = disabled
	return c
}

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

// WithCustomID returns a new ChannelSelectMenuComponent with the provided customID
func (c ChannelSelectMenuComponent) WithCustomID(customID string) ChannelSelectMenuComponent {
	c.CustomID = customID
	return c
}

// WithPlaceholder returns a new ChannelSelectMenuComponent with the provided placeholder
func (c ChannelSelectMenuComponent) WithPlaceholder(placeholder string) ChannelSelectMenuComponent {
	c.Placeholder = placeholder
	return c
}

// WithMinValues returns a new ChannelSelectMenuComponent with the provided minValue
func (c ChannelSelectMenuComponent) WithMinValues(minValue int) ChannelSelectMenuComponent {
	c.MinValues = &minValue
	return c
}

// WithMaxValues returns a new ChannelSelectMenuComponent with the provided maxValue
func (c ChannelSelectMenuComponent) WithMaxValues(maxValue int) ChannelSelectMenuComponent {
	c.MaxValues = maxValue
	return c
}

// AsEnabled returns a new ChannelSelectMenuComponent but enabled
func (c ChannelSelectMenuComponent) AsEnabled() ChannelSelectMenuComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new ChannelSelectMenuComponent but disabled
func (c ChannelSelectMenuComponent) AsDisabled() ChannelSelectMenuComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new ChannelSelectMenuComponent with the provided disabled
func (c ChannelSelectMenuComponent) WithDisabled(disabled bool) ChannelSelectMenuComponent {
	c.Disabled = disabled
	return c
}

// WithChannelTypes returns a new ChannelSelectMenuComponent with the provided channelTypes
func (c ChannelSelectMenuComponent) WithChannelTypes(channelTypes ...ComponentType) ChannelSelectMenuComponent {
	c.ChannelTypes = channelTypes
	return c
}
