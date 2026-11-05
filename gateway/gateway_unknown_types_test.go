package gateway

import (
	"testing"

	"github.com/disgoorg/disgo/discord"
)

// A dispatch that fails to decode makes the gateway reconnect and replay it, so
// a channel type added by Discord must not fail the event.
func TestUnmarshalEventData_UnknownChannelType(t *testing.T) {
	channel := `{"id":"2","type":99,"guild_id":"1","name":"new-kind"}`

	tests := []struct {
		eventType EventType
		data      string
	}{
		{EventTypeGuildCreate, `{"id":"1","channels":[{"id":"3","type":0,"name":"general"},` + channel + `]}`},
		{EventTypeChannelCreate, channel},
		{EventTypeChannelUpdate, channel},
		{EventTypeChannelDelete, channel},
	}
	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			e, err := UnmarshalEventData([]byte(tt.data), tt.eventType)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var ch discord.GuildChannel
			switch e := e.(type) {
			case EventGuildCreate:
				if _, ok := e.Channels[0].(discord.GuildTextChannel); !ok {
					t.Errorf("expected GuildTextChannel, got %T", e.Channels[0])
				}
				ch = e.Channels[1]
			case EventChannelCreate:
				ch = e.GuildChannel
			case EventChannelUpdate:
				ch = e.GuildChannel
			case EventChannelDelete:
				ch = e.GuildChannel
			default:
				t.Fatalf("unexpected event %T", e)
			}
			if _, ok := ch.(discord.UnknownChannel); !ok {
				t.Errorf("expected UnknownChannel, got %T", ch)
			}
		})
	}
}
