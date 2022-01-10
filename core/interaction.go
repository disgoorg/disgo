package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type InteractionFilter func(interaction Interaction) bool

// Interaction represents a generic Interaction received from discord
type Interaction interface {
	Type() discord.InteractionType
	interaction()
	Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error
}

type BaseInteraction struct {
	ID              discord.Snowflake
	ApplicationID   discord.Snowflake
	Token           string
	Version         int
	GuildID         *discord.Snowflake
	ChannelID       discord.Snowflake
	Member          *Member
	User            *User
	ResponseChannel chan<- discord.InteractionResponse
	Acknowledged    bool
	Bot             *Bot
}

func (i *BaseInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	if i.Acknowledged {
		return discord.ErrInteractionAlreadyReplied
	}
	i.Acknowledged = true

	response := discord.InteractionResponse{
		Type: callbackType,
		Data: callbackData,
	}

	if i.ResponseChannel != nil {
		i.ResponseChannel <- response
		return nil
	}

	return i.Bot.RestServices.InteractionService().CreateInteractionResponse(i.ID, i.Token, response, opts...)
}

// Guild returns the Guild from the Caches
func (i *BaseInteraction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.Guilds().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *BaseInteraction) Channel() MessageChannel {
	if ch := i.Bot.Caches.Channels().Get(i.ChannelID); ch != nil {
		return ch.(MessageChannel)
	}
	return nil
}

type ReplyInteraction struct {
	*BaseInteraction
}

func (i ReplyInteraction) GetResponse(opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().GetInteractionResponse(i.ApplicationID, i.Token, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (i ReplyInteraction) UpdateResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateInteractionResponse(i.ApplicationID, i.Token, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (i ReplyInteraction) DeleteResponse(opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.InteractionService().DeleteInteractionResponse(i.ApplicationID, i.Token, opts...)
}

func (i ReplyInteraction) Create(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeChannelMessageWithSource, messageCreate, opts...)
}

func (i ReplyInteraction) DeferCreate(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionCallbackData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return i.Respond(discord.InteractionCallbackTypeDeferredChannelMessageWithSource, data, opts...)
}

func (i ReplyInteraction) GetFollowupMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().GetFollowupMessage(i.ApplicationID, i.Token, messageID, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (i ReplyInteraction) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().CreateFollowupMessage(i.ApplicationID, i.Token, messageCreate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (i ReplyInteraction) UpdateFollowupMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := i.Bot.RestServices.InteractionService().UpdateFollowupMessage(i.ApplicationID, i.Token, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder.CreateMessage(*message, CacheStrategyNoWs), nil
}

func (i ReplyInteraction) DeleteFollowupMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return i.Bot.RestServices.InteractionService().DeleteFollowupMessage(i.ApplicationID, i.Token, messageID, opts...)
}
