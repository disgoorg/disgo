package core

type UserCommandInteraction struct {
	*ContextCommandInteraction
	Data *UserCommandInteractionData `json:"data,omitempty"`
}

func (i *UserCommandInteraction) TargetUser() *User {
	return i.Resolved().Users[i.TargetID()]
}

func (i *UserCommandInteraction) TargetMember() *Member {
	return i.Resolved().Members[i.TargetID()]
}

type UserCommandInteractionData struct {
	*ContextCommandInteractionData
}
