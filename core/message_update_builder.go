package core

import (
	"fmt"
	"io"

	"github.com/DisgoOrg/disgo/discord"
)

// MessageUpdateBuilder helper to build MessageUpdate easier
type MessageUpdateBuilder struct {
	discord.MessageUpdate
}

// NewMessageUpdateBuilder creates a new MessageUpdateBuilder to be built later
func NewMessageUpdateBuilder() *MessageUpdateBuilder {
	return &MessageUpdateBuilder{
		MessageUpdate: discord.MessageUpdate{
			AllowedMentions: &DefaultAllowedMentions,
		},
	}
}

// SetContent sets content of the Message
func (b *MessageUpdateBuilder) SetContent(content string) *MessageUpdateBuilder {
	b.Content = &content
	return b
}

// SetContentf sets content of the Message
func (b *MessageUpdateBuilder) SetContentf(content string, a ...interface{}) *MessageUpdateBuilder {
	return b.SetContent(fmt.Sprintf(content, a...))
}

// ClearContent removes content of the Message
func (b *MessageUpdateBuilder) ClearContent() *MessageUpdateBuilder {
	return b.SetContent("")
}

// SetEmbeds sets the discord.Embed(s) of the Message
func (b *MessageUpdateBuilder) SetEmbeds(embeds ...discord.Embed) *MessageUpdateBuilder {
	if b.Embeds == nil {
		b.Embeds = new([]discord.Embed)
	}
	*b.Embeds = embeds
	return b
}

