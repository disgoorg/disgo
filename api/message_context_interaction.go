package api

type MessageContextInteraction struct {
	*GenericContextInteraction
	Data *MessageContextInteractionData `json:"data,omitempty"`
}

func (i *MessageContextInteraction) Message() *Message {
	return i.Data.Resolved.Messages[i.Data.TargetID]
}

type MessageContextInteractionData struct {
	*GenericContextInteractionData
}
