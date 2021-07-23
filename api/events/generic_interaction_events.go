package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	*GenericEvent
	Interaction *api.Interaction
}

// Respond replies to the api.Interaction with the provided api.InteractionResponse
func (e *GenericInteractionEvent) Respond(responseType api.InteractionResponseType, data interface{}) error {
	return e.Interaction.Respond(responseType, data)
}

// DeferReply replies to the api.SlashCommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (e *GenericInteractionEvent) DeferReply(ephemeral bool) error {
	return e.Interaction.DeferReply(ephemeral)
}

// Reply replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.MessageCreate
func (e *GenericInteractionEvent) Reply(messageCreate api.MessageCreate) error {
	return e.Interaction.Reply(messageCreate)
}

// EditOriginal edits the original api.InteractionResponse
func (e *GenericInteractionEvent) EditOriginal(messageUpdate api.MessageUpdate) (*api.Message, error) {
	return e.Interaction.EditOriginal(messageUpdate)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e *GenericInteractionEvent) DeleteOriginal() error {
	return e.Interaction.DeleteOriginal()
}

// SendFollowup used to send a followup api.MessageCreate to an api.Interaction
func (e *GenericInteractionEvent) SendFollowup(messageCreate api.MessageCreate) (*api.Message, error) {
	return e.Interaction.SendFollowup(messageCreate)
}

// EditFollowup used to edit a followup api.Message from an api.Interaction
func (e *GenericInteractionEvent) EditFollowup(messageID api.Snowflake, messageUpdate api.MessageUpdate) (*api.Message, error) {
	return e.Interaction.EditFollowup(messageID, messageUpdate)
}

// DeleteFollowup used to delete a followup api.Message from an api.Interaction
func (e *GenericInteractionEvent) DeleteFollowup(messageID api.Snowflake) error {
	return e.Interaction.DeleteFollowup(messageID)
}
