package utils

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// NewSelectCollector gives you a channel to receive on and a function to close the collector
func NewSelectCollector(disgo api.Disgo, channelID api.Snowflake, guildID api.Snowflake, messageID api.Snowflake, filter SelectFilter) (chan *api.SelectMenuInteraction, func()) {
	ch := make(chan *api.SelectMenuInteraction)

	col := &SelectCollector{
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

// NewSelectCollectorFromMessage is an overload of NewSelectCollector that takes an api.Message for information
//goland:noinspection GoUnusedExportedFunction
func NewSelectCollectorFromMessage(message *api.Message, filter SelectFilter) (chan *api.SelectMenuInteraction, func()) {
	return NewSelectCollector(message.Disgo, message.ChannelID, message.ID, *message.GuildID, filter)
}

// SelectFilter used to filter api.MessageReaction in a ReactionCollector
type SelectFilter func(reaction *api.SelectMenuInteraction) bool

// SelectCollector used to collect api.SelectMenuInteraction(s) from an api.Message using a SelectFilter function
type SelectCollector struct {
	Channel   chan *api.SelectMenuInteraction
	Filter    SelectFilter
	Close     func()
	ChannelID api.Snowflake
	GuildID   api.Snowflake
	MessageID api.Snowflake
}

// OnEvent used to get events for the SelectCollector
func (r *SelectCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.SelectMenuSubmitEvent); ok {
		if !r.Filter(event.SelectMenuInteraction) {
			return
		}

		r.Channel <- event.SelectMenuInteraction
	} else if event, ok := e.(*events.GuildChannelDeleteEvent); ok && event.ChannelID == r.ChannelID {
		r.Close()
	} else if event, ok := e.(events.GuildLeaveEvent); ok && event.GuildID == r.GuildID {
		r.Close()
	} else if event, ok := e.(events.MessageDeleteEvent); ok && event.MessageID == r.MessageID {
		r.Close()
	}
}
