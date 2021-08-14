package core

import "github.com/DisgoOrg/disgo/discord"

type SlashCommandInteraction struct {
	*ApplicationCommandInteraction
	Data *SlashCommandInteractionData
}

// SubCommandName the subcommand name of the ApplicationCommand which got used. May be nil
func (i *SlashCommandInteraction) SubCommandName() *string {
	return i.Data.SubCommandName
}

// SubCommandGroupName the subcommand group name of the ApplicationCommand which got used. May be nil
func (i *SlashCommandInteraction) SubCommandGroupName() *string {
	return i.Data.SubCommandGroupName
}

// CommandPath returns the ApplicationCommand path
func (i *SlashCommandInteraction) CommandPath() string {
	path := i.CommandName()
	if name := i.SubCommandName(); name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName(); name != nil {
		path += "/" + *name
	}
	return path
}

// Options returns the parsed ApplicationCommandOption which the ApplicationCommand got used with
func (i *SlashCommandInteraction) Options() []ApplicationCommandOption {
	return i.Data.Options
}

// Option returns an Option by name
func (i *SlashCommandInteraction) Option(name string) *ApplicationCommandOption {
	options := i.OptionN(name)
	if len(options) == 0 {
		return nil
	}
	return &options[0]
}

// OptionN returns Option(s) by name
func (i *SlashCommandInteraction) OptionN(name string) []ApplicationCommandOption {
	options := make([]ApplicationCommandOption, 0)
	for _, option := range i.Options() {
		if option.Name == name {
			options = append(options, option)
		}
	}
	return options
}

// OptionsT returns Option(s) by api.CommandOptionType
func (i *SlashCommandInteraction) OptionsT(optionType discord.ApplicationCommandOptionType) []ApplicationCommandOption {
	options := make([]ApplicationCommandOption, 0)
	for _, option := range i.Options() {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

type SlashCommandInteractionData struct {
	*ApplicationCommandInteractionData
	SubCommandName      *string
	SubCommandGroupName *string
	Options             []ApplicationCommandOption
}

// ApplicationCommandOption holds info about an ApplicationCommandOption.Value
type ApplicationCommandOption struct {
	Resolved *Resolved
	Name     string
	Type     discord.ApplicationCommandOptionType
	Value    interface{}
}

// String returns the ApplicationCommandOption.Value as string
func (o ApplicationCommandOption) String() string {
	return o.Value.(string)
}

// Int returns the ApplicationCommandOption.Value as int
func (o ApplicationCommandOption) Int() int {
	return o.Value.(int)
}

// Float64 returns the ApplicationCommandOption.Value as float64
func (o ApplicationCommandOption) Float64() float64 {
	return o.Value.(float64)
}

// Float32 returns the ApplicationCommandOption.Value as float32
func (o ApplicationCommandOption) Float32() float32 {
	return o.Value.(float32)
}

// Bool returns the ApplicationCommandOption.Value as bool
func (o ApplicationCommandOption) Bool() bool {
	return o.Value.(bool)
}

// Snowflake returns the ApplicationCommandOption.Value as Snowflake
func (o ApplicationCommandOption) Snowflake() discord.Snowflake {
	return discord.Snowflake(o.String())
}

// User returns the ApplicationCommandOption.Value as User
func (o ApplicationCommandOption) User() *User {
	return o.Resolved.Users[o.Snowflake()]
}

// Member returns the ApplicationCommandOption.Value as Member
func (o ApplicationCommandOption) Member() *Member {
	return o.Resolved.Members[o.Snowflake()]
}

// Role returns the ApplicationCommandOption.Value as Role
func (o ApplicationCommandOption) Role() *Role {
	return o.Resolved.Roles[o.Snowflake()]
}

// Channel returns the ApplicationCommandOption.Value as Channel
func (o ApplicationCommandOption) Channel() Channel {
	return o.Channel()
}

// MessageChannel returns the ApplicationCommandOption.Value as MessageChannel
func (o ApplicationCommandOption) MessageChannel() MessageChannel {
	channel := o.Channel()
	if channel == nil || !channel.IsMessageChannel() {
		return nil
	}
	return channel.(MessageChannel)
}

// GuildChannel returns the ApplicationCommandOption.Value as GuildChannel
func (o ApplicationCommandOption) GuildChannel() GuildChannel {
	channel := o.Channel()
	if channel == nil || !channel.IsGuildChannel() {
		return nil
	}
	return channel.(GuildChannel)
}

// VoiceChannel returns the ApplicationCommandOption.Value as VoiceChannel
func (o ApplicationCommandOption) VoiceChannel() VoiceChannel {
	channel := o.Channel()
	if channel == nil || !channel.IsVoiceChannel() {
		return nil
	}
	return channel.(VoiceChannel)
}

// TextChannel returns the ApplicationCommandOption.Value as TextChannel
func (o ApplicationCommandOption) TextChannel() TextChannel {
	channel := o.Channel()
	if channel == nil || channel.IsTextChannel() {
		return nil
	}
	return channel.(TextChannel)
}

// Category returns the ApplicationCommandOption.Value as Category
func (o ApplicationCommandOption) Category() Category {
	channel := o.Channel()
	if channel == nil || channel.IsCategory() {
		return nil
	}
	return channel.(Category)
}

// StoreChannel returns the ApplicationCommandOption.Value as StoreChannel
func (o ApplicationCommandOption) StoreChannel() StoreChannel {
	channel := o.Channel()
	if channel == nil || channel.IsStoreChannel() {
		return nil
	}
	return channel.(StoreChannel)
}

// StageChannel returns the ApplicationCommandOption.Value as StageChannel
func (o ApplicationCommandOption) StageChannel() StageChannel {
	channel := o.Channel()
	if channel == nil || channel.IsStageChannel() {
		return nil
	}
	return channel.(StageChannel)
}
