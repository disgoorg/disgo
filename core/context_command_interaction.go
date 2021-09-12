package core

import "github.com/DisgoOrg/disgo/discord"

type ContextCommandInteraction struct {
	*ApplicationCommandInteraction
	ContextCommandInteractionData
}

type ContextCommandInteractionData struct {
	TargetID discord.Snowflake
}
