package utils

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// NewReactionAddCollector gives you a channel to receive on and a function to close the collector
func NewReactionAddCollector(disgo api.Disgo, channelID api.Snowflake, guildID api.Snowflake, messageID api.Snowflake, filter ReactionFilter) (chan *api.MessageReaction, func()) {
	ch := make(chan *api.MessageReaction)

	col := &ReactionAddCollector{
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

// NewReactionAddCollectorFromMessage is an overload of NewReactionCollector that takes an api.Message for information
//goland:noinspection GoUnusedExportedFunction
func NewReactionAddCollectorFromMessage(message *api.Message, filter ReactionFilter) (chan *api.MessageReaction, func()) {
	return NewReactionAddCollector(message.Disgo, message.ChannelID, message.ID, *message.GuildID, filter)
}

// ReactionFilter used to filter api.MessageReaction in a ReactionCollector
type ReactionFilter func(reaction *events.MessageReactionAddEvent) bool

// ReactionAddCollector used to collect api.MessageReaction(s) from an api.Message using a ReactionFilter function
type ReactionAddCollector struct {
	Channel   chan *api.MessageReaction
	Filter    ReactionFilter
	Close     func()
	ChannelID api.Snowflake
	GuildID   api.Snowflake
	MessageID api.Snowflake
}

// OnEvent used to get events for the ReactionCollector
func (r *ReactionAddCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.MessageReactionAddEvent); ok {
		if !r.Filter(event) {
			return
		}

		r.Channel <- &event.MessageReaction
	} else if event, ok := e.(*events.GuildChannelDeleteEvent); ok && event.ChannelID == r.ChannelID {
		r.Close()
	} else if event, ok := e.(events.GuildLeaveEvent); ok && event.GuildID == r.GuildID {
		r.Close()
	} else if event, ok := e.(events.MessageDeleteEvent); ok && event.MessageID == r.MessageID {
		r.Close()
	}
}
