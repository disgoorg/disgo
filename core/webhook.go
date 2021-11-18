package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Webhook interface {
	discord.Webhook
}

type IncomingWebhook struct {
	discord.IncomingWebhook
	Bot  *Bot
	User *User
}

func (h *IncomingWebhook) URL() string {
	compiledRoute, _ := route.WebhookURL.Compile(nil, h.ID(), h.Token)
	return compiledRoute.URL()
}

func (h *IncomingWebhook) Update(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*IncomingWebhook, error) {
	webhook, err := h.Bot.RestServices.WebhookService().UpdateWebhook(h.ID(), webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.Bot.EntityBuilder.CreateWebhook(webhook, CacheStrategyNoWs).(*IncomingWebhook), nil
}

func (h *IncomingWebhook) Delete(opts ...rest.RequestOpt) error {
	return h.Bot.RestServices.WebhookService().DeleteWebhook(h.ID(), opts...)
}

type ChannelFollowerWebhook struct {
	discord.ChannelFollowerWebhook
	Bot  *Bot
	User *User
}

type ApplicationWebhook struct {
	discord.ApplicationWebhook
	Bot *Bot
}
