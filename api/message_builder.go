package api

import "fmt"

// MessageBuilder helper to build Message(s) easier
type MessageBuilder struct {
	MessageCreate
}

// NewMessageBuilder creates a new MessageBuilder to be built later
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		MessageCreate: MessageCreate{
			Nonce:           "test nonce",
			AllowedMentions: &DefaultMessageAllowedMentions,
		},
	}
}

// NewMessageBuilderWithEmbed creates a new MessageBuilder with an Embed to be built later
func NewMessageBuilderWithEmbed(embed *Embed) *MessageBuilder {
	return NewMessageBuilder().SetEmbed(embed)
}

// NewMessageBuilderWithContent creates a new MessageBuilder with a content to be built later
func NewMessageBuilderWithContent(content string) *MessageBuilder {
	return NewMessageBuilder().SetContent(content)
}

// SetContent sets content of the Message
func (b *MessageBuilder) SetContent(content string) *MessageBuilder {
	b.Content = content
	return b
}

// SetContentf sets content of the Message
func (b *MessageBuilder) SetContentf(content string, a ...interface{}) *MessageBuilder {
	b.Content = fmt.Sprintf(content, a...)
	return b
}

// SetTTS sets the text to speech of the Message
func (b *MessageBuilder) SetTTS(tts bool) *MessageBuilder {
	b.TTS = tts
	return b
}

// SetEmbed sets the Embed of the Message
func (b *MessageBuilder) SetEmbed(embed *Embed) *MessageBuilder {
	b.Embed = embed
	return b
}

// SetComponents sets the Component(s) of the Message
func (b *MessageBuilder) SetComponents(components ...Component) *MessageBuilder {
	b.Components = components
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *MessageBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *MessageBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *MessageBuilder) ClearAllowedMentions() *MessageBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetMessageReference allows you to specify a MessageReference to reply to
func (b *MessageBuilder) SetMessageReference(messageReference *MessageReference) *MessageBuilder {
	b.MessageReference = messageReference
	return b
}

// SetMessageReferenceByMessageID allows you to specify a Message ID to reply to
func (b *MessageBuilder) SetMessageReferenceByMessageID(messageID Snowflake) *MessageBuilder {
	if b.MessageReference == nil {
		b.MessageReference = &MessageReference{}
	}
	b.MessageReference.MessageID = &messageID
	return b
}

// Build builds the MessageBuilder to a MessageCreate struct
func (b *MessageBuilder) Build() *MessageCreate {
	return &b.MessageCreate
}
