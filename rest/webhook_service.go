package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var (
	_ Service        = (*webhookServiceImpl)(nil)
	_ WebhookService = (*webhookServiceImpl)(nil)
)

func NewWebhookService(restClient Client) WebhookService {
	return &webhookServiceImpl{restClient: restClient}
}

type WebhookService interface {
	Service
	GetWebhook(webhookID discord.Snowflake, opts ...RequestOpt) (discord.Webhook, error)
	UpdateWebhook(webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (discord.Webhook, error)
	DeleteWebhook(webhookID discord.Snowflake, opts ...RequestOpt) error

	GetWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) (discord.Webhook, error)
	UpdateWebhookWithToken(webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdateWithToken, opts ...RequestOpt) (discord.Webhook, error)
	DeleteWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) error

	CreateMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	CreateMessageSlack(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	CreateMessageGitHub(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	UpdateMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error)
	DeleteMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, threadID discord.Snowflake, opts ...RequestOpt) error
}

type webhookServiceImpl struct {
	restClient Client
}

func (s *webhookServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *webhookServiceImpl) GetWebhook(webhookID discord.Snowflake, opts ...RequestOpt) (webhook discord.Webhook, err error) {
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

func (s *webhookServiceImpl) UpdateWebhook(webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (webhook discord.Webhook, err error) {
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

func (s *webhookServiceImpl) DeleteWebhook(webhookID discord.Snowflake, opts ...RequestOpt) (err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhook.Compile(nil, webhookID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *webhookServiceImpl) GetWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) (webhook discord.Webhook, err error) {
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

func (s *webhookServiceImpl) UpdateWebhookWithToken(webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdateWithToken, opts ...RequestOpt) (webhook discord.Webhook, err error) {
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

func (s *webhookServiceImpl) DeleteWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) (err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.DeleteWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *webhookServiceImpl) createMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, apiRoute *route.APIRoute, opts []RequestOpt) (message *discord.Message, err error) {
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

func (s *webhookServiceImpl) CreateMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessage, opts)
}

func (s *webhookServiceImpl) CreateMessageSlack(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageSlack, opts)
}

func (s *webhookServiceImpl) CreateMessageGitHub(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageGitHub, opts)
}

func (s *webhookServiceImpl) UpdateMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID discord.Snowflake, opts ...RequestOpt) (message *discord.Message, err error) {
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

func (s *webhookServiceImpl) DeleteMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, threadID discord.Snowflake, opts ...RequestOpt) (err error) {
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
