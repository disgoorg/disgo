package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Interaction struct {
	discord.UnmarshalInteraction
	Disgo           Disgo
	User            *User
	Member          *Member
	ResponseChannel chan discord.InteractionResponse
	Replied         bool
	Data            *InteractionData
}

// Respond responds to the Interaction with the provided InteractionResponse
func (i *Interaction) Respond(ctx context.Context, responseType discord.InteractionResponseType, data interface{}) rest.Error {
	response := discord.InteractionResponse{
		Type: responseType,
		Data: data,
	}
	if i.Replied {
		return rest.NewError(nil, discord.ErrInteractionAlreadyReplied)
	}
	i.Replied = true

	if !i.FromGateway() {
		i.ResponseChannel <- response
		return nil
	}

	return i.Disgo.RestServices().InteractionService().CreateInteractionResponse(ctx, i.ID, i.Token, response)
}

// DeferCreate replies to the Interaction with discord.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (i *Interaction) DeferCreate(ctx context.Context, ephemeral bool) rest.Error {
	var messageCreate interface{}
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return i.Respond(ctx, discord.InteractionResponseTypeDeferredChannelMessageWithSource, messageCreate)
}

// Create replies to the Interaction with discord.InteractionResponseTypeChannelMessageWithSource & discord.MessageCreate
func (i *Interaction) Create(ctx context.Context, messageCreate discord.MessageCreate) rest.Error {
	return i.Respond(ctx, discord.InteractionResponseTypeChannelMessageWithSource, messageCreate)
}

// UpdateOriginal edits the original InteractionResponse
func (i *Interaction) UpdateOriginal(ctx context.Context, messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().UpdateInteractionResponse(ctx, i.Disgo.ApplicationID(), i.Token, messageUpdate)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteOriginal deletes the original InteractionResponse
func (i *Interaction) DeleteOriginal(ctx context.Context, ) rest.Error {
	return i.Disgo.RestServices().InteractionService().DeleteInteractionResponse(ctx, i.Disgo.ApplicationID(), i.Token)
}

// CreateFollowup is used to send an MessageCreate to an Interaction
func (i *Interaction) CreateFollowup(ctx context.Context, messageCreate discord.MessageCreate) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().CreateFollowupMessage(ctx, i.Disgo.ApplicationID(), i.Token, messageCreate)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// UpdateFollowup is used to edit a Message from an Interaction
func (i *Interaction) UpdateFollowup(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().UpdateFollowupMessage(ctx, i.Disgo.ApplicationID(), i.Token, messageID, messageUpdate)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteFollowup used to delete a Message from an Interaction
func (i *Interaction) DeleteFollowup(ctx context.Context, messageID discord.Snowflake) rest.Error {
	return i.Disgo.RestServices().InteractionService().DeleteFollowupMessage(ctx, i.Disgo.ApplicationID(), i.Token, messageID)
}

// FromGateway returns is the Interaction came in via gateway.Gateway or httpserver.Server
func (i *Interaction) FromGateway() bool {
	return i.ResponseChannel == nil
}

// Guild returns the Guild from the Cache
func (i *Interaction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Disgo.Cache().GuildCache().Get(*i.GuildID)
}

// DMChannel returns the DMChannel from the Cache
func (i *Interaction) DMChannel() DMChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().DMChannelCache().Get(*i.ChannelID)
}

// MessageChannel returns the MessageChannel from the Cache
func (i *Interaction) MessageChannel() MessageChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().ChannelCache().MessageChannel(*i.ChannelID)
}

// TextChannel returns the TextChannel from the Cache
func (i *Interaction) TextChannel() TextChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().TextChannelCache().Get(*i.ChannelID)
}

// GuildChannel returns the GuildChannel from the Cache
func (i *Interaction) GuildChannel() GuildChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().ChannelCache().GuildChannel(*i.ChannelID)
}

type InteractionData struct {
	discord.UnmarshalInteractionData
}
