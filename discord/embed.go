package discord

import "time"

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

// The EmbedResource of an Embed.Image/Embed.Thumbnail/Embed.Video
type EmbedResource struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

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
