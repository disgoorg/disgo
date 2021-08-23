package events

import (
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
func (e GenericInteractionEvent) Respond(responseType discord.InteractionResponseType, data interface{}) error {
	return e.Interaction.Respond(responseType, data)
}

// DeferCreate replies to the core.SlashCommandInteraction with discord.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e GenericInteractionEvent) DeferCreate(ephemeral bool) error {
	return e.Interaction.DeferCreate(ephemeral)
}

// Create replies to the api.Interaction with discord.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e GenericInteractionEvent) Create(messageCreate discord.MessageCreate) error {
	return e.Interaction.Create(messageCreate)
}

// GetOriginal gets the original discord.InteractionResponse
func (e GenericInteractionEvent) GetOriginal(opts ...rest.RequestOpt) (*core.Message, rest.Error) {
	return e.Interaction.GetOriginal(opts...)
}

// UpdateOriginal edits the original discord.InteractionResponse
func (e GenericInteractionEvent) UpdateOriginal(messageUpdate discord.MessageUpdate) (*core.Message, rest.Error) {
	return e.Interaction.UpdateOriginal(messageUpdate)
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (e GenericInteractionEvent) DeleteOriginal(opts ...rest.RequestOpt) rest.Error {
	return e.Interaction.DeleteOriginal(opts...)
}

// CreateFollowup is used to send a followup discord.MessageCreate to an api.Interaction
func (e GenericInteractionEvent) CreateFollowup(messageCreate discord.MessageCreate) (*core.Message, rest.Error) {
	return e.Interaction.CreateFollowup(messageCreate)
}

// UpdateFollowup is used to edit a followup discord.Message from an api.Interaction
func (e GenericInteractionEvent) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*core.Message, rest.Error) {
	return e.Interaction.UpdateFollowup(messageID, messageUpdate)
}

// DeleteFollowup used to delete a followup discord.Message from a core.Interaction
func (e GenericInteractionEvent) DeleteFollowup(messageID discord.Snowflake) rest.Error {
	return e.Interaction.DeleteFollowup(messageID)
}
