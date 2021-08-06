package utils

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// NewButtonCollector gives you a channel to receive on and a function to close the collector
func NewButtonClickCollector(disgo api.Disgo, channelID api.Snowflake, guildID api.Snowflake, messageID api.Snowflake, filter ButtonFilter) (chan *api.ButtonInteraction, func()) {
	ch := make(chan *api.ButtonInteraction)

	col := &ButtonCollector{
		Filter:    filter,
		Channel:   ch,
		ChannelID: channelID,
		GuildID:   guildID,
		MessageID: messageID,
	}

	cls := func() {
		close(ch)
		disgo.EventManager().RemoveEventListener(col)
	}

	col.Close = cls

	disgo.EventManager().AddEventListeners(col)

	return ch, cls
}

// NewButtonCollectorFromMessage is an overload of NewButtonCollector that takes an api.Message for information
//goland:noinspection GoUnusedExportedFunction
func NewButtonClickCollectorFromMessage(message *api.Message, filter ButtonFilter) (chan *api.ButtonInteraction, func()) {
	return NewButtonCollector(message.Disgo, message.ChannelID, message.ID, *message.GuildID, filter)
}

// ButtonFilter used to filter api.ButtonInteraction for ButtonCollector
type ButtonFilter func(reaction *api.ButtonInteraction) bool

// ButtonCollector used to collect api.ButtonInteraction(s) from an api.Message using a ButtonFilter function
type ButtonClickCollector struct {
	Channel   chan *api.ButtonInteraction
	Filter    ButtonFilter
	Close     func()
	ChannelID api.Snowflake
	GuildID   api.Snowflake
	MessageID api.Snowflake
}

// OnEvent used to get events for the ButtonCollector
func (r *ButtonCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.ButtonClickEvent); ok {
		if !r.Filter(event.ButtonInteraction) {
			return
		}

		r.Channel <- event.ButtonInteraction
	} else if event, ok := e.(*events.GuildChannelDeleteEvent); ok && event.ChannelID == r.ChannelID {
		r.Close()
	} else if event, ok := e.(*events.GuildLeaveEvent); ok && event.GuildID == r.GuildID {
		r.Close()
	} else if event, ok := e.(*events.MessageDeleteEvent); ok && event.MessageID == r.MessageID {
		r.Close()
	}
}
