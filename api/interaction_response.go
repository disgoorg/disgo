package api

// InteractionResponseType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionResponseType int

// Constants for the InteractionResponseType(s)
const (
	InteractionResponseTypePong InteractionResponseType = iota + 1
	_
	_
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeDeferredChannelMessageWithSource
	_
	InteractionResponseTypeButtonResponse
)

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionResponseType `json:"type"`
	Data interface{}             `json:"data,omitempty"`
}
