package discord

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/disgoorg/validate"
)

const (
	MessageContentMaxLength = 2000
	MessageNonceMaxLength   = 25
	MessageMaxEmbeds        = 10
	MessageMaxComponents    = 5
	MessageMaxStickers      = 3
	MessageMaxFiles         = 10
)

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Nonce            string               `json:"nonce,omitempty"`
	Content          string               `json:"content,omitempty"`
	TTS              bool                 `json:"tts,omitempty"`
	Embeds           []Embed              `json:"embeds,omitempty"`
	Components       []ContainerComponent `json:"components,omitempty"`
	StickerIDs       []snowflake.ID       `json:"sticker_ids,omitempty"`
	Files            []*File              `json:"-"`
	Attachments      []AttachmentCreate   `json:"attachments,omitempty"`
	AllowedMentions  *AllowedMentions     `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference    `json:"message_reference,omitempty"`
	Flags            MessageFlags         `json:"flags,omitempty"`
}

func (MessageCreate) interactionCallbackData() {}

func (m MessageCreate) Validate() error {
	return validate.Validate(
		validate.Value(m.Nonce, validate.StringRange(0, MessageNonceMaxLength)),
		validate.Value(m.Content, validate.StringRange(0, MessageContentMaxLength)),
		validate.Value(m.Embeds, validate.SliceMaxLen[Embed](MessageMaxEmbeds)),
		validate.Slice(m.Embeds),
		validate.Value(m.Components, validate.SliceMaxLen[ContainerComponent](MessageMaxComponents)),
		validate.Slice(m.Components),
		validate.Value(m.StickerIDs, validate.SliceMaxLen[snowflake.ID](MessageMaxStickers)),
		validate.Slice(m.StickerIDs),
		validate.Value(m.Files, validate.SliceMaxLen[*File](MessageMaxFiles)),
		validate.Slice(m.Files),
		validate.Value(m.Attachments, validate.SliceMaxLen[AttachmentCreate](MessageMaxFiles)),
		validate.Slice(m.Attachments),
		validate.Value(m.Flags, validate.AllowedFlags[MessageFlags](MessageFlagCrossposted, MessageFlagIsCrosspost, MessageFlagSuppressEmbeds, MessageFlagSourceMessageDeleted, MessageFlagUrgent)),
	)
}

// ToBody returns the MessageCreate ready for body
func (m MessageCreate) ToBody() (any, error) {
	if len(m.Files) > 0 {
		m.Attachments = parseAttachments(m.Files)
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

func (m MessageCreate) ToResponseBody(response InteractionResponse) (any, error) {
	if len(m.Files) > 0 {
		m.Attachments = parseAttachments(m.Files)
		response.Data = m
		return PayloadWithFiles(response, m.Files...)
	}
	return response, nil
}
