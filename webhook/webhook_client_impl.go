package webhook

import (
	"net/http"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// New returns a new Client
func New(client *http.Client, httpClient rest.HTTPClient, logger log.Logger, id discord.Snowflake, token string) Client {
	if client == nil {
		client = http.DefaultClient
	}
	if httpClient == nil {
		// TODO: user agent
		httpClient = rest.NewHTTPClient(logger, client, "")
	}
	if logger == nil {
		logger = log.Default()
	}
	webhookClient := &webhookClientImpl{
		logger:                 logger,
		webhookService:         rest.NewWebhookService(httpClient),
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

func (h *webhookClientImpl) HTTPClient() rest.HTTPClient {
	return h.webhookService.HTTPClient()
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

func (h *webhookClientImpl) GetWebhook() (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().GetWebhookWithToken(h.id, h.token)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) UpdateWebhook(webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().UpdateWebhookWithToken(h.id, h.token, webhookUpdate)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) DeleteWebhook() rest.Error {
	return h.WebhookService().DeleteWebhookWithToken(h.id, h.token)
}

func (h *webhookClientImpl) CreateMessageInThread(messageCreate discord.MessageCreate, threadID discord.Snowflake) (*Message, rest.Error) {
	message, err := h.WebhookService().CreateMessage(h.id, h.token, messageCreate, true, threadID)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateMessage(*message), nil
}

func (h *webhookClientImpl) CreateMessage(messageCreate discord.MessageCreate) (*Message, rest.Error) {
	return h.CreateMessageInThread(messageCreate, "")
}

func (h *webhookClientImpl) CreateContent(content string) (*Message, rest.Error) {
	return h.CreateMessage(discord.MessageCreate{Content: content})
}

func (h *webhookClientImpl) CreateEmbeds(embeds ...discord.Embed) (*Message, rest.Error) {
	return h.CreateMessage(discord.MessageCreate{Embeds: embeds})
}

func (h *webhookClientImpl) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*Message, rest.Error) {
	message, err := h.WebhookService().UpdateMessage(h.id, h.token, messageID, messageUpdate)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateMessage(*message), nil
}

func (h *webhookClientImpl) UpdateContent(messageID discord.Snowflake, content string) (*Message, rest.Error) {
	return h.UpdateMessage(messageID, discord.MessageUpdate{Content: &content})
}

func (h *webhookClientImpl) UpdateEmbeds(messageID discord.Snowflake, embeds ...discord.Embed) (*Message, rest.Error) {
	return h.UpdateMessage(messageID, discord.MessageUpdate{Embeds: embeds})
}

func (h *webhookClientImpl) DeleteMessage(messageID discord.Snowflake) rest.Error {
	return h.WebhookService().DeleteMessage(h.id, h.token, messageID)
}

func (h *webhookClientImpl) ID() discord.Snowflake {
	return h.id
}
func (h *webhookClientImpl) Token() string {
	return h.token
}

func (h *webhookClientImpl) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, h.id, h.token)
	return compiledRoute.Route()
}
