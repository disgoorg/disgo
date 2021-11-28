package webhook

import (
	"fmt"

	"github.com/DisgoOrg/disgo/discord"
)

// EntityBuilder is used to transform discord package entities into webhook package entities which hold a reference to the webhook client
type EntityBuilder interface {
	// WebhookClient returns the underlying webhook client used by this EntityBuilder
	WebhookClient() *Client

	// CreateMessage returns a new webhook.Message from the discord.Message

	CreateMessage(message discord.Message) *Message

	// CreateWebhook returns a new webhook.Webhook from the discord.Webhook
	CreateWebhook(webhook discord.Webhook) *Webhook
}

// NewEntityBuilder returns a new default EntityBuilder
func NewEntityBuilder(webhookClient *Client) EntityBuilder {
	return &entityBuilderImpl{
		webhookClient: webhookClient,
	}
}

type entityBuilderImpl struct {
	webhookClient *Client
}

func (b *entityBuilderImpl) WebhookClient() *Client {
	return b.webhookClient
}

func (b *entityBuilderImpl) CreateMessage(message discord.Message) *Message {
	return &Message{
		Message:       message,
		WebhookClient: b.WebhookClient(),
	}
}

func (b *entityBuilderImpl) CreateWebhook(webhook discord.Webhook) *Webhook {
	if w, ok := webhook.(discord.IncomingWebhook); ok {
		return &Webhook{
			IncomingWebhook: w,
			WebhookClient:   b.WebhookClient(),
		}
	}
	panic(fmt.Sprintf("invalid webhook type %d received", webhook.Type()))
}
