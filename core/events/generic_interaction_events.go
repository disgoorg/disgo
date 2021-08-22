package events

import (
	"context"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	*GenericEvent
	Interaction *core.Interaction
}

// Respond replies to the core.Interaction with the provided api.InteractionResponse
func (e GenericInteractionEvent) Respond(ctx context.Context, responseType discord.InteractionResponseType, data interface{}) error {
	return e.Interaction.Respond(ctx, responseType, data)
}

// DeferCreate replies to the core.SlashCommandInteraction with discord.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e GenericInteractionEvent) DeferCreate(ctx context.Context, ephemeral bool) error {
	return e.Interaction.DeferCreate(ctx, ephemeral)
}

// Create replies to the api.Interaction with discord.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e GenericInteractionEvent) Create(ctx context.Context, messageCreate discord.MessageCreate) error {
	return e.Interaction.Create(ctx, messageCreate)
}

// GetOriginal gets the original discord.InteractionResponse
func (e GenericInteractionEvent) GetOriginal(ctx context.Context) (*core.Message, rest.Error) {
	return e.Interaction.GetOriginal(ctx)
}

// UpdateOriginal edits the original discord.InteractionResponse
func (e GenericInteractionEvent) UpdateOriginal(ctx context.Context, messageUpdate discord.MessageUpdate) (*core.Message, rest.Error) {
	return e.Interaction.UpdateOriginal(ctx, messageUpdate)
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (e GenericInteractionEvent) DeleteOriginal(ctx context.Context) rest.Error {
	return e.Interaction.DeleteOriginal(ctx)
}

// CreateFollowup is used to send a followup discord.MessageCreate to an api.Interaction
func (e GenericInteractionEvent) CreateFollowup(ctx context.Context, messageCreate discord.MessageCreate) (*core.Message, rest.Error) {
	return e.Interaction.CreateFollowup(ctx, messageCreate)
}

// UpdateFollowup is used to edit a followup discord.Message from an api.Interaction
func (e GenericInteractionEvent) UpdateFollowup(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*core.Message, rest.Error) {
	return e.Interaction.UpdateFollowup(ctx, messageID, messageUpdate)
}

// DeleteFollowup used to delete a followup discord.Message from a core.Interaction
func (e GenericInteractionEvent) DeleteFollowup(ctx context.Context, messageID discord.Snowflake) rest.Error {
	return e.Interaction.DeleteFollowup(ctx, messageID)
}
