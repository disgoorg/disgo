package api

type UserContextInteraction struct {
	*GenericContextInteraction
	Data *UserContextInteractionData `json:"data,omitempty"`
}

func (i *UserContextInteraction) User() *User {
	return i.Data.Resolved.Users[i.Data.TargetID]
}

func (i *UserContextInteraction) Member() *Member {
	return i.Data.Resolved.Members[i.Data.TargetID]
}

type UserContextInteractionData struct {
	*GenericContextInteractionData
}
