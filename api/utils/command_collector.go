package utils

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// NewCommandCollectorFromMember is an overload of NewCommandCollector that takes an api.Member for information
//goland:noinspection GoUnusedExportedFunction
func NewCommandCollectorFromMember(member *api.Member, filter CommandFilter) (chan *api.CommandInteraction, func()) {
	return NewCommandCollector(member.Disgo, member.GuildID, member.User.ID, filter)
}

// NewCommandCollector gives you a channel to receive on and a function to close the collector
func NewCommandCollector(disgo api.Disgo, guildID api.Snowflake, memberID api.Snowflake, filter CommandFilter) (chan *api.CommandInteraction, func()) {
	ch := make(chan *api.CommandInteraction)

	col := &CommandCollector{
		Filter:   filter,
		Channel:  ch,
		GuildID:  guildID,
		MemberID: memberID,
	}

	cls := func() {
		close(ch)
		disgo.EventManager().RemoveEventListener(col)
	}

	col.Close = cls

	disgo.EventManager().AddEventListeners(col)

	return ch, cls
}

// CommandFilter used to filter api.CommandInteraction in a CommandCollector
type CommandFilter func(reaction *api.CommandInteraction) bool

// CommandCollector used to collect api.CommandInteraction(s) using a CommandFilter function
type CommandCollector struct {
	Channel  chan *api.CommandInteraction
	Filter   CommandFilter
	Close    func()
	GuildID  api.Snowflake
	MemberID api.Snowflake
}

// OnEvent used to get events for the CommandCollector
func (c *CommandCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.CommandEvent); ok {
		if !c.Filter(event.CommandInteraction) {
			return
		}

		c.Channel <- event.CommandInteraction
	} else if event, ok := e.(events.GuildLeaveEvent); ok && event.GuildID == c.GuildID {
		c.Close()
	}
}
