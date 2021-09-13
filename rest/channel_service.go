package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ ChannelService = (*channelServiceImpl)(nil)

func NewChannelService(restClient Client) ChannelService {
	return &channelServiceImpl{restClient: restClient}
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

type channelServiceImpl struct {
	restClient Client
}

func (s *channelServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *channelServiceImpl) GetChannel(channelID discord.Snowflake, opts ...RequestOpt) (channel *discord.Channel, rErr Error) {
	compiledRoute, err := route.GetChannel.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &channel, opts...)
	return
}

func (s *channelServiceImpl) UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (channel *discord.Channel, rErr Error) {
	compiledRoute, err := route.UpdateChannel.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, channelUpdate, &channel, opts...)
	return
}

func (s *channelServiceImpl) DeleteChannel(channelID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteChannel.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetWebhooks(channelID discord.Snowflake, opts ...RequestOpt) (webhooks []discord.Webhook, rErr Error) {
	compiledRoute, err := route.GetChannelWebhooks.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *channelServiceImpl) CreateWebhook(channelID discord.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.CreateWebhook.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, webhookCreate, &webhook, opts...)
	return
}

func (s *channelServiceImpl) UpdatePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) Error {
	compiledRoute, err := route.UpdatePermissionOverride.Compile(nil, channelID, overwriteID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, permissionOverwrite, nil, opts...)
}

func (s *channelServiceImpl) DeletePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeletePermissionOverride.Compile(nil, channelID, overwriteID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) SendTyping(channelID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.SendTyping.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.GetMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *channelServiceImpl) CreateMessage(channelID discord.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	body, err := messageCreate.ToBody()
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *channelServiceImpl) UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	body, err := messageUpdate.ToBody()
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *channelServiceImpl) DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) BulkDeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.BulkDeleteMessages.Compile(nil, channelID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, discord.MessageBulkDelete{Messages: messageIDs}, nil, opts...)
}

func (s *channelServiceImpl) CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *channelServiceImpl) AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error {
	compiledRoute, err := route.AddReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveOwnReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveUserReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveAllReactions(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveAllReactions.Compile(nil, channelID, messageID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveAllReactionsForEmoji(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error {
	compiledRoute, err := route.RemoveAllReactionsForEmoji.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
