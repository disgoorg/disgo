package webhook

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// DefaultAllowedMentions gives you the default AllowedMentions for a Message
var DefaultAllowedMentions = discord.AllowedMentions{
	Parse:       []discord.AllowedMentionType{discord.AllowedMentionTypeUsers, discord.AllowedMentionTypeRoles, discord.AllowedMentionTypeEveryone},
	Roles:       []discord.Snowflake{},
	Users:       []discord.Snowflake{},
	RepliedUser: true,
}

// NewClient returns a new Client
//goland:noinspection GoUnusedExportedFunction
func NewClient(id discord.Snowflake, token string, opts ...ConfigOpt) *Client {
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

type Client struct {
	ID    discord.Snowflake
	Token string
	Config
}

func (h *Client) GetWebhook(opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.WebhookService.GetWebhookWithToken(h.ID, h.Token, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateWebhook(*webhook), nil
}

func (h *Client) UpdateWebhook(webhookUpdate discord.WebhookUpdate, opts ...rest.RequestOpt) (*Webhook, error) {
	webhook, err := h.WebhookService.UpdateWebhookWithToken(h.ID, h.Token, webhookUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateWebhook(*webhook), nil
}

func (h *Client) DeleteWebhook(opts ...rest.RequestOpt) error {
	return h.WebhookService.DeleteWebhookWithToken(h.ID, h.Token, opts...)
}

func (h *Client) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID discord.Snowflake, opts ...rest.RequestOpt) (*Message, error) {
	message, err := h.WebhookService.CreateMessage(h.ID, h.Token, messageCreate, true, threadID, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateMessage(*message), nil
}

func (h *Client) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessageInThread(messageCreate, "", opts...)
}

func (h *Client) CreateContent(content string, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Content: content}, opts...)
}

func (h *Client) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error) {
	return h.CreateMessage(discord.WebhookMessageCreate{Embeds: embeds}, opts...)
}

func (h *Client) UpdateMessage(messageID discord.Snowflake, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*Message, error) {
	message, err := h.WebhookService.UpdateMessage(h.ID, h.Token, messageID, messageUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return h.EntityBuilder.CreateMessage(*message), nil
}

func (h *Client) UpdateContent(messageID discord.Snowflake, content string, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Content: &content}, opts...)
}

func (h *Client) UpdateEmbeds(messageID discord.Snowflake, embeds []discord.Embed, opts ...rest.RequestOpt) (*Message, error) {
	return h.UpdateMessage(messageID, discord.WebhookMessageUpdate{Embeds: &embeds}, opts...)
}

func (h *Client) DeleteMessage(messageID discord.Snowflake, opts ...rest.RequestOpt) error {
	return h.WebhookService.DeleteMessage(h.ID, h.Token, messageID, opts...)
}

func (h *Client) URL() string {
	compiledRoute, _ := route.GetWebhook.Compile(nil, h.ID, h.Token)
	return compiledRoute.URL()
}
