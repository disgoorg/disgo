package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewMessageCollectorByChannel(channel *Channel, filter MessageFilter) (chan *Message, func()) {
	var guildID *discord.Snowflake = nil
	if guildChannel := channel; channel.IsGuildChannel() {
		guildID = guildChannel.GuildID
	}
	return NewMessageCollector(channel.Bot, channel.ID, guildID, filter)
}

// NewMessageCollector gives you a channel to receive on and a function to close the collector
func NewMessageCollector(disgo *Bot, channelID discord.Snowflake, guildID *discord.Snowflake, filter MessageFilter) (chan *Message, func()) {
	ch := make(chan *Message)

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

// MessageFilter used to filter core.Message(s) in a MessageCollector
type MessageFilter func(message *Message) bool

// MessageCollector collects core.Message(s) using a MessageFilter function
type MessageCollector struct {
	Filter    MessageFilter
	Channel   chan *Message
	Close     func()
	ChannelID discord.Snowflake
	GuildID   *discord.Snowflake
}

// OnEvent used to get events for the MessageCollector
func (c *MessageCollector) OnEvent(e interface{}) {
	switch event := e.(type) {
	case *MessageCreateEvent:
		if !c.Filter(event.Message) {
			return
		}
		c.Channel <- event.Message

	case *GuildChannelDeleteEvent:
		if event.ChannelID == c.ChannelID {
			c.Close()
		}

	case *GuildLeaveEvent:
		if c.GuildID != nil && event.GuildID == *c.GuildID {
			c.Close()
		}
	}
}
