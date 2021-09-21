package webhook

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type EntityBuilder interface {
	WebhookClient() *Client
	CreateMessage(message discord.Message) *Message
	CreateComponents(unmarshalComponents []discord.Component) []core.Component
	CreateWebhook(webhook discord.Webhook) *Webhook
}

func NewEntityBuilder(webhookClient *Client) EntityBuilder {
	return &entityBuilderImpl{
		webhookClient: webhookClient,
	}
}

type entityBuilderImpl struct {
	webhookClient *Client
}

func (b *entityBuilderImpl) WebhookClient() *Client {
	return b.webhookClient
}

func (b *entityBuilderImpl) CreateMessage(message discord.Message) *Message {
	webhookMessage := &Message{
		Message:       message,
		WebhookClient: b.WebhookClient(),
	}
	if len(message.Components) > 0 {
		webhookMessage.Components = b.CreateComponents(message.Components)
	}
	return webhookMessage
}

func (b *entityBuilderImpl) CreateComponents(unmarshalComponents []discord.Component) []core.Component {
	components := make([]core.Component, len(unmarshalComponents))
	for i, component := range unmarshalComponents {
		switch component.Type {
		case discord.ComponentTypeActionRow:
			actionRow := core.ActionRow{
				Component: component,
			}
			if len(component.Components) > 0 {
				actionRow.Components = b.CreateComponents(component.Components)
			}
			components[i] = actionRow

		case discord.ComponentTypeButton:
			components[i] = core.Button{
				Component: component,
			}

		case discord.ComponentTypeSelectMenu:
			components[i] = core.SelectMenu{
				Component: component,
			}
		}
	}
	return components
}

func (b *entityBuilderImpl) CreateWebhook(webhook discord.Webhook) *Webhook {
	return &Webhook{
		Webhook:       webhook,
		WebhookClient: b.WebhookClient(),
	}
}
