package api

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Content          string            `json:"content,omitempty"`
	TTS              bool              `json:"tts,omitempty"`
	Embed            *Embed            `json:"embed,omitempty"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
}

// MessageBuilder helper to build Message(s) easier
type MessageBuilder struct {
	MessageCreate
}

// NewMessageBuilder creates a new MessageBuilder to be built later
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		MessageCreate: MessageCreate{
			AllowedMentions: &DefaultMessageAllowedMentions,
		},
	}
}

// SetContent sets content of the Message
func (b *MessageBuilder) SetContent(content string) *MessageBuilder {
	b.Content = content
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

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *MessageBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *MessageBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the Message to nothing
func (b *MessageBuilder) SetAllowedMentionsEmpty() *MessageBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

// SetMessageReference allows you to specify a MessageReference to reply to
func (b *MessageBuilder) SetMessageReference(messageReference *MessageReference) *MessageBuilder {
	b.MessageReference = messageReference
	return b
}

// SetMessageReferenceByMessageID allows you to specify a Message ID to reply to
func (b *MessageBuilder) SetMessageReferenceByMessageID(messageID Snowflake) *MessageBuilder {
	b.MessageReference = &MessageReference{
		MessageID:       &messageID,
	}
	return b
}

// Build builds the MessageBuilder to a MessageCreate struct
func (b *MessageBuilder) Build() MessageCreate {
	return b.MessageCreate
}
