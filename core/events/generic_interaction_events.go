package events

import (
	"context"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	*GenericEvent
	Interaction *core.Interaction
}

// Respond replies to the api.Interaction with the provided api.InteractionResponse
func (e GenericInteractionEvent) Respond(ctx context.Context, responseType discord.InteractionResponseType, data interface{}) error {
	return e.Interaction.Respond(ctx, responseType, data)
}

// DeferReply replies to the api.SlashCommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e GenericInteractionEvent) DeferReply(ctx context.Context, ephemeral bool) error {
	return e.Interaction.DeferReply(ctx, ephemeral)
}

// Reply replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e GenericInteractionEvent) Reply(ctx context.Context, messageCreate discord.MessageCreate) error {
	return e.Interaction.Reply(ctx, messageCreate)
}

// UpdateOriginal edits the original api.InteractionResponse
func (e GenericInteractionEvent) UpdateOriginal(ctx context.Context, messageUpdate discord.MessageUpdate) (*core.Message, error) {
	return e.Interaction.UpdateOriginal(ctx, messageUpdate)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e GenericInteractionEvent) DeleteOriginal(ctx context.Context, ) error {
	return e.Interaction.DeleteOriginal(ctx)
}

// CreateFollowup is used to send a followup api.MessageCreate to an api.Interaction
func (e GenericInteractionEvent) CreateFollowup(ctx context.Context, messageCreate discord.MessageCreate) (*core.Message, error) {
	return e.Interaction.CreateFollowup(ctx, messageCreate)
}

// EditFollowup is used to edit a followup api.Message from an api.Interaction
func (e GenericInteractionEvent) EditFollowup(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*core.Message, error) {
	return e.Interaction.UpdateFollowup(ctx, messageID, messageUpdate)
}

// DeleteFollowup used to delete a followup api.Message from an api.Interaction
func (e GenericInteractionEvent) DeleteFollowup(ctx context.Context, messageID discord.Snowflake) error {
	return e.Interaction.DeleteFollowup(ctx, messageID)
}
