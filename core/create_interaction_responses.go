package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type CreateInteractionResponses struct {
	*Interaction
}

// DeferCreate replies to the Interaction with discord.InteractionCallbackTypeDeferredChannelMessageWithSource and shows a loading state
func (i *CreateInteractionResponses) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) rest.Error {
	var messageCreate interface{}
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return i.Respond(discord.InteractionCallbackTypeDeferredChannelMessageWithSource, messageCreate, opts...)
}

// Create replies to the Interaction with discord.InteractionCallbackTypeChannelMessageWithSource & discord.MessageCreate
func (i *CreateInteractionResponses) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) rest.Error {
	return i.Respond(discord.InteractionCallbackTypeChannelMessageWithSource, messageCreate, opts...)
}

// GetOriginal gets the original discord.InteractionResponse
func (i *CreateInteractionResponses) GetOriginal(opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().GetInteractionResponse(i.Bot.ApplicationID, i.Token, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateOriginal edits the original discord.InteractionResponse
func (i *CreateInteractionResponses) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateInteractionResponse(i.Bot.ApplicationID, i.Token, messageUpdate, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (i *CreateInteractionResponses) DeleteOriginal(opts ...rest.RequestOpt) rest.Error {
	return i.Bot.RestServices.InteractionService().DeleteInteractionResponse(i.Bot.ApplicationID, i.Token, opts...)
}

// CreateFollowup is used to send a discord.MessageCreate to an Interaction
func (i *CreateInteractionResponses) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().CreateFollowupMessage(i.Bot.ApplicationID, i.Token, messageCreate, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateFollowup is used to edit a Message from an Interaction
func (i *CreateInteractionResponses) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateFollowupMessage(i.Bot.ApplicationID, i.Token, messageID, messageUpdate, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteFollowup used to delete a Message from an Interaction
func (i *CreateInteractionResponses) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return i.Bot.RestServices.InteractionService().DeleteFollowupMessage(i.Bot.ApplicationID, i.Token, messageID, opts...)
}

