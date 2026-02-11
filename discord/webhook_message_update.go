package discord

import (
	"fmt"
	"io"
	"slices"

	"github.com/disgoorg/snowflake/v2"
)

// NewWebhookMessageUpdate returns a new WebhookMessageUpdate with no fields set.
func NewWebhookMessageUpdate() WebhookMessageUpdate {
	return WebhookMessageUpdate{}
}

// NewWebhookMessageUpdateV2 returns a new WebhookMessageUpdate with MessageFlagIsComponentsV2 set and no other fields set.
func NewWebhookMessageUpdateV2(components []LayoutComponent) WebhookMessageUpdate {
	flags := MessageFlagIsComponentsV2
	return WebhookMessageUpdate{
		Flags:      &flags,
		Components: &components,
	}
}

// WebhookMessageUpdate is used to edit a Message.
type WebhookMessageUpdate struct {
	Content         *string             `json:"content,omitempty"`
	Embeds          *[]Embed            `json:"embeds,omitempty"`
	Components      *[]LayoutComponent  `json:"components,omitempty"`
	Attachments     *[]AttachmentUpdate `json:"attachments,omitempty"`
	Files           []*File             `json:"-"`
	AllowedMentions *AllowedMentions    `json:"allowed_mentions,omitempty"`
	Poll            *PollCreate         `json:"poll,omitempty"`
	// Flags are the MessageFlags of the message.
	// Be careful not to override the current flags when editing messages from other users - this will result in a permission error.
	// Use MessageFlags.Add for flags like discord.MessageFlagIsComponentsV2.
	Flags *MessageFlags `json:"flags,omitempty"`
}

