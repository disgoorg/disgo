package events

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/json"

	"github.com/disgoorg/disgo/discord"
)

func HandleRawEvent(client bot.Client, gatewayEventType discord.GatewayEventType, sequenceNumber int, respondFunc httpserver.RespondFunc, reader io.Reader) io.Reader {
	if client.EventManager().RawEventsEnabled() {
		var buf bytes.Buffer
		data, err := ioutil.ReadAll(io.TeeReader(reader, &buf))
		if err != nil {
			client.Logger().Error("error reading raw payload from event")
		}
		client.EventManager().DispatchEvent(&RawEvent{
			GenericEvent: NewGenericEvent(client, sequenceNumber),
			Type:         gatewayEventType,
			RawPayload:   data,
			RespondFunc:  respondFunc,
		})

		return &buf
	}
	return reader
}

// RawEvent is called for any discord.GatewayEventType we receive if enabled in the bot.Config
type RawEvent struct {
	*GenericEvent
	Type        discord.GatewayEventType
	RawPayload  json.RawMessage
	RespondFunc httpserver.RespondFunc
}
