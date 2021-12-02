package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service        = (*channelServiceImpl)(nil)
	_ ChannelService = (*channelServiceImpl)(nil)
)

func NewChannelService(restClient Client) ChannelService {
	return &channelServiceImpl{restClient: restClient}
}

type ChannelService interface {
	Service
	GetChannel(channelID discord.Snowflake, opts ...RequestOpt) (discord.Channel, error)
	UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (discord.Channel, error)
	DeleteChannel(channelID discord.Snowflake, opts ...RequestOpt) error

	GetWebhooks(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Webhook, error)
	CreateWebhook(channelID discord.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (discord.Webhook, error)

	GetPermissionOverwrites(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.PermissionOverwrite, error)
	GetPermissionOverwrite(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) (*discord.PermissionOverwrite, error)
	UpdatePermissionOverwrite(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) error
	DeletePermissionOverwrite(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) error

	SendTyping(channelID discord.Snowflake, opts ...RequestOpt) error

	GetMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	GetMessages(channelID discord.Snowflake, around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...RequestOpt) ([]discord.Message, error)
	CreateMessage(channelID discord.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (*discord.Message, error)
	UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, error)
	DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error
	BulkDeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...RequestOpt) error
	CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error)

	GetReactions(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) ([]discord.User, error)
	AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) error
	RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) error
	RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...RequestOpt) error
	RemoveAllReactions(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error
	RemoveAllReactionsForEmoji(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) error

	GetPinnedMessages(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Message, error)
	PinMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error
	UnpinMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error
	// TODO: add missing endpoints
}

type channelServiceImpl struct {
	restClient Client
}

func (s *channelServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *channelServiceImpl) GetChannel(channelID discord.Snowflake, opts ...RequestOpt) (channel discord.Channel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetChannel.Compile(nil, channelID)
	if err != nil {
		return
	}
	var ch discord.UnmarshalChannel
	err = s.restClient.DoBot(compiledRoute, nil, &ch, opts...)
	if err == nil {
		channel = ch.Channel
	}
	return
}

func (s *channelServiceImpl) UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (channel discord.Channel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateChannel.Compile(nil, channelID)
	if err != nil {
		return
	}
	var ch discord.UnmarshalChannel
	err = s.restClient.DoBot(compiledRoute, channelUpdate, &ch, opts...)
	if err == nil {
		channel = ch.Channel
	}
	return
}

func (s *channelServiceImpl) DeleteChannel(channelID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteChannel.Compile(nil, channelID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetWebhooks(channelID discord.Snowflake, opts ...RequestOpt) (webhooks []discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetChannelWebhooks.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *channelServiceImpl) CreateWebhook(channelID discord.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateWebhook.Compile(nil, channelID)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.restClient.DoBot(compiledRoute, webhookCreate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *channelServiceImpl) GetPermissionOverwrites(channelID discord.Snowflake, opts ...RequestOpt) (overwrites []discord.PermissionOverwrite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetPermissionOverwrites.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &overwrites, opts...)
	return
}

func (s *channelServiceImpl) GetPermissionOverwrite(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) (overwrite *discord.PermissionOverwrite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetPermissionOverwrite.Compile(nil, channelID, overwriteID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &overwrite, opts...)
	return
}

func (s *channelServiceImpl) UpdatePermissionOverwrite(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdatePermissionOverwrite.Compile(nil, channelID, overwriteID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, permissionOverwrite, nil, opts...)
}

func (s *channelServiceImpl) DeletePermissionOverwrite(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeletePermissionOverwrite.Compile(nil, channelID, overwriteID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) SendTyping(channelID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.SendTyping.Compile(nil, channelID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &message, opts...)
	return
}

func (s *channelServiceImpl) GetMessages(channelID discord.Snowflake, around discord.Snowflake, before discord.Snowflake, after discord.Snowflake, limit int, opts ...RequestOpt) (messages []discord.Message, err error) {
	values := route.QueryValues{}
	if around != "" {
		values["around"] = around
	}
	if before != "" {
		values["before"] = before
	}
	if after != "" {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMessages.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &messages, opts...)
	return
}

func (s *channelServiceImpl) CreateMessage(channelID discord.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return
	}
	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, body, &message, opts...)
	return
}

func (s *channelServiceImpl) UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, body, &message, opts...)
	return
}

func (s *channelServiceImpl) DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) BulkDeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.BulkDeleteMessages.Compile(nil, channelID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, discord.MessageBulkDelete{Messages: messageIDs}, nil, opts...)
}

func (s *channelServiceImpl) CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &message, opts...)
	return
}

func (s *channelServiceImpl) GetReactions(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) (users []discord.User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetReactions.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &users, opts...)
	return
}

func (s *channelServiceImpl) AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) error {
	compiledRoute, err := route.AddReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveOwnReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveUserReaction.Compile(nil, channelID, messageID, emoji, userID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveAllReactions(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveAllReactions.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveAllReactionsForEmoji(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveAllReactionsForEmoji.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetPinnedMessages(channelID discord.Snowflake, opts ...RequestOpt) (messages []discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetPinnedMessages.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.DoBot(compiledRoute, nil, &messages, opts...)
	return
}

func (s *channelServiceImpl) PinMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.PinMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) UnpinMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.UnpinMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.DoBot(compiledRoute, nil, nil, opts...)
}
