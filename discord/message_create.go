package discord

import "github.com/DisgoOrg/snowflake"

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Nonce            string                `json:"nonce,omitempty"`
	Content          string                `json:"content,omitempty"`
	TTS              bool                  `json:"tts,omitempty"`
	Embeds           []Embed               `json:"embeds,omitempty"`
	Components       []ContainerComponent  `json:"components,omitempty"`
	StickerIDs       []snowflake.Snowflake `json:"sticker_ids,omitempty"`
	Files            []*File               `json:"-"`
	AllowedMentions  *AllowedMentions      `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference     `json:"message_reference,omitempty"`
	Flags            MessageFlags          `json:"flags,omitempty"`
}

func (MessageCreate) interactionCallbackData() {}

// ToBody returns the MessageCreate ready for body
func (m MessageCreate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

func (m MessageCreate) ToResponseBody(response InteractionResponse) (interface{}, error) {
	if len(m.Files) > 0 {
		return PayloadWithFiles(response, m.Files...)
	}
	return response, nil
}
