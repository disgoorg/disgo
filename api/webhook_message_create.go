package api

import "fmt"

// WebhookMessageCreate is used to add additional messages to an Interaction after you've responded initially
type WebhookMessageCreate struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content,omitempty"`
	Embeds          []Embed          `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags     `json:"flags,omitempty"`
}

// WebhookMessageCreateBuilder allows you to create an WebhookMessageCreate with ease
type WebhookMessageCreateBuilder struct {
	WebhookMessageCreate
}

// NewWebhookMessageCreateBuilder returns a new WebhookMessageCreateBuilder
func NewWebhookMessageCreateBuilder() *WebhookMessageCreateBuilder {
	return &WebhookMessageCreateBuilder{
		WebhookMessageCreate: WebhookMessageCreate{
			AllowedMentions: &DefaultInteractionAllowedMentions,
		},
	}
}

// SetTTS sets if the WebhookMessageCreate is a tts message
func (b *WebhookMessageCreateBuilder) SetTTS(tts bool) *WebhookMessageCreateBuilder {
	b.TTS = tts
	return b
}

// SetContent sets the content of the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) SetContent(content string) *WebhookMessageCreateBuilder {
	b.Content = content
	return b
}

// SetContentf sets the content of the WebhookMessageCreate with format
func (b *WebhookMessageCreateBuilder) SetContentf(content string, a ...interface{}) *WebhookMessageCreateBuilder {
	b.Content = fmt.Sprintf(content, a...)
	return b
}

// SetEmbeds sets the embeds of the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) SetEmbeds(embeds ...Embed) *WebhookMessageCreateBuilder {
	b.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) AddEmbeds(embeds ...Embed) *WebhookMessageCreateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) ClearEmbeds() *WebhookMessageCreateBuilder {
	b.Embeds = []Embed{}
	return b
}

// RemoveEmbed removes an embed from the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) RemoveEmbed(index int) *WebhookMessageCreateBuilder {
	if b != nil && len(b.Embeds) > index {
		b.Embeds = append(b.Embeds[:index], b.Embeds[index+1:]...)
	}
	return b
}

// SetComponents sets the Component(s) of the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) SetComponents(components ...Component) *WebhookMessageCreateBuilder {
	b.Components = components
	return b
}

// AddComponents adds the Component(s) to the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) AddComponents(components ...Component) *WebhookMessageCreateBuilder {
	b.Components = append(b.Components, components...)
	return b
}

// SetAllowedMentions sets the allowed mentions of the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *WebhookMessageCreateBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the WebhookMessageCreate to nothing
func (b *WebhookMessageCreateBuilder) SetAllowedMentionsEmpty() *WebhookMessageCreateBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the message flags of the WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) SetFlags(flags MessageFlags) *WebhookMessageCreateBuilder {
	b.Flags = flags
	return b
}

// SetEphemeral adds/removes MessageFlagEphemeral to the message flags
func (b *WebhookMessageCreateBuilder) SetEphemeral(ephemeral bool) *WebhookMessageCreateBuilder {
	if ephemeral {
		b.Flags = b.Flags.Add(MessageFlagEphemeral)
	} else {
		b.Flags = b.Flags.Remove(MessageFlagEphemeral)
	}
	return b
}

// Build returns your built WebhookMessageCreate
func (b *WebhookMessageCreateBuilder) Build() WebhookMessageCreate {
	return b.WebhookMessageCreate
}
