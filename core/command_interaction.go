package core

import "github.com/DisgoOrg/disgo/discord"

type CommandInteraction struct {
	*Interaction
	Data *CommandInteractionData
}

// CommandID returns the ID of the Command which got used
func (i *CommandInteraction) CommandID() discord.Snowflake {
	return i.Data.ID
}

// CommandName the name of the Command which got used
func (i *CommandInteraction) CommandName() string {
	return i.Data.CommandName
}

// SubCommandName the subcommand name of the Command which got used. May be nil
func (i *CommandInteraction) SubCommandName() *string {
	return i.Data.SubCommandName
}

// SubCommandGroupName the subcommand group name of the Command which got used. May be nil
func (i *CommandInteraction) SubCommandGroupName() *string {
	return i.Data.SubCommandGroupName
}

// CommandPath returns the Command path
func (i *CommandInteraction) CommandPath() string {
	path := i.CommandName()
	if name := i.SubCommandName(); name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName(); name != nil {
		path += "/" + *name
	}
	return path
}

// Options returns the parsed CommandOption which the Command got used with
func (i *CommandInteraction) Options() []CommandOption {
	return i.Data.Options
}

type Resolved struct {
	discord.Resolved
	Users    map[discord.Snowflake]*User
	Members  map[discord.Snowflake]*Member
	Roles    map[discord.Snowflake]*Role
	Channels map[discord.Snowflake]Channel
}

type CommandInteractionData struct {
	*InteractionData
	Resolved            *Resolved
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             []CommandOption
}

// CommandOption holds info about an CommandOption.Value
type CommandOption struct {
	Resolved *Resolved
	Name     string
	Type     discord.CommandOptionType
	Value    interface{}
}

// String returns the CommandOption.Value as string
func (o CommandOption) String() string {
	return o.Value.(string)
}

// Int returns the CommandOption.Value as int
func (o CommandOption) Int() int {
	return o.Value.(int)
}

// Float64 returns the CommandOption.Value as float64
func (o CommandOption) Float64() float64 {
	return o.Value.(float64)
}

// Float32 returns the CommandOption.Value as float32
func (o CommandOption) Float32() float32 {
	return o.Value.(float32)
}

// Bool returns the CommandOption.Value as bool
func (o CommandOption) Bool() bool {
	return o.Value.(bool)
}

// Snowflake returns the CommandOption.Value as Snowflake
func (o CommandOption) Snowflake() discord.Snowflake {
	return discord.Snowflake(o.String())
}

// User returns the CommandOption.Value as User
func (o CommandOption) User() *User {
	return o.Resolved.Users[o.Snowflake()]
}

// Member returns the CommandOption.Value as Member
func (o CommandOption) Member() *Member {
	return o.Resolved.Members[o.Snowflake()]
}

// Role returns the CommandOption.Value as Role
func (o CommandOption) Role() *Role {
	return o.Resolved.Roles[o.Snowflake()]
}

// Channel returns the CommandOption.Value as Channel
func (o CommandOption) Channel() Channel {
	return o.Channel()
}

// MessageChannel returns the CommandOption.Value as MessageChannel
func (o CommandOption) MessageChannel() MessageChannel {
	channel := o.Channel()
	if channel == nil || !channel.IsMessageChannel() {
		return nil
	}
	return channel.(MessageChannel)
}

// GuildChannel returns the CommandOption.Value as GuildChannel
func (o CommandOption) GuildChannel() GuildChannel {
	channel := o.Channel()
	if channel == nil || !channel.IsGuildChannel() {
		return nil
	}
	return channel.(GuildChannel)
}

// VoiceChannel returns the CommandOption.Value as VoiceChannel
func (o CommandOption) VoiceChannel() VoiceChannel {
	channel := o.Channel()
	if channel == nil || !channel.IsVoiceChannel() {
		return nil
	}
	return channel.(VoiceChannel)
}

// TextChannel returns the CommandOption.Value as TextChannel
func (o CommandOption) TextChannel() TextChannel {
	channel := o.Channel()
	if channel == nil || channel.IsTextChannel() {
		return nil
	}
	return channel.(TextChannel)
}

// Category returns the CommandOption.Value as Category
func (o CommandOption) Category() Category {
	channel := o.Channel()
	if channel == nil || channel.IsCategory() {
		return nil
	}
	return channel.(Category)
}

// StoreChannel returns the CommandOption.Value as StoreChannel
func (o CommandOption) StoreChannel() StoreChannel {
	channel := o.Channel()
	if channel == nil || channel.IsStoreChannel() {
		return nil
	}
	return channel.(StoreChannel)
}

// StageChannel returns the CommandOption.Value as StageChannel
func (o CommandOption) StageChannel() StageChannel {
	channel := o.Channel()
	if channel == nil || channel.IsStageChannel() {
		return nil
	}
	return channel.(StageChannel)
}
