package webhook

import (
	"fmt"
	"io"

	"github.com/DisgoOrg/disgo/discord"
)

// MessageCreateBuilder helper to build Message(s) easier
type MessageCreateBuilder struct {
	discord.WebhookMessageCreate
}

// NewMessageCreateBuilder creates a new MessageCreateBuilder to be built later
//goland:noinspection GoUnusedExportedFunction
func NewMessageCreateBuilder() *MessageCreateBuilder {
	return &MessageCreateBuilder{
		WebhookMessageCreate: discord.WebhookMessageCreate{
			AllowedMentions: &DefaultAllowedMentions,
		},
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

func (b *MessageCreateBuilder) SetUsername(username string) *MessageCreateBuilder {
	b.Username = username
	return b
}

func (b *MessageCreateBuilder) SetAvatarURL(url string) *MessageCreateBuilder {
	b.AvatarURL = url
	return b
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

// SetContainerComponents sets the discord.ContainerComponent(s) of the Message
func (b *MessageCreateBuilder) SetContainerComponents(containerComponents ...discord.ContainerComponent) *MessageCreateBuilder {
	b.Components = containerComponents
	return b
}

// SetContainerComponent sets the provided discord.InteractiveComponent at the index of discord.InteractiveComponent(s)
func (b *MessageCreateBuilder) SetContainerComponent(i int, container discord.ContainerComponent) *MessageCreateBuilder {
	if len(b.Components) > i {
		b.Components[i] = container
	}
	return b
}

// AddActionRow adds a new discord.ActionRowComponent with the provided discord.InteractiveComponent(s) to the Message
func (b *MessageCreateBuilder) AddActionRow(components ...discord.InteractiveComponent) *MessageCreateBuilder {
	b.Components = append(b.Components, discord.ActionRowComponent(components))
	return b
}

// AddContainerComponents adds the discord.ContainerComponent(s) to the Message
func (b *MessageCreateBuilder) AddContainerComponents(containers ...discord.ContainerComponent) *MessageCreateBuilder {
	b.Components = append(b.Components, containers...)
	return b
}

// RemoveContainerComponent removes a discord.ActionRowComponent from the Message
func (b *MessageCreateBuilder) RemoveContainerComponent(i int) *MessageCreateBuilder {
	if len(b.Components) > i {
		b.Components = append(b.Components[:i], b.Components[i+1:]...)
	}
	return b
}

// ClearContainerComponents removes all the discord.ContainerComponent(s) of the Message
func (b *MessageCreateBuilder) ClearContainerComponents() *MessageCreateBuilder {
	b.Components = []discord.ContainerComponent{}
	return b
}

// SetFiles sets the File(s) for this MessageCreate
func (b *MessageCreateBuilder) SetFiles(files ...*discord.File) *MessageCreateBuilder {
	b.Files = files
	return b
}

// SetFile sets the discord.File at the index for this discord.MessageCreate
func (b *MessageCreateBuilder) SetFile(i int, file *discord.File) *MessageCreateBuilder {
	if len(b.Files) > i {
		b.Files[i] = file
	}
	return b
}

// AddFiles adds the discord.File(s) to the discord.MessageCreate
func (b *MessageCreateBuilder) AddFiles(files ...*discord.File) *MessageCreateBuilder {
	b.Files = append(b.Files, files...)
	return b
}

// AddFile adds a discord.File to the discord.MessageCreate
func (b *MessageCreateBuilder) AddFile(name string, reader io.Reader, flags ...discord.FileFlags) *MessageCreateBuilder {
	b.Files = append(b.Files, discord.NewFile(name, reader, flags...))
	return b
}

// ClearFiles removes all discord.File(s) of this discord.MessageCreate
func (b *MessageCreateBuilder) ClearFiles() *MessageCreateBuilder {
	b.Files = []*discord.File{}
	return b
}

// RemoveFile removes the discord.File at this index
func (b *MessageCreateBuilder) RemoveFile(i int) *MessageCreateBuilder {
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

// Build builds the MessageCreateBuilder to a MessageCreate struct
func (b *MessageCreateBuilder) Build() discord.WebhookMessageCreate {
	b.WebhookMessageCreate.Components = b.Components
	return b.WebhookMessageCreate
}
