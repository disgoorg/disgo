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
	GetChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, Error)
	UpdateChannel(ctx context.Context, channelID discord.Snowflake, channelUpdate discord.ChannelUpdate) (*discord.Channel, Error)
	DeleteChannel(ctx context.Context, channelID discord.Snowflake) Error

	GetWebhooks(ctx context.Context, channelID discord.Snowflake) ([]discord.Webhook, Error)
	CreateWebhook(ctx context.Context, channelID discord.Snowflake, update discord.WebhookCreate) (*discord.Webhook, Error)

	UpdatePermissionOverride(ctx context.Context, channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate) Error
	DeletePermissionOverride(ctx context.Context, channelID discord.Snowflake, overwriteID discord.Snowflake) Error

	SendTyping(ctx context.Context, channelID discord.Snowflake) Error

	GetMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)
	CreateMessage(ctx context.Context, channelID discord.Snowflake, message discord.MessageCreate) (*discord.Message, Error)
	UpdateMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake) Error
	BulkDeleteMessages(ctx context.Context, channelID discord.Snowflake, messageIDs ...discord.Snowflake) Error
	CrosspostMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)

	AddReaction(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveOwnReaction(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveUserReaction(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake) Error
}
