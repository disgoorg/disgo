package api

import "fmt"

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Nonce            string            `json:"nonce,omitempty"`
	Content          string            `json:"content,omitempty"`
	Components       []Component       `json:"components,omitempty"`
	TTS              bool              `json:"tts,omitempty"`
	Embed            *Embed            `json:"embed,omitempty"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
}

// MessageCreateBuilder helper to build Message(s) easier
type MessageCreateBuilder struct {
	MessageCreate
}

// NewMessageCreateBuilder creates a new MessageCreateBuilder to be built later
func NewMessageCreateBuilder() *MessageCreateBuilder {
	return &MessageCreateBuilder{
		MessageCreate: MessageCreate{
			AllowedMentions: &DefaultMessageAllowedMentions,
		},
	}
}

// NewMessageBuilderByMessage returns a new MessageCreateBuilder and takes an existing Message
func NewMessageBuilderByMessage(message *Message) *MessageCreateBuilder {
	msg := MessageCreate{
		TTS:             message.TTS,
		Components:      message.Components,
		AllowedMentions: &DefaultInteractionAllowedMentions,
	}
	if message.Content != nil {
		msg.Content = *message.Content
	}
	if len(message.Embeds) > 0 {
		msg.Embed = message.Embeds[0]
	}
	return &MessageCreateBuilder{
		MessageCreate: msg,
	}
}

// NewMessageCreateBuilderWithEmbed creates a new MessageCreateBuilder with an Embed to be built later
func NewMessageCreateBuilderWithEmbed(embed *Embed) *MessageCreateBuilder {
	return NewMessageCreateBuilder().SetEmbed(embed)
}

// NewMessageCreateBuilderWithContent creates a new MessageCreateBuilder with a content to be built later
func NewMessageCreateBuilderWithContent(content string) *MessageCreateBuilder {
	return NewMessageCreateBuilder().SetContent(content)
}

// SetContent sets content of the Message
func (b *MessageCreateBuilder) SetContent(content string) *MessageCreateBuilder {
	b.Content = content
	return b
}

// SetContentf sets content of the Message
func (b *MessageCreateBuilder) SetContentf(content string, a ...interface{}) *MessageCreateBuilder {
	b.Content = fmt.Sprintf(content, a...)
	return b
}

// SetTTS sets the text to speech of the Message
func (b *MessageCreateBuilder) SetTTS(tts bool) *MessageCreateBuilder {
	b.TTS = tts
	return b
}

// SetEmbed sets the Embed of the Message
func (b *MessageCreateBuilder) SetEmbed(embed *Embed) *MessageCreateBuilder {
	b.Embed = embed
	return b
}

// SetComponents sets the Component(s) of the Message
func (b *MessageCreateBuilder) SetComponents(components ...Component) *MessageCreateBuilder {
	b.Components = components
	return b
}

// AddComponents adds the Component(s) to the Message
func (b *MessageCreateBuilder) AddComponents(components ...Component) *MessageCreateBuilder {
	b.Components = append(b.Components, components...)
	return b
}

// ClearComponents removes all of the Component(s) of the Message
func (b *MessageCreateBuilder) ClearComponents() *MessageCreateBuilder {
	if b != nil {
		b.Components = []Component{}
	}
	return b
}

// RemoveComponent removes a Component from the Message
func (b *MessageCreateBuilder) RemoveComponent(i int) *MessageCreateBuilder {
	if b != nil && len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *MessageCreateBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *MessageCreateBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *MessageCreateBuilder) ClearAllowedMentions() *MessageCreateBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetMessageReference allows you to specify a MessageReference to reply to
func (b *MessageCreateBuilder) SetMessageReference(messageReference *MessageReference) *MessageCreateBuilder {
	b.MessageReference = messageReference
	return b
}

// SetMessageReferenceByMessageID allows you to specify a Message ID to reply to
func (b *MessageCreateBuilder) SetMessageReferenceByMessageID(messageID Snowflake) *MessageCreateBuilder {
	if b.MessageReference == nil {
		b.MessageReference = &MessageReference{}
	}
	b.MessageReference.MessageID = &messageID
	return b
}

// Build builds the MessageCreateBuilder to a MessageCreate struct
func (b *MessageCreateBuilder) Build() MessageCreate {
	return b.MessageCreate
}
