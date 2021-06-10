package api

import (
	"encoding/json"
	"fmt"
)

type updateFlags int

const (
	updateFlagContent = 1 << iota
	updateFlagComponents
	updateFlagEmbed
	updateFlagFlags
	updateFlagAllowedMentions
)

// MessageUpdate is used to edit a Message
type MessageUpdate struct {
	Content         string           `json:"content"`
	Embeds          []Embed          `json:"embeds"`
	Components      []Component      `json:"components"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions"`
	Flags           MessageFlags     `json:"flags"`
	updateFlags     updateFlags
}

func (u MessageUpdate) isUpdated(flag updateFlags) bool {
	return (u.updateFlags & flag) == flag
}

// MarshalJSON marshals the MessageUpdate into json
func (u MessageUpdate) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	if u.isUpdated(updateFlagContent) {
		data["content"] = u.Content
	}
	if u.isUpdated(updateFlagEmbed) {
		data["embeds"] = u.Embeds
	}
	if u.isUpdated(updateFlagComponents) {
		data["components"] = u.Components
	}
	if u.isUpdated(updateFlagAllowedMentions) {
		data["allowed_mentions"] = u.AllowedMentions
	}
	if u.isUpdated(updateFlagFlags) {
		data["flags"] = u.Flags
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
	b.updateFlags |= updateFlagEmbed
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *MessageUpdateBuilder) AddEmbeds(embeds ...Embed) *MessageUpdateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	b.updateFlags |= updateFlagEmbed
	return b
}

// ClearEmbeds removes all of the embeds from the Message
func (b *MessageUpdateBuilder) ClearEmbeds() *MessageUpdateBuilder {
	b.Embeds = []Embed{}
	b.updateFlags |= updateFlagEmbed
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *MessageUpdateBuilder) RemoveEmbed(index int) *MessageUpdateBuilder {
	if b != nil && len(b.Embeds) > index {
		b.Embeds = append(b.Embeds[:index], b.Embeds[index+1:]...)
	}
	b.updateFlags |= updateFlagEmbed
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
	if b != nil && len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	b.updateFlags |= updateFlagComponents
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
