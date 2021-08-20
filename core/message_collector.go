package core

import (
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// NewMessageCollector gives you a channel to receive on and a function to close the collector
func NewMessageCollector(disgo Disgo, channelID discord.Snowflake, guildID *discord.Snowflake, filter MessageFilter) (chan *Message, func()) {
	ch := make(chan *Message)

	col := &MessageCollector{
		Filter:    filter,
		Channel:   ch,
		ChannelID: channelID,
		GuildID:   guildID,
	}

	cls := func() {
		close(ch)
		disgo.EventManager().RemoveEventListeners(col)
	}

	col.Close = cls

	disgo.EventManager().AddEventListeners(col)

	return ch, cls
}

// MessageFilter used to filter api.Message(s) in a MessageCollector
type MessageFilter func(msg *Message) bool

// MessageCollector collects api.Message(s) using a MessageFilter function
type MessageCollector struct {
	Filter  MessageFilter
	Channel chan *Message
	Close   func()
	ChannelID discord.Snowflake
	GuildID   *discord.Snowflake
}

// OnEvent used to get events for the MessageCollector
func (c *MessageCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.MessageCreateEvent); ok {
		if !c.Filter(event.Message) {
			return
		}

		c.Channel <- event.Message
	} else if event, ok := e.(*events.GuildChannelDeleteEvent); ok && event.ChannelID == c.ChannelID {
		c.Close()
	} else if event, ok := e.(events.GuildLeaveEvent); ok && c.GuildID != nil && event.GuildID == *c.GuildID {
		c.Close()
	}
}
