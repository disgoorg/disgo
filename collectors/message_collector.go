package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

func NewMessageCollectorByChannel(channel *core.Channel, filter core.MessageFilter) (<-chan *core.Message, func()) {
	var guildID *discord.Snowflake = nil
	if guildChannel := channel; channel.IsGuildChannel() {
		guildID = guildChannel.GuildID
	}
	return NewMessageCollector(channel.Bot, channel.ID, guildID, filter)
}

// NewMessageCollector gives you a channel to receive on and a function to close the collector
func NewMessageCollector(disgo *core.Bot, channelID discord.Snowflake, guildID *discord.Snowflake, filter core.MessageFilter) (<-chan *core.Message, func()) {
	ch := make(chan *core.Message)

	collector := &MessageCollector{
		Filter:    filter,
		Channel:   ch,
		ChannelID: channelID,
		GuildID:   guildID,
	}
	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(collector)
	}
	collector.Close = cls
	disgo.EventManager.AddEventListeners(collector)

	return ch, cls
}

// MessageCollector collects Message(s) using a MessageFilter function
type MessageCollector struct {
	Filter    core.MessageFilter
	Channel   chan<- *core.Message
	Close     func()
	ChannelID discord.Snowflake
	GuildID   *discord.Snowflake
}

// OnEvent used to get events for the MessageCollector
func (c *MessageCollector) OnEvent(e interface{}) {
	switch event := e.(type) {
	case *events.MessageCreateEvent:
		if !c.Filter(event.Message) {
			return
		}
		c.Channel <- event.Message

	case *events.GuildChannelDeleteEvent:
		if event.ChannelID == c.ChannelID {
			c.Close()
		}

	case *events.GuildLeaveEvent:
		if c.GuildID != nil && event.GuildID == *c.GuildID {
			c.Close()
		}
	}
}
