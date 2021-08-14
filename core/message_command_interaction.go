package core

type MessageCommandInteraction struct {
	*ContextCommandInteraction
	Data *MessageCommandInteractionData `json:"data,omitempty"`
}

func (i *MessageCommandInteraction) TargetMessage() *Message {
	return i.Resolved().Messages[i.TargetID()]
}

type MessageCommandInteractionData struct {
	*ContextCommandInteractionData
}
