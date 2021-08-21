package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ WebhookService = (*WebhookServiceImpl)(nil)

func NewWebhookService(client Client) WebhookService {
	return &WebhookServiceImpl{restClient: client}
}

type WebhookService interface {
	Service
	GetWebhook(ctx context.Context, webhookID discord.Snowflake) (*discord.Webhook, Error)
	UpdateWebhook(ctx context.Context, webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate) (*discord.Webhook, Error)
	DeleteWebhook(ctx context.Context, webhookID discord.Snowflake) Error

	GetWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string) (*discord.Webhook, Error)
	UpdateWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdate) (*discord.Webhook, Error)
	DeleteWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string) Error

	CreateMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	CreateMessageSlack(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	CreateMessageGitHub(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	UpdateMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate) (*discord.Message, Error)
	DeleteMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake) Error
}

type WebhookServiceImpl struct {
	restClient Client
}

func (s *WebhookServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *WebhookServiceImpl) GetWebhook(ctx context.Context, webhookID discord.Snowflake) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.GetWebhook.Compile(nil, webhookID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &webhook)
	return
}

func (s *WebhookServiceImpl) UpdateWebhook(ctx context.Context, webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.UpdateWebhook.Compile(nil, webhookID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, webhookUpdate, &webhook)
	return
}

func (s *WebhookServiceImpl) DeleteWebhook(ctx context.Context, webhookID discord.Snowflake) (rErr Error) {
	compiledRoute, err := route.DeleteWebhook.Compile(nil, webhookID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, nil)
	return
}

func (s *WebhookServiceImpl) GetWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.GetWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, &webhook)
	return
}

func (s *WebhookServiceImpl) UpdateWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdate) (webhook *discord.Webhook, rErr Error) {
	compiledRoute, err := route.UpdateWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, webhookUpdate, &webhook)
	return
}

func (s *WebhookServiceImpl) DeleteWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string) (rErr Error) {
	compiledRoute, err := route.DeleteWebhookWithToken.Compile(nil, webhookID, webhookToken)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, nil)
	return
}

func (s *WebhookServiceImpl) createMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake, apiRoute *route.APIRoute) (message *discord.Message, rErr Error) {
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
		rErr = s.restClient.Do(ctx, compiledRoute, body, &message)
	} else {
		rErr = s.restClient.Do(ctx, compiledRoute, body, nil)
	}
	return
}

func (s *WebhookServiceImpl) CreateMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.WebhookMessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error) {
	return s.createMessage(ctx, webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessage)
}

func (s *WebhookServiceImpl) CreateMessageSlack(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake) (*discord.Message, Error) {
	return s.createMessage(ctx, webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageSlack)
}

func (s *WebhookServiceImpl) CreateMessageGitHub(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.Payload, wait bool, threadID discord.Snowflake) (*discord.Message, Error) {
	return s.createMessage(ctx, webhookID, webhookToken, messageCreate, wait, threadID, route.CreateWebhookMessageGitHub)
}

func (s *WebhookServiceImpl) UpdateMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate) (message *discord.Message, rErr Error) {
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

	rErr = s.restClient.Do(ctx, compiledRoute, body, &message)
	return
}

func (s *WebhookServiceImpl) DeleteMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake) (rErr Error) {
	compiledRoute, err := route.DeleteWebhookMessage.Compile(nil, webhookID, webhookToken, messageID)
	if err != nil {
		rErr = NewError(nil, err)
		return
	}
	rErr = s.restClient.Do(ctx, compiledRoute, nil, nil)
	return
}
