package api

type MessageCreate struct {
	Content          string            `json:"content,omitempty"`
	TTS              bool              `json:"tts,omitempty"`
	Embed            *Embed            `json:"embed,omitempty"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
}

type MessageBuilder struct {
	MessageCreate
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		MessageCreate: MessageCreate{
			AllowedMentions: &DefaultMessageAllowedMentions,
		},
	}
}

func (b *MessageBuilder) SetContent(content string) *MessageBuilder {
	b.Content = content
	return b
}

func (b *MessageBuilder) SetTTS(tts bool) *MessageBuilder {
	b.TTS = tts
	return b
}

func (b *MessageBuilder) SetEmbed(embed *Embed) *MessageBuilder {
	b.Embed = embed
	return b
}

func (b *MessageBuilder) SetAllowedMentions(allowedMentions *AllowedMentions) *MessageBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// SetAllowedMentionsEmpty sets the allowed mentions of the Message to nothing
func (b *MessageBuilder) SetAllowedMentionsEmpty() *MessageBuilder {
	return b.SetAllowedMentions(&AllowedMentions{})
}

func (b *MessageBuilder) SetMessageReference(messageReference *MessageReference) *MessageBuilder {
	b.MessageReference = messageReference
	return b
}

func (b *MessageBuilder) SetMessageReferenceByMessageID(messageID Snowflake) *MessageBuilder {
	b.MessageReference = &MessageReference{
		MessageID:       &messageID,
	}
	return b
}

func (b *MessageBuilder) Build() *MessageCreate {
	return &b.MessageCreate
}
