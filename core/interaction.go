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
	Replied         bool
	Data            *InteractionData
}

// Respond responds to the Interaction with the provided InteractionResponse
func (i *Interaction) Respond(responseType discord.InteractionResponseType, data interface{}) rest.Error {
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

	return i.Disgo.RestServices().InteractionService().CreateInteractionResponse(i.ID, i.Token, response)
}

// DeferReply replies to the Interaction with InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (i *Interaction) DeferReply(ephemeral bool) rest.Error {
	var messageCreate interface{}
	if ephemeral {
		messageCreate = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return i.Respond(discord.InteractionResponseTypeDeferredChannelMessageWithSource, messageCreate)
}

// Reply replies to the Interaction with InteractionResponseTypeDeferredChannelMessageWithSource & MessageCreate
func (i *Interaction) Reply(messageCreate discord.MessageCreate) rest.Error {
	return i.Respond(discord.InteractionResponseTypeChannelMessageWithSource, messageCreate)
}

// EditOriginal edits the original InteractionResponse
func (i *Interaction) EditOriginal(messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().UpdateInteractionResponse(i.Disgo.ApplicationID(), i.Token, messageUpdate)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteOriginal deletes the original InteractionResponse
func (i *Interaction) DeleteOriginal() rest.Error {
	return i.Disgo.RestServices().InteractionService().DeleteInteractionResponse(i.Disgo.ApplicationID(), i.Token)
}

// SendFollowup used to send an MessageCreate to an Interaction
func (i *Interaction) SendFollowup(messageCreate discord.MessageCreate) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().CreateFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageCreate)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// EditFollowup used to edit a Message from an Interaction
func (i *Interaction) EditFollowup(messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := i.Disgo.RestServices().InteractionService().UpdateFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID, messageUpdate)
	if err != nil {

	}
	return i.Disgo.EntityBuilder().CreateMessage(*message, CacheStrategyNoWs), nil
}

// DeleteFollowup used to delete a Message from an Interaction
func (i *Interaction) DeleteFollowup(messageID discord.Snowflake) rest.Error {
	return i.Disgo.RestServices().InteractionService().DeleteFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID)
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
