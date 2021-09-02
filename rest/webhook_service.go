package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ WebhookService = (*WebhookServiceImpl)(nil)

func NewWebhookService(restClient Client) WebhookService {
	return &WebhookServiceImpl{restClient: restClient}
}

type WebhookService interface {
	Service
	GetWebhook(webhookID discord.Snowflake, opts ...RequestOpt) (*discord.Webhook, Error)
	UpdateWebhook(webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (*discord.Webhook, Error)
	DeleteWebhook(webhookID discord.Snowflake, opts ...RequestOpt) Error

	GetWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) (*discord.Webhook, Error)
	UpdateWebhookWithToken(webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (*discord.Webhook, Error)
	DeleteWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) Error

	CreateMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)
	CreateMessageSlack(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)
	CreateMessageGitHub(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error)
	UpdateMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...RequestOpt) (*discord.Message, Error)
	DeleteMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, opts ...RequestOpt) Error
}

type WebhookServiceImpl struct {
	restClient Client
}

func (s *WebhookServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *WebhookServiceImpl) GetWebhook(webhookID discord.Snowflake, opts ...RequestOpt) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.GetWebhook.Compile(nil, webhookID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, &webhook, opts...)
	return
}

func (s *WebhookServiceImpl) UpdateWebhook(webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.UpdateWebhook.Compile(nil, webhookID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, webhookUpdate, &webhook, opts...)
	return
}

func (s *WebhookServiceImpl) DeleteWebhook(webhookID discord.Snowflake, opts ...RequestOpt) (rErr Error) {
	compiledRoute, err := route.DeleteWebhook.Compile(nil, webhookID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *WebhookServiceImpl) GetWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.GetWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, &webhook, opts...)
	return
}

func (s *WebhookServiceImpl) UpdateWebhookWithToken(webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdate, opts ...RequestOpt) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.UpdateWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, webhookUpdate, &webhook, opts...)
	return
}

func (s *WebhookServiceImpl) DeleteWebhookWithToken(webhookID discord.Snowflake, webhookToken string, opts ...RequestOpt) (rErr Error) {
	compiledRoute, err := route.DeleteWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}

func (s *WebhookServiceImpl) createMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, apiRoute *route.APIRoute, opts []RequestOpt) (message *discord.Message, rErr Error) {
	params := route.QueryValues{}
	if wait {
		params["wait"] = true
	}
	if threadID != "" {
		params["thread_id"] = threadID
	}
	compiledRoute, err := apiRoute.Compile(params, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}

	body, err := messageCreate.ToBody()
	if err != nil {
		rErr = NewError(nil, err)
		return
	}

	if wait {
		rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	} else {
		rErr = s.restClient.Do(compiledRoute, body, nil, opts...)
	}
	return
}

func (s *WebhookServiceImpl) CreateMessage(webhookID discord.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessage, opts)
}

func (s *WebhookServiceImpl) CreateMessageSlack(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageSlack, opts)
}

func (s *WebhookServiceImpl) CreateMessageGitHub(webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, opts ...RequestOpt) (*discord.Message, Error) {
	return s.createMessage(webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageGitHub, opts)
}

func (s *WebhookServiceImpl) UpdateMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...RequestOpt) (message *discord.Message, rErr Error) {
	compiledRoute, err := route.UpdateWebhookMessage.Compile(nil, webhookID, webhookToken, messageID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}

	body, err := messageUpdate.ToBody()
	if err != nil {
		rErr = NewError(nil, err)
		return
	}

	rErr = s.restClient.Do(compiledRoute, body, &message, opts...)
	return
}

func (s *WebhookServiceImpl) DeleteMessage(webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, opts ...RequestOpt) (rErr Error) {
	compiledRoute, err := route.DeleteWebhookMessage.Compile(nil, webhookID, webhookToken, messageID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(compiledRoute, nil, nil, opts...)
	return
}
