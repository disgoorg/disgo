package discord

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Nonce            string            `json:"nonce,omitempty"`
	Content          string            `json:"content,omitempty"`
	TTS              bool              `json:"tts,omitempty"`
	Embeds           []Embed           `json:"embeds,omitempty"`
	Components       []Component       `json:"components,omitempty"`
	StickerIDs       []Snowflake       `json:"sticker_ids,omitempty"`
	Files            []*File           `json:"-"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Flags            MessageFlags      `json:"flags,omitempty"`
}

// ToBody returns the MessageCreate ready for body
func (m MessageCreate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}
