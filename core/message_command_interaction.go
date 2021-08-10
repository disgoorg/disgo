package core

type MessageCommandInteraction struct {
	*ContextCommandInteraction
	Data *MessageCommandInteractionData `json:"data,omitempty"`
}

func (i *MessageCommandInteraction) Message() *Message {
	return i.Resolved().Messages[i.TargetID()]
}

type MessageCommandInteractionData struct {
	*ContextCommandInteractionData
}
