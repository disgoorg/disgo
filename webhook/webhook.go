package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// Webhook can be used to update or delete the Webhook
type Webhook struct {
	discord.IncomingWebhook
	WebhookClient Client
}

// Update is used to update the Webhook
func (h *Webhook) Update(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*Webhook, error) {
	return h.WebhookClient.UpdateWebhook(webhookUpdate, opts...)
}

// Delete is used to delete the Webhook
func (h *Webhook) Delete(opts ...rest.RequestOpt) error {
	return h.WebhookClient.DeleteWebhook(opts...)
}
