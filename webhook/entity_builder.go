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
	return &EntityBuilderImpl{
		webhookClient: webhookClient,
	}
}

type EntityBuilderImpl struct {
	webhookClient *Client
}

func (b *EntityBuilderImpl) WebhookClient() *Client {
	return b.webhookClient
}

func (b *EntityBuilderImpl) CreateMessage(message discord.Message) *Message {
	webhookMessage := &Message{
		Message:       message,
		WebhookClient: b.WebhookClient(),
	}
	if len(message.Components) > 0 {
		webhookMessage.Components = b.CreateComponents(message.Components)
	}
	return nil
}

func (b *EntityBuilderImpl) CreateComponents(unmarshalComponents []discord.Component) []core.Component {
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

func (b *EntityBuilderImpl) CreateWebhook(webhook discord.Webhook) *Webhook {
	return &Webhook{
		Webhook:       webhook,
		WebhookClient: b.WebhookClient(),
	}
}
