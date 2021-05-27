package api

import "fmt"

// InteractionResponseType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionResponseType int

// Constants for the InteractionResponseType(s)
const (
	InteractionResponseTypePong InteractionResponseType = iota + 1
	_
	_
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
	InteractionResponseTypeDeferredUpdateMessage
	InteractionResponseTypeUpdateMessage
)

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionResponseType  `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

// The InteractionResponseData is used to specify the message_events options when creating an InteractionResponse
type InteractionResponseData struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags     `json:"flags,omitempty"`
}

// InteractionResponseBuilder allows you to create an InteractionResponse with ease
type InteractionResponseBuilder struct {
	InteractionResponse
}

// NewInteractionResponseBuilder returns a new InteractionResponseBuilder
func NewInteractionResponseBuilder() *InteractionResponseBuilder {
	return &InteractionResponseBuilder{
		InteractionResponse: InteractionResponse{
			Data: &InteractionResponseData{
				AllowedMentions: &DefaultInteractionAllowedMentions,
			},
		},
	}
}

// NewInteractionResponseBuilderByMessage returns a new InteractionResponseBuilder and takes an existing Message
func NewInteractionResponseBuilderByMessage(message *Message) *InteractionResponseBuilder {
	rs := &InteractionResponseData{
		TTS:             message.TTS,
		Embeds:          message.Embeds,
		Components:      message.Components,
		AllowedMentions: &DefaultInteractionAllowedMentions,
		Flags:           message.Flags,
	}
	if message.Content != nil {
		rs.Content = message.Content
	}
	return &InteractionResponseBuilder{
		InteractionResponse: InteractionResponse{
			Data: rs,
		},
	}
}

// SetType sets if the InteractionResponseType of this InteractionResponse
func (b *InteractionResponseBuilder) SetType(responseType InteractionResponseType) *InteractionResponseBuilder {
	b.Type = responseType
	return b
}

// SetTTS sets if the InteractionResponse is a tts message
func (b *InteractionResponseBuilder) SetTTS(tts bool) *InteractionResponseBuilder {
	b.Data.TTS = tts
	return b
}

// SetContent sets the content of the InteractionResponse
func (b *InteractionResponseBuilder) SetContent(content string) *InteractionResponseBuilder {
	b.Data.Content = &content
	return b
}

// SetContentf sets the content of the InteractionResponse with format
func (b *InteractionResponseBuilder) SetContentf(content string, a ...interface{}) *InteractionResponseBuilder {
	contentf := fmt.Sprintf(content, a...)
	b.Data.Content = &contentf
	return b
}

// ClearContent sets the content of the InteractionResponse to nil
func (b *InteractionResponseBuilder) ClearContent() *InteractionResponseBuilder {
	b.Data.Content = nil
	return b
}

// SetEmbeds sets the embeds of the InteractionResponse
func (b *InteractionResponseBuilder) SetEmbeds(embeds ...*Embed) *InteractionResponseBuilder {
	b.Data.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the InteractionResponse
func (b *InteractionResponseBuilder) AddEmbeds(embeds ...*Embed) *InteractionResponseBuilder {
	b.Data.Embeds = append(b.Data.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the InteractionResponse
func (b *InteractionResponseBuilder) ClearEmbeds() *InteractionResponseBuilder {
	if b != nil {
		b.Data.Embeds = []*Embed{}
	}
	return b
}

// RemoveEmbed removes an embed from the InteractionResponse
func (b *InteractionResponseBuilder) RemoveEmbed(i int) *InteractionResponseBuilder {
	if b != nil && len(b.Data.Embeds) > i {
		b.Data.Embeds = append(b.Data.Embeds[:i], b.Data.Embeds[i+1:]...)
	}
	return b
}

// SetComponents sets the Component(s) of the InteractionResponse
func (b *InteractionResponseBuilder) SetComponents(components ...Component) *InteractionResponseBuilder {
	b.Data.Components = components
	return b
}

// AddComponents adds the Component(s) to the InteractionResponse
func (b *InteractionResponseBuilder) AddComponents(components ...Component) *InteractionResponseBuilder {
	b.Data.Components = append(b.Data.Components, components...)
	return b
}

// ClearComponents removes all of the Component(s) of the InteractionResponse
func (b *InteractionResponseBuilder) ClearComponents() *InteractionResponseBuilder {
	if b != nil {
		b.Data.Components = []Component{}
	}
	return b
}

// RemoveComponent removes a Component from the InteractionResponse
func (b *InteractionResponseBuilder) RemoveComponent(i int) *InteractionResponseBuilder {
	if b != nil && len(b.Data.Components) > i {
		b.Data.Components = append(b.Data.Components[:i], b.Data.Components[i+1:]...)
	}
	return b
}

// SetAllowedMentions sets the allowed mentions of the InteractionResponse
func (b *InteractionResponseBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *InteractionResponseBuilder {
	b.Data.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the InteractionResponse to nothing
func (b *InteractionResponseBuilder) SetAllowedMentionsEmpty() *InteractionResponseBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the message flags of the InteractionResponse
func (b *InteractionResponseBuilder) SetFlags(flags MessageFlags) *InteractionResponseBuilder {
	b.Data.Flags = flags
	return b
}

// SetEphemeral adds/removes MessageFlagEphemeral to the message flags
func (b *InteractionResponseBuilder) SetEphemeral(ephemeral bool) *InteractionResponseBuilder {
	if ephemeral {
		b.Data.Flags &= MessageFlagEphemeral

	} else {
		b.Data.Flags |= MessageFlagEphemeral
	}
	return b
}

// Build returns your built InteractionResponse
func (b *InteractionResponseBuilder) Build() *InteractionResponse {
	return &b.InteractionResponse
}

// BuildData returns your built InteractionResponseData
func (b *InteractionResponseBuilder) BuildData() *InteractionResponseData {
	return b.Data
}
