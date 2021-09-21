package utils

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// NewReactionAddCollector gives you a channel to receive on and a function to close the collector
func NewReactionAddCollector(disgo *core.Bot, channelID discord.Snowflake, guildID discord.Snowflake, messageID discord.Snowflake, filter ReactionFilter) (chan *core.MessageReactionAddEvent, func()) {
	ch := make(chan *core.MessageReactionAddEvent)

	col := &ReactionAddCollector{
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

// NewReactionAddCollectorFromMessage is an overload of NewReactionCollector that takes an api.Message for information
//goland:noinspection GoUnusedExportedFunction
func NewReactionAddCollectorFromMessage(message *core.Message, filter ReactionFilter) (chan *core.MessageReactionAddEvent, func()) {
	return NewReactionAddCollector(message.Bot, message.ChannelID, message.ID, *message.GuildID, filter)
}

// ReactionFilter used to filter api.MessageReaction in a ReactionCollector
type ReactionFilter func(reaction *core.MessageReactionAddEvent) bool

// ReactionAddCollector used to collect api.MessageReaction(s) from an api.Message using a ReactionFilter function
type ReactionAddCollector struct {
	Channel   chan *core.MessageReactionAddEvent
	Filter    ReactionFilter
	Close     func()
	ChannelID discord.Snowflake
	GuildID   discord.Snowflake
	MessageID discord.Snowflake
}

// OnEvent used to get events for the ReactionCollector
func (r *ReactionAddCollector) OnEvent(e interface{}) {
	if event, ok := e.(*core.MessageReactionAddEvent); ok {
		if !r.Filter(event) {
			return
		}

		r.Channel <- event
	} else if event, ok := e.(*core.GuildChannelDeleteEvent); ok && event.ChannelID == r.ChannelID {
		r.Close()
	} else if event, ok := e.(core.GuildLeaveEvent); ok && event.GuildID == r.GuildID {
		r.Close()
	} else if event, ok := e.(core.MessageDeleteEvent); ok && event.MessageID == r.MessageID {
		r.Close()
	}
}
