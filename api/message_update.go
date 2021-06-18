package api

import (
	"fmt"
	"io"

	"github.com/DisgoOrg/restclient"
)

// MessageUpdate is used to edit a Message
type MessageUpdate struct {
	Content         *string           `json:"content,omitempty"`
	Embeds          []Embed           `json:"embeds,omitempty"`
	Components      []Component       `json:"components,omitempty"`
	Attachments     []Attachment      `json:"attachments,omitempty"`
	Files           []restclient.File `json:"-"`
	AllowedMentions *AllowedMentions  `json:"allowed_mentions,omitempty"`
	Flags           *MessageFlags     `json:"flags,omitempty"`
}

// ToBody returns the MessageUpdate ready for body
func (m MessageUpdate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return restclient.PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

// MessageUpdateBuilder helper to build MessageUpdate easier
type MessageUpdateBuilder struct {
	MessageUpdate
}

// NewMessageUpdateBuilder creates a new MessageUpdateBuilder to be built later
func NewMessageUpdateBuilder() *MessageUpdateBuilder {
	return &MessageUpdateBuilder{
		MessageUpdate: MessageUpdate{
			AllowedMentions: &DefaultMessageAllowedMentions,
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

// SetEmbeds sets the embeds of the Message
func (b *MessageUpdateBuilder) SetEmbeds(embeds ...Embed) *MessageUpdateBuilder {
	b.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *MessageUpdateBuilder) AddEmbeds(embeds ...Embed) *MessageUpdateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the Message
func (b *MessageUpdateBuilder) ClearEmbeds() *MessageUpdateBuilder {
	b.Embeds = []Embed{}
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *MessageUpdateBuilder) RemoveEmbed(index int) *MessageUpdateBuilder {
	if b != nil && len(b.Embeds) > index {
		b.Embeds = append(b.Embeds[:index], b.Embeds[index+1:]...)
	}
	return b
}

// SetComponents sets the Component(s) of the Message
func (b *MessageUpdateBuilder) SetComponents(components ...Component) *MessageUpdateBuilder {
	b.Components = components
	return b
}

// AddComponents adds the Component(s) to the Message
func (b *MessageUpdateBuilder) AddComponents(components ...Component) *MessageUpdateBuilder {
	b.Components = append(b.Components, components...)
	return b
}

// ClearComponents removes all of the Component(s) of the Message
func (b *MessageUpdateBuilder) ClearComponents() *MessageUpdateBuilder {
	b.Components = []Component{}
	return b
}

// RemoveComponent removes a Component from the Message
func (b *MessageUpdateBuilder) RemoveComponent(i int) *MessageUpdateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// SetFiles sets the files for this Message
func (b *MessageUpdateBuilder) SetFiles(files ...restclient.File) *MessageUpdateBuilder {
	b.Files = files
	return b
}

// AddFiles adds the files to the Message
func (b *MessageUpdateBuilder) AddFiles(files ...restclient.File) *MessageUpdateBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a file to the Message
func (b *MessageUpdateBuilder) AddFile(name string, reader io.Reader, flags ...restclient.FileFlags) *MessageUpdateBuilder {
	b.Files = append(b.Files, restclient.File{
		Name:   name,
		Reader: reader,
		Flags:  restclient.FileFlagNone.Add(flags...),
	})
	return b
}

// ClearFiles removes all files of this Message
func (b *MessageUpdateBuilder) ClearFiles() *MessageUpdateBuilder {
	b.Files = []restclient.File{}
	return b
}

// RemoveFiles removes the file at this index
func (b *MessageUpdateBuilder) RemoveFiles(i int) *MessageUpdateBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	return b
}

// RetainAttachments removes all Attachment(s) from this Message except the ones provided
func (b *MessageUpdateBuilder) RetainAttachments(attachments ...Attachment) *MessageUpdateBuilder {
	b.Attachments = append(b.Attachments, attachments...)
	return b
}

// RetainAttachmentsByID removes all Attachment(s) from this Message except the ones provided
func (b *MessageUpdateBuilder) RetainAttachmentsByID(attachmentIDs ...Snowflake) *MessageUpdateBuilder {
	for _, attachmentID := range attachmentIDs {
		b.Attachments = append(b.Attachments, Attachment{
			ID: attachmentID,
		})
	}
	return b
}

// SetFiles sets the files for this Message
func (b *MessageUpdateBuilder) SetFiles(files ...restclient.File) *MessageUpdateBuilder {
	b.Files = files
	b.updateFlags |= updateFlagFiles
	return b
}

// AddFiles adds the files to the Message
func (b *MessageUpdateBuilder) AddFiles(files ...restclient.File) *MessageUpdateBuilder {
	b.Files = append(b.Files, files...)
	b.updateFlags |= updateFlagFiles
	return b
}

// AddFile adds a file to the Message
func (b *MessageUpdateBuilder) AddFile(name string, reader io.Reader, flags ...restclient.FileFlags) *MessageUpdateBuilder {
	b.Files = append(b.Files, restclient.File{
		Name:   name,
		Reader: reader,
		Flags:  restclient.FileFlagNone.Add(flags...),
	})
	b.updateFlags |= updateFlagFiles
	return b
}

// ClearFiles removes all files of this Message
func (b *MessageUpdateBuilder) ClearFiles() *MessageUpdateBuilder {
	b.Files = []restclient.File{}
	b.updateFlags |= updateFlagFiles
	return b
}

// RemoveFiles removes the file at this index
func (b *MessageUpdateBuilder) RemoveFiles(i int) *MessageUpdateBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	b.updateFlags |= updateFlagFiles
	return b
}

// RetainAttachments removes all Attachment(s) from this Message except the ones provided
func (b *MessageUpdateBuilder) RetainAttachments(attachments ...Attachment) *MessageUpdateBuilder {
	b.Attachments = append(b.Attachments, attachments...)
	b.updateFlags |= updateFlagRetainAttachment
	return b
}

// RetainAttachmentsByID removes all Attachment(s) from this Message except the ones provided
func (b *MessageUpdateBuilder) RetainAttachmentsByID(attachmentIDs ...Snowflake) *MessageUpdateBuilder {
	for _, attachmentID := range attachmentIDs {
		b.Attachments = append(b.Attachments, Attachment{
			ID: attachmentID,
		})
	}
	b.updateFlags |= updateFlagRetainAttachment
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *MessageUpdateBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *MessageUpdateBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *MessageUpdateBuilder) ClearAllowedMentions() *MessageUpdateBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the MessageFlags of the Message
func (b *MessageUpdateBuilder) SetFlags(flags ...MessageFlags) *MessageUpdateBuilder {
	*b.Flags = MessageFlagNone.Add(flags...)
	return b
}

// AddFlags adds the MessageFlags of the Message
func (b *MessageUpdateBuilder) AddFlags(flags ...MessageFlags) *MessageUpdateBuilder {
	*b.Flags = b.Flags.Add(flags...)
	return b
}

// RemoveFlags removes the MessageFlags of the Message
func (b *MessageUpdateBuilder) RemoveFlags(flags ...MessageFlags) *MessageUpdateBuilder {
	*b.Flags = b.Flags.Remove(flags...)
	return b
}

// ClearFlags clears the MessageFlags of the Message
func (b *MessageUpdateBuilder) ClearFlags() *MessageUpdateBuilder {
	return b.SetFlags(MessageFlagNone)
}

// Build builds the MessageUpdateBuilder to a MessageUpdate struct
func (b *MessageUpdateBuilder) Build() MessageUpdate {
	return b.MessageUpdate
}
