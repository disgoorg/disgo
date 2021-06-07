package api

import (
	"encoding/json"
	"fmt"
)

// WebhookMessageUpdate is used to edit a WebhookMessage
type WebhookMessageUpdate struct {
	Content         string           `json:"content"`
	Components      []Component      `json:"components"`
	Embeds          []*Embed         `json:"embeds"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions"`
	Flags           MessageFlags     `json:"flags"`
	updateFlags     updateFlags
}

func (u WebhookMessageUpdate) isUpdated(flag updateFlags) bool {
	return (u.updateFlags & flag) == flag
}

// MarshalJSON marshals the WebhookMessageUpdate into json
func (u WebhookMessageUpdate) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	if u.isUpdated(updateFlagContent) {
		data["content"] = u.Content
	}
	if u.isUpdated(updateFlagComponents) {
		data["components"] = u.Components
	}
	if u.isUpdated(updateFlagEmbed) {
		data["embeds"] = u.Embeds
	}
	if u.isUpdated(updateFlagAllowedMentions) {
		data["allowed_mentions"] = u.AllowedMentions
	}
	if u.isUpdated(updateFlagFlags) {
		data["flags"] = u.Flags
	}

	return json.Marshal(data)
}

// WebhookMessageUpdateBuilder helper to build WebhookMessageUpdate easier
type WebhookMessageUpdateBuilder struct {
	WebhookMessageUpdate
}

// NewWebhookMessageUpdateBuilder creates a new WebhookMessageUpdateBuilder to be built later
func NewWebhookMessageUpdateBuilder() *WebhookMessageUpdateBuilder {
	return &WebhookMessageUpdateBuilder{
		WebhookMessageUpdate: WebhookMessageUpdate{
			AllowedMentions: &DefaultInteractionAllowedMentions,
		},
	}
}

// SetContent sets content of the Message
func (b *WebhookMessageUpdateBuilder) SetContent(content string) *WebhookMessageUpdateBuilder {
	b.Content = content
	b.updateFlags |= updateFlagContent
	return b
}

// SetContentf sets content of the Message
func (b *WebhookMessageUpdateBuilder) SetContentf(content string, a ...interface{}) *WebhookMessageUpdateBuilder {
	return b.SetContent(fmt.Sprintf(content, a...))
}

// SetEmbeds sets the embeds of the WebhookMessageUpdate
func (b *WebhookMessageUpdateBuilder) SetEmbeds(embeds ...*Embed) *WebhookMessageUpdateBuilder {
	b.Embeds = embeds
	b.updateFlags |= updateFlagEmbed
	return b
}

// AddEmbeds adds multiple embeds to the WebhookMessageUpdate
func (b *WebhookMessageUpdateBuilder) AddEmbeds(embeds ...*Embed) *WebhookMessageUpdateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	b.updateFlags |= updateFlagEmbed
	return b
}

// ClearEmbeds removes all of the embeds from the WebhookMessageUpdate
func (b *WebhookMessageUpdateBuilder) ClearEmbeds() *WebhookMessageUpdateBuilder {
	b.Embeds = []*Embed{}
	b.updateFlags |= updateFlagEmbed
	return b
}

// RemoveEmbed removes an embed from the WebhookMessageUpdate
func (b *WebhookMessageUpdateBuilder) RemoveEmbed(index int) *WebhookMessageUpdateBuilder {
	if b != nil && len(b.Embeds) > index {
		b.Embeds = append(b.Embeds[:index], b.Embeds[index+1:]...)
	}
	b.updateFlags |= updateFlagEmbed
	return b
}

// SetComponents sets the Component(s) of the Message
func (b *WebhookMessageUpdateBuilder) SetComponents(components ...Component) *WebhookMessageUpdateBuilder {
	b.Components = components
	b.updateFlags |= updateFlagComponents
	return b
}

// AddComponents adds the Component(s) to the Message
func (b *WebhookMessageUpdateBuilder) AddComponents(components ...Component) *WebhookMessageUpdateBuilder {
	b.Components = append(b.Components, components...)
	b.updateFlags |= updateFlagComponents
	return b
}

// ClearComponents removes all of the Component(s) of the Message
func (b *WebhookMessageUpdateBuilder) ClearComponents() *WebhookMessageUpdateBuilder {
	b.Components = []Component{}
	b.updateFlags |= updateFlagComponents
	return b
}

// RemoveComponent removes a Component from the Message
func (b *WebhookMessageUpdateBuilder) RemoveComponent(i int) *WebhookMessageUpdateBuilder {
	if b != nil && len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	b.updateFlags |= updateFlagComponents
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *WebhookMessageUpdateBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *WebhookMessageUpdateBuilder {
	b.AllowedMentions = allowedMentions
	b.updateFlags |= updateFlagAllowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *WebhookMessageUpdateBuilder) ClearAllowedMentions() *WebhookMessageUpdateBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// Build builds the WebhookMessageUpdateBuilder to a WebhookMessageUpdate struct
func (b *WebhookMessageUpdateBuilder) Build() WebhookMessageUpdate {
	return b.WebhookMessageUpdate
}
