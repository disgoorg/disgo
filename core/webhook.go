package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/disgo/webhook"
)

type Webhook interface {
	discord.Webhook
}

type IncomingWebhook struct {
	discord.IncomingWebhook
	Bot  Bot
	User *User
}

func (h *IncomingWebhook) URL() string {
	compiledRoute, _ := route.WebhookURL.Compile(nil, h.ID(), h.Token)
	return compiledRoute.URL()
}

func (h *IncomingWebhook) Update(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*IncomingWebhook, error) {
	wh, err := h.Bot.RestServices().WebhookService().UpdateWebhook(h.ID(), webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.Bot.EntityBuilder().CreateWebhook(wh, CacheStrategyNoWs).(*IncomingWebhook), nil
}

func (h *IncomingWebhook) Delete(opts ...rest.RequestOpt) error {
	return h.Bot.RestServices().WebhookService().DeleteWebhook(h.ID(), opts...)
}

func (h *IncomingWebhook) NewWebhookClient(opts ...webhook.ConfigOpt) *webhook.Client {
	return webhook.NewClient(h.ID(), h.Token, append(opts, webhook.WithRestClient(h.Bot.RestServices().RestClient()), webhook.WithLogger(h.Bot.Logger))...)
}

type ChannelFollowerWebhook struct {
	discord.ChannelFollowerWebhook
	Bot  Bot
	User *User
}

func (h *ChannelFollowerWebhook) Update(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*ChannelFollowerWebhook, error) {
	wh, err := h.Bot.RestServices().WebhookService().UpdateWebhook(h.ID(), webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.Bot.EntityBuilder().CreateWebhook(wh, CacheStrategyNoWs).(*ChannelFollowerWebhook), nil
}

func (h *ChannelFollowerWebhook) Delete(opts ...rest.RequestOpt) error {
	return h.Bot.RestServices().WebhookService().DeleteWebhook(h.ID(), opts...)
}

type ApplicationWebhook struct {
	discord.ApplicationWebhook
	Bot Bot
}
