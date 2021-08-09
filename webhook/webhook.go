package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Webhook struct {
	discord.Webhook
	WebhookClient Client
}

func (h *Webhook) Update(webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error) {
	return h.WebhookClient.UpdateWebhook(webhookUpdate)
}

func (h *Webhook) Delete() rest.Error {
	return h.WebhookClient.DeleteWebhook()
}