// ToBody returns the WebhookMessageUpdate ready for body.
func (m WebhookMessageUpdate) ToBody() (any, error) {
	if len(m.Files) > 0 {
		for _, attachmentCreate := range parseAttachments(m.Files) {
			if m.Attachments == nil {
				m.Attachments = new([]AttachmentUpdate)
			}
			*m.Attachments = append(*m.Attachments, attachmentCreate)
		}
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

// WithContent returns a new WebhookMessageUpdate with the provided content.
func (m WebhookMessageUpdate) WithContent(content string) WebhookMessageUpdate {
	m.Content = &content
	return m
}

// WithContentf returns a new WebhookMessageUpdate with the formatted content.
func (m WebhookMessageUpdate) WithContentf(content string, a ...any) WebhookMessageUpdate {
	return m.WithContent(fmt.Sprintf(content, a...))
}

// ClearContent returns a new WebhookMessageUpdate with no content.
func (m WebhookMessageUpdate) ClearContent() WebhookMessageUpdate {
	return m.WithContent("")
}

// WithEmbeds returns a new WebhookMessageUpdate with the provided Embed(s).
func (m WebhookMessageUpdate) WithEmbeds(embeds ...Embed) WebhookMessageUpdate {
	m.Embeds = &embeds
	return m
}

// WithEmbed returns a new WebhookMessageUpdate with the provided Embed at the index.
func (m WebhookMessageUpdate) WithEmbed(i int, embed Embed) WebhookMessageUpdate {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	if len(*m.Embeds) > i {
		newEmbeds := slices.Clone(*m.Embeds)
		newEmbeds[i] = embed
		m.Embeds = &newEmbeds
	}
	return m
}

// AddEmbeds returns a new WebhookMessageUpdate with the provided embeds added.
func (m WebhookMessageUpdate) AddEmbeds(embeds ...Embed) WebhookMessageUpdate {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	newEmbeds := append(slices.Clone(*m.Embeds), embeds...)
	m.Embeds = &newEmbeds
	return m
}

// RemoveEmbed returns a new WebhookMessageUpdate with the embed at the index removed.
func (m WebhookMessageUpdate) RemoveEmbed(i int) WebhookMessageUpdate {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	if len(*m.Embeds) > i {
		newEmbeds := slices.Delete(slices.Clone(*m.Embeds), i, i+1)
		m.Embeds = &newEmbeds
	}
	return m
}

// ClearEmbeds returns a new WebhookMessageUpdate with no embeds.
func (m WebhookMessageUpdate) ClearEmbeds() WebhookMessageUpdate {
	m.Embeds = &[]Embed{}
	return m
}

// WithComponents returns a new WebhookMessageUpdate with the provided LayoutComponent(s).
func (m WebhookMessageUpdate) WithComponents(components ...LayoutComponent) WebhookMessageUpdate {
	m.Components = &components
	return m
}

// UpdateComponent returns a new WebhookMessageUpdate with the provided LayoutComponent at the index.
func (m WebhookMessageUpdate) UpdateComponent(id int, component LayoutComponent) WebhookMessageUpdate {
	if m.Components == nil {
		return m
	}
	for i, cc := range *m.Components {
		if cc.GetID() == id {
			newComponents := slices.Clone(*m.Components)
			newComponents[i] = component
			m.Components = &newComponents
			return m
		}
	}
	return m
}

// AddActionRow returns a new WebhookMessageUpdate with a new ActionRowComponent containing the provided InteractiveComponent(s) added.
func (m WebhookMessageUpdate) AddActionRow(components ...InteractiveComponent) WebhookMessageUpdate {
	if m.Components == nil {
		m.Components = &[]LayoutComponent{}
	}
	newComponents := append(slices.Clone(*m.Components), NewActionRow(components...))
	m.Components = &newComponents
	return m
}

// AddComponents returns a new WebhookMessageUpdate with the provided LayoutComponent(s) added.
func (m WebhookMessageUpdate) AddComponents(containers ...LayoutComponent) WebhookMessageUpdate {
	if m.Components == nil {
		m.Components = &[]LayoutComponent{}
	}
	newComponents := append(slices.Clone(*m.Components), containers...)
	m.Components = &newComponents
	return m
}

// RemoveComponent returns a new WebhookMessageUpdate with the LayoutComponent with the provided ID removed.
func (m WebhookMessageUpdate) RemoveComponent(id int) WebhookMessageUpdate {
	if m.Components == nil {
		return m
	}
	for i, cc := range *m.Components {
		if cc.GetID() == id {
			newComponents := slices.Delete(slices.Clone(*m.Components), i, i+1)
			m.Components = &newComponents
			return m
		}
	}
	return m
}

// ClearComponents returns a new WebhookMessageUpdate with no LayoutComponent(s).
func (m WebhookMessageUpdate) ClearComponents() WebhookMessageUpdate {
	m.Components = &[]LayoutComponent{}
	return m
}

// WithFiles returns a new WebhookMessageUpdate with the provided File(s).
func (m WebhookMessageUpdate) WithFiles(files ...*File) WebhookMessageUpdate {
	m.Files = files
	return m
}

// UpdateFile returns a new WebhookMessageUpdate with the File at the index.
func (m WebhookMessageUpdate) UpdateFile(i int, file *File) WebhookMessageUpdate {
	if len(m.Files) > i {
		m.Files = slices.Clone(m.Files)
		m.Files[i] = file
	}
	return m
}

// AddFiles returns a new WebhookMessageUpdate with the File(s) added.
func (m WebhookMessageUpdate) AddFiles(files ...*File) WebhookMessageUpdate {
	m.Files = append(m.Files, files...)
	return m
}

// AddFile returns a new WebhookMessageUpdate with a File added.
func (m WebhookMessageUpdate) AddFile(name string, description string, reader io.Reader, flags ...FileFlags) WebhookMessageUpdate {
	m.Files = append(m.Files, NewFile(name, description, reader, flags...))
	return m
}

// ClearFiles returns a new WebhookMessageUpdate with no File(s).
func (m WebhookMessageUpdate) ClearFiles() WebhookMessageUpdate {
	m.Files = []*File{}
	return m
}

// RemoveFile returns a new WebhookMessageUpdate with the File at the index removed.
func (m WebhookMessageUpdate) RemoveFile(i int) WebhookMessageUpdate {
	if len(m.Files) > i {
		m.Files = slices.Delete(slices.Clone(m.Files), i, i+1)
	}
	return m
}

// RetainAttachments returns a new WebhookMessageUpdate that retains the provided Attachment(s).
func (m WebhookMessageUpdate) RetainAttachments(attachments ...Attachment) WebhookMessageUpdate {
	if m.Attachments == nil {
		m.Attachments = &[]AttachmentUpdate{}
	}
	newAttachments := slices.Clone(*m.Attachments)
	for _, attachment := range attachments {
		newAttachments = append(newAttachments, AttachmentKeep{
			ID: attachment.ID,
		})
	}
	m.Attachments = &newAttachments
	return m
}

// RetainAttachmentsByID returns a new WebhookMessageUpdate that retains the Attachment(s) with the provided IDs.
func (m WebhookMessageUpdate) RetainAttachmentsByID(attachmentIDs ...snowflake.ID) WebhookMessageUpdate {
	if m.Attachments == nil {
		m.Attachments = &[]AttachmentUpdate{}
	}
	newAttachments := slices.Clone(*m.Attachments)
	for _, attachmentID := range attachmentIDs {
		newAttachments = append(newAttachments, AttachmentKeep{
			ID: attachmentID,
		})
	}
	m.Attachments = &newAttachments
	return m
}

// WithAllowedMentions returns a new WebhookMessageUpdate with the provided AllowedMentions.
func (m WebhookMessageUpdate) WithAllowedMentions(allowedMentions *AllowedMentions) WebhookMessageUpdate {
	m.AllowedMentions = allowedMentions
	return m
}

// ClearAllowedMentions returns a new WebhookMessageUpdate with no AllowedMentions.
func (m WebhookMessageUpdate) ClearAllowedMentions() WebhookMessageUpdate {
	return m.WithAllowedMentions(nil)
}

// WithPoll returns a new WebhookMessageUpdate with the provided Poll.
func (m WebhookMessageUpdate) WithPoll(poll PollCreate) WebhookMessageUpdate {
	m.Poll = &poll
	return m
}

// ClearPoll returns a new WebhookMessageUpdate with no Poll.
func (m WebhookMessageUpdate) ClearPoll() WebhookMessageUpdate {
	m.Poll = nil
	return m
}

// WithFlags returns a new WebhookMessageUpdate with the provided MessageFlags.
// Be careful not to override the current flags when editing messages from other users - this will result in a permission error.
// Use WithSuppressEmbeds or AddFlags for flags like MessageFlagSuppressEmbeds.
func (m WebhookMessageUpdate) WithFlags(flags MessageFlags) WebhookMessageUpdate {
	m.Flags = &flags
	return m
}

// AddFlags returns a new WebhookMessageUpdate with the provided MessageFlags added.
func (m WebhookMessageUpdate) AddFlags(flags ...MessageFlags) WebhookMessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	newFlags := m.Flags.Add(flags...)
	m.Flags = &newFlags
	return m
}

// RemoveFlags returns a new WebhookMessageUpdate with the provided MessageFlags removed.
func (m WebhookMessageUpdate) RemoveFlags(flags ...MessageFlags) WebhookMessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	newFlags := m.Flags.Remove(flags...)
	m.Flags = &newFlags
	return m
}

