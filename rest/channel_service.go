package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewChannelService(client Client) ChannelService {
	return nil
}

type ChannelService interface {
	Service
	GetChannel(channelID discord.Snowflake) (*discord.Channel, Error)
	UpdateChannel(channelID discord.Snowflake, channelUpdate discord.ChannelUpdate) (*discord.Channel, Error)
	DeleteChannel(channelID discord.Snowflake) Error

	GetWebhooks(channelID discord.Snowflake) ([]discord.Webhook, Error)
	CreateWebhook(channelID discord.Snowflake, update discord.WebhookCreate) (*discord.Webhook, Error)

	UpdatePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate) Error
	DeletePermissionOverride(channelID discord.Snowflake, overwriteID discord.Snowflake) Error

	SendTyping(channelID discord.Snowflake) Error

	GetMessage(channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)
	CreateMessage(channelID discord.Snowflake, message discord.MessageCreate) (*discord.Message, Error)
	UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake) Error
	BulkDeleteMessages(channelID discord.Snowflake, messageIDs ...discord.Snowflake) Error
	CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)

	AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake) Error
}
