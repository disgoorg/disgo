package webhook

import (
	"net/http"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// New returns a new Client
//goland:noinspection GoUnusedExportedFunction
func New(restClient rest.Client, logger log.Logger, id discord.Snowflake, token string) Client {
	if restClient == nil {
		config := &rest.DefaultConfig
		config.Headers = http.Header{"authorization": []string{}}
		restClient = rest.NewClient(logger, nil, nil, config)
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

func (h *webhookClientImpl) GetWebhook(opts ...rest.RequestOpt) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().GetWebhookWithToken(h.id, h.token, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) UpdateWebhook(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().UpdateWebhookWithToken(h.id, h.token, webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) DeleteWebhook(opts ...rest.RequestOpt) rest.Error {
	return h.WebhookService().DeleteWebhookWithToken(h.id, h.token)
}

func (h *webhookClientImpl) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := h.WebhookService().CreateMessage(h.id, h.token, messageCreate, true, threadID, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateMessage(*message), nil
}

func (h *webhookClientImpl) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	return h.CreateMessageInThread(messageCreate, "", opts...)
}

func (h *webhookClientImpl) CreateContent(content string, opts ...rest.RequestOpt) (*Message, rest.Error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Content: content}, opts...)
}

func (h *webhookClientImpl) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, rest.Error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, opts...)
}

func (h *webhookClientImpl) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := h.WebhookService().UpdateMessage(h.id, h.token, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateMessage(*message), nil
}

func (h *webhookClientImpl) UpdateContent(messageID discord.Snowflake, content string, opts ...rest.RequestOpt) (*Message, rest.Error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, opts...)
}

func (h *webhookClientImpl) UpdateEmbeds(messageID discord.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, rest.Error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: embeds}, opts...)
}

func (h *webhookClientImpl) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return h.WebhookService().DeleteMessage(h.id, h.token, messageID, opts...)
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
