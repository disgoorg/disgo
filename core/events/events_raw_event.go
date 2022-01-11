package events

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/json"

	"github.com/DisgoOrg/disgo/discord"
)

func HandleRawEvent(bot core.Bot, gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) io.Reader {
	if bot.EventManager().Config().RawEventsEnabled {
		var buf bytes.Buffer
		data, err := ioutil.ReadAll(io.TeeReader(reader, &buf))
		if err != nil {
			bot.Logger().Error("error reading raw payload from event")
		}
		bot.EventManager().Dispatch(&RawEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			Type:         gatewayEventType,
			RawPayload:   data,
		})

		return &buf
	}
	return reader
}

// RawEvent is called for any discord.GatewayEventType we receive if enabled in the bot.Config
type RawEvent struct {
	*GenericEvent
	Type       discord.GatewayEventType
	RawPayload json.RawMessage
}
