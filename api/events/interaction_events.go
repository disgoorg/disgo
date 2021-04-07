package events

import (
	"errors"

	"github.com/DisgoOrg/disgo/api"
)

// GenericInteractionEvent generic api.Interaction event
type GenericInteractionEvent struct {
	GenericEvent
	Interaction api.Interaction
}

// Guild returns the api.Guild from the api.Cache
func (e GenericInteractionEvent) Guild() *api.Guild {
	if e.Interaction.GuildID == nil {
		return nil
	}
	return e.Disgo().Cache().Guild(*e.Interaction.GuildID)
}

// DMChannel returns the api.DMChannel from the api.Cache
func (e GenericInteractionEvent) DMChannel() *api.DMChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().DMChannel(*e.Interaction.ChannelID)
}

// MessageChannel returns the api.MessageChannel from the api.Cache
func (e GenericInteractionEvent) MessageChannel() *api.MessageChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().MessageChannel(*e.Interaction.ChannelID)
}

// TextChannel returns the api.TextChannel from the api.Cache
func (e GenericInteractionEvent) TextChannel() *api.TextChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().TextChannel(*e.Interaction.ChannelID)
}

// GuildChannel returns the api.GuildChannel from the api.Cache
func (e GenericInteractionEvent) GuildChannel() *api.GuildChannel {
	if e.Interaction.ChannelID == nil {
		return nil
	}
	return e.Disgo().Cache().GuildChannel(*e.Interaction.ChannelID)
}

// SlashCommandEvent indicates a slash api.Command was ran in a api.Guild
type SlashCommandEvent struct {
	GenericInteractionEvent
	ResponseChannel     chan interface{}
	FromWebhook         bool
	CommandID           api.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             []*Option
	Replied             bool
}

// Option holds info about an Option.Value
type Option struct {
	Resolved *api.Resolved
	Name     string
	Type     api.CommandOptionType
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

// CommandPath returns the api.Command path
func (e SlashCommandEvent) CommandPath() string {
	path := e.CommandName
	if e.SubCommandName != nil {
		path += "/" + *e.SubCommandName
	}
	if e.SubCommandGroupName != nil {
		path += "/" + *e.SubCommandGroupName
	}
	return path
}

// Option returns an Option by name
func (e SlashCommandEvent) Option(name string) *Option {
	options := e.OptionN(name)
	if len(options) == 0 {
		return nil
	}
	return options[0]
}

// OptionN returns Option(s) by name
func (e SlashCommandEvent) OptionN(name string) []*Option {
	options := make([]*Option, 0)
	for _, option := range e.Options {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsT returns Option(s) by api.CommandOptionType
func (e SlashCommandEvent) OptionsT(optionType api.CommandOptionType) []*Option {
	options := make([]*Option, 0)
	for _, option := range e.Options {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

// Acknowledge replies to the api.Interaction with api.InteractionResponseTypeDeferredChannelMessageWithSource
func (e *SlashCommandEvent) Acknowledge() error {
	return e.Reply(api.NewInteractionResponseBuilder().SetType(api.InteractionResponseTypeDeferredChannelMessageWithSource).Build())
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

	return e.Disgo().RestClient().SendInteractionResponse(e.Interaction.ID, e.Interaction.Token, response)
}

// EditOriginal edits the original api.InteractionResponse
func (e *SlashCommandEvent) EditOriginal(followupMessage api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().EditInteractionResponse(e.Disgo().SelfUserID(), e.Interaction.Token, followupMessage)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (e *SlashCommandEvent) DeleteOriginal() error {
	return e.Disgo().RestClient().DeleteInteractionResponse(e.Disgo().SelfUserID(), e.Interaction.Token)
}

// SendFollowup used to send a api.FollowupMessage to an api.Interaction
func (e *SlashCommandEvent) SendFollowup(followupMessage api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().SendFollowupMessage(e.Disgo().SelfUserID(), e.Interaction.Token, followupMessage)
}

// EditFollowup used to edit a api.FollowupMessage from an api.Interaction
func (e *SlashCommandEvent) EditFollowup(messageID api.Snowflake, followupMessage api.FollowupMessage) (*api.Message, error) {
	return e.Disgo().RestClient().EditFollowupMessage(e.Disgo().SelfUserID(), e.Interaction.Token, messageID, followupMessage)
}

// DeleteFollowup used to delete a api.FollowupMessage from an api.Interaction
func (e *SlashCommandEvent) DeleteFollowup(messageID api.Snowflake) error {
	return e.Disgo().RestClient().DeleteFollowupMessage(e.Disgo().SelfUserID(), e.Interaction.Token, messageID)
}
