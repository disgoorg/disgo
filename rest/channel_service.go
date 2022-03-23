package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
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
	GetChannel(channelID snowflake.Snowflake, opts ...RequestOpt) (discord.Channel, error)
	UpdateChannel(channelID snowflake.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (discord.Channel, error)
	DeleteChannel(channelID snowflake.Snowflake, opts ...RequestOpt) error

	GetWebhooks(channelID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Webhook, error)
	CreateWebhook(channelID snowflake.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (discord.Webhook, error)

	GetPermissionOverwrites(channelID snowflake.Snowflake, opts ...RequestOpt) ([]discord.PermissionOverwrite, error)
	GetPermissionOverwrite(channelID snowflake.Snowflake, overwriteID snowflake.Snowflake, opts ...RequestOpt) (*discord.PermissionOverwrite, error)
	UpdatePermissionOverwrite(channelID snowflake.Snowflake, overwriteID snowflake.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) error
	DeletePermissionOverwrite(channelID snowflake.Snowflake, overwriteID snowflake.Snowflake, opts ...RequestOpt) error

	SendTyping(channelID snowflake.Snowflake, opts ...RequestOpt) error

	GetMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	GetMessages(channelID snowflake.Snowflake, around snowflake.Snowflake, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) ([]discord.Message, error)
	CreateMessage(channelID snowflake.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (*discord.Message, error)
	UpdateMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, error)
	DeleteMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error
	BulkDeleteMessages(channelID snowflake.Snowflake, messageIDs []snowflake.Snowflake, opts ...RequestOpt) error
	CrosspostMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error)

	GetReactions(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) ([]discord.User, error)
	AddReaction(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) error
	RemoveOwnReaction(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) error
	RemoveUserReaction(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, userID snowflake.Snowflake, opts ...RequestOpt) error
	RemoveAllReactions(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error
	RemoveAllReactionsForEmoji(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) error

	GetPinnedMessages(channelID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Message, error)
	PinMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error
	UnpinMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error
	// TODO: add missing endpoints
}

type channelServiceImpl struct {
	restClient Client
}

func (s *channelServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *channelServiceImpl) GetChannel(channelID snowflake.Snowflake, opts ...RequestOpt) (channel discord.Channel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetChannel.Compile(nil, channelID)
	if err != nil {
		return
	}
	var ch discord.UnmarshalChannel
	err = s.restClient.Do(compiledRoute, nil, &ch, opts...)
	if err == nil {
		channel = ch.Channel
	}
	return
}

func (s *channelServiceImpl) UpdateChannel(channelID snowflake.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (channel discord.Channel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateChannel.Compile(nil, channelID)
	if err != nil {
		return
	}
	var ch discord.UnmarshalChannel
	err = s.restClient.Do(compiledRoute, channelUpdate, &ch, opts...)
	if err == nil {
		channel = ch.Channel
	}
	return
}

func (s *channelServiceImpl) DeleteChannel(channelID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteChannel.Compile(nil, channelID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetWebhooks(channelID snowflake.Snowflake, opts ...RequestOpt) (webhooks []discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetChannelWebhooks.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *channelServiceImpl) CreateWebhook(channelID snowflake.Snowflake, webhookCreate discord.WebhookCreate, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateWebhook.Compile(nil, channelID)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.restClient.Do(compiledRoute, webhookCreate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *channelServiceImpl) GetPermissionOverwrites(channelID snowflake.Snowflake, opts ...RequestOpt) (overwrites []discord.PermissionOverwrite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetPermissionOverwrites.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &overwrites, opts...)
	return
}

func (s *channelServiceImpl) GetPermissionOverwrite(channelID snowflake.Snowflake, overwriteID snowflake.Snowflake, opts ...RequestOpt) (overwrite *discord.PermissionOverwrite, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetPermissionOverwrite.Compile(nil, channelID, overwriteID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &overwrite, opts...)
	return
}

func (s *channelServiceImpl) UpdatePermissionOverwrite(channelID snowflake.Snowflake, overwriteID snowflake.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdatePermissionOverwrite.Compile(nil, channelID, overwriteID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, permissionOverwrite, nil, opts...)
}

func (s *channelServiceImpl) DeletePermissionOverwrite(channelID snowflake.Snowflake, overwriteID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeletePermissionOverwrite.Compile(nil, channelID, overwriteID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) SendTyping(channelID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.SendTyping.Compile(nil, channelID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *channelServiceImpl) GetMessages(channelID snowflake.Snowflake, around snowflake.Snowflake, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) (messages []discord.Message, err error) {
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
	compiledRoute, err = route.GetMessages.Compile(values, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &messages, opts...)
	return
}

func (s *channelServiceImpl) CreateMessage(channelID snowflake.Snowflake, messageCreate discord.MessageCreate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateMessage.Compile(nil, channelID)
	if err != nil {
		return
	}
	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *channelServiceImpl) UpdateMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *channelServiceImpl) DeleteMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) BulkDeleteMessages(channelID snowflake.Snowflake, messageIDs []snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.BulkDeleteMessages.Compile(nil, channelID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, discord.MessageBulkDelete{Messages: messageIDs}, nil, opts...)
}

func (s *channelServiceImpl) CrosspostMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) (message *discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CrosspostMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &message, opts...)
	return
}

func (s *channelServiceImpl) GetReactions(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) (users []discord.User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetReactions.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &users, opts...)
	return
}

func (s *channelServiceImpl) AddReaction(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) error {
	compiledRoute, err := route.AddReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveOwnReaction(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveOwnReaction.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveUserReaction(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, userID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveUserReaction.Compile(nil, channelID, messageID, emoji, userID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveAllReactions(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveAllReactions.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) RemoveAllReactionsForEmoji(channelID snowflake.Snowflake, messageID snowflake.Snowflake, emoji string, opts ...RequestOpt) error {
	compiledRoute, err := route.RemoveAllReactionsForEmoji.Compile(nil, channelID, messageID, emoji)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) GetPinnedMessages(channelID snowflake.Snowflake, opts ...RequestOpt) (messages []discord.Message, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetPinnedMessages.Compile(nil, channelID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &messages, opts...)
	return
}

func (s *channelServiceImpl) PinMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.PinMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}

func (s *channelServiceImpl) UnpinMessage(channelID snowflake.Snowflake, messageID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.UnpinMessage.Compile(nil, channelID, messageID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
