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
	Components      []Component      `json:"components"`
	Embed           *Embed           `json:"embed"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions"`
	Flags           MessageFlags     `json:"flags"`
	updateFlags     updateFlags
}

func (u MessageUpdate) isUpdated(flag updateFlags) bool {
	return (u.updateFlags & flag) == flag
}

func (u MessageUpdate) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	if u.isUpdated(updateFlagContent) {
		data["content"] = u.Content
	}
	if u.isUpdated(updateFlagComponents) {
		data["components"] = u.Components
	}
	if u.isUpdated(updateFlagEmbed) {
		data["embed"] = u.Embed
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

// SetEmbed sets the Embed of the Message
func (b *MessageUpdateBuilder) SetEmbed(embed *Embed) *MessageUpdateBuilder {
	b.Embed = embed
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

// Build builds the MessageUpdateBuilder to a MessageUpdate struct
func (b *MessageUpdateBuilder) Build() MessageUpdate {
	return b.MessageUpdate
}
