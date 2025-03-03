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

// SetLayoutComponents sets the discord.LayoutComponent(s) of the ModalCreate
func (b *ModalCreateBuilder) SetLayoutComponents(LayoutComponents ...LayoutComponent) *ModalCreateBuilder {
	b.Components = LayoutComponents
	return b
}

// SetLayoutComponent sets the provided discord.InteractiveComponent at the index of discord.InteractiveComponent(s)
func (b *ModalCreateBuilder) SetLayoutComponent(i int, container LayoutComponent) *ModalCreateBuilder {
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

// AddLayoutComponents adds the discord.LayoutComponent(s) to the ModalCreate
func (b *ModalCreateBuilder) AddLayoutComponents(containers ...LayoutComponent) *ModalCreateBuilder {
	b.Components = append(b.Components, containers...)
	return b
}

// RemoveLayoutComponent removes a discord.ActionRowComponent from the ModalCreate
func (b *ModalCreateBuilder) RemoveLayoutComponent(i int) *ModalCreateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// ClearLayoutComponents removes all the discord.LayoutComponent(s) of the ModalCreate
func (b *ModalCreateBuilder) ClearLayoutComponents() *ModalCreateBuilder {
	b.Components = []LayoutComponent{}
	return b
}

// Build builds the ModalCreateBuilder to a ModalCreate struct
func (b *ModalCreateBuilder) Build() ModalCreate {
	return b.ModalCreate
}
