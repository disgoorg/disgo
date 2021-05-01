package api

import "fmt"

// The CommandResponseData is used to specify the message_events options when creating an InteractionResponse
type CommandResponseData struct {
	TTS             *bool            `json:"tts,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags     `json:"flags,omitempty"`
}

// CommandResponseBuilder allows you to create an InteractionResponse with ease
type CommandResponseBuilder struct {
	CommandResponseData
}

// NewCommandResponseBuilder returns a new CommandResponseBuilder
func NewCommandResponseBuilder() *CommandResponseBuilder {
	return &CommandResponseBuilder{
		CommandResponseData: CommandResponseData{
			AllowedMentions: &DefaultInteractionAllowedMentions,
		},
	}
}

// SetTTS sets if the InteractionResponse is a tts message
func (b *CommandResponseBuilder) SetTTS(tts bool) *CommandResponseBuilder {
	b.TTS = &tts
	return b
}

// SetContent sets the content of the InteractionResponse
func (b *CommandResponseBuilder) SetContent(content string) *CommandResponseBuilder {
	b.Content = &content
	return b
}

// SetContentf sets the content of the InteractionResponse with format
func (b *CommandResponseBuilder) SetContentf(content string, a ...interface{}) *CommandResponseBuilder {
	contentf := fmt.Sprintf(content, a...)
	b.Content = &contentf
	return b
}

// SetEmbeds sets the embeds of the InteractionResponse
func (b *CommandResponseBuilder) SetEmbeds(embeds ...*Embed) *CommandResponseBuilder {
	b.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the InteractionResponse
func (b *CommandResponseBuilder) AddEmbeds(embeds ...*Embed) *CommandResponseBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the InteractionResponse
func (b *CommandResponseBuilder) ClearEmbeds() *CommandResponseBuilder {
	if b != nil {
		b.Embeds = []*Embed{}
	}
	return b
}

// RemoveEmbed removes an embed from the InteractionResponse
func (b *CommandResponseBuilder) RemoveEmbed(i int) *CommandResponseBuilder {
	if b != nil && len(b.Embeds) > i {
		b.Embeds = append(b.Embeds[:i], b.Embeds[i+1:]...)
	}
	return b
}

// SetComponents sets the Component(s) of the InteractionResponse
func (b *CommandResponseBuilder) SetComponents(components ...Component) *CommandResponseBuilder {
	b.Components = components
	return b
}

// SetAllowedMentions sets the allowed mentions of the InteractionResponse
func (b *CommandResponseBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *CommandResponseBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the InteractionResponse to nothing
func (b *CommandResponseBuilder) SetAllowedMentionsEmpty() *CommandResponseBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetFlags sets the message flags of the InteractionResponse
func (b *CommandResponseBuilder) SetFlags(flags MessageFlags) *CommandResponseBuilder {
	b.Flags = flags
	return b
}

// SetEphemeral adds/removes MessageFlagEphemeral to the message flags
func (b *CommandResponseBuilder) SetEphemeral(ephemeral bool) *CommandResponseBuilder {
	if ephemeral {
		b.Flags = MessageFlagEphemeral

	} else {
		b.Flags |= MessageFlagEphemeral
	}
	return b
}

// Build returns your built InteractionResponse
func (b *CommandResponseBuilder) Build() *InteractionResponse {
	return &InteractionResponse{
		Type: InteractionResponseTypeChannelMessageWithSource,
		Data: &b.CommandResponseData,
	}
}
