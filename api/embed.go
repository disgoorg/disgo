package api

import "time"

type EmbedType string

const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGifV    EmbedType = "rich"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

// Embed allows you to send embeds to discord
type Embed struct {
	Title       string         `json:"title,omitempty"`
	Type        EmbedType      `json:"type,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Timestamp   *time.Time     `json:"timestamp,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
}

type EmbedFooter struct {
}

type EmbedImage struct {
}

type EmbedThumbnail struct {
}

type EmbedVideo struct {
}

type EmbedProvider struct {
}

type EmbedAuthor struct {
}

type EmbedField struct {
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{}
}

type EmbedBuilder struct {
	Embed
}

func (b *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	b.Description = description
	return b
}

func (b *EmbedBuilder) Build() Embed {
	return b.Embed
}