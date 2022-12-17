package handler

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

type CommandEvent struct {
	*events.ApplicationCommandInteractionCreate
}

// CreateMessage responds to the interaction with a new message.
func (e *CommandEvent) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *CommandEvent) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// CreateModal responds to the interaction with a new modal.
func (e *CommandEvent) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

func (e *CommandEvent) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *CommandEvent) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *CommandEvent) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *CommandEvent) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *CommandEvent) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *CommandEvent) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest().UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *CommandEvent) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest().DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}
