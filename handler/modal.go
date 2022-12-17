package handler

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

type ModalEvent struct {
	*events.ModalSubmitInteractionCreate
}

// CreateMessage responds to the interaction with a new message.
func (e *ModalEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ModalEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ModalEvent) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ModalEvent) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

func (e *ModalEvent) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ModalEvent) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *ModalEvent) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ModalEvent) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *ModalEvent) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *ModalEvent) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *ModalEvent) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}
