package discord

// InteractionCallbackType indicates the type of slash command response, whether it's responding immediately or deferring to edit your response later
type InteractionCallbackType int

// Constants for the InteractionCallbackType(s)
const (
	InteractionCallbackTypePong InteractionCallbackType = iota + 1
	_
	_
	InteractionCallbackTypeChannelMessageWithSource
	InteractionCallbackTypeDeferredChannelMessageWithSource
	InteractionCallbackTypeDeferredUpdateMessage
	InteractionCallbackTypeUpdateMessage
	InteractionCallbackTypeAutocompleteResult
)

// InteractionResponse is how you answer interactions. If an answer is not sent within 3 seconds of receiving it, the interaction is failed, and you will be unable to respond to it.
type InteractionResponse struct {
	Type InteractionCallbackType `json:"type"`
	Data InteractionCallbackData `json:"data,omitempty"`
}

type dataType int

const (
	dataTypeMessageCreate = iota
	dataTypeMessageUpdate
	dataTypeAutocompleteResult
)

type InteractionCallbackData interface {
	dataType() dataType
	ToBody() (interface{}, error)
}

// ToBody returns the InteractionResponse ready for body
func (r *InteractionResponse) ToBody() (interface{}, error) {
	if r.Data == nil {
		return r, nil
	}
	return r.ToBody()
}
