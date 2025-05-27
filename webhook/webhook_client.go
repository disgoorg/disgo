package webhook

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

var ErrInvalidWebhookURL = errors.New("invalid webhook URL")

// New creates a new Client with the given ID, Token and ConfigOpt(s).
func New(id snowflake.ID, token string, opts ...ConfigOpt) *Client {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &Client{
		ID:         id,
		Token:      token,
		Rest:       cfg.Webhooks,
		RestClient: cfg.RestClient,
		logger:     cfg.Logger,
	}
}

// NewWithURL creates a new Client by parsing the given webhookURL for the ID and Token.
func NewWithURL(webhookURL string, opts ...ConfigOpt) (*Client, error) {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
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

type Client struct {
	ID         snowflake.ID
	Token      string
	Rest       rest.Webhooks
	RestClient rest.Client
	logger     *slog.Logger
}

func (c *Client) URL() string {
	return discord.WebhookURL(c.ID, c.Token)
}

func (c *Client) Close(ctx context.Context) {
	c.RestClient.Close(ctx)
}

func (c *Client) GetWebhook(opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.Rest.GetWebhookWithToken(c.ID, c.Token, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

func (c *Client) UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.Rest.UpdateWebhookWithToken(c.ID, c.Token, webhookUpdate, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

func (c *Client) DeleteWebhook(opts ...rest.RequestOpt) error {
	return c.Rest.DeleteWebhookWithToken(c.ID, c.Token, opts...)
}

func (c *Client) GetMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest.GetWebhookMessage(c.ID, c.Token, messageID, opts...)
}

func (c *Client) CreateMessage(messageCreate discord.WebhookMessageCreate, params rest.CreateWebhookMessageParams, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest.CreateWebhookMessage(c.ID, c.Token, messageCreate, params, opts...)
}

func (c *Client) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(messageCreate, rest.CreateWebhookMessageParams{Wait: true, ThreadID: threadID}, opts...)
}

func (c *Client) CreateContent(content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Content: content}, rest.CreateWebhookMessageParams{}, opts...)
}

func (c *Client) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, rest.CreateWebhookMessageParams{}, opts...)
}

func (c *Client) UpdateMessage(messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, params rest.UpdateWebhookMessageParams, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest.UpdateWebhookMessage(c.ID, c.Token, messageID, messageUpdate, params, opts...)
}

func (c *Client) UpdateMessageInThread(messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.Rest.UpdateWebhookMessage(c.ID, c.Token, messageID, messageUpdate, rest.UpdateWebhookMessageParams{ThreadID: threadID}, opts...)
}

func (c *Client) UpdateContent(messageID snowflake.ID, content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, rest.UpdateWebhookMessageParams{}, opts...)
}

func (c *Client) UpdateEmbeds(messageID snowflake.ID, embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, rest.UpdateWebhookMessageParams{}, opts...)
}

func (c *Client) DeleteMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return c.DeleteMessageInThread(messageID, 0, opts...)
}

func (c *Client) DeleteMessageInThread(messageID snowflake.ID, threadID snowflake.ID, opts ...rest.RequestOpt) error {
	return c.Rest.DeleteWebhookMessage(c.ID, c.Token, messageID, threadID, opts...)
}
