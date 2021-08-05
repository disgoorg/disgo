package util

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// NewMessageCollectorFromChannel is an overload of NewMessageCollector that takes an api.MessageChannel for information
//goland:noinspection GoUnusedExportedFunction
func NewMessageCollectorFromChannel(channel *api.MessageChannel, filter MessageFilter) (chan *api.Message, func()) {
	return NewMessageCollector(channel.Disgo, channel.ID, *channel.GuildID, filter)
}

// NewMessageCollector gives you a channel to receive on and a function to close the collector
func NewMessageCollector(disgo api.Disgo, channelID api.Snowflake, guildID api.Snowflake, filter MessageFilter) (chan *api.Message, func()) {
	ch := make(chan *api.Message)

	col := &MessageCollector{
		Filter:    filter,
		Channel:   ch,
		ChannelID: channelID,
		GuildID:   guildID,
	}

	cls := func() {
		close(ch)
		disgo.EventManager().RemoveEventListener(col)
	}

	col.Close = cls

	disgo.EventManager().AddEventListeners(col)

	return ch, cls
}

// MessageFilter used to filter api.Message(s) in a MessageCollector
type MessageFilter func(msg *api.Message) bool

// MessageCollector collects api.Message(s) using a MessageFilter function
type MessageCollector struct {
	Filter    MessageFilter
	Channel   chan *api.Message
	Close     func()
	ChannelID api.Snowflake
	GuildID   api.Snowflake
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
	} else if event, ok := e.(events.GuildVoiceLeaveEvent); ok && event.GuildID == c.GuildID {
		c.Close()
	}
}
