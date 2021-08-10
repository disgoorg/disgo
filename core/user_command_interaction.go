package core

type UserCommandInteraction struct {
	*ContextCommandInteraction
	Data *UserCommandInteractionData `json:"data,omitempty"`
}

func (i *UserCommandInteraction) User() *User {
	return i.Resolved().Users[i.TargetID()]
}

func (i *UserCommandInteraction) Member() *Member {
	return i.Resolved().Members[i.TargetID()]
}

type UserCommandInteractionData struct {
	*ContextCommandInteractionData
}
