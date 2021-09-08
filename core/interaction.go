package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Interaction struct {
	discord.Interaction
	Bot             *Bot
	User            *User
	Member          *Member
	ResponseChannel chan discord.InteractionResponse
	Responded       bool
}

// Respond responds to the Interaction with the provided discord.InteractionResponse
func (i *Interaction) Respond(responseType discord.InteractionResponseType, data interface{}, opts ...rest.RequestOpt) rest.Error {
	response := discord.InteractionResponse{
		Type: responseType,
		Data: data,
	}
	if i.Responded {
		return rest.NewError(nil, discord.ErrInteractionAlreadyReplied)
	}
	i.Responded = true

	if !i.FromGateway() {
		i.ResponseChannel <- response
		return nil
	}

	return i.Bot.RestServices.InteractionService().CreateInteractionResponse(i.ID, i.Token, response, opts...)
}

// DeferCreate replies to the Interaction with discord.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (i *Interaction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) rest.Error {
	var messageCreate interface{}
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return i.Respond(discord.InteractionResponseTypeDeferredChannelMessageWithSource, messageCreate, opts...)
}

// Create replies to the Interaction with discord.InteractionResponseTypeChannelMessageWithSource & discord.MessageCreate
func (i *Interaction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) rest.Error {
	return i.Respond(discord.InteractionResponseTypeChannelMessageWithSource, messageCreate, opts...)
}

// GetOriginal gets the original discord.InteractionResponse
func (i *Interaction) GetOriginal(opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().GetInteractionResponse(i.Bot.ApplicationID, i.Token, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateOriginal edits the original discord.InteractionResponse
func (i *Interaction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateInteractionResponse(i.Bot.ApplicationID, i.Token, messageUpdate, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (i *Interaction) DeleteOriginal(opts ...rest.RequestOpt) rest.Error {
	return i.Bot.RestServices.InteractionService().DeleteInteractionResponse(i.Bot.ApplicationID, i.Token, opts...)
}

// CreateFollowup is used to send a discord.MessageCreate to an Interaction
func (i *Interaction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().CreateFollowupMessage(i.Bot.ApplicationID, i.Token, messageCreate, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateFollowup is used to edit a Message from an Interaction
func (i *Interaction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateFollowupMessage(i.Bot.ApplicationID, i.Token, messageID, messageUpdate, opts...)
	if err != nil {

	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteFollowup used to delete a Message from an Interaction
func (i *Interaction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return i.Bot.RestServices.InteractionService().DeleteFollowupMessage(i.Bot.ApplicationID, i.Token, messageID, opts...)
}

// FromGateway returns is the Interaction came in via gateway.Gateway or httpserver.Server
func (i *Interaction) FromGateway() bool {
	return i.ResponseChannel == nil
}

// Guild returns the Guild from the Caches
func (i *Interaction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.GuildCache().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *Interaction) Channel() *Channel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Bot.Caches.ChannelCache().Get(*i.ChannelID)
}
