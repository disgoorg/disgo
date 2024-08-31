package handler

import (
	"context"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

type InteractionEvent struct {
	*events.InteractionCreate
	Vars map[string]string
	Ctx  context.Context
}

// CreateMessage responds to the interaction with a new message.
func (e *InteractionEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *InteractionEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *InteractionEvent) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *InteractionEvent) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// Deprecated: Respond with a discord.ButtonStylePremium button instead.
// PremiumRequired responds to the interaction with an upgrade button if available.
func (e *InteractionEvent) PremiumRequired(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypePremiumRequired, nil, opts...)
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *InteractionEvent) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}

// Modal responds to the interaction with a new modal.
func (e *InteractionEvent) Modal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

func (e *InteractionEvent) AutocompleteResult(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

func (e *InteractionEvent) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *InteractionEvent) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *InteractionEvent) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *InteractionEvent) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *InteractionEvent) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *InteractionEvent) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *InteractionEvent) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}
