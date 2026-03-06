package discord

import (
	"fmt"
	"time"
)

// EmbedType is the type of Embed
type EmbedType string

// Constants for EmbedType
const (
	EmbedTypeRich                  EmbedType = "rich"
	EmbedTypeImage                 EmbedType = "image"
	EmbedTypeVideo                 EmbedType = "video"
	EmbedTypeGifV                  EmbedType = "gifv"
	EmbedTypeArticle               EmbedType = "article"
	EmbedTypeLink                  EmbedType = "link"
	EmbedTypeAutoModerationMessage EmbedType = "auto_moderation_message"
	EmbedTypePollResult            EmbedType = "poll_result"
)

// NewEmbed returns a new Embed struct with no fields set.
func NewEmbed() Embed {
	return Embed{}
}

// Embed allows you to send embeds to discord
type Embed struct {
	Title       string         `json:"title,omitempty"`
	Type        EmbedType      `json:"type,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Timestamp   *time.Time     `json:"timestamp,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Image       *EmbedResource `json:"image,omitempty"`
	Thumbnail   *EmbedResource `json:"thumbnail,omitempty"`
	Video       *EmbedResource `json:"video,omitempty"`
	Provider    *EmbedProvider `json:"provider,omitempty"`
	Author      *EmbedAuthor   `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
	Flags       EmbedFlags     `json:"flags,omitempty"`
}

// WithTitle Withs the title of the Embed
func (e Embed) WithTitle(title string) Embed {
	e.Title = title
	return e
}

// WithTitlef Withs the title of the Embed with format
func (e Embed) WithTitlef(title string, a ...any) Embed {
	return e.WithTitle(fmt.Sprintf(title, a...))
}

// WithDescription Withs the description of the Embed
func (e Embed) WithDescription(description string) Embed {
	e.Description = description
	return e
}

// WithDescriptionf Withs the description of the Embed with format
func (e Embed) WithDescriptionf(description string, a ...any) Embed {
	return e.WithDescription(fmt.Sprintf(description, a...))
}

// WithEmbedAuthor Withs the author of the Embed using an EmbedAuthor struct
func (e Embed) WithEmbedAuthor(author *EmbedAuthor) Embed {
	e.Author = author
	return e
}

// WithAuthor Withs the author of the Embed with all properties
func (e Embed) WithAuthor(name string, url string, iconURL string) Embed {
	if e.Author == nil {
		e.Author = &EmbedAuthor{}
	}
	e.Author.Name = name
	e.Author.URL = url
	e.Author.IconURL = iconURL
	return e
}

// WithAuthorName Withs the author name of the Embed
func (e Embed) WithAuthorName(name string) Embed {
	if e.Author == nil {
		e.Author = &EmbedAuthor{}
	}
	e.Author.Name = name
	return e
}

// WithAuthorNamef Withs the author name of the Embed with format
func (e Embed) WithAuthorNamef(name string, a ...any) Embed {
	return e.WithAuthorName(fmt.Sprintf(name, a...))
}

// WithAuthorURL Withs the author URL of the Embed
func (e Embed) WithAuthorURL(url string) Embed {
	if e.Author == nil {
		e.Author = &EmbedAuthor{}
	}
	e.Author.URL = url
	return e
}

// WithAuthorURLf Withs the author URL of the Embed with format
func (e Embed) WithAuthorURLf(url string, a ...any) Embed {
	return e.WithAuthorURL(fmt.Sprintf(url, a...))
}

// WithAuthorIcon Withs the author icon of the Embed
func (e Embed) WithAuthorIcon(iconURL string) Embed {
	if e.Author == nil {
		e.Author = &EmbedAuthor{}
	}
	e.Author.IconURL = iconURL
	return e
}

// WithAuthorIconf Withs the author icon of the Embed with format
func (e Embed) WithAuthorIconf(iconURL string, a ...any) Embed {
	return e.WithAuthorIcon(fmt.Sprintf(iconURL, a...))
}

// WithColor Withs the color of the Embed
// The color should be an integer representation of a hexadecimal color code (e.g. 0xFF0000 for red)
func (e Embed) WithColor(color int) Embed {
	e.Color = color
	return e
}

// WithEmbedFooter Withs the footer of the Embed
func (e Embed) WithEmbedFooter(footer *EmbedFooter) Embed {
	e.Footer = footer
	return e
}

// WithFooter Withs the footer icon of the Embed
func (e Embed) WithFooter(text string, iconURL string) Embed {
	if e.Footer == nil {
		e.Footer = &EmbedFooter{}
	}
	e.Footer.Text = text
	e.Footer.IconURL = iconURL
	return e
}

// WithFooterText Withs the footer text of the Embed
func (e Embed) WithFooterText(text string) Embed {
	if e.Footer == nil {
		e.Footer = &EmbedFooter{}
	}
	e.Footer.Text = text
	return e
}

// WithFooterTextf Withs the footer text of the Embed with format
func (e Embed) WithFooterTextf(text string, a ...any) Embed {
	return e.WithFooterText(fmt.Sprintf(text, a...))
}

// WithFooterIcon Withs the footer icon of the Embed
func (e Embed) WithFooterIcon(iconURL string) Embed {
	if e.Footer == nil {
		e.Footer = &EmbedFooter{}
	}
	e.Footer.IconURL = iconURL
	return e
}

// WithFooterIconf Withs the footer icon of the Embed
func (e Embed) WithFooterIconf(iconURL string, a ...any) Embed {
	return e.WithFooterIcon(fmt.Sprintf(iconURL, a...))
}

// WithImage Withs the image of the Embed
func (e Embed) WithImage(url string) Embed {
	if e.Image == nil {
		e.Image = &EmbedResource{}
	}
	e.Image.URL = url
	return e
}

// WithImagef Withs the image of the Embed with format
func (e Embed) WithImagef(url string, a ...any) Embed {
	return e.WithImage(fmt.Sprintf(url, a...))
}

// WithThumbnail Withs the thumbnail of the Embed
func (e Embed) WithThumbnail(url string) Embed {
	if e.Thumbnail == nil {
		e.Thumbnail = &EmbedResource{}
	}
	e.Thumbnail.URL = url
	return e
}

// WithThumbnailf Withs the thumbnail of the Embed with format
func (e Embed) WithThumbnailf(url string, a ...any) Embed {
	return e.WithThumbnail(fmt.Sprintf(url, a...))
}

// WithURL Withs the URL of the Embed
func (e Embed) WithURL(url string) Embed {
	e.URL = url
	return e
}

// WithURLf Withs the URL of the Embed with format
func (e Embed) WithURLf(url string, a ...any) Embed {
	return e.WithURL(fmt.Sprintf(url, a...))
}

// WithTimestamp Withs the timestamp of the Embed
func (e Embed) WithTimestamp(time time.Time) Embed {
	e.Timestamp = &time
	return e
}

// AddField adds a field to the Embed by name and value
func (e Embed) AddField(name string, value string, inline bool) Embed {
	e.Fields = append(e.Fields, EmbedField{Name: name, Value: value, Inline: &inline})
	return e
}

// WithField Withs a field to the Embed by name and value
func (e Embed) WithField(i int, name string, value string, inline bool) Embed {
	if len(e.Fields) > i {
		e.Fields[i] = EmbedField{Name: name, Value: value, Inline: &inline}
	}
	return e
}

// AddFields adds multiple fields to the Embed
func (e Embed) AddFields(fields ...EmbedField) Embed {
	e.Fields = append(e.Fields, fields...)
	return e
}

// WithFields Withs fields of the Embed
func (e Embed) WithFields(fields ...EmbedField) Embed {
	e.Fields = fields
	return e
}

// ClearFields removes all the fields from the Embed
func (e Embed) ClearFields() Embed {
	e.Fields = []EmbedField{}
	return e
}

// RemoveField removes a field from the Embed
func (e Embed) RemoveField(i int) Embed {
	if len(e.Fields) > i {
		e.Fields = append(e.Fields[:i], e.Fields[i+1:]...)
	}
	return e
}

func (e Embed) FindField(fieldFindFunc func(field EmbedField) bool) (EmbedField, bool) {
	for _, field := range e.Fields {
		if fieldFindFunc(field) {
			return field, true
		}
	}
	return EmbedField{}, false
}

func (e Embed) FindAllFields(fieldFindFunc func(field EmbedField) bool) []EmbedField {
	var fields []EmbedField
	for _, field := range e.Fields {
		if fieldFindFunc(field) {
			fields = append(fields, field)
		}
	}
	return fields
}

type EmbedFlags int

const (
	EmbedFlagIsContentInventoryEntry EmbedFlags = 1 << (iota + 5)
	EmbedFlagsNone                   EmbedFlags = 0
)

// The EmbedResource of an Embed.Image/Embed.Thumbnail/Embed.Video
type EmbedResource struct {
	URL                string             `json:"url,omitempty"`
	ProxyURL           string             `json:"proxy_url,omitempty"`
	Height             int                `json:"height,omitempty"`
	Width              int                `json:"width,omitempty"`
	Placeholder        string             `json:"placeholder,omitempty"`
	PlaceholderVersion string             `json:"placeholder_version,omitempty"`
	Description        string             `json:"description,omitempty"`
	Flags              EmbedResourceFlags `json:"flags,omitempty"`
}

type EmbedResourceFlags int

const (
	EmbedResourceFlagIsAnimated EmbedResourceFlags = 1 << 5
	EmbedResourceFlagsNone      EmbedResourceFlags = 0
)

// The EmbedProvider of an Embed
type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// The EmbedAuthor of an Embed
type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// The EmbedFooter of an Embed
type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedField (s) of an Embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

type EmbedFieldPollResult string

const (
	EmbedFieldPollResultQuestionText              EmbedFieldPollResult = "poll_question_text"
	EmbedFieldPollResultVictorAnswerVotes         EmbedFieldPollResult = "victor_answer_votes"
	EmbedFieldPollResultTotalVotes                EmbedFieldPollResult = "total_votes"
	EmbedFieldPollResultVictorAnswerID            EmbedFieldPollResult = "victor_answer_id"
	EmbedFieldPollResultVictorAnswerText          EmbedFieldPollResult = "victor_answer_text"
	EmbedFieldPollResultVictorAnswerEmojiID       EmbedFieldPollResult = "victor_answer_emoji_id"
	EmbedFieldPollResultVictorAnswerEmojiName     EmbedFieldPollResult = "victor_answer_emoji_name"
	EmbedFieldPollResultVictorAnswerEmojiAnimated EmbedFieldPollResult = "victor_answer_emoji_animated"
)
