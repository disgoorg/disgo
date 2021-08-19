package webhook

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Webhook struct {
	discord.Webhook
	WebhookClient Client
}

func (h *Webhook) Update(ctx context.Context, webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error) {
	return h.WebhookClient.UpdateWebhook(ctx, webhookUpdate)
}

func (h *Webhook) Delete(ctx context.Context) rest.Error {
	return h.WebhookClient.DeleteWebhook(ctx)
}
