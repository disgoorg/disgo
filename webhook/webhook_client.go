package webhook

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// NewClient returns a new Client
//goland:noinspection GoUnusedExportedFunction
func NewClient(id discord.Snowflake, token string, opts ...ConfigOpt) *Client {
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

	webhookClient := &Client{
		ID:    id,
		Token: token,
	}

	if config.EntityBuilder == nil {
		config.EntityBuilder = NewEntityBuilder(webhookClient)
	}

	webhookClient.Config = *config
	return webhookClient
}

// Client is used to interact with the discord webhook api
type Client struct {
	ID    discord.Snowflake
	Token string
	Config
}

// GetWebhook fetches the current webhook from discord
func (h *Client) GetWebhook(opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.WebhookService.GetWebhookWithToken(h.ID, h.Token, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateWebhook(webhook), nil
}

// UpdateWebhook updates the current webhook
func (h *Client) UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.WebhookService.UpdateWebhookWithToken(h.ID, h.Token, webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateWebhook(webhook), nil
}

// DeleteWebhook deletes the current webhook
func (h *Client) DeleteWebhook(opts ...rest.RequestOpt) error {
	return h.WebhookService.DeleteWebhookWithToken(h.ID, h.Token, opts...)
}

// CreateMessageInThread creates a new Message in the provided thread
func (h *Client) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := h.WebhookService.CreateMessage(h.ID, h.Token, messageCreate, true, threadID, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateMessage(*message), nil
}

// CreateMessage creates a new message from the discord.WebhookMessageCreate
func (h *Client) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessageInThread(messageCreate, "", opts...)
}

// CreateContent creates a new message from the provided content
func (h *Client) CreateContent(content string, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Content: content}, opts...)
}

// CreateEmbeds creates a new message from the provided embeds
func (h *Client) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, opts...)
}

// UpdateMessage updates an already sent webhook message with the discord.WebhookMessageUpdate
func (h *Client) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessageInThread(messageID, messageUpdate, "", opts...)
}

// UpdateMessageInThread updates an already sent webhook message with the discord.WebhookMessageUpdate in a thread
func (h *Client) UpdateMessageInThread(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := h.WebhookService.UpdateMessage(h.ID, h.Token, messageID, messageUpdate, threadID, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateMessage(*message), nil
}

// UpdateContent updates an already sent webhook message with the content
func (h *Client) UpdateContent(messageID discord.Snowflake, content string, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, opts...)
}

// UpdateEmbeds updates an already sent webhook message with the embeds
func (h *Client) UpdateEmbeds(messageID discord.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, opts...)
}

// DeleteMessage deletes an already sent webhook message
func (h *Client) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return h.DeleteMessageInThread(messageID, "", opts...)
}

// DeleteMessageInThread deletes an already sent webhook message in a thread
func (h *Client) DeleteMessageInThread(messageID discord.Snowflake, threadID discord.Snowflake, opts ...rest.RequestOpt) error {
	return h.WebhookService.DeleteMessage(h.ID, h.Token, messageID, threadID, opts...)
}

// URL returns the full webhook URL
func (h *Client) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, h.ID, h.Token)
	return compiledRoute.URL()
}

// Close closes all connections the webhook client has open
func (h *Client) Close(ctx context.Context) error {
	return h.RestClient.Close(ctx)
}
