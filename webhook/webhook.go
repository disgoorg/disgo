package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Webhook struct {
	discord.Webhook
	WebhookClient *Client
}

func (h *Webhook) Update(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, error) {
	return h.WebhookClient.UpdateWebhook(webhookUpdate, opts...)
}

func (h *Webhook) Delete(opts ...rest.RequestOpt) error {
	return h.WebhookClient.DeleteWebhook(opts...)
}
