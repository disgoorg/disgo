package webhook

import (
	"context"
	"net/http"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// New returns a new Client
func New(client *http.Client, restClient rest.Client, logger log.Logger, id discord.Snowflake, token string) Client {
	if client == nil {
		client = http.DefaultClient
	}
	if restClient == nil {
		// TODO: user agent
		config := &rest.DefaultConfig
		config.Headers = http.Header{"authorization": []string{}}
		restClient = rest.NewClient(logger, client, nil, config)
	}
	if logger == nil {
		logger = log.Default()
	}
	webhookClient := &webhookClientImpl{
		logger:                 logger,
		webhookService:         rest.NewWebhookService(restClient),
		defaultAllowedMentions: &DefaultAllowedMentions,
		id:                     id,
		token:                  token,
	}
	webhookClient.entityBuilder = NewEntityBuilder(webhookClient)

	webhookClient.webhookService = nil
	return webhookClient
}

type webhookClientImpl struct {
	logger                 log.Logger
	webhookService         rest.WebhookService
	entityBuilder          EntityBuilder
	defaultAllowedMentions *discord.AllowedMentions
	id                     discord.Snowflake
	token                  string
}

func (h *webhookClientImpl) Logger() log.Logger {
	return h.logger
}

func (h *webhookClientImpl) RestClient() rest.Client {
	return h.webhookService.RestClient()
}
func (h *webhookClientImpl) WebhookService() rest.WebhookService {
	return h.webhookService
}
func (h *webhookClientImpl) EntityBuilder() EntityBuilder {
	return h.entityBuilder
}
func (h *webhookClientImpl) DefaultAllowedMentions() *discord.AllowedMentions {
	return h.defaultAllowedMentions
}
func (h *webhookClientImpl) SetDefaultAllowedMentions(allowedMentions *discord.AllowedMentions) {
	h.defaultAllowedMentions = allowedMentions
}

func (h *webhookClientImpl) GetWebhook(ctx context.Context) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().GetWebhookWithToken(ctx, h.id, h.token)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) UpdateWebhook(ctx context.Context, webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().UpdateWebhookWithToken(ctx, h.id, h.token, webhookUpdate)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) DeleteWebhook(ctx context.Context) rest.Error {
	return h.WebhookService().DeleteWebhookWithToken(ctx, h.id, h.token)
}

func (h *webhookClientImpl) CreateMessageInThread(ctx context.Context, messageCreate discord.MessageCreate, threadID discord.Snowflake) (*Message, rest.Error) {
	message, err := h.WebhookService().CreateMessage(ctx, h.id, h.token, messageCreate, true, threadID)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateMessage(*message), nil
}

func (h *webhookClientImpl) CreateMessage(ctx context.Context, messageCreate discord.MessageCreate) (*Message, rest.Error) {
	return h.CreateMessageInThread(ctx, messageCreate, "")
}

func (h *webhookClientImpl) CreateContent(ctx context.Context, content string) (*Message, rest.Error) {
	return h.CreateMessage(ctx, discord.MessageCreate{Content: content})
}

func (h *webhookClientImpl) CreateEmbeds(ctx context.Context, embeds ...discord.Embed) (*Message, rest.Error) {
	return h.CreateMessage(ctx, discord.MessageCreate{Embeds: embeds})
}

func (h *webhookClientImpl) UpdateMessage(ctx context.Context, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := h.WebhookService().UpdateMessage(ctx, h.id, h.token, messageID, messageUpdate)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateMessage(*message), nil
}

func (h *webhookClientImpl) UpdateContent(ctx context.Context, messageID discord.Snowflake, content string) (*Message, rest.Error) {
	return h.UpdateMessage(ctx, messageID, discord.MessageUpdate{Content: &content})
}

func (h *webhookClientImpl) UpdateEmbeds(ctx context.Context, messageID discord.Snowflake, embeds ...discord.Embed) (*Message, rest.Error) {
	return h.UpdateMessage(ctx, messageID, discord.MessageUpdate{Embeds: embeds})
}

func (h *webhookClientImpl) DeleteMessage(ctx context.Context, messageID discord.Snowflake) rest.Error {
	return h.WebhookService().DeleteMessage(ctx, h.id, h.token, messageID)
}

func (h *webhookClientImpl) ID() discord.Snowflake {
	return h.id
}
func (h *webhookClientImpl) Token() string {
	return h.token
}

func (h *webhookClientImpl) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, h.id, h.token)
	return compiledRoute.URL()
}
