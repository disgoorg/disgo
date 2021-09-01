package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// New returns a new Client
//goland:noinspection GoUnusedExportedFunction
func New(id discord.Snowflake, token string, opts ...ConfigOpt) Client {
	config := &DefaultConfig
	config.Apply(opts)

	if config.Logger == nil {
		config.Logger = log.Default()
	}

	if config.RestClient == nil {
		config.RestClient = rest.NewClient(config.RestClientConfig)
	}
	if config.WebhookService == nil {
		config.WebhookService = rest.NewWebhookService(config.RestClient)
	}
	if config.DefaultAllowedMentions == nil {
		config.DefaultAllowedMentions = &DefaultAllowedMentions
	}

	webhookClient := &webhookClientImpl{
		id: id,
		token: token,
	}

	if config.EntityBuilder == nil {
		config.EntityBuilder = NewEntityBuilder(webhookClient)
	}

	webhookClient.config = *config
	return webhookClient
}

type webhookClientImpl struct {
	id discord.Snowflake
	token string
	config Config
}

func (h *webhookClientImpl) Logger() log.Logger {
	return h.config.Logger
}

func (h *webhookClientImpl) RestClient() rest.Client {
	return h.config.WebhookService.RestClient()
}
func (h *webhookClientImpl) WebhookService() rest.WebhookService {
	return h.config.WebhookService
}
func (h *webhookClientImpl) EntityBuilder() EntityBuilder {
	return h.config.EntityBuilder
}
func (h *webhookClientImpl) DefaultAllowedMentions() *discord.AllowedMentions {
	return h.config.DefaultAllowedMentions
}

func (h *webhookClientImpl) GetWebhook(opts ...rest.RequestOpt) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().GetWebhookWithToken(h.ID(), h.Token(), opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) UpdateWebhook(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, rest.Error) {
	webhook, err := h.WebhookService().UpdateWebhookWithToken(h.ID(), h.Token(), webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder().CreateWebhook(*webhook), nil
}

func (h *webhookClientImpl) DeleteWebhook(opts ...rest.RequestOpt) rest.Error {
	return h.WebhookService().DeleteWebhookWithToken(h.ID(), h.Token(), opts...)
}

func (h *webhookClientImpl) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, rest.Error) {
	message, err := h.WebhookService().CreateMessage(h.ID(), h.Token(), messageCreate, true, threadID, opts...)
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
	message, err := h.WebhookService().UpdateMessage(h.ID(), h.Token(), messageID, messageUpdate, opts...)
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
	return h.WebhookService().DeleteMessage(h.ID(), h.Token(), messageID, opts...)
}

func (h *webhookClientImpl) ID() discord.Snowflake {
	return h.id
}
func (h *webhookClientImpl) Token() string {
	return h.token
}

func (h *webhookClientImpl) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, h.ID(), h.Token())
	return compiledRoute.URL()
}
