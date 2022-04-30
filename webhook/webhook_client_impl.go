package webhook

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

// NewClient returns a new Client
func NewClient(id snowflake.Snowflake, token string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &clientImpl{
		id:     id,
		token:  token,
		config: *config,
	}
}

type clientImpl struct {
	id     snowflake.Snowflake
	token  string
	config Config
}

func (c *clientImpl) ID() snowflake.Snowflake {
	return c.id
}

func (c *clientImpl) Token() string {
	return c.token
}

func (c *clientImpl) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, c.ID, c.Token)
	return compiledRoute.URL()
}

func (c *clientImpl) Close(ctx context.Context) {
	c.config.RestClient.Close(ctx)
}

func (c *clientImpl) Rest() rest.Webhooks {
	return c.config.Webhooks
}

func (c *clientImpl) GetWebhook(opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.Rest().GetWebhookWithToken(c.id, c.token, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

func (c *clientImpl) UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.Rest().UpdateWebhookWithToken(c.id, c.token, webhookUpdate, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

func (c *clientImpl) DeleteWebhook(opts ...rest.RequestOpt) error {
	return c.Rest().DeleteWebhookWithToken(c.id, c.token, opts...)
}

func (c *clientImpl) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID snowflake.Snowflake, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest().CreateWebhookMessage(c.id, c.token, messageCreate, true, threadID, opts...)
}

func (c *clientImpl) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessageInThread(messageCreate, "", opts...)
}

func (c *clientImpl) CreateContent(content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Content: content}, opts...)
}

func (c *clientImpl) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, opts...)
}

func (c *clientImpl) UpdateMessage(messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessageInThread(messageID, messageUpdate, "", opts...)
}

func (c *clientImpl) UpdateMessageInThread(messageID snowflake.Snowflake, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.Snowflake, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest().UpdateWebhookMessage(c.id, c.token, messageID, messageUpdate, threadID, opts...)
}

func (c *clientImpl) UpdateContent(messageID snowflake.Snowflake, content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, opts...)
}

func (c *clientImpl) UpdateEmbeds(messageID snowflake.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, opts...)
}

func (c *clientImpl) DeleteMessage(messageID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return c.DeleteMessageInThread(messageID, "", opts...)
}

func (c *clientImpl) DeleteMessageInThread(messageID snowflake.Snowflake, threadID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return c.Rest().DeleteWebhookMessage(c.id, c.token, messageID, threadID, opts...)
}
