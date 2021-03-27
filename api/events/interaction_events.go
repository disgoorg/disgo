package events

import (
	"errors"

	"github.com/DiscoOrg/disgo/api"
)

type GenericInteractionEvent struct {
	api.Event
	api.Interaction
}

func (e GenericInteractionEvent) Guild() *api.Guild {
	if e.GuildID == nil {
		return nil
	}
	return e.Disgo.Cache().Guild(*e.GuildID)
}

func (e GenericInteractionEvent) DMChannel() *api.DMChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().DMChannel(*e.ChannelID)
}
func (e GenericInteractionEvent) MessageChannel() *api.MessageChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().MessageChannel(*e.ChannelID)
}
func (e GenericInteractionEvent) TextChannel() *api.TextChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().TextChannel(*e.ChannelID)
}
func (e GenericInteractionEvent) GuildChannel() *api.GuildChannel {
	if e.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().GuildChannel(*e.ChannelID)
}

type SlashCommandEvent struct {
	GenericInteractionEvent
	CommandID       api.Snowflake
	Name            string
	SubCommandName  *string
	SubCommandGroup *string
	Options         []*Option
	Replied         bool
}

type Option struct {
	Resolved *api.Resolved
	Name     string
	Type     api.CommandOptionType
	Value    interface{}
}

func (d Option) String() string {
	return d.Value.(string)
}

func (d Option) Bool() bool {
	return d.Value.(bool)
}

func (d Option) Snowflake() api.Snowflake {
	return api.Snowflake(d.String())
}

func (d Option) User() *api.User {
	return d.Resolved.Users[d.Snowflake()]
}

func (d Option) Member() *api.Member {
	return d.Resolved.Members[d.Snowflake()]
}

func (d Option) Role() *api.Role {
	return d.Resolved.Roles[d.Snowflake()]
}

func (d Option) Channel() *api.Channel {
	return d.Resolved.Channels[d.Snowflake()]
}

func (d Option) MessageChannel() *api.MessageChannel {
	channel := d.Channel()
	if channel == nil || (channel.Type != api.ChannelTypeText && channel.Type != api.ChannelTypeNews) {
		return nil
	}
	return &api.MessageChannel{Channel: *channel}
}

func (d Option) GuildChannel() *api.GuildChannel {
	channel := d.Channel()
	if channel == nil || (channel.Type != api.ChannelTypeText && channel.Type != api.ChannelTypeNews && channel.Type != api.ChannelTypeCategory && channel.Type != api.ChannelTypeStore && channel.Type != api.ChannelTypeVoice) {
		return nil
	}
	return &api.GuildChannel{Channel: *channel}
}

func (d Option) VoiceChannel() *api.VoiceChannel {
	channel := d.Channel()
	if channel == nil || channel.Type != api.ChannelTypeVoice {
		return nil
	}
	return &api.VoiceChannel{GuildChannel: api.GuildChannel{Channel: *channel}}
}

func (d Option) TextChannel() *api.TextChannel {
	channel := d.Channel()
	if channel == nil || (channel.Type != api.ChannelTypeText && channel.Type != api.ChannelTypeNews) {
		return nil
	}
	return &api.TextChannel{GuildChannel: api.GuildChannel{Channel: *channel}, MessageChannel: api.MessageChannel{Channel: *channel}}
}

func (d Option) Category() *api.Category {
	channel := d.Channel()
	if channel == nil || channel.Type != api.ChannelTypeCategory {
		return nil
	}
	return &api.Category{GuildChannel: api.GuildChannel{Channel: *channel}}
}

func (d Option) StoreChannel() *api.StoreChannel {
	channel := d.Channel()
	if channel == nil || channel.Type != api.ChannelTypeStore {
		return nil
	}
	return &api.StoreChannel{GuildChannel: api.GuildChannel{Channel: *channel}}
}

func (e SlashCommandEvent) CommandPath() string {
	path := e.Name
	if e.SubCommandName != nil {
		path += "/" + *e.SubCommandName
	}
	if e.SubCommandGroup != nil {
		path += "/" + *e.SubCommandGroup
	}
	return path
}

func (e SlashCommandEvent) OptionByName(name string) *Option {
	options := e.OptionsByName(name)
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

func (e SlashCommandEvent) OptionsByName(name string) []*Option {
	options := make([]*Option, 0)
	for _, option := range e.Options {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

func (e SlashCommandEvent) OptionsByType(optionType api.CommandOptionType) []*Option {
	options := make([]*Option, 0)
	for _, option := range e.Options {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

func (e SlashCommandEvent) Reply(response api.InteractionResponse) error {
	if e.Replied {
		return errors.New("you already replied to this interaction")
	}
	e.Replied = true
	return e.Disgo.RestClient().SendInteractionResponse(e.Interaction.ID, e.Token, response)
}
