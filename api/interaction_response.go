package api

type InteractionResponseType int

const (
	InteractionResponseTypePong           = iota + 1
	InteractionResponseTypeAcknowledge    // DEPRECATED
	InteractionResponseTypeChannelMessage // DEPRECATED
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
)

type InteractionResponse struct {
	Type InteractionResponseType  `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

type InteractionResponseData struct {
	TTS             bool        `json:"tts,omitempty"`
	Content         string      `json:"content,omitempty"`
	Embeds          []Embed     `json:"embeds,omitempty"`
	AllowedMentions interface{} `json:"allowed_mentions,omitempty"`
	Flags           int         `json:"flags,omitempty"`
}

type FollowupMessage struct {
}

type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User User            `json:"user"`
}
