package webhook

import (
	"context"

	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/rest/route"
)

// NewClient returns a new Client
//goland:noinspection GoUnusedExportedFunction
func NewClient(id snowflake.Snowflake, token string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &ClientImpl{
		id:     id,
		token:  token,
		config: *config,
	}
}

// ClientImpl is used to interact with the discord webhook api
type ClientImpl struct {
	id     snowflake.Snowflake
	token  string
	config Config
}

func (c *ClientImpl) ID() snowflake.Snowflake {
	return c.id
}

func (c *ClientImpl) Token() string {
	return c.token
}

// URL returns the full webhook URL
func (c *ClientImpl) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, c.ID, c.Token)
	return compiledRoute.URL()
}

// Close closes all connections the webhook client has open
func (c *ClientImpl) Close(ctx context.Context) {
	c.config.RestClient.Close(ctx)
}

func (c *ClientImpl) WebhookService() rest.WebhookService {
	return c.config.WebhookService
}

// GetWebhook fetches the current webhook from discord
func (c *ClientImpl) GetWebhook(opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.WebhookService().GetWebhookWithToken(c.id, c.token, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

// UpdateWebhook updates the current webhook
func (c *ClientImpl) UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.WebhookService().UpdateWebhookWithToken(c.id, c.token, webhookUpdate, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

// DeleteWebhook deletes the current webhook
func (c *ClientImpl) DeleteWebhook(opts ...rest.RequestOpt) error {
	return c.WebhookService().DeleteWebhookWithToken(c.id, c.token, opts...)
}

// CreateMessageInThread creates a new Message in the provided thread
func (c *ClientImpl) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID snowflake.Snowflake, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.WebhookService().CreateMessage(c.id, c.token, messageCreate, true, threadID, opts...)
}

// CreateMessage creates a new message from the discord.WebhookMessageCreate
func (c *ClientImpl) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessageInThread(messageCreate, "", opts...)
}

// CreateContent creates a new message from the provided content
func (c *ClientImpl) CreateContent(content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Content: content}, opts...)
}

// CreateEmbeds creates a new message from the provided embeds
func (c *ClientImpl) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, opts...)
}

// UpdateMessage updates an already sent webhook message with the discord.WebhookMessageUpdate
func (c *ClientImpl) UpdateMessage(messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessageInThread(messageID, messageUpdate, "", opts...)
}

// UpdateMessageInThread updates an already sent webhook message with the discord.WebhookMessageUpdate in a thread
func (c *ClientImpl) UpdateMessageInThread(messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.Snowflake, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.WebhookService().UpdateMessage(c.id, c.token, messageID, messageUpdate, threadID, opts...)
}

// UpdateContent updates an already sent webhook message with the content
func (c *ClientImpl) UpdateContent(messageID snowflake.Snowflake, content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, opts...)
}

// UpdateEmbeds updates an already sent webhook message with the embeds
func (c *ClientImpl) UpdateEmbeds(messageID snowflake.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, opts...)
}

// DeleteMessage deletes an already sent webhook message
func (c *ClientImpl) DeleteMessage(messageID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return c.DeleteMessageInThread(messageID, "", opts...)
}

// DeleteMessageInThread deletes an already sent webhook message in a thread
func (c *ClientImpl) DeleteMessageInThread(messageID snowflake.Snowflake, threadID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return c.WebhookService().DeleteMessage(c.id, c.token, messageID, threadID, opts...)
}
