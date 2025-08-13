package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/httpgateway"
)

func GetHTTPGatewayHandler() bot.HTTPGatewayEventHandler {
	return &httpGatewayHandler{}
}

type httpGatewayHandler struct{}

func (h *httpGatewayHandler) HandleHTTPGatewayEvent(client *bot.Client, ack func(), message httpgateway.Message) {
	switch event := message.Event.(type) {
	case httpgateway.EventPing:
		handleHTTPGatewayPing(client, ack, message)
	case httpgateway.Event:
		handleHTTPGatewayEvent(client, ack, message, event)
	case httpgateway.EventUnknown:
		handleHTTPGatewayEventUnknown(client, ack, message, event)
	default:
		client.Logger.Warn("received unknown http gateway event", slog.Any("event", event))
	}
}

func handleHTTPGatewayEvent(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.Event) {
	switch data := event.Data.(type) {
	case httpgateway.EventDataApplicationAuthorized:
		handleHTTPGatewayEventDataApplicationAuthorized(client, ack, message, event, data)
	case httpgateway.EventDataApplicationDeauthorized:
		handleHTTPGatewayEventDataApplicationDeauthorized(client, ack, message, event, data)
	case httpgateway.EventDataEntitlementCreate:
		handleHTTPGatewayEventDataEntitlementCreate(client, ack, message, event, data)
	case httpgateway.EventDataQuestUserEnrollment:
		handleHTTPGatewayEventDataQuestUserEnrollment(client, ack, message, event, data)
	case httpgateway.EventDataUnknown:
		handleHTTPGatewayEventDataUnknown(client, ack, message, event, data)
	default:
		client.Logger.Warn("received unknown http gateway event data", slog.Any("data", data))
	}
}

func handleHTTPGatewayPing(client *bot.Client, ack func(), message httpgateway.Message) {
	ack()
	client.Logger.Debug("received http gateway ping", slog.String("application_id", message.ApplicationID.String()), slog.Int("version", message.Version))
}

func handleHTTPGatewayEventDataApplicationAuthorized(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.Event, data httpgateway.EventDataApplicationAuthorized) {
}

func handleHTTPGatewayEventDataApplicationDeauthorized(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.Event, data httpgateway.EventDataApplicationDeauthorized) {

}

func handleHTTPGatewayEventDataEntitlementCreate(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.Event, data httpgateway.EventDataEntitlementCreate) {

}

func handleHTTPGatewayEventDataQuestUserEnrollment(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.Event, data httpgateway.EventDataQuestUserEnrollment) {

}

func handleHTTPGatewayEventDataUnknown(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.Event, data httpgateway.EventDataUnknown) {

}

func handleHTTPGatewayEventUnknown(client *bot.Client, ack func(), message httpgateway.Message, event httpgateway.EventUnknown) {

}
