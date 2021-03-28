package events

import (
	"errors"

	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	api.Event
	Interaction api.Interaction
}

// Guild returns the api.Guild from the api.Cache
func (e GenericInteractionEvent) Guild() *api.Guild {
	if e.Interaction.GuildID == nil {
		return nil
	}
	return e.Disgo.Cache().Guild(*e.Interaction.GuildID)
}

// DMChannel returns the api.DMChannel from the api.Cache
func (e GenericInteractionEvent) DMChannel() *api.DMChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().DMChannel(*e.Interaction.ChannelID)
}

// MessageChannel returns the api.MessageChannel from the api.Cache
func (e GenericInteractionEvent) MessageChannel() *api.MessageChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().MessageChannel(*e.Interaction.ChannelID)
}

// TextChannel returns the api.TextChannel from the api.Cache
func (e GenericInteractionEvent) TextChannel() *api.TextChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().TextChannel(*e.Interaction.ChannelID)
}

// GuildChannel returns the api.GuildChannel from the api.Cache
func (e GenericInteractionEvent) GuildChannel() *api.GuildChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo.Cache().GuildChannel(*e.Interaction.ChannelID)
}

// SlashCommandEvent indicates a slash api.SlashCommand was ran in a api.Guild
type SlashCommandEvent struct {
	GenericInteractionEvent
	ResponseChannel chan interface{}
	FromWebhook     bool
	CommandID       api.Snowflake
	Name            string
	SubCommandName  *string
	SubCommandGroup *string
	Options         []*Option
	Replied         bool
}

// Option holds info about an Option.Value
type Option struct {
	Resolved *api.Resolved
	Name     string
	Type     api.SlashCommandOptionType
	Value    interface{}
}

// String returns the Option.Value as string
func (o Option) String() string {
	return o.Value.(string)
}

// Bool returns the Option.Value as bool
func (o Option) Bool() bool {
	return o.Value.(bool)
}

// Snowflake returns the Option.Value as api.Snowflake
func (o Option) Snowflake() api.Snowflake {
	return api.Snowflake(o.String())
}

// User returns the Option.Value as api.User
func (o Option) User() *api.User {
	return o.Resolved.Users[o.Snowflake()]
}

// Member returns the Option.Value as api.Member
func (o Option) Member() *api.Member {
	return o.Resolved.Members[o.Snowflake()]
}

// Role returns the Option.Value as api.Role
func (o Option) Role() *api.Role {
	return o.Resolved.Roles[o.Snowflake()]
}

// Channel returns the Option.Value as api.Channel
func (o Option) Channel() *api.Channel {
	return o.Resolved.Channels[o.Snowflake()]
}

// MessageChannel returns the Option.Value as api.MessageChannel
func (o Option) MessageChannel() *api.MessageChannel {
	channel := o.Channel()
	if channel == nil || (channel.Type != api.ChannelTypeText && channel.Type != api.ChannelTypeNews) {
		return nil
	}
	return &api.MessageChannel{Channel: *channel}
}

// GuildChannel returns the Option.Value as api.GuildChannel
func (o Option) GuildChannel() *api.GuildChannel {
	channel := o.Channel()
	if channel == nil || (channel.Type != api.ChannelTypeText && channel.Type != api.ChannelTypeNews && channel.Type != api.ChannelTypeCategory && channel.Type != api.ChannelTypeStore && channel.Type != api.ChannelTypeVoice) {
		return nil
	}
	return &api.GuildChannel{Channel: *channel}
}

// VoiceChannel returns the Option.Value as api.VoiceChannel
func (o Option) VoiceChannel() *api.VoiceChannel {
	channel := o.Channel()
	if channel == nil || channel.Type != api.ChannelTypeVoice {
		return nil
	}
	return &api.VoiceChannel{GuildChannel: api.GuildChannel{Channel: *channel}}
}

// TextChannel returns the Option.Value as api.TextChannel
func (o Option) TextChannel() *api.TextChannel {
	channel := o.Channel()
	if channel == nil || (channel.Type != api.ChannelTypeText && channel.Type != api.ChannelTypeNews) {
		return nil
	}
	return &api.TextChannel{GuildChannel: api.GuildChannel{Channel: *channel}, MessageChannel: api.MessageChannel{Channel: *channel}}
}

// Category returns the Option.Value as api.Category
func (o Option) Category() *api.Category {
	channel := o.Channel()
	if channel == nil || channel.Type != api.ChannelTypeCategory {
		return nil
	}
	return &api.Category{GuildChannel: api.GuildChannel{Channel: *channel}}
}

// StoreChannel returns the Option.Value as api.StoreChannel
func (o Option) StoreChannel() *api.StoreChannel {
	channel := o.Channel()
	if channel == nil || channel.Type != api.ChannelTypeStore {
		return nil
	}
	return &api.StoreChannel{GuildChannel: api.GuildChannel{Channel: *channel}}
}

// CommandPath returns the api.SlashCommand path
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

// OptionByName returns an Option by name
func (e SlashCommandEvent) OptionByName(name string) *Option {
	options := e.OptionsByName(name)
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

// OptionsByName returns Option(s) by name
func (e SlashCommandEvent) OptionsByName(name string) []*Option {
	options := make([]*Option, 0)
	for _, option := range e.Options {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsByType returns Option(s) by api.SlashCommandOptionType
func (e SlashCommandEvent) OptionsByType(optionType api.SlashCommandOptionType) []*Option {
	options := make([]*Option, 0)
	for _, option := range e.Options {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

// Reply replies to the api.Interaction with the provided api.InteractionResponse
func (e *SlashCommandEvent) Reply(response api.InteractionResponse) error {
	if e.Replied {
		return errors.New("you already replied to this interaction")
	}
	e.Replied = true

	if e.FromWebhook {
		e.ResponseChannel <- response
		return nil
	}

	return e.Disgo.RestClient().SendInteractionResponse(e.Interaction.ID, e.Interaction.Token, response)
}
