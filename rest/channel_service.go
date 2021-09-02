package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ ChannelService = (*ChannelServiceImpl)(nil)

func NewChannelService(restClient Client) ChannelService {
	return &ChannelServiceImpl{restClient: restClient}
}

type ChannelService interface {
	Service
	GetChannel(channelID discord.Snowflake, opts ...RequestOpt) (*discord.Channel, Error)
	UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (*discord.Channel, Error)
	DeleteChannel(channelID discord.Snowflake, opts ...RequestOpt) Error

	GetWebhooks(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Webhook, Error)
	CreateWebhook(channelID discord.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (*discord.Webhook, Error)

	UpdatePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) Error
	DeletePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) Error

	SendTyping(channelID discord.Snowflake, opts ...RequestOpt) Error

	GetMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)
	CreateMessage(channelID discord.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (*discord.Message, Error)
	UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, Error)
	DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error
	BulkDeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...RequestOpt) Error
	CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)

	AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error
	RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error
	RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...RequestOpt) Error
	RemoveAllReactions(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error
	RemoveAllReactionsForEmoji(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error
}

type ChannelServiceImpl struct {
	restClient Client
}

func (s *ChannelServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *ChannelServiceImpl) GetChannel(channelID discord.Snowflake, opts ...RequestOpt) (channel *discord.Channel, rErr Error) {
	compiledRoute, err := route.GetChannel.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &channel, opts...)
	return
}

func (s *ChannelServiceImpl) UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (channel *discord.Channel, rErr Error) {
	compiledRoute, err := route.UpdateChannel.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, channelUpdate, &channel, opts...)
	return
}

func (s *ChannelServiceImpl) DeleteChannel(channelID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteChannel.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) GetWebhooks(channelID discord.Snowflake, opts ...RequestOpt) (webhooks []discord.Webhook, rErr Error) {
	compiledRoute, err := route.GetChannelWebhooks.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *ChannelServiceImpl) CreateWebhook(channelID discord.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.CreateWebhook.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, webhookCreate, &webhook, opts...)
	return
}

func (s *ChannelServiceImpl) UpdatePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) Error {
	compiledRoute, err := route.UpdatePermissionOverride.Compile(nil, channelID, overwriteID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, permissionOverwrite, nil, opts...)
}

func (s *ChannelServiceImpl) DeletePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeletePermissionOverride.Compile(nil, channelID, overwriteID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) SendTyping(channelID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.SendTyping.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) GetMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.GetMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *ChannelServiceImpl) CreateMessage(channelID discord.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, messageCreate, &message, opts...)
	return
}

func (s *ChannelServiceImpl) UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, messageUpdate, &message, opts...)
	return
}

func (s *ChannelServiceImpl) DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) BulkDeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.BulkDeleteMessages.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, discord.MessageBulkDelete{Messages: messageIDs}, nil, opts...)
}

func (s *ChannelServiceImpl) CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *ChannelServiceImpl) AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error {
	compiledRoute, err := route.AddReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveOwnReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveUserReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) RemoveAllReactions(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveAllReactions.Compile(nil, channelID, messageID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *ChannelServiceImpl) RemoveAllReactionsForEmoji(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveAllReactionsForEmoji.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
