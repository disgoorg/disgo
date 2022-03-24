package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Webhooks = (*webhookImpl)(nil)

func NewWebhooks(restClient Client) Webhooks {
	return &webhookImpl{restClient: restClient}
}

type Webhooks interface {
	GetWebhook(webhookID snowflake.Snowflake, opts ...RequestOpt) (discord.Webhook, error)
	UpdateWebhook(webhookID snowflake.Snowflake, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (discord.Webhook, error)
	DeleteWebhook(webhookID snowflake.Snowflake, opts ...RequestOpt) error

	GetWebhookWithToken(webhookID snowflake.Snowflake, webhookToken string, opts ...RequestOpt) (discord.Webhook, error)
	UpdateWebhookWithToken(webhookID snowflake.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdateWithToken, opts ...RequestOpt) (discord.Webhook, error)
	DeleteWebhookWithToken(webhookID snowflake.Snowflake, webhookToken string, opts ...RequestOpt) error

	CreateMessage(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	CreateMessageSlack(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	CreateMessageGitHub(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	UpdateMessage(webhookID snowflake.Snowflake, webhookToken string, messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	DeleteMessage(webhookID snowflake.Snowflake, webhookToken string, messageID snowflake.Snowflake, threadID snowflake.Snowflake, opts ...RequestOpt) error
}

type webhookImpl struct {
	restClient Client
}

func (s *webhookImpl) GetWebhook(webhookID snowflake.Snowflake, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.restClient.Do(compiledRoute, nil, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) UpdateWebhook(webhookID snowflake.Snowflake, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.restClient.Do(compiledRoute, webhookUpdate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) DeleteWebhook(webhookID snowflake.Snowflake, opts ...RequestOpt) (err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *webhookImpl) GetWebhookWithToken(webhookID snowflake.Snowflake, webhookToken string, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.restClient.Do(compiledRoute, nil, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) UpdateWebhookWithToken(webhookID snowflake.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdateWithToken, opts ...RequestOpt) (webhook discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}

	var unmarshalWebhook discord.UnmarshalWebhook
	err = s.restClient.Do(compiledRoute, webhookUpdate, &unmarshalWebhook, opts...)
	if err == nil {
		webhook = unmarshalWebhook.Webhook
	}
	return
}

func (s *webhookImpl) DeleteWebhookWithToken(webhookID snowflake.Snowflake, webhookToken string, opts ...RequestOpt) (err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *webhookImpl) createMessage(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.Snowflake, apiRoute *route.APIRoute, opts []RequestOpt) (message *discord.Message, err error) {
	params := route.QueryValues{}
	if wait {
		params["wait"] = true
	}
	if threadID != "" {
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
		err = s.restClient.Do(compiledRoute, body, &message, opts...)
	} else {
		err = s.restClient.Do(compiledRoute, body, nil, opts...)
	}
	return
}

func (s *webhookImpl) CreateMessage(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessage, opts)
}

func (s *webhookImpl) CreateMessageSlack(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageSlack, opts)
}

func (s *webhookImpl) CreateMessageGitHub(webhookID snowflake.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID snowflake.Snowflake, opts ...RequestOpt) (*discord.Message, error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageGitHub, opts)
}

func (s *webhookImpl) UpdateMessage(webhookID snowflake.Snowflake, webhookToken string, messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.Snowflake, opts ...RequestOpt) (message *discord.Message, err error) {
	params := route.QueryValues{}
	if threadID != "" {
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

	err = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *webhookImpl) DeleteMessage(webhookID snowflake.Snowflake, webhookToken string, messageID snowflake.Snowflake, threadID snowflake.Snowflake, opts ...RequestOpt) (err error) {
	params := route.QueryValues{}
	if threadID != "" {
		params["thread_id"] = threadID
	}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhookMessage.Compile(params, webhookID, webhookToken, messageID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}