// ClearFlags returns a new WebhookMessageUpdate with no MessageFlags.
func (m WebhookMessageUpdate) ClearFlags() WebhookMessageUpdate {
	return m.WithFlags(MessageFlagsNone)
}

// WithSuppressEmbeds returns a new WebhookMessageUpdate with MessageFlagSuppressEmbeds added/removed.
func (m WebhookMessageUpdate) WithSuppressEmbeds(suppressEmbeds bool) WebhookMessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	var newFlags MessageFlags
	if suppressEmbeds {
		newFlags = m.Flags.Add(MessageFlagSuppressEmbeds)
	} else {
		newFlags = m.Flags.Remove(MessageFlagSuppressEmbeds)
	}
	m.Flags = &newFlags
	return m
}

// WithIsComponentsV2 returns a new WebhookMessageUpdate with MessageFlagIsComponentsV2 added/removed.
// Once a message with the flag has been sent, it cannot be removed by editing the message.
func (m WebhookMessageUpdate) WithIsComponentsV2(isComponentV2 bool) WebhookMessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	var newFlags MessageFlags
	if isComponentV2 {
		newFlags = m.Flags.Add(MessageFlagIsComponentsV2)
	} else {
		newFlags = m.Flags.Remove(MessageFlagIsComponentsV2)
	}
	m.Flags = &newFlags
	return m
}
