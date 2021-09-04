package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Interaction struct {
	discord.UnmarshalInteraction
	Disgo           Disgo
	User            *User
	Member          *Member
	ResponseChannel chan discord.InteractionResponse
	Responded       bool
	Data            *InteractionData
}

// Respond responds to the Interaction with the provided InteractionResponse
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

	return i.Disgo.RestServices().InteractionService().CreateInteractionResponse(i.ID, i.Token, response, opts...)
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
	message, err := i.Disgo.RestServices().InteractionService().GetInteractionResponse(i.Disgo.ApplicationID(), i.Token, opts...)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateOriginal edits the original discord.InteractionResponse
func (i *Interaction) UpdateOriginal(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().UpdateInteractionResponse(i.Disgo.ApplicationID(), i.Token, messageUpdate, opts...)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteOriginal deletes the original discord.InteractionResponse
func (i *Interaction) DeleteOriginal(opts ...rest.RequestOpt) rest.Error {
	return i.Disgo.RestServices().InteractionService().DeleteInteractionResponse(i.Disgo.ApplicationID(), i.Token, opts...)
}

// CreateFollowup is used to send a discord.MessageCreate to an Interaction
func (i *Interaction) CreateFollowup(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().CreateFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageCreate, opts...)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateFollowup is used to edit a Message from an Interaction
func (i *Interaction) UpdateFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().UpdateFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID, messageUpdate, opts...)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteFollowup used to delete a Message from an Interaction
func (i *Interaction) DeleteFollowup(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return i.Disgo.RestServices().InteractionService().DeleteFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID, opts...)
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
	return i.Disgo.Caches().GuildCache().Get(*i.GuildID)
}

// DMChannel returns the DMChannel from the Caches
func (i *Interaction) DMChannel() DMChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Caches().DMChannelCache().Get(*i.ChannelID)
}

// MessageChannel returns the MessageChannel from the Caches
func (i *Interaction) MessageChannel() MessageChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Caches().ChannelCache().GetMessageChannel(*i.ChannelID)
}

// TextChannel returns the TextChannel from the Caches
func (i *Interaction) TextChannel() TextChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Caches().TextChannelCache().Get(*i.ChannelID)
}

// GuildChannel returns the GuildChannel from the Caches
func (i *Interaction) GuildChannel() GuildChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Caches().ChannelCache().GetGuildChannel(*i.ChannelID)
}

type InteractionData struct {
	discord.UnmarshalInteractionData
}
