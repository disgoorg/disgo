package api

type GenericContextInteraction struct {
	*GenericCommandInteraction
	Data *GenericContextInteractionData `json:"data,omitempty"`
}

func (i *GenericContextInteraction) TargetID() Snowflake {
	return i.Data.TargetID
}

type GenericContextInteractionData struct {
	*GenericCommandInteractionData
	TargetID Snowflake `json:"target_id"`
}
