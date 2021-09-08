package core

type MessageCommandInteraction struct {
	*ContextCommandInteraction
}

func (i *MessageCommandInteraction) TargetMessage() *Message {
	return i.Resolved.Messages[i.TargetID]
}
