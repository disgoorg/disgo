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

func (h *Webhook) URL() string {
	if h.Type != discord.WebhookTypeIncoming {
		return ""
	}
	compiledRoute, _ := route.WebhookURL.Compile(nil, h.ID, h.Token)
	return compiledRoute.URL()
}

func (h *Webhook) Update(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.Bot.RestServices.WebhookService().UpdateWebhook(h.ID, webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.Bot.EntityBuilder.CreateWebhook(*webhook), nil
}

func (h *Webhook) Delete(opts ...rest.RequestOpt) error {
	return h.Bot.RestServices.WebhookService().DeleteWebhook(h.ID, opts...)
}
