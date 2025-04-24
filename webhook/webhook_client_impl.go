package webhook

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// NewWithURL creates a new Client by parsing the given webhookURL for the ID and token.
func NewWithURL(webhookURL string, opts ...ConfigOpt) (Client, error) {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return nil, err
	}

	parts := strings.FieldsFunc(u.Path, func(r rune) bool { return r == '/' })
	if len(parts) != 4 {
		return nil, ErrInvalidWebhookURL
	}

	token := parts[3]
	id, err := snowflake.Parse(parts[2])
	if err != nil {
		return nil, err
	}

	return New(id, token, opts...), nil
}

// New creates a new Client with the given ID, token and ConfigOpt(s).
func New(id snowflake.ID, token string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "webhook"))

	return &clientImpl{
		id:     id,
		token:  token,
		config: *config,
	}
}

type clientImpl struct {
	id     snowflake.ID
	token  string
	config Config
}

func (c *clientImpl) ID() snowflake.ID {
	return c.id
}

func (c *clientImpl) Token() string {
	return c.token
}

func (c *clientImpl) URL() string {
	return discord.WebhookURL(c.id, c.token)
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

func (c *clientImpl) GetMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest().GetWebhookMessage(c.id, c.token, messageID, opts...)
}

func (c *clientImpl) CreateMessage(messageCreate discord.WebhookMessageCreate, params rest.CreateWebhookMessageParams, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest().CreateWebhookMessage(c.id, c.token, messageCreate, params, opts...)
}

func (c *clientImpl) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(messageCreate, rest.CreateWebhookMessageParams{Wait: true, ThreadID: threadID}, opts...)
}

func (c *clientImpl) CreateContent(content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Content: content}, rest.CreateWebhookMessageParams{}, opts...)
}

func (c *clientImpl) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, rest.CreateWebhookMessageParams{}, opts...)
}

func (c *clientImpl) UpdateMessage(messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, params rest.UpdateWebhookMessageParams, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest().UpdateWebhookMessage(c.id, c.token, messageID, messageUpdate, params, opts...)
}

func (c *clientImpl) UpdateMessageInThread(messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, messageUpdate, rest.UpdateWebhookMessageParams{ThreadID: threadID}, opts...)
}

func (c *clientImpl) UpdateContent(messageID snowflake.ID, content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, rest.UpdateWebhookMessageParams{}, opts...)
}

func (c *clientImpl) UpdateEmbeds(messageID snowflake.ID, embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, rest.UpdateWebhookMessageParams{}, opts...)
}

func (c *clientImpl) DeleteMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return c.DeleteMessageInThread(messageID, 0, opts...)
}

func (c *clientImpl) DeleteMessageInThread(messageID snowflake.ID, threadID snowflake.ID, opts ...rest.RequestOpt) error {
	return c.Rest().DeleteWebhookMessage(c.id, c.token, messageID, threadID, opts...)
}
