package webhook

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

type Client interface {
	// ID returns the configured Webhook id
	ID() discord.Snowflake
	// Token returns the configured Webhook token
	Token() string
	// URL returns the full Webhook URL
	URL() string
	// Config returns the configured Config
	Config() Config
	// Close closes all connections the Webhook Client has open
	Close(ctx context.Context) error

	// GetWebhook fetches the current Webhook from discord
	GetWebhook(opts ...rest.RequestOpt) (*Webhook, error)
	// UpdateWebhook updates the current Webhook
	UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*Webhook, error)
	// DeleteWebhook deletes the current Webhook
	DeleteWebhook(opts ...rest.RequestOpt) error

	// CreateMessage creates a new Message from the discord.WebhookMessageCreate
	CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*Message, error)
	// CreateMessageInThread creates a new Message from the discord.WebhookMessageCreate in the provided thread
	CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error)
	// CreateContent creates a new Message from the provided content
	CreateContent(content string, opts ...rest.RequestOpt) (*Message, error)
	// CreateEmbeds creates a new Message from the provided discord.Embed(s)
	CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error)

	// UpdateMessage updates an already sent Webhook Message with the discord.WebhookMessageUpdate
	UpdateMessage(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, error)
	// UpdateMessageInThread updates an already sent Webhook Message with the discord.WebhookMessageUpdate in the provided thread
	UpdateMessageInThread(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error)
	// UpdateContent updates an already sent Webhook Message with the content
	UpdateContent(messageID discord.Snowflake, content string, opts ...rest.RequestOpt) (*Message, error)
	// UpdateEmbeds updates an already sent Webhook Message with the discord.Embed(s)
	UpdateEmbeds(messageID discord.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error)

	// DeleteMessage deletes an already sent Webhook Message
	DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error
	// DeleteMessageInThread deletes an already sent Webhook Message in the provided thread
	DeleteMessageInThread(messageID discord.Snowflake, threadID discord.Snowflake, opts ...rest.RequestOpt) error
}

// NewClient returns a new Client
//goland:noinspection GoUnusedExportedFunction
func NewClient(id discord.Snowflake, token string, opts ...ConfigOpt) Client {
	config := &DefaultConfig
	config.Apply(opts)

	if config.Logger == nil {
		config.Logger = log.Default()
	}

	if config.WebhookService == nil {
		if config.RestClient == nil {
			config.RestClient = rest.NewClient(config.RestClientConfig)
		}
		config.WebhookService = rest.NewWebhookService(config.RestClient)
	}
	if config.DefaultAllowedMentions == nil {
		config.DefaultAllowedMentions = &discord.DefaultAllowedMentions
	}

	webhookClient := &clientImpl{
		id:    id,
		token: token,
	}

	if config.EntityBuilder == nil {
		config.EntityBuilder = NewEntityBuilder(webhookClient)
	}

	webhookClient.config = *config
	return webhookClient
}

// clientImpl is used to interact with the discord webhook api
type clientImpl struct {
	id     discord.Snowflake
	token  string
	config Config
}

// ID returns the configured Webhook id
func (h *clientImpl) ID() discord.Snowflake {
	return h.id
}

// Token returns the configured Webhook token
func (h *clientImpl) Token() string {
	return h.token
}

// URL returns the full Webhook URL
func (h *clientImpl) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, h.id, h.token)
	return compiledRoute.URL()
}

// Config returns the configured Config
func (h *clientImpl) Config() Config {
	return h.config
}

// Close closes all connections the Webhook Client has open
func (h *clientImpl) Close(ctx context.Context) error {
	return h.config.RestClient.Close(ctx)
}

// GetWebhook fetches the current Webhook from discord
func (h *clientImpl) GetWebhook(opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.config.WebhookService.GetWebhookWithToken(h.id, h.token, opts...)
	if err != nil {
		return nil, err
	}
	return h.config.EntityBuilder.CreateWebhook(webhook), nil
}

// UpdateWebhook updates the current Webhook
func (h *clientImpl) UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.config.WebhookService.UpdateWebhookWithToken(h.id, h.token, webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.config.EntityBuilder.CreateWebhook(webhook), nil
}

// DeleteWebhook deletes the current Webhook
func (h *clientImpl) DeleteWebhook(opts ...rest.RequestOpt) error {
	return h.config.WebhookService.DeleteWebhookWithToken(h.id, h.token, opts...)
}

// CreateMessage creates a new Message from the discord.WebhookMessageCreate
func (h *clientImpl) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessageInThread(messageCreate, "", opts...)
}

// CreateMessageInThread creates a new Message from the discord.WebhookMessageCreate in the provided thread
func (h *clientImpl) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := h.config.WebhookService.CreateMessage(h.id, h.token, messageCreate, true, threadID, opts...)
	if err != nil {
		return nil, err
	}
	return h.config.EntityBuilder.CreateMessage(*message), nil
}

// CreateContent creates a new Message from the provided content
func (h *clientImpl) CreateContent(content string, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Content: content}, opts...)
}

// CreateEmbeds creates a new Message from the provided discord.Embed(s)
func (h *clientImpl) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, opts...)
}

// UpdateMessage updates an already sent Webhook Message with the discord.WebhookMessageUpdate
func (h *clientImpl) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessageInThread(messageID, messageUpdate, "", opts...)
}

// UpdateMessageInThread updates an already sent Webhook Message with the discord.WebhookMessageUpdate in the provided thread
func (h *clientImpl) UpdateMessageInThread(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := h.config.WebhookService.UpdateMessage(h.id, h.token, messageID, messageUpdate, threadID, opts...)
	if err != nil {
		return nil, err
	}
	return h.config.EntityBuilder.CreateMessage(*message), nil
}

// UpdateContent updates an already sent Webhook Message with the content
func (h *clientImpl) UpdateContent(messageID discord.Snowflake, content string, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, opts...)
}

// UpdateEmbeds updates an already sent Webhook Message with the discord.Embed(s)
func (h *clientImpl) UpdateEmbeds(messageID discord.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, opts...)
}

// DeleteMessage deletes an already sent Webhook Message
func (h *clientImpl) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return h.DeleteMessageInThread(messageID, "", opts...)
}

// DeleteMessageInThread deletes an already sent Webhook Message in the provided thread
func (h *clientImpl) DeleteMessageInThread(messageID discord.Snowflake, threadID discord.Snowflake, opts ...rest.RequestOpt) error {
	return h.config.WebhookService.DeleteMessage(h.id, h.token, messageID, threadID, opts...)
}
