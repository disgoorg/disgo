package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Webhooks = (*webhookImpl)(nil)

func NewWebhooks(client Client) Webhooks {
	return &webhookImpl{client: client}
}

type Webhooks interface {
	GetWebhook(webhookID snowflake.ID, opts ...RequestOpt) (discord.Webhook, error)
	UpdateWebhook(webhookID snowflake.ID, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (discord.Webhook, error)
	DeleteWebhook(webhookID snowflake.ID, opts ...RequestOpt) error

	GetWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) (discord.Webhook, error)
	UpdateWebhookWithToken(webhookID snowflake.ID, webhookToken string, webhookUpdate discord.WebhookUpdateWithToken, opts ...RequestOpt) (discord.Webhook, error)
	DeleteWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) error

	CreateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error)
	CreateWebhookMessageSlack(webhookID snowflake.ID, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error)
	CreateWebhookMessageGitHub(webhookID snowflake.ID, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error)
	UpdateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error)
	DeleteWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, threadID snowflake.ID, opts ...RequestOpt) error
}

type webhookImpl struct {
	client Client
}

func (s *webhookImpl) GetWebhook(webhookID snowflake.ID, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.client.Do(compiledRoute, nil, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) UpdateWebhook(webhookID snowflake.ID, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.client.Do(compiledRoute, webhookUpdate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) DeleteWebhook(webhookID snowflake.ID, opts ...RequestOpt) (err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *webhookImpl) GetWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.client.Do(compiledRoute, nil, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) UpdateWebhookWithToken(webhookID snowflake.ID, webhookToken string, webhookUpdate discord.WebhookUpdateWithToken, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.client.Do(compiledRoute, webhookUpdate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) DeleteWebhookWithToken(webhookID snowflake.ID, webhookToken string, opts ...RequestOpt) (err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *webhookImpl) createWebhookMessage(webhookID snowflake.ID, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.ID, apiRoute *route.APIRoute, opts []RequestOpt) (message *discord.Message, err error) {
	params := route.QueryValues{}
	if wait {
		params["wait"] = true
	}
	if threadID != 0 {
		params["thread_id"] = threadID
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = apiRoute.Compile(params, webhookID, webhookToken)
	if err != nil {
		return
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		return
	}

	if wait {
		err = s.client.Do(compiledRoute, body, &message, opts...)
	} else {
		err = s.client.Do(compiledRoute, body, nil, opts...)
	}
	return
}

func (s *webhookImpl) CreateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error) {
	return s.createWebhookMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessage, opts)
}

func (s *webhookImpl) CreateWebhookMessageSlack(webhookID snowflake.ID, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error) {
	return s.createWebhookMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageSlack, opts)
}

func (s *webhookImpl) CreateWebhookMessageGitHub(webhookID snowflake.ID, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.ID, opts ...RequestOpt) (*discord.Message, error) {
	return s.createWebhookMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageGitHub, opts)
}

func (s *webhookImpl) UpdateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.ID, opts ...RequestOpt) (message *discord.Message, err error) {
	params := route.QueryValues{}
	if threadID != 0 {
		params["thread_id"] = threadID
	}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateWebhookMessage.Compile(params, webhookID, webhookToken, messageID)
	if err != nil {
		return
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *webhookImpl) DeleteWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, threadID snowflake.ID, opts ...RequestOpt) (err error) {
	params := route.QueryValues{}
	if threadID != 0 {
		params["thread_id"] = threadID
	}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhookMessage.Compile(params, webhookID, webhookToken, messageID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, nil, opts...)
	return
}
