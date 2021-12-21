package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

// NewMessageReactionRemoveCollector gives you a channel to receive on and a function to close the collector
func NewMessageReactionRemoveCollector(disgo *core.Bot, filter core.MessageReactionRemoveFilter) (<-chan *core.MessageReactionRemove, func()) {
	ch := make(chan *core.MessageReactionRemove)

	col := &MessageReactionRemoveCollector{
		Filter: filter,
		Chan:   ch,
	}
	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(col)
	}
	col.Close = cls
	disgo.EventManager.RemoveEventListeners(col)

	return ch, cls
}

// MessageReactionRemoveCollector used to collect discord.MessageReaction(s) from a core.Message using an events.MessageReactionRemoveFilter function
type MessageReactionRemoveCollector struct {
	Filter core.MessageReactionRemoveFilter
	Chan   chan *core.MessageReactionRemove
	Close  func()
}

// OnEvent used to get events for the ReactionCollector
func (c *MessageReactionRemoveCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.MessageReactionRemoveEvent); ok {
		messageReactionRemove := &core.MessageReactionRemove{
			UserID:    event.UserID,
			ChannelID: event.ChannelID,
			MessageID: event.MessageID,
			GuildID:   event.GuildID,
			Emoji:     event.Emoji,
		}
		if !c.Filter(messageReactionRemove) {
			return
		}
		c.Chan <- messageReactionRemove
	}
}
