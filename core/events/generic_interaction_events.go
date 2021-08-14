package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	*GenericEvent
	Interaction *core.Interaction
}

// Respond replies to the api.Interaction with the provided api.InteractionResponse
func (e GenericInteractionEvent) Respond(responseType discord.InteractionResponseType, data interface{}) error {
	return e.Interaction.Respond(responseType, data)
}

// DeferReply replies to the api.SlashCommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e GenericInteractionEvent) DeferReply(ephemeral bool) error {
	return e.Interaction.DeferReply(ephemeral)
}

// Reply replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e GenericInteractionEvent) Reply(messageCreate discord.MessageCreate) error {
	return e.Interaction.Reply(messageCreate)
}

// UpdateOriginal edits the original api.InteractionResponse
func (e GenericInteractionEvent) UpdateOriginal(messageUpdate discord.MessageUpdate) (*core.Message, error) {
	return e.Interaction.UpdateOriginal(messageUpdate)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e GenericInteractionEvent) DeleteOriginal() error {
	return e.Interaction.DeleteOriginal()
}

// CreateFollowup is used to send a followup api.MessageCreate to an api.Interaction
func (e GenericInteractionEvent) CreateFollowup(messageCreate discord.MessageCreate) (*core.Message, error) {
	return e.Interaction.CreateFollowup(messageCreate)
}

// EditFollowup is used to edit a followup api.Message from an api.Interaction
func (e GenericInteractionEvent) EditFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*core.Message, error) {
	return e.Interaction.UpdateFollowup(messageID, messageUpdate)
}

// DeleteFollowup used to delete a followup api.Message from an api.Interaction
func (e GenericInteractionEvent) DeleteFollowup(messageID discord.Snowflake) error {
	return e.Interaction.DeleteFollowup(messageID)
}
