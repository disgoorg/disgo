package api

import "fmt"

// The SlashCommandResponse is used to specify the message_events options when creating an InteractionResponse
type SlashCommandResponse struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags     `json:"flags,omitempty"`
}

// SlashCommandResponseBuilder allows you to create an InteractionResponse with ease
type SlashCommandResponseBuilder struct {
	SlashCommandResponse
}

// NewSlashCommandResponseBuilder returns a new SlashCommandResponseBuilder
func NewSlashCommandResponseBuilder() *SlashCommandResponseBuilder {
	return &SlashCommandResponseBuilder{
		SlashCommandResponse: SlashCommandResponse{
			AllowedMentions: &DefaultInteractionAllowedMentions,
		},
	}
}

// SetTTS sets if the InteractionResponse is a tts message
func (b *SlashCommandResponseBuilder) SetTTS(tts bool) *SlashCommandResponseBuilder {
	b.TTS = tts
	return b
}

// SetContent sets the content of the InteractionResponse
func (b *SlashCommandResponseBuilder) SetContent(content string) *SlashCommandResponseBuilder {
	b.Content = content
	return b
}

// SetContentf sets the content of the InteractionResponse with format
func (b *SlashCommandResponseBuilder) SetContentf(content string, a ...interface{}) *SlashCommandResponseBuilder {
	b.Content = fmt.Sprintf(content, a...)
	return b
}

// SetEmbeds sets the embeds of the InteractionResponse
func (b *SlashCommandResponseBuilder) SetEmbeds(embeds ...*Embed) *SlashCommandResponseBuilder {
	b.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the InteractionResponse
func (b *SlashCommandResponseBuilder) AddEmbeds(embeds ...*Embed) *SlashCommandResponseBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the InteractionResponse
func (b *SlashCommandResponseBuilder) ClearEmbeds() *SlashCommandResponseBuilder {
	if b != nil {
		b.Embeds = []*Embed{}
	}
	return b
}

// RemoveEmbed removes an embed from the InteractionResponse
func (b *SlashCommandResponseBuilder) RemoveEmbed(i int) *SlashCommandResponseBuilder {
	if b != nil && len(b.Embeds) > i {
		b.Embeds = append(b.Embeds[:i], b.Embeds[i+1:]...)
	}
	return b
}

// SetComponents sets the Component(s) of the InteractionResponse
func (b *SlashCommandResponseBuilder) SetComponents(components ...Component) *SlashCommandResponseBuilder {
	b.Components = components
	return b
}

// AddComponents adds the Component(s) to the InteractionResponse
func (b *SlashCommandResponseBuilder) AddComponents(components ...Component) *SlashCommandResponseBuilder {
	b.Components = append(b.Components, components...)
	return b
}

// SetAllowedMentions sets the allowed mentions of the InteractionResponse
func (b *SlashCommandResponseBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *SlashCommandResponseBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the InteractionResponse to nothing
func (b *SlashCommandResponseBuilder) SetAllowedMentionsEmpty() *SlashCommandResponseBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the message flags of the InteractionResponse
func (b *SlashCommandResponseBuilder) SetFlags(flags MessageFlags) *SlashCommandResponseBuilder {
	b.Flags = flags
	return b
}

// SetEphemeral adds/removes MessageFlagEphemeral to the message flags
func (b *SlashCommandResponseBuilder) SetEphemeral(ephemeral bool) *SlashCommandResponseBuilder {
	if ephemeral {
		b.Flags &= MessageFlagEphemeral

	} else {
		b.Flags |= MessageFlagEphemeral
	}
	return b
}

// Build returns your built InteractionResponse
func (b *SlashCommandResponseBuilder) Build() *InteractionResponse {
	return &InteractionResponse{
		Type: InteractionResponseTypeChannelMessageWithSource,
		Data: &b.SlashCommandResponse,
	}
}
