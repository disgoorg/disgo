package api

import "fmt"

// NewEmbedBuilder returns a new embed builder
func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{}
}

// EmbedBuilder allows you to create embeds and use methods to set values
type EmbedBuilder struct {
	Embed
}

// SetTitle sets the title of the EmbedBuilder
func (b *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	b.Title = &title
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

// SetEmbedAuthor sets the author of the EmbedBuilder using an EmbedAuthor struct
func (b *EmbedBuilder) SetEmbedAuthor(author *EmbedAuthor) *EmbedBuilder {
	b.Author = author
	return b
}

// SetAuthor sets the author of the EmbedBuilder without an Icon URL
func (b *EmbedBuilder) SetAuthor(name string, url string) *EmbedBuilder {
	b.Author = &EmbedAuthor{
		Name: &name,
		URL:  &url,
	}
	return b
}

// SetAuthorI sets the author of the EmbedBuilder with all properties
func (b *EmbedBuilder) SetAuthorI(name string, url string, iconURL string) *EmbedBuilder {
	b.Author = &EmbedAuthor{
		Name:    &name,
		URL:     &url,
		IconURL: &iconURL,
	}
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
func (b *EmbedBuilder) SetFooterBy(text string, iconURL string) *EmbedBuilder {
	b.Footer = &EmbedFooter{
		Text:    text,
		IconURL: &iconURL,
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
func (b *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	b.URL = &url
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
