package rest

import "github.com/DisgoOrg/disgo/discord"

type WebhookService interface {
	Service
	GetWebhook(webhookID discord.Snowflake) (*discord.Webhook, Error)
	UpdateWebhook(webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate) (*discord.Webhook, Error)
	DeleteWebhook(webhookID discord.Snowflake) Error

	GetWebhookWithToken(webhookID discord.Snowflake, webhookToken string) (*discord.Webhook, Error)
	UpdateWebhookWithToken(webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdate) (*discord.Webhook, Error)
	DeleteWebhookWithToken(webhookID discord.Snowflake, webhookToken string) Error

	CreateMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.MessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	CreateMessageSlack(webhookID discord.Snowflake, webhookToken string, messageCreate discord.MessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	CreateMessageGitHub(webhookID discord.Snowflake, webhookToken string, messageCreate discord.MessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	UpdateMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake) Error
}

func NewWebhookService(client HTTPClient) WebhookService {
	return nil
}
