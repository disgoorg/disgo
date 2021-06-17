package api

import (
	"fmt"
	"io"

	"github.com/DisgoOrg/restclient"
)

// MessageCreate is the struct to create a new Message with
type MessageCreate struct {
	Nonce            string            `json:"nonce,omitempty"`
	Content          string            `json:"content,omitempty"`
	TTS              bool              `json:"tts,omitempty"`
	Embeds           []Embed           `json:"embeds,omitempty"`
	Components       []Component       `json:"components,omitempty"`
	Files            []restclient.File `json:"-"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Flags            MessageFlags      `json:"flags,omitempty"`
}

// ToBody returns the MessageCreate ready for body
func (m MessageCreate) ToBody() (interface{}, error) {
	if len(m.Files) > 0 {
		return restclient.PayloadWithFiles(m, m.Files...)
	}
	return m, nil
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
		Embeds:          message.Embeds,
		AllowedMentions: &DefaultMessageAllowedMentions,
	}
	if message.Content != nil {
		msg.Content = *message.Content
	}
	return &MessageCreateBuilder{
		MessageCreate: msg,
	}
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

// ClearContent removes content of the Message
func (b *MessageCreateBuilder) ClearContent() *MessageCreateBuilder {
	return b.SetContent("")
}

// SetTTS sets the text to speech of the Message
func (b *MessageCreateBuilder) SetTTS(tts bool) *MessageCreateBuilder {
	b.TTS = tts
	return b
}

// SetEmbeds sets the embeds of the Message
func (b *MessageCreateBuilder) SetEmbeds(embeds ...Embed) *MessageCreateBuilder {
	b.Embeds = embeds
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *MessageCreateBuilder) AddEmbeds(embeds ...Embed) *MessageCreateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the Message
func (b *MessageCreateBuilder) ClearEmbeds() *MessageCreateBuilder {
	b.Embeds = []Embed{}
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *MessageCreateBuilder) RemoveEmbed(index int) *MessageCreateBuilder {
	if b != nil && len(b.Embeds) > index {
		b.Embeds = append(b.Embeds[:index], b.Embeds[index+1:]...)
	}
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
	b.Components = []Component{}
	return b
}

// RemoveComponent removes a Component from the Message
func (b *MessageCreateBuilder) RemoveComponent(i int) *MessageCreateBuilder {
	if b != nil && len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// SetFiles sets the files for this WebhookMessageCreate
func (b *MessageCreateBuilder) SetFiles(files ...restclient.File) *MessageCreateBuilder {
	b.Files = files
	return b
}

// AddFiles adds the files to the WebhookMessageCreate
func (b *MessageCreateBuilder) AddFiles(files ...restclient.File) *MessageCreateBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a file to the WebhookMessageCreate
func (b *MessageCreateBuilder) AddFile(name string, reader io.Reader, flags ...restclient.FileFlags) *MessageCreateBuilder {
	b.Files = append(b.Files, restclient.File{
		Name:   name,
		Reader: reader,
		Flags:  restclient.FileFlagNone.Add(flags...),
	})
	return b
}

// ClearFiles removes all files of this WebhookMessageCreate
func (b *MessageCreateBuilder) ClearFiles() *MessageCreateBuilder {
	b.Files = []restclient.File{}
	return b
}

// RemoveFiles removes the file at this index
func (b *MessageCreateBuilder) RemoveFiles(i int) *MessageCreateBuilder {
	if b != nil && len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
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
	return b.SetAllowedMentions(nil)
}

// SetMessageReference allows you to specify a MessageReference to reply to
func (b *MessageCreateBuilder) SetMessageReference(messageReference *MessageReference) *MessageCreateBuilder {
	b.MessageReference = messageReference
	return b
}

// SetMessageReferenceByID allows you to specify a Message ID to reply to
func (b *MessageCreateBuilder) SetMessageReferenceByID(messageID Snowflake) *MessageCreateBuilder {
	if b.MessageReference == nil {
		b.MessageReference = &MessageReference{}
	}
	b.MessageReference.MessageID = &messageID
	return b
}

// SetFlags sets the message flags of the Message
func (b *MessageCreateBuilder) SetFlags(flags MessageFlags) *MessageCreateBuilder {
	b.Flags = flags
	return b
}

// AddFlags adds the MessageFlags of the Message
func (b *MessageCreateBuilder) AddFlags(flags ...MessageFlags) *MessageCreateBuilder {
	b.Flags = b.Flags.Add(flags...)
	return b
}

// RemoveFlags removes the MessageFlags of the Message
func (b *MessageCreateBuilder) RemoveFlags(flags ...MessageFlags) *MessageCreateBuilder {
	b.Flags = b.Flags.Remove(flags...)
	return b
}

// ClearFlags clears the MessageFlags of the Message
func (b *MessageCreateBuilder) ClearFlags() *MessageCreateBuilder {
	return b.SetFlags(MessageFlagNone)
}

// SetEphemeral adds/removes MessageFlagEphemeral to the Message flags
func (b *MessageCreateBuilder) SetEphemeral(ephemeral bool) *MessageCreateBuilder {
	if ephemeral {
		b.Flags = b.Flags.Add(MessageFlagEphemeral)
	} else {
		b.Flags = b.Flags.Remove(MessageFlagEphemeral)
	}
	return b
}

// Build builds the MessageCreateBuilder to a MessageCreate struct
func (b *MessageCreateBuilder) Build() MessageCreate {
	return b.MessageCreate
}
