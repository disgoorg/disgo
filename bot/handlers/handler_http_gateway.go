package handlers

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/webhookevent"
)

func GetHTTPGatewayHandler() bot.HTTPGatewayEventHandler {
	return &httpGatewayHandler{}
}

type httpGatewayHandler struct{}

func (h *httpGatewayHandler) HandleHTTPGatewayEvent(client *bot.Client, message webhookevent.Message) {
	switch event := message.Event.(type) {
	case webhookevent.EventPing:
		handleHTTPGatewayPing(client, message)
	case webhookevent.Event:
		handleHTTPGatewayEvent(client, message, event)
	case webhookevent.EventUnknown:
		// nothing to do
	default:
		client.Logger.Warn("received unknown http gateway event", slog.Any("event", event))
	}
}

func handleHTTPGatewayEvent(client *bot.Client, message webhookevent.Message, event webhookevent.Event) {
	switch data := event.Data.(type) {
	case webhookevent.EventDataApplicationAuthorized:
		handleHTTPGatewayEventDataApplicationAuthorized(client, message, event, data)
	case webhookevent.EventDataApplicationDeauthorized:
		handleHTTPGatewayEventDataApplicationDeauthorized(client, message, event, data)
	case webhookevent.EventDataEntitlementCreate:
		handleHTTPGatewayEventDataEntitlementCreate(client, message, event, data)
	case webhookevent.EventDataQuestUserEnrollment:
		handleHTTPGatewayEventDataQuestUserEnrollment(client, message, event, data)
	case webhookevent.EventDataRaw:
		handleHTTPGatewayEventDataRaw(client, message, event, data)
	case webhookevent.EventDataUnknown:
		// nothing to do
	default:
		client.Logger.Warn("received unknown http gateway event data", slog.Any("data", data))
	}
}

func handleHTTPGatewayPing(client *bot.Client, message webhookevent.Message) {
	client.EventManager.DispatchEvent(&events.WebhookPing{
		Event:        events.NewEvent(client),
		WebhookEvent: events.NewWebhookEvent(message.Version, message.ApplicationID, time.Now()),
	})
}

func handleHTTPGatewayEventDataApplicationAuthorized(client *bot.Client, message webhookevent.Message, event webhookevent.Event, data webhookevent.EventDataApplicationAuthorized) {
	client.EventManager.DispatchEvent(&events.ApplicationAuthorized{
		Event:           events.NewEvent(client),
		WebhookEvent:    events.NewWebhookEvent(message.Version, message.ApplicationID, event.Timestamp),
		IntegrationType: data.IntegrationType,
		User:            data.User,
		Scopes:          data.Scopes,
		Guild:           data.Guild,
	})
}

func handleHTTPGatewayEventDataApplicationDeauthorized(client *bot.Client, message webhookevent.Message, event webhookevent.Event, data webhookevent.EventDataApplicationDeauthorized) {
	client.EventManager.DispatchEvent(&events.ApplicationDeauthorized{
		Event:        events.NewEvent(client),
		WebhookEvent: events.NewWebhookEvent(message.Version, message.ApplicationID, event.Timestamp),
		User:         data.User,
	})
}

func handleHTTPGatewayEventDataEntitlementCreate(client *bot.Client, message webhookevent.Message, event webhookevent.Event, data webhookevent.EventDataEntitlementCreate) {
	client.EventManager.DispatchEvent(&events.EntitlementCreate{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(-1, -1),
		WebhookEvent: events.NewWebhookEvent(message.Version, message.ApplicationID, event.Timestamp),
		Entitlement:  data.Entitlement,
	})
}

func handleHTTPGatewayEventDataRaw(client *bot.Client, message webhookevent.Message, event webhookevent.Event, data webhookevent.EventDataRaw) {
	client.EventManager.DispatchEvent(&events.WebhookRaw{
		Event:        events.NewEvent(client),
		WebhookEvent: events.NewWebhookEvent(message.Version, message.ApplicationID, event.Timestamp),
		EventDataRaw: data,
	})
}

func handleHTTPGatewayEventDataQuestUserEnrollment(client *bot.Client, message webhookevent.Message, event webhookevent.Event, data webhookevent.EventDataQuestUserEnrollment) {
	// This event cannot be received by apps at this time. It's documented because it appears on the Webhooks settings page.
}
