package utils

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// NewButtonClickCollector gives you a channel to receive on and a function to close the collector
func NewButtonClickCollector(disgo *core.Bot, channelID discord.Snowflake, guildID discord.Snowflake, messageID discord.Snowflake, filter ButtonFilter) (chan *core.ButtonInteraction, func()) {
	ch := make(chan *core.ButtonInteraction)

	col := &ButtonClickCollector{
		Filter:    filter,
		Channel:   ch,
		ChannelID: channelID,
		GuildID:   guildID,
		MessageID: messageID,
	}

	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(col)
	}

	col.Close = cls

	disgo.EventManager.AddEventListeners(col)

	return ch, cls
}

// NewButtonClickCollectorFromMessage is an overload of NewButtonCollector that takes an api.Message for information
//goland:noinspection GoUnusedExportedFunction
func NewButtonClickCollectorFromMessage(message *core.Message, filter ButtonFilter) (chan *core.ButtonInteraction, func()) {
	return NewButtonClickCollector(message.Bot, message.ChannelID, message.ID, *message.GuildID, filter)
}

// ButtonFilter used to filter api.ButtonInteraction for ButtonCollector
type ButtonFilter func(reaction *core.ButtonInteraction) bool

// ButtonClickCollector used to collect api.ButtonInteraction(s) from an api.Message using a ButtonFilter function
type ButtonClickCollector struct {
	Channel   chan *core.ButtonInteraction
	Filter    ButtonFilter
	Close     func()
	ChannelID discord.Snowflake
	GuildID   discord.Snowflake
	MessageID discord.Snowflake
}

// OnEvent used to get events for the ButtonCollector
func (b *ButtonClickCollector) OnEvent(e interface{}) {
	if event, ok := e.(*core.ButtonClickEvent); ok {
		if !b.Filter(event.ButtonInteraction) {
			return
		}

		b.Channel <- event.ButtonInteraction
	} else if event, ok := e.(*core.GuildChannelDeleteEvent); ok && event.ChannelID == b.ChannelID {
		b.Close()
	} else if event, ok := e.(core.GuildLeaveEvent); ok && event.GuildID == b.GuildID {
		b.Close()
	} else if event, ok := e.(core.MessageDeleteEvent); ok && event.MessageID == b.MessageID {
		b.Close()
	}
}
