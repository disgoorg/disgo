package api

// InteractionResponseType indicates the type of slash command response, whether it's responding immediately or
// deferring to edit your response later
type InteractionResponseType int

// Constants for the InteractionResponseType(s)
const (
	InteractionResponseTypePong           = iota + 1
	InteractionResponseTypeAcknowledge    // DEPRECATED
	InteractionResponseTypeChannelMessage // DEPRECATED
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
)

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the
// interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionResponseType  `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

// The InteractionResponseData is used to specify the message_events options when creating an InteractionResponse
type InteractionResponseData struct {
	TTS             bool        `json:"tts,omitempty"`
	Content         string      `json:"content,omitempty"`
	Embeds          []Embed     `json:"embeds,omitempty"`
	AllowedMentions interface{} `json:"allowed_mentions,omitempty"`
	Flags           int         `json:"flags,omitempty"`
}

// FollowupMessage is used to add additional messages to an Interaction after you've responded initially
type FollowupMessage struct {
	// Todo: fill this
}

// MessageInteraction is sent on the Message object when the message_events is a response to an interaction
type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User User            `json:"user"`
}