// SetEmbed sets the provided discord.Embed at the index of the Message
func (b *MessageUpdateBuilder) SetEmbed(i int, embed discord.Embed) *MessageUpdateBuilder {
	if b.Embeds == nil {
		b.Embeds = new([]discord.Embed)
	}
	if len(*b.Embeds) > i {
		(*b.Embeds)[i] = embed
	}
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *MessageUpdateBuilder) AddEmbeds(embeds ...discord.Embed) *MessageUpdateBuilder {
	if b.Embeds == nil {
		b.Embeds = new([]discord.Embed)
	}
	*b.Embeds = append(*b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all the embeds from the Message
func (b *MessageUpdateBuilder) ClearEmbeds() *MessageUpdateBuilder {
	b.Embeds = &[]discord.Embed{}
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *MessageUpdateBuilder) RemoveEmbed(i int) *MessageUpdateBuilder {
	if b.Embeds == nil {
		b.Embeds = new([]discord.Embed)
	}
	if len(*b.Embeds) > i {
		*b.Embeds = append((*b.Embeds)[:i], (*b.Embeds)[i+1:]...)
	}
	return b
}

// SetActionRows sets the discord.ActionRowComponent(s) of the Message
func (b *MessageUpdateBuilder) SetActionRows(actionRows ...discord.ActionRowComponent) *MessageUpdateBuilder {
	if b.Components == nil {
		b.Components = new([]discord.ActionRowComponent)
	}
	*b.Components = actionRows
	return b
}

// SetActionRow sets the provided discord.ActionRowComponent at the index of discord.Component(s)
func (b *MessageUpdateBuilder) SetActionRow(i int, actionRow discord.ActionRowComponent) *MessageUpdateBuilder {
	if b.Components == nil {
		b.Components = new([]discord.ActionRowComponent)
	}
	if len(*b.Components) > i {
		(*b.Components)[i] = actionRow
	}
	return b
}

// AddActionRow adds a new discord.ActionRowComponent with the provided discord.Component(s) to the Message
func (b *MessageUpdateBuilder) AddActionRow(components ...discord.Component) *MessageUpdateBuilder {
	if b.Components == nil {
		b.Components = new([]discord.ActionRowComponent)
	}
	*b.Components = append(*b.Components, components)
	return b
}

// AddActionRows adds the discord.ActionRowComponent(s) to the Message
func (b *MessageUpdateBuilder) AddActionRows(actionRows ...discord.ActionRowComponent) *MessageUpdateBuilder {
	if b.Components == nil {
		b.Components = new([]discord.ActionRowComponent)
	}
	*b.Components = append(*b.Components, actionRows...)
	return b
}

// RemoveActionRow removes a discord.ActionRowComponent from the Message
func (b *MessageUpdateBuilder) RemoveActionRow(i int) *MessageUpdateBuilder {
	if b.Components == nil {
		b.Components = new([]discord.ActionRowComponent)
	}
	if len(*b.Components) > i {
		*b.Components = append((*b.Components)[:i], (*b.Components)[i+1:]...)
	}
	return b
}

// ClearActionRows removes all the discord.ActionRowComponent(s) of the Message
func (b *MessageUpdateBuilder) ClearActionRows() *MessageUpdateBuilder {
	b.Components = &[]discord.ActionRowComponent{}
	return b
}

// SetFiles sets the new discord.File(s) for this discord.MessageUpdate
func (b *MessageUpdateBuilder) SetFiles(files ...*discord.File) *MessageUpdateBuilder {
	b.Files = files
	return b
}

// SetFile sets the new discord.File at the index for this discord.MessageUpdate
func (b *MessageUpdateBuilder) SetFile(i int, file *discord.File) *MessageUpdateBuilder {
	if len(b.Files) > i {
		b.Files[i] = file
	}
	return b
}

// AddFiles adds the new discord.File(s) to the discord.MessageUpdate
func (b *MessageUpdateBuilder) AddFiles(files ...*discord.File) *MessageUpdateBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a new discord.File to the discord.MessageUpdate
func (b *MessageUpdateBuilder) AddFile(name string, reader io.Reader, flags ...discord.FileFlags) *MessageUpdateBuilder {
	b.Files = append(b.Files, discord.NewFile(name, reader, flags...))
	return b
}

// ClearFiles removes all new files of this discord.MessageUpdate
func (b *MessageUpdateBuilder) ClearFiles() *MessageUpdateBuilder {
	b.Files = []*discord.File{}
	return b
}

// RemoveFile removes the new discord.File at this index
func (b *MessageUpdateBuilder) RemoveFile(i int) *MessageUpdateBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	return b
}

// RetainAttachments removes all Attachment(s) from this Message except the ones provided
func (b *MessageUpdateBuilder) RetainAttachments(attachments ...discord.Attachment) *MessageUpdateBuilder {
	if b.Attachments == nil {
		b.Attachments = new([]discord.Attachment)
	}
	*b.Attachments = append(*b.Attachments, attachments...)
	return b
}

// RetainAttachmentsByID removes all Attachment(s) from this Message except the ones provided
func (b *MessageUpdateBuilder) RetainAttachmentsByID(attachmentIDs ...discord.Snowflake) *MessageUpdateBuilder {
	if b.Attachments == nil {
		b.Attachments = new([]discord.Attachment)
	}
	for _, attachmentID := range attachmentIDs {
		*b.Attachments = append(*b.Attachments, discord.Attachment{ID: attachmentID})
	}
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *MessageUpdateBuilder) SetAllowedMentions(allowedMentions *discord.AllowedMentions) *MessageUpdateBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *MessageUpdateBuilder) ClearAllowedMentions() *MessageUpdateBuilder {
	return b.SetAllowedMentions(nil)
}

// SetFlags sets the message flags of the Message
func (b *MessageUpdateBuilder) SetFlags(flags discord.MessageFlags) *MessageUpdateBuilder {
	if b.Flags == nil {
		b.Flags = new(discord.MessageFlags)
	}
	*b.Flags = flags
	return b
}

// AddFlags adds the MessageFlags of the Message
func (b *MessageUpdateBuilder) AddFlags(flags ...discord.MessageFlags) *MessageUpdateBuilder {
	if b.Flags == nil {
		b.Flags = new(discord.MessageFlags)
	}
	*b.Flags = b.Flags.Add(flags...)
	return b
}

// RemoveFlags removes the MessageFlags of the Message
func (b *MessageUpdateBuilder) RemoveFlags(flags ...discord.MessageFlags) *MessageUpdateBuilder {
	if b.Flags == nil {
		b.Flags = new(discord.MessageFlags)
	}
	*b.Flags = b.Flags.Remove(flags...)
	return b
}

// ClearFlags clears the MessageFlags of the Message
func (b *MessageUpdateBuilder) ClearFlags() *MessageUpdateBuilder {
	return b.SetFlags(discord.MessageFlagNone)
}

// Build builds the MessageUpdateBuilder to a MessageUpdate struct
func (b *MessageUpdateBuilder) Build() discord.MessageUpdate {
	return b.MessageUpdate
}
