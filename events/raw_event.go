package events

import (
	"bytes"
	"io"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/json"
)

// HandleRawEvent handles raw events and dispatches the raw event.
func HandleRawEvent(client bot.Client, gatewayEventType gateway.EventType, sequenceNumber int, shardID int, respondFunc httpserver.RespondFunc, reader io.Reader) io.Reader {
	if client.EventManager().RawEventsEnabled() {
		var buf bytes.Buffer
		data, err := io.ReadAll(io.TeeReader(reader, &buf))
		if err != nil {
			client.Logger().Error("error reading raw payload from event")
		}
		client.EventManager().DispatchEvent(&Raw{
			GenericEvent: NewGenericEvent(client, sequenceNumber, shardID),
			Type:         gatewayEventType,
			RawPayload:   data,
			RespondFunc:  respondFunc,
		})

		return &buf
	}
	return reader
}

// Raw is called for any discord.Gateway Type we receive if enabled in the bot.Config
type Raw struct {
	*GenericEvent
	Type        gateway.EventType
	RawPayload  json.RawMessage
	RespondFunc httpserver.RespondFunc
}
