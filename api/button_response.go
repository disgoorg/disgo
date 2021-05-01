package api

import "fmt"

// The ButtonResponseData ...
type ButtonResponseData struct {
	Content         *string          `json:"content"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           *MessageFlags    `json:"flags,omitempty"`
}

// ButtonResponseBuilder allows you to create an InteractionResponse with ease
type ButtonResponseBuilder struct {
	ButtonResponseData
}

// NewButtonResponseBuilder returns a new ButtonResponseBuilder
func NewButtonResponseBuilder() *ButtonResponseBuilder {
	return &ButtonResponseBuilder{
		ButtonResponseData: ButtonResponseData{
			AllowedMentions: &DefaultInteractionAllowedMentions,
		},
	}
}

// SetContent sets the content of the InteractionResponse
func (b *ButtonResponseBuilder) SetContent(content string) *ButtonResponseBuilder {
	b.Content = &content
	return b
}

// SetContentf sets the content of the InteractionResponse with format
func (b *ButtonResponseBuilder) SetContentf(content string, a ...interface{}) *ButtonResponseBuilder {
	contentf := fmt.Sprintf(content, a...)
	b.Content = &contentf
	return b
}

// SetEmbeds sets the embeds of the InteractionResponse
func (b *ButtonResponseBuilder) SetEmbeds(embeds ...*Embed) *ButtonResponseBuilder {
	b.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the InteractionResponse
func (b *ButtonResponseBuilder) AddEmbeds(embeds ...*Embed) *ButtonResponseBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the InteractionResponse
func (b *ButtonResponseBuilder) ClearEmbeds() *ButtonResponseBuilder {
	if b != nil {
		b.Embeds = []*Embed{}
	}
	return b
}

// RemoveEmbed removes an embed from the InteractionResponse
func (b *ButtonResponseBuilder) RemoveEmbed(i int) *ButtonResponseBuilder {
	if b != nil && len(b.Embeds) > i {
		b.Embeds = append(b.Embeds[:i], b.Embeds[i+1:]...)
	}
	return b
}

// SetComponents sets the Component(s) of the InteractionResponse
func (b *ButtonResponseBuilder) SetComponents(components ...Component) *ButtonResponseBuilder {
	b.Components = components
	return b
}

// SetAllowedMentions sets the allowed mentions of the InteractionResponse
func (b *ButtonResponseBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *ButtonResponseBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the InteractionResponse to nothing
func (b *ButtonResponseBuilder) SetAllowedMentionsEmpty() *ButtonResponseBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the message flags of the InteractionResponse
func (b *ButtonResponseBuilder) SetFlags(flags MessageFlags) *ButtonResponseBuilder {
	b.Flags = &flags
	return b
}

// SetEphemeral adds/removes MessageFlagEphemeral to the message flags
func (b *ButtonResponseBuilder) SetEphemeral(ephemeral bool) *ButtonResponseBuilder {
	if ephemeral {
		*b.Flags = MessageFlagEphemeral

	} else {
		*b.Flags |= MessageFlagEphemeral
	}
	return b
}

// Build returns your built InteractionResponse
func (b *ButtonResponseBuilder) Build() *InteractionResponse {
	return &InteractionResponse{
		Type: InteractionResponseTypeButtonResponse,
		Data: &b.ButtonResponseData,
	}
}
