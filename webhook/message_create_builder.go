package webhook

import (
	"fmt"
	"io"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// MessageCreateBuilder helper to build Message(s) easier
type MessageCreateBuilder struct {
	discord.MessageCreate
	Components []core.Component
}

// NewMessageCreateBuilder creates a new MessageCreateBuilder to be built later
func NewMessageCreateBuilder() *MessageCreateBuilder {
	return &MessageCreateBuilder{
		MessageCreate: discord.MessageCreate{
			AllowedMentions: &DefaultAllowedMentions,
		},
	}
}

// NewMessageCreateBuilderByMessage returns a new MessageCreateBuilder and takes an existing Message
//goland:noinspection GoUnusedExportedFunction
func NewMessageCreateBuilderByMessage(message *Message) *MessageCreateBuilder {
	msg := discord.MessageCreate{
		TTS:             message.TTS,
		Components:      message.Components,
		Embeds:          message.Embeds,
		AllowedMentions: &DefaultAllowedMentions,
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
	return b.SetContent(fmt.Sprintf(content, a...))
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

// SetEmbeds sets the Embed(s) of the Message
func (b *MessageCreateBuilder) SetEmbeds(embeds ...discord.Embed) *MessageCreateBuilder {
	b.Embeds = embeds
	return b
}

// SetEmbed sets the provided Embed at the index of the Message
func (b *MessageCreateBuilder) SetEmbed(i int, embed discord.Embed) *MessageCreateBuilder {
	if len(b.Embeds) > i {
		b.Embeds[i] = embed
	}
	return b
}

// AddEmbeds adds multiple embeds to the Message
func (b *MessageCreateBuilder) AddEmbeds(embeds ...discord.Embed) *MessageCreateBuilder {
	b.Embeds = append(b.Embeds, embeds...)
	return b
}

// ClearEmbeds removes all of the embeds from the Message
func (b *MessageCreateBuilder) ClearEmbeds() *MessageCreateBuilder {
	b.Embeds = []discord.Embed{}
	return b
}

// RemoveEmbed removes an embed from the Message
func (b *MessageCreateBuilder) RemoveEmbed(i int) *MessageCreateBuilder {
	if len(b.Embeds) > i {
		b.Embeds = append(b.Embeds[:i], b.Embeds[i+1:]...)
	}
	return b
}

// SetActionRows sets the ActionRow(s) of the Message
func (b *MessageCreateBuilder) SetActionRows(actionRows ...core.ActionRow) *MessageCreateBuilder {
	b.Components = actionRowsToComponents(actionRows)
	return b
}

// SetActionRow sets the provided ActionRow at the index of Component(s)
func (b *MessageCreateBuilder) SetActionRow(i int, actionRow core.ActionRow) *MessageCreateBuilder {
	if len(b.Components) > i {
		b.Components[i] = actionRow
	}
	return b
}

// AddActionRow adds a new ActionRow with the provided Component(s) to the Message
func (b *MessageCreateBuilder) AddActionRow(components ...core.Component) *MessageCreateBuilder {
	b.Components = append(b.Components, core.NewActionRow(components...))
	return b
}

// AddActionRows adds the ActionRow(s) to the Message
func (b *MessageCreateBuilder) AddActionRows(actionRows ...core.ActionRow) *MessageCreateBuilder {
	b.Components = append(b.Components, actionRowsToComponents(actionRows)...)
	return b
}

// RemoveActionRow removes a ActionRow from the Message
func (b *MessageCreateBuilder) RemoveActionRow(i int) *MessageCreateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// ClearActionRows removes all the ActionRow(s) of the Message
func (b *MessageCreateBuilder) ClearActionRows() *MessageCreateBuilder {
	b.Components = []core.Component{}
	return b
}

// SetFiles sets the File(s) for this MessageCreate
func (b *MessageCreateBuilder) SetFiles(files ...discord.File) *MessageCreateBuilder {
	b.Files = files
	return b
}

// SetFile sets the File at the index for this MessageCreate
func (b *MessageCreateBuilder) SetFile(i int, file discord.File) *MessageCreateBuilder {
	if len(b.Files) > i {
		b.Files[i] = file
	}
	return b
}

// AddFiles adds the File(s) to the MessageCreate
func (b *MessageCreateBuilder) AddFiles(files ...discord.File) *MessageCreateBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a File to the MessageCreate
func (b *MessageCreateBuilder) AddFile(name string, reader io.Reader, flags ...discord.FileFlags) *MessageCreateBuilder {
	b.Files = append(b.Files, discord.File{
		Name:   name,
		Reader: reader,
		Flags:  discord.FileFlagNone.Add(flags...),
	})
	return b
}

// ClearFiles removes all files of this MessageCreate
func (b *MessageCreateBuilder) ClearFiles() *MessageCreateBuilder {
	b.Files = []discord.File{}
	return b
}

// RemoveFiles removes the file at this index
func (b *MessageCreateBuilder) RemoveFiles(i int) *MessageCreateBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	return b
}

// SetAllowedMentions sets the AllowedMentions of the Message
func (b *MessageCreateBuilder) SetAllowedMentions(allowedMentions *discord.AllowedMentions) *MessageCreateBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the allowed mentions of the Message
func (b *MessageCreateBuilder) ClearAllowedMentions() *MessageCreateBuilder {
	return b.SetAllowedMentions(nil)
}

// SetMessageReference allows you to specify a MessageReference to reply to
func (b *MessageCreateBuilder) SetMessageReference(messageReference *discord.MessageReference) *MessageCreateBuilder {
	b.MessageReference = messageReference
	return b
}

// SetMessageReferenceByID allows you to specify a Message ID to reply to
func (b *MessageCreateBuilder) SetMessageReferenceByID(messageID discord.Snowflake) *MessageCreateBuilder {
	if b.MessageReference == nil {
		b.MessageReference = &discord.MessageReference{}
	}
	b.MessageReference.MessageID = &messageID
	return b
}

// SetFlags sets the message flags of the Message
func (b *MessageCreateBuilder) SetFlags(flags discord.MessageFlags) *MessageCreateBuilder {
	b.Flags = flags
	return b
}

// AddFlags adds the MessageFlags of the Message
func (b *MessageCreateBuilder) AddFlags(flags ...discord.MessageFlags) *MessageCreateBuilder {
	b.Flags = b.Flags.Add(flags...)
	return b
}

// RemoveFlags removes the MessageFlags of the Message
func (b *MessageCreateBuilder) RemoveFlags(flags ...discord.MessageFlags) *MessageCreateBuilder {
	b.Flags = b.Flags.Remove(flags...)
	return b
}

// ClearFlags clears the MessageFlags of the Message
func (b *MessageCreateBuilder) ClearFlags() *MessageCreateBuilder {
	return b.SetFlags(discord.MessageFlagNone)
}

// SetEphemeral adds/removes MessageFlagEphemeral to the Message flags
func (b *MessageCreateBuilder) SetEphemeral(ephemeral bool) *MessageCreateBuilder {
	if ephemeral {
		b.Flags = b.Flags.Add(discord.MessageFlagEphemeral)
	} else {
		b.Flags = b.Flags.Remove(discord.MessageFlagEphemeral)
	}
	return b
}

// Build builds the MessageCreateBuilder to a MessageCreate struct
func (b *MessageCreateBuilder) Build() discord.MessageCreate {
	b.MessageCreate.Components = b.Components
	return b.MessageCreate
}

func actionRowsToComponents(actionRows []core.ActionRow) []core.Component {
	components := make([]core.Component, len(actionRows))
	for i := range actionRows {
		components[i] = components[i]
	}
	return components
}
