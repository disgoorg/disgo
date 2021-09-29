package core

import (
	"fmt"
	"io"

	"github.com/DisgoOrg/disgo/discord"
)

// MessageCreateBuilder helper to build Message(s) easier
type MessageCreateBuilder struct {
	discord.MessageCreate
	Components []Component
}

// NewMessageCreateBuilder creates a new MessageCreateBuilder to be built later
//goland:noinspection GoUnusedExportedFunction
func NewMessageCreateBuilder() *MessageCreateBuilder {
	return &MessageCreateBuilder{
		MessageCreate: discord.MessageCreate{
			AllowedMentions: &DefaultAllowedMentions,
		},
	}
}

// SetContent sets the content of the Message
func (b *MessageCreateBuilder) SetContent(content string) *MessageCreateBuilder {
	b.Content = content
	return b
}

// SetContentf sets the content of the Message but with format
func (b *MessageCreateBuilder) SetContentf(content string, a ...interface{}) *MessageCreateBuilder {
	return b.SetContent(fmt.Sprintf(content, a...))
}

// SetTTS sets whether the Message should be text to speech
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

// ClearEmbeds removes all the embeds from the Message
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
func (b *MessageCreateBuilder) SetActionRows(actionRows ...ActionRow) *MessageCreateBuilder {
	b.Components = actionRowsToComponents(actionRows)
	return b
}

// SetActionRow sets the provided ActionRow at the index of Component(s)
func (b *MessageCreateBuilder) SetActionRow(i int, actionRow ActionRow) *MessageCreateBuilder {
	if len(b.Components) > i {
		b.Components[i] = actionRow
	}
	return b
}

// AddActionRow adds a new ActionRow with the provided Component(s) to the Message
func (b *MessageCreateBuilder) AddActionRow(components ...Component) *MessageCreateBuilder {
	b.Components = append(b.Components, NewActionRow(components...))
	return b
}

// AddActionRows adds the ActionRow(s) to the Message
func (b *MessageCreateBuilder) AddActionRows(actionRows ...ActionRow) *MessageCreateBuilder {
	b.Components = append(b.Components, actionRowsToComponents(actionRows)...)
	return b
}

// RemoveActionRow removes an ActionRow from the Message
func (b *MessageCreateBuilder) RemoveActionRow(i int) *MessageCreateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// ClearActionRows removes all the ActionRow(s) of the Message
func (b *MessageCreateBuilder) ClearActionRows() *MessageCreateBuilder {
	b.Components = []Component{}
	return b
}

// AddStickers adds provided stickers to the Message
func (b *MessageCreateBuilder) AddStickers(stickerIds ...discord.Snowflake) *MessageCreateBuilder {
	b.StickerIDs = append(b.StickerIDs, stickerIds...)
	return b
}

// SetStickers sets the stickers of the Message
func (b *MessageCreateBuilder) SetStickers(stickerIds ...discord.Snowflake) *MessageCreateBuilder {
	b.StickerIDs = stickerIds
	return b
}

// ClearStickers removes all Sticker(s) from the Message
func (b *MessageCreateBuilder) ClearStickers() *MessageCreateBuilder {
	b.StickerIDs = []discord.Snowflake{}
	return b
}

// SetFiles sets the File(s) for this Message
func (b *MessageCreateBuilder) SetFiles(files ...*discord.File) *MessageCreateBuilder {
	b.Files = files
	return b
}

// SetFile sets the File at the index for the Message
func (b *MessageCreateBuilder) SetFile(i int, file *discord.File) *MessageCreateBuilder {
	if len(b.Files) > i {
		b.Files[i] = file
	}
	return b
}

// AddFiles adds the discord.File(s) to the Message
func (b *MessageCreateBuilder) AddFiles(files ...*discord.File) *MessageCreateBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a discord.File to the Message
func (b *MessageCreateBuilder) AddFile(name string, reader io.Reader, flags ...discord.FileFlags) *MessageCreateBuilder {
	b.Files = append(b.Files, discord.NewFile(name, reader, flags...))
	return b
}

// ClearFiles removes all discord.File(s) of the Message
func (b *MessageCreateBuilder) ClearFiles() *MessageCreateBuilder {
	b.Files = []*discord.File{}
	return b
}

// RemoveFiles removes the discord.File at this index
func (b *MessageCreateBuilder) RemoveFiles(i int) *MessageCreateBuilder {
	if len(b.Files) > i {
		b.Files = append(b.Files[:i], b.Files[i+1:]...)
	}
	return b
}

// SetAllowedMentions sets the discord.AllowedMentions of the Message
func (b *MessageCreateBuilder) SetAllowedMentions(allowedMentions *discord.AllowedMentions) *MessageCreateBuilder {
	b.AllowedMentions = allowedMentions
	return b
}

// ClearAllowedMentions clears the discord.AllowedMentions of the Message
func (b *MessageCreateBuilder) ClearAllowedMentions() *MessageCreateBuilder {
	return b.SetAllowedMentions(nil)
}

// SetMessageReference allows you to specify a discord.MessageReference to reply to
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

// SetFlags sets the discord.MessageFlags of the Message
func (b *MessageCreateBuilder) SetFlags(flags discord.MessageFlags) *MessageCreateBuilder {
	b.Flags = flags
	return b
}

// AddFlags adds the discord.MessageFlags of the Message
func (b *MessageCreateBuilder) AddFlags(flags ...discord.MessageFlags) *MessageCreateBuilder {
	b.Flags = b.Flags.Add(flags...)
	return b
}

// RemoveFlags removes the discord.MessageFlags of the Message
func (b *MessageCreateBuilder) RemoveFlags(flags ...discord.MessageFlags) *MessageCreateBuilder {
	b.Flags = b.Flags.Remove(flags...)
	return b
}

// ClearFlags clears the discord.MessageFlags of the Message
func (b *MessageCreateBuilder) ClearFlags() *MessageCreateBuilder {
	return b.SetFlags(discord.MessageFlagNone)
}

// SetEphemeral adds/removes discord.MessageFlagEphemeral to the Message flags
func (b *MessageCreateBuilder) SetEphemeral(ephemeral bool) *MessageCreateBuilder {
	if ephemeral {
		b.Flags = b.Flags.Add(discord.MessageFlagEphemeral)
	} else {
		b.Flags = b.Flags.Remove(discord.MessageFlagEphemeral)
	}
	return b
}

// Build builds the MessageCreateBuilder to a discord.MessageCreate struct
func (b *MessageCreateBuilder) Build() discord.MessageCreate {
	b.MessageCreate.Components = b.Components
	return b.MessageCreate
}

func actionRowsToComponents(actionRows []ActionRow) []Component {
	components := make([]Component, len(actionRows))
	for i := range actionRows {
		components[i] = actionRows[i]
	}
	return components
}
