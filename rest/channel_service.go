package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewChannelService(client Client) ChannelService {
	return nil
}

type ChannelService interface {
	Service
	GetChannel(channelID discord.Snowflake, opts ...RequestOpt) (*discord.Channel, Error)
	UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate, opts ...RequestOpt) (*discord.Channel, Error)
	DeleteChannel(channelID discord.Snowflake, opts ...RequestOpt) Error

	GetWebhooks(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Webhook, Error)
	CreateWebhook(channelID discord.Snowflake, update discord.WebhookCreate, opts ...RequestOpt) (*discord.Webhook, Error)

	UpdatePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate, opts ...RequestOpt) Error
	DeletePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, opts ...RequestOpt) Error

	SendTyping(channelID discord.Snowflake, opts ...RequestOpt) Error

	GetMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)
	CreateMessage(channelID discord.Snowflake, message discord.MessageCreate, opts ...RequestOpt) (*discord.Message, Error)
	UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate, opts ...RequestOpt) (*discord.Message, Error)
	DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) Error
	BulkDeleteMessages(channelID discord.Snowflake, messageIDs []discord.Snowflake, opts ...RequestOpt) Error
	CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)

	AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error
	RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, opts ...RequestOpt) Error
	RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake, opts ...RequestOpt) Error
}
