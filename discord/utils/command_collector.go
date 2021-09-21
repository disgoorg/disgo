package utils

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// NewCommandCollectorFromMember is an overload of NewCommandCollector that takes an api.Member for information
//goland:noinspection GoUnusedExportedFunction
func NewCommandCollectorFromMember(member *core.Member, filter CommandFilter) (chan *core.ApplicationCommandInteraction, func()) {
	return NewCommandCollector(member.Bot, member.GuildID, member.User.ID, filter)
}

// NewCommandCollector gives you a channel to receive on and a function to close the collector
func NewCommandCollector(disgo *core.Bot, guildID discord.Snowflake, memberID discord.Snowflake, filter CommandFilter) (chan *core.ApplicationCommandInteraction, func()) {
	ch := make(chan *core.ApplicationCommandInteraction)

	col := &CommandCollector{
		Filter:   filter,
		Channel:  ch,
		GuildID:  guildID,
		MemberID: memberID,
	}

	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(col)
	}

	col.Close = cls

	disgo.EventManager.AddEventListeners(col)

	return ch, cls
}

// CommandFilter used to filter api.CommandInteraction in a CommandCollector
type CommandFilter func(reaction *core.ApplicationCommandInteraction) bool

// CommandCollector used to collect api.CommandInteraction(s) using a CommandFilter function
type CommandCollector struct {
	Channel  chan *core.ApplicationCommandInteraction
	Filter   CommandFilter
	Close    func()
	GuildID  discord.Snowflake
	MemberID discord.Snowflake
}

// OnEvent used to get events for the CommandCollector
func (c *CommandCollector) OnEvent(e interface{}) {
	if event, ok := e.(*core.SlashCommandEvent); ok {
		if !c.Filter(event.ApplicationCommandInteraction) {
			return
		}

		c.Channel <- event.ApplicationCommandInteraction
	} else if event, ok := e.(core.GuildLeaveEvent); ok && event.GuildID == c.GuildID {
		c.Close()
	}
}
