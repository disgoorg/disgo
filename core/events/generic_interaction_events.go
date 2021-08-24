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
func (e GenericInteractionEvent) Respond(responseType discord.InteractionResponseType, data interface{}, opts ...rest.RequestOpt) error {
	return e.Interaction.Respond(responseType, data, opts...)
}

// DeferCreate replies to the core.SlashCommandInteraction with discord.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e GenericInteractionEvent) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	return e.Interaction.DeferCreate(ephemeral, opts...)
}

// Create replies to the api.Interaction with discord.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e GenericInteractionEvent) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Interaction.Create(messageCreate, opts...)
}

// GetOriginal gets the original discord.InteractionResponse
func (e GenericInteractionEvent) GetOriginal(opts ...rest.RequestOpt) (*core.Message, rest.Error) {
	return e.Interaction.GetOriginal(opts...)
}

// UpdateOriginal edits the original discord.InteractionResponse
func (e GenericInteractionEvent) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*core.Message, rest.Error) {
	return e.Interaction.UpdateOriginal(messageUpdate, opts...)
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (e GenericInteractionEvent) DeleteOriginal(opts ...rest.RequestOpt) rest.Error {
	return e.Interaction.DeleteOriginal(opts...)
}

// CreateFollowup is used to send a followup discord.MessageCreate to an api.Interaction
func (e GenericInteractionEvent) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*core.Message, rest.Error) {
	return e.Interaction.CreateFollowup(messageCreate, opts...)
}

// UpdateFollowup is used to edit a followup discord.Message from an api.Interaction
func (e GenericInteractionEvent) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*core.Message, rest.Error) {
	return e.Interaction.UpdateFollowup(messageID, messageUpdate, opts...)
}

// DeleteFollowup used to delete a followup discord.Message from a core.Interaction
func (e GenericInteractionEvent) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return e.Interaction.DeleteFollowup(messageID, opts...)
}
