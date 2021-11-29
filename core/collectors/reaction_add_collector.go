package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

// NewMessageReactionAddCollector gives you a channel to receive on and a function to close the collector
func NewMessageReactionAddCollector(disgo *core.Bot, filter core.MessageReactionAddFilter) (<-chan *core.MessageReactionAdd, func()) {
	ch := make(chan *core.MessageReactionAdd)

	col := &MessageReactionAddCollector{
		Filter: filter,
		Chan:   ch,
	}
	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(col)
	}
	col.Close = cls
	disgo.EventManager.AddEventListeners(col)

	return ch, cls
}

// MessageReactionAddCollector used to collect discord.MessageReaction(s) from a core.Message using an events.MessageReactionAddFilter function
type MessageReactionAddCollector struct {
	Filter core.MessageReactionAddFilter
	Chan   chan *core.MessageReactionAdd
	Close  func()
}

// OnEvent used to get events for the ReactionCollector
func (c *MessageReactionAddCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.MessageReactionAddEvent); ok {
		messageReactionAdd := &core.MessageReactionAdd{
			UserID:    event.UserID,
			ChannelID: event.ChannelID,
			MessageID: event.MessageID,
			GuildID:   event.GuildID,
			Member:    event.Member,
			Emoji:     event.Emoji,
		}
		if !c.Filter(messageReactionAdd) {
			return
		}
		c.Chan <- messageReactionAdd
	}
}
