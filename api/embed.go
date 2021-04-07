package api

import (
	"fmt"
	"time"
)

// EmbedType is the type of an Embed
type EmbedType string

// Constants for EmbedType
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
	Title       *string        `json:"title,omitempty"`
	Type        *EmbedType     `json:"type,omitempty"`
	Description *string        `json:"description,omitempty"`
	URL         *string        `json:"url,omitempty"`
	Timestamp   *time.Time     `json:"timestamp,omitempty"`
	Color       *int           `json:"color,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Image       *EmbedResource `json:"image,omitempty"`
	Thumbnail   *EmbedResource `json:"thumbnail,omitempty"`
	Video       *EmbedResource `json:"video,omitempty"`
	Provider    *EmbedProvider `json:"provider,omitempty"`
	Author      *EmbedAuthor   `json:"author,omitempty"`
	Fields      []*EmbedField  `json:"fields,omitempty"`
}

// The EmbedFooter of an Embed
type EmbedFooter struct {
	Text         string  `json:"text"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

// The EmbedResource of an Embed.Image/Embed.Thumbnail/Embed.Video
type EmbedResource struct {
	URL      *string `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

// The EmbedProvider of an Embed
type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// The EmbedAuthor of an Embed
type EmbedAuthor struct {
	Name         *string `json:"name,omitempty"`
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

// EmbedField (s) of an Embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

// NewEmbedBuilder returns a new embed builder
func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{}
}

// EmbedBuilder allows you to create embeds and use methods to set values
type EmbedBuilder struct {
	Embed
}

// SetTitle sets the title of the EmbedBuilder
func (b *EmbedBuilder) SetTitle(title *string) *EmbedBuilder {
	b.Title = title
	return b
}

// SetDescription sets the description of the EmbedBuilder
func (b *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	b.Description = &description
	return b
}

// SetDescriptionf sets the description of the EmbedBuilder with format
func (b *EmbedBuilder) SetDescriptionf(description string, a ...interface{}) *EmbedBuilder {
	descriptionf := fmt.Sprintf(description, a...)
	b.Description = &descriptionf
	return b
}

// SetAuthor sets the author of the EmbedBuilder
func (b *EmbedBuilder) SetAuthor(author *EmbedAuthor) *EmbedBuilder {
	b.Author = author
	return b
}

// SetColor sets the color of the EmbedBuilder
func (b *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	b.Color = &color
	return b
}

// SetFooter sets the footer of the EmbedBuilder
func (b *EmbedBuilder) SetFooter(footer *EmbedFooter) *EmbedBuilder {
	b.Footer = footer
	return b
}

// SetFooterBy sets the footer of the EmbedBuilder by text and iconURL
func (b *EmbedBuilder) SetFooterBy(text string, iconURL *string) *EmbedBuilder {
	b.Footer = &EmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return b
}

// SetImage sets the image of the EmbedBuilder
func (b *EmbedBuilder) SetImage(i *string) *EmbedBuilder {
	b.Image = &EmbedResource{
		URL: i,
	}
	return b
}

// SetThumbnail sets the thumbnail of the EmbedBuilder
func (b *EmbedBuilder) SetThumbnail(i *string) *EmbedBuilder {
	b.Thumbnail = &EmbedResource{
		URL: i,
	}
	return b
}

// SetURL sets the URL of the EmbedBuilder
func (b *EmbedBuilder) SetURL(u *string) *EmbedBuilder {
	b.URL = u
	return b
}

// AddField adds a field to the EmbedBuilder by name and value
func (b *EmbedBuilder) AddField(name string, value string, inline bool) *EmbedBuilder {
	b.Fields = append(b.Fields, &EmbedField{name, value, &inline})
	return b
}

// SetField sets a field to the EmbedBuilder by name and value
func (b *EmbedBuilder) SetField(index int, name string, value string, inline bool) *EmbedBuilder {
	if len(b.Fields) > index {
		b.Fields[index] = &EmbedField{name, value, &inline}
	}
	return b
}

// AddFields adds multiple fields to the EmbedBuilder
func (b *EmbedBuilder) AddFields(f *EmbedField, fs ...*EmbedField) *EmbedBuilder {
	b.Fields = append(b.Fields, f)
	b.Fields = append(b.Fields, fs...)
	return b
}

// SetFields sets fields of the EmbedBuilder
func (b *EmbedBuilder) SetFields(fs ...*EmbedField) *EmbedBuilder {
	b.Fields = fs
	return b
}

// ClearFields removes all of the fields from the EmbedBuilder
func (b *EmbedBuilder) ClearFields() *EmbedBuilder {
	b.Fields = []*EmbedField{}
	return b
}

// RemoveField removes a field from the EmbedBuilder
func (b *EmbedBuilder) RemoveField(index int) *EmbedBuilder {
	if len(b.Fields) > index {
		b.Fields = append(b.Fields[:index], b.Fields[index+1:]...)
	}
	return b
}

// Build returns your built Embed
func (b *EmbedBuilder) Build() Embed {
	return b.Embed
}
