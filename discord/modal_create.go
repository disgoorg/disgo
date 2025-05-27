package discord

var _ InteractionResponseData = (*ModalCreate)(nil)

type ModalCreate struct {
	CustomID   string            `json:"custom_id"`
	Title      string            `json:"title"`
	Components []LayoutComponent `json:"components"`
}

func (ModalCreate) interactionCallbackData() {}

// NewModalCreateBuilder creates a new ModalCreateBuilder to be built later
func NewModalCreateBuilder() *ModalCreateBuilder {
	return &ModalCreateBuilder{}
}

type ModalCreateBuilder struct {
	ModalCreate
}

// SetCustomID sets the CustomID of the ModalCreate
func (b *ModalCreateBuilder) SetCustomID(customID string) *ModalCreateBuilder {
	b.CustomID = customID
	return b
}

// SetTitle sets the title of the ModalCreate
func (b *ModalCreateBuilder) SetTitle(title string) *ModalCreateBuilder {
	b.Title = title
	return b
}

// SetComponents sets the discord.LayoutComponent(s) of the ModalCreate
func (b *ModalCreateBuilder) SetComponents(components ...LayoutComponent) *ModalCreateBuilder {
	b.Components = components
	return b
}

// SetComponent sets the provided discord.LayoutComponent at the index of discord.LayoutComponent(s)
func (b *ModalCreateBuilder) SetComponent(i int, container LayoutComponent) *ModalCreateBuilder {
	if len(b.Components) > i {
		b.Components[i] = container
	}
	return b
}

// AddActionRow adds a new discord.ActionRowComponent with the provided discord.InteractiveComponent(s) to the ModalCreate
func (b *ModalCreateBuilder) AddActionRow(components ...InteractiveComponent) *ModalCreateBuilder {
	b.Components = append(b.Components, ActionRowComponent{Components: components})
	return b
}

// AddComponents adds the discord.LayoutComponent(s) to the ModalCreate
func (b *ModalCreateBuilder) AddComponents(containers ...LayoutComponent) *ModalCreateBuilder {
	b.Components = append(b.Components, containers...)
	return b
}

// RemoveComponent removes a discord.LayoutComponent from the ModalCreate
func (b *ModalCreateBuilder) RemoveComponent(i int) *ModalCreateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// ClearComponents removes all the discord.LayoutComponent(s) of the ModalCreate
func (b *ModalCreateBuilder) ClearComponents() *ModalCreateBuilder {
	b.Components = []LayoutComponent{}
	return b
}

// Build builds the ModalCreateBuilder to a ModalCreate struct
func (b *ModalCreateBuilder) Build() ModalCreate {
	return b.ModalCreate
}
