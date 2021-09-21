package utils

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// NewSelectMenuSubmitCollector gives you a channel to receive on and a function to close the collector
func NewSelectMenuSubmitCollector(disgo *core.Bot, channelID discord.Snowflake, guildID discord.Snowflake, messageID discord.Snowflake, filter SelectFilter) (chan *core.SelectMenuInteraction, func()) {
	ch := make(chan *core.SelectMenuInteraction)

	col := &SelectMenuSubmitCollector{
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

// NewSelectMenuSubmitCollectorFromMessage is an overload of NewSelectCollector that takes an api.Message for information
//goland:noinspection GoUnusedExportedFunction
func NewSelectMenuSubmitCollectorFromMessage(message *core.Message, filter SelectFilter) (chan *core.SelectMenuInteraction, func()) {
	return NewSelectMenuSubmitCollector(message.Bot, message.ChannelID, message.ID, *message.GuildID, filter)
}

// SelectFilter used to filter api.MessageReaction in a SelectCollector
type SelectFilter func(reaction *core.SelectMenuInteraction) bool

// SelectMenuSubmitCollector used to collect api.SelectMenuInteraction(s) from an api.Message using a SelectFilter function
type SelectMenuSubmitCollector struct {
	Channel   chan *core.SelectMenuInteraction
	Filter    SelectFilter
	Close     func()
	ChannelID discord.Snowflake
	GuildID   discord.Snowflake
	MessageID discord.Snowflake
}

// OnEvent used to get events for the SelectCollector
func (s *SelectMenuSubmitCollector) OnEvent(e interface{}) {
	if event, ok := e.(*core.SelectMenuSubmitEvent); ok {
		if !s.Filter(event.SelectMenuInteraction) {
			return
		}

		s.Channel <- event.SelectMenuInteraction
	} else if event, ok := e.(*core.GuildChannelDeleteEvent); ok && event.ChannelID == s.ChannelID {
		s.Close()
	} else if event, ok := e.(core.GuildLeaveEvent); ok && event.GuildID == s.GuildID {
		s.Close()
	} else if event, ok := e.(core.MessageDeleteEvent); ok && event.MessageID == s.MessageID {
		s.Close()
	}
}
