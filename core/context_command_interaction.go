package core

import "github.com/DisgoOrg/disgo/discord"

type ContextCommandInteraction struct {
	*ApplicationCommandInteraction
	Data *ContextCommandInteractionData `json:"data,omitempty"`
}

func (i *ContextCommandInteraction) TargetID() discord.Snowflake {
	return i.Data.TargetID
}

type ContextCommandInteractionData struct {
	*ApplicationCommandInteractionData
}
