package discord

import (
	"fmt"
	"io"
	"slices"

	"github.com/disgoorg/snowflake/v2"
)

// NewMessageCreate returns a new MessageCreate to create a message without any fields set.
func NewMessageCreate() MessageCreate {
	return MessageCreate{}
}

// NewMessageCreateV2 returns a new MessageCreate with [MessageFlagIsComponentsV2] flag set & allows to directly pass components.
func NewMessageCreateV2(components ...LayoutComponent) MessageCreate {
	return MessageCreate{
		Flags:      MessageFlagIsComponentsV2,
		Components: components,
	}
}

// MessageCreate is the struct to create a new Message with.
type MessageCreate struct {
	Nonce            string             `json:"nonce,omitempty"`
	Content          string             `json:"content,omitempty"`
	TTS              bool               `json:"tts,omitempty"`
	Embeds           []Embed            `json:"embeds,omitempty"`
	Components       []LayoutComponent  `json:"components,omitempty"`
	StickerIDs       []snowflake.ID     `json:"sticker_ids,omitempty"`
	Files            []*File            `json:"-"`
	Attachments      []AttachmentCreate `json:"attachments,omitempty"`
	AllowedMentions  *AllowedMentions   `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference  `json:"message_reference,omitempty"`
	Flags            MessageFlags       `json:"flags,omitempty"`
	EnforceNonce     bool               `json:"enforce_nonce,omitempty"`
	Poll             *PollCreate        `json:"poll,omitempty"`
}

func (MessageCreate) interactionCallbackData() {}

// ToBody returns the MessageCreate ready for body.
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

// WithContent returns a new MessageCreate with the provided content.
func (m MessageCreate) WithContent(content string) MessageCreate {
	m.Content = content
	return m
}

// WithContentf returns a new MessageCreate with the formatted content.
func (m MessageCreate) WithContentf(content string, a ...any) MessageCreate {
	return m.WithContent(fmt.Sprintf(content, a...))
}

// WithNonce returns a new MessageCreate with the provided nonce.
func (m MessageCreate) WithNonce(nonce string) MessageCreate {
	m.Nonce = nonce
	return m
}

// WithEnforceNonce returns a new MessageCreate with the provided enforce nonce setting (use with WithNonce).
func (m MessageCreate) WithEnforceNonce(enforce bool) MessageCreate {
	m.EnforceNonce = enforce
	return m
}

// WithTTS returns a new MessageCreate with the provided TTS setting.
func (m MessageCreate) WithTTS(tts bool) MessageCreate {
	m.TTS = tts
	return m
}

// WithEmbeds returns a new MessageCreate with the provided Embed(s).
func (m MessageCreate) WithEmbeds(embeds ...Embed) MessageCreate {
	m.Embeds = embeds
	return m
}

// WithEmbed returns a new MessageCreate with the provided Embed at the index.
func (m MessageCreate) WithEmbed(i int, embed Embed) MessageCreate {
	if len(m.Embeds) > i {
		m.Embeds = slices.Insert(m.Embeds, i, embed)
	}
	return m
}

// AddEmbeds returns a new MessageCreate with the provided embeds added.
func (m MessageCreate) AddEmbeds(embeds ...Embed) MessageCreate {
	m.Embeds = append(m.Embeds, embeds...)
	return m
}

// ClearEmbeds returns a new MessageCreate with no embeds.
func (m MessageCreate) ClearEmbeds() MessageCreate {
	m.Embeds = []Embed{}
	return m
}

// RemoveEmbed returns a new MessageCreate with the embed at the index removed.
func (m MessageCreate) RemoveEmbed(i int) MessageCreate {
	if len(m.Embeds) > i {
		m.Embeds = slices.Delete(slices.Clone(m.Embeds), i, i+1)
	}
	return m
}

// WithComponents returns a new MessageCreate with the provided LayoutComponent(s).
func (m MessageCreate) WithComponents(components ...LayoutComponent) MessageCreate {
	m.Components = components
	return m
}

// UpdateComponent returns a new MessageCreate with the provided LayoutComponent at the index.
func (m MessageCreate) UpdateComponent(id int, component LayoutComponent) MessageCreate {
	for i, cc := range m.Components {
		if cc.GetID() == id {
			m.Components = slices.Clone(m.Components)
			m.Components[i] = component
			return m
		}
	}
	return m
}

// AddComponents returns a new MessageCreate with the provided LayoutComponent(s) added.
func (m MessageCreate) AddComponents(containers ...LayoutComponent) MessageCreate {
	m.Components = append(m.Components, containers...)
	return m
}

// RemoveComponent returns a new MessageCreate with the LayoutComponent at the index removed.
func (m MessageCreate) RemoveComponent(id int) MessageCreate {
	for i, cc := range m.Components {
		if cc.GetID() == id {
			m.Components = slices.Delete(slices.Clone(m.Components), i, i+1)
			return m
		}
	}
	return m
}

// AddActionRow returns a new MessageCreate with a new ActionRowComponent containing the provided InteractiveComponent(s) added.
func (m MessageCreate) AddActionRow(components ...InteractiveComponent) MessageCreate {
	m.Components = append(m.Components, NewActionRow(components...))
	return m
}

// ClearComponents returns a new MessageCreate with no LayoutComponent(s).
func (m MessageCreate) ClearComponents() MessageCreate {
	m.Components = []LayoutComponent{}
	return m
}

// WithStickers returns a new MessageCreate with the provided stickers.
func (m MessageCreate) WithStickers(stickerIds ...snowflake.ID) MessageCreate {
	m.StickerIDs = stickerIds
	return m
}

// AddStickers returns a new MessageCreate with the provided stickers added.
func (m MessageCreate) AddStickers(stickerIds ...snowflake.ID) MessageCreate {
	m.StickerIDs = append(m.StickerIDs, stickerIds...)
	return m
}

// RemoveSticker returns a new MessageCreate with the provided sticker removed.
func (m MessageCreate) RemoveSticker(stickerId snowflake.ID) MessageCreate {
	m.StickerIDs = slices.DeleteFunc(slices.Clone(m.StickerIDs), func(id snowflake.ID) bool {
		return id == stickerId
	})
	return m
}

// ClearStickers returns a new MessageCreate with no Sticker(s).
func (m MessageCreate) ClearStickers() MessageCreate {
	m.StickerIDs = []snowflake.ID{}
	return m
}

// WithFiles returns a new MessageCreate with the provided File(s).
func (m MessageCreate) WithFiles(files ...*File) MessageCreate {
	m.Files = files
	return m
}

// UpdateFile returns a new MessageCreate with the provided File at the index.
func (m MessageCreate) UpdateFile(i int, file *File) MessageCreate {
	if len(m.Files) > i {
		m.Files = slices.Clone(m.Files)
		m.Files[i] = file
	}
	return m
}

// AddFiles returns a new MessageCreate with the File(s) added.
func (m MessageCreate) AddFiles(files ...*File) MessageCreate {
	m.Files = append(m.Files, files...)
	return m
}

// AddFile returns a new MessageCreate with a File added.
func (m MessageCreate) AddFile(name string, description string, reader io.Reader, flags ...FileFlags) MessageCreate {
	m.Files = append(m.Files, NewFile(name, description, reader, flags...))
	return m
}

// RemoveFile returns a new MessageCreate with the File at the index removed.
func (m MessageCreate) RemoveFile(i int) MessageCreate {
	if len(m.Files) > i {
		m.Files = slices.Delete(slices.Clone(m.Files), i, i+1)
	}
	return m
}

// ClearFiles returns a new MessageCreate with no File(s).
func (m MessageCreate) ClearFiles() MessageCreate {
	m.Files = []*File{}
	return m
}

// WithAllowedMentions returns a new MessageCreate with the provided AllowedMentions.
func (m MessageCreate) WithAllowedMentions(allowedMentions *AllowedMentions) MessageCreate {
	m.AllowedMentions = allowedMentions
	return m
}

// ClearAllowedMentions returns a new MessageCreate with no AllowedMentions.
func (m MessageCreate) ClearAllowedMentions() MessageCreate {
	return m.WithAllowedMentions(nil)
}

// WithMessageReference returns a new MessageCreate with the provided MessageReference to reply to.
func (m MessageCreate) WithMessageReference(messageReference *MessageReference) MessageCreate {
	m.MessageReference = messageReference
	return m
}

// WithMessageReferenceByID returns a new MessageCreate with a MessageReference to the provided Message ID to reply to.
func (m MessageCreate) WithMessageReferenceByID(messageID snowflake.ID) MessageCreate {
	m.MessageReference = &MessageReference{
		MessageID: &messageID,
	}
	return m
}

// WithFlags returns a new MessageCreate with the provided message flags.
func (m MessageCreate) WithFlags(flags ...MessageFlags) MessageCreate {
	m.Flags = m.Flags.Add(flags...)
	return m
}

// AddFlags returns a new MessageCreate with the provided MessageFlags added.
func (m MessageCreate) AddFlags(flags ...MessageFlags) MessageCreate {
	m.Flags = m.Flags.Add(flags...)
	return m
}

// RemoveFlags returns a new MessageCreate with the provided MessageFlags removed.
func (m MessageCreate) RemoveFlags(flags ...MessageFlags) MessageCreate {
	m.Flags = m.Flags.Remove(flags...)
	return m
}

// ClearFlags returns a new MessageCreate with no MessageFlags.
func (m MessageCreate) ClearFlags() MessageCreate {
	return m.WithFlags(MessageFlagsNone)
}

// WithEphemeral returns a new MessageCreate with MessageFlagEphemeral added/removed.
func (m MessageCreate) WithEphemeral(ephemeral bool) MessageCreate {
	if ephemeral {
		m.Flags = m.Flags.Add(MessageFlagEphemeral)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagEphemeral)
	}
	return m
}

// WithIsComponentsV2 returns a new MessageCreate with MessageFlagIsComponentsV2 added/removed.
// Once a message with the flag has been sent, it cannot be removed by editing the message.
func (m MessageCreate) WithIsComponentsV2(isComponentV2 bool) MessageCreate {
	if isComponentV2 {
		m.Flags = m.Flags.Add(MessageFlagIsComponentsV2)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagIsComponentsV2)
	}
	return m
}

// WithSuppressEmbeds returns a new MessageCreate with MessageFlagSuppressEmbeds added/removed.
func (m MessageCreate) WithSuppressEmbeds(suppressEmbeds bool) MessageCreate {
	if suppressEmbeds {
		m.Flags = m.Flags.Add(MessageFlagSuppressEmbeds)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagSuppressEmbeds)
	}
	return m
}

// WithSuppressNotifications returns a new MessageCreate with MessageFlagSuppressNotifications added/removed.
func (m MessageCreate) WithSuppressNotifications(suppressNotifications bool) MessageCreate {
	if suppressNotifications {
		m.Flags = m.Flags.Add(MessageFlagSuppressNotifications)
	} else {
		m.Flags = m.Flags.Remove(MessageFlagSuppressNotifications)
	}
	return m
}

// WithPoll returns a new MessageCreate with the provided Poll.
func (m MessageCreate) WithPoll(poll PollCreate) MessageCreate {
	m.Poll = &poll
	return m
}

// ClearPoll returns a new MessageCreate with no Poll.
func (m MessageCreate) ClearPoll() MessageCreate {
	m.Poll = nil
	return m
}
