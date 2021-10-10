package core

type UserCommandInteractionFilter func(userCommandInteraction *UserCommandInteraction) bool

type UserCommandInteraction struct {
	*ContextCommandInteraction
}

func (i *UserCommandInteraction) TargetUser() *User {
	return i.Resolved.Users[i.TargetID]
}

func (i *UserCommandInteraction) TargetMember() *Member {
	return i.Resolved.Members[i.TargetID]
}
