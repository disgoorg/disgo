package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DisgoOrg/restclient"
)

type updateFlags int

const (
	updateFlagContent = 1 << iota
	updateFlagComponents
	updateFlagEmbeds
	updateFlagFiles
	updateFlagRetainAttachment
	updateFlagAllowedMentions
	updateFlagFlags
)

// MessageUpdate is used to edit a Message
type MessageUpdate struct {
	Content         string            `json:"content"`
	Embeds          []Embed           `json:"embeds"`
	Components      []Component       `json:"components"`
	Attachments     []Attachment      `json:"attachments"`
	Files           []restclient.File `json:"-"`
	AllowedMentions *AllowedMentions  `json:"allowed_mentions"`
	Flags           MessageFlags      `json:"flags"`
	updateFlags     updateFlags
}

func (m MessageUpdate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 && m.isUpdated(updateFlagFiles) {
		return restclient.PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

func (m MessageUpdate) isUpdated(flag updateFlags) bool {
	return (m.updateFlags & flag) == flag
}

// MarshalJSON marshals the MessageUpdate into json
func (m MessageUpdate) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	if m.isUpdated(updateFlagContent) {
		data["content"] = m.Content
	}
	if m.isUpdated(updateFlagEmbeds) {
		data["embeds"] = m.Embeds
	}
	if m.isUpdated(updateFlagComponents) {
		data["components"] = m.Components
	}
	if m.isUpdated(updateFlagRetainAttachment) {
		data["attachments"] = m.Attachments
	}
	if m.isUpdated(updateFlagAllowedMentions) {
		data["allowed_mentions"] = m.AllowedMentions
	}
	if m.isUpdated(updateFlagFlags) {
		data["flags"] = m.Flags
	}

	return json.Marshal(data)
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
	b.Content = content
	b.updateFlags |= updateFlagContent
	return b
}

// SetContentf sets content of the Message
func (b *MessageUpdateBuilder) SetContentf(content string, a ...interface{}) *MessageUpdateBuilder {
	return b.SetContent(fmt.Sprintf(content, a...))
}

// SetEmbeds sets the embeds of the Message
func (b *MessageUpdateBuilder) SetEmbeds(embeds ...Embed) *MessageUpdateBuilder {
	b.Embeds = embeds
	b.updateFlags |= updateFlagEmbeds
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *MessageUpdateBuilder) AddEmbeds(embeds ...Embed) *MessageUpdateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	b.updateFlags |= updateFlagEmbeds
	return b
}

// ClearEmbeds removes all of the embeds from the Message
func (b *MessageUpdateBuilder) ClearEmbeds() *MessageUpdateBuilder {
	b.Embeds = []Embed{}
	b.updateFlags |= updateFlagEmbeds
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *MessageUpdateBuilder) RemoveEmbed(index int) *MessageUpdateBuilder {
	if b != nil && len(b.Embeds) > index {
		b.Embeds = append(b.Embeds[:index], b.Embeds[index+1:]...)
	}
	b.updateFlags |= updateFlagEmbeds
	return b
}

// SetComponents sets the Component(s) of the Message
func (b *MessageUpdateBuilder) SetComponents(components ...Component) *MessageUpdateBuilder {
	b.Components = components
	b.updateFlags |= updateFlagComponents
	return b
}

// AddComponents adds the Component(s) to the Message
func (b *MessageUpdateBuilder) AddComponents(components ...Component) *MessageUpdateBuilder {
	b.Components = append(b.Components, components...)
	b.updateFlags |= updateFlagComponents
	return b
}

// ClearComponents removes all of the Component(s) of the Message
func (b *MessageUpdateBuilder) ClearComponents() *MessageUpdateBuilder {
	b.Components = []Component{}
	b.updateFlags |= updateFlagComponents
	return b
}

// RemoveComponent removes a Component from the Message
func (b *MessageUpdateBuilder) RemoveComponent(i int) *MessageUpdateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	b.updateFlags |= updateFlagComponents
	return b
}

func (b *MessageUpdateBuilder) SetFiles(files ...restclient.File) *MessageUpdateBuilder {
	b.Files = files
	b.updateFlags |= updateFlagFiles
	return b
}

func (b *MessageUpdateBuilder) AddFiles(files ...restclient.File) *MessageUpdateBuilder {
	b.Files = append(b.Files, files...)
	b.updateFlags |= updateFlagFiles
	return b
}

func (b *MessageUpdateBuilder) AddFile(name string, reader io.Reader, flags ...restclient.FileFlags) *MessageUpdateBuilder {
	b.Files = append(b.Files, restclient.File{
		Name:   name,
		Reader: reader,
		Flags:  restclient.FileFlagNone.Add(flags...),
	})
	b.updateFlags |= updateFlagFiles
	return b
}

func (b *MessageUpdateBuilder) ClearFiles() *MessageUpdateBuilder {
	b.Files = []restclient.File{}
	b.updateFlags |= updateFlagFiles
	return b
}

func (b *MessageUpdateBuilder) RemoveFiles(i int) *MessageUpdateBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	b.updateFlags |= updateFlagFiles
	return b
}

func (b *MessageUpdateBuilder) RetainAttachments(attachments ...Attachment) *MessageUpdateBuilder {
	b.Attachments = append(b.Attachments, attachments...)
	b.updateFlags |= updateFlagRetainAttachment
	return b
}

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
	b.updateFlags |= updateFlagAllowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *MessageUpdateBuilder) ClearAllowedMentions() *MessageUpdateBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the MessageFlags of the Message
func (b *MessageUpdateBuilder) SetFlags(flags MessageFlags) *MessageUpdateBuilder {
	b.Flags = flags
	return b
}

// Build builds the MessageUpdateBuilder to a MessageUpdate struct
func (b *MessageUpdateBuilder) Build() MessageUpdate {
	return b.MessageUpdate
}
