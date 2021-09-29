package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Webhook struct {
	discord.Webhook
	Bot *Bot
}

// URL returns the URL of this Webhook
func (h *Webhook) URL() string {
	if h.Type != discord.WebhookTypeIncoming {
		return ""
	}
	compiledRoute, _ := route.WebhookURL.Compile(nil, h.ID, h.Token)
	return compiledRoute.URL()
}

// Update updates this Webhook with the properties provided in discord.WebhookUpdate
func (h *Webhook) Update(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, rest.Error) {
	webhook, err := h.Bot.RestServices.WebhookService().UpdateWebhook(h.ID, webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.Bot.EntityBuilder.CreateWebhook(*webhook), nil
}

// Delete deletes this Webhook
func (h *Webhook) Delete(opts ...rest.RequestOpt) rest.Error {
	return h.Bot.RestServices.WebhookService().DeleteWebhook(h.ID, opts...)
}
