package api

// CommandInteraction is a specific Interaction when using Command(s)
type CommandInteraction struct {
	*Interaction
	Data *CommandInteractionData `json:"data,omitempty"`
}

// CommandID returns the ID of the api.Command which got used
func (i *CommandInteraction) CommandID() Snowflake {
	return i.Data.ID
}

// CommandName the name of the api.Command which got used
func (i *CommandInteraction) CommandName() string {
	return i.Data.CommandName
}

// SubCommandName the subcommand name of the api.Command which got used. May be nil
func (i *CommandInteraction) SubCommandName() *string {
	return i.Data.SubCommandName
}

// SubCommandGroupName the subcommand group name of the api.Command which got used. May be nil
func (i *CommandInteraction) SubCommandGroupName() *string {
	return i.Data.SubCommandGroupName
}

// CommandPath returns the api.Command path
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

// Options returns the parsed Option which the Command got used with
func (i *CommandInteraction) Options() []*Option {
	return i.Data.Options
}

// CommandInteractionData is the command data payload
type CommandInteractionData struct {
	ID                  Snowflake    `json:"id"`
	CommandName         string       `json:"name"`
	SubCommandName      *string      `json:"-"`
	SubCommandGroupName *string      `json:"-"`
	Resolved            *Resolved    `json:"resolved,omitempty"`
	RawOptions          []OptionData `json:"options,omitempty"`
	Options             []*Option    `json:"-"`
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[Snowflake]*User        `json:"users,omitempty"`
	Members  map[Snowflake]*Member      `json:"members,omitempty"`
	Roles    map[Snowflake]*Role        `json:"roles,omitempty"`
	Channels map[Snowflake]*ChannelImpl `json:"channels,omitempty"`
}

// OptionData is used for options or subcommands in your slash commands
type OptionData struct {
	Name    string            `json:"name"`
	Type    CommandOptionType `json:"type"`
	Value   interface{}       `json:"value,omitempty"`
	Options []OptionData      `json:"options,omitempty"`
}

// Option holds info about an Option.Value
type Option struct {
	Resolved *Resolved
	Name     string
	Type     CommandOptionType
	Value    interface{}
}

// String returns the Option.Value as string
func (o Option) String() string {
	return o.Value.(string)
}

// Integer returns the Option.Value as int
func (o Option) Integer() int {
	return o.Value.(int)
}

// Bool returns the Option.Value as bool
func (o Option) Bool() bool {
	return o.Value.(bool)
}

// Snowflake returns the Option.Value as Snowflake
func (o Option) Snowflake() Snowflake {
	return Snowflake(o.String())
}

// User returns the Option.Value as User
func (o Option) User() *User {
	return o.Resolved.Users[o.Snowflake()]
}

// Member returns the Option.Value as Member
func (o Option) Member() *Member {
	return o.Resolved.Members[o.Snowflake()]
}

// Role returns the Option.Value as Role
func (o Option) Role() *Role {
	return o.Resolved.Roles[o.Snowflake()]
}

// Channel returns the Option.Value as Channel
func (o Option) Channel() Channel {
	return o.Resolved.Channels[o.Snowflake()]
}

func (o Option) channelImpl() *ChannelImpl {
	return o.Resolved.Channels[o.Snowflake()]
}

// MessageChannel returns the Option.Value as MessageChannel
func (o Option) MessageChannel() MessageChannel {
	channel := o.channelImpl()
	if channel == nil || channel.MessageChannel() {
		return nil
	}
	return channel
}

// GuildChannel returns the Option.Value as GuildChannel
func (o Option) GuildChannel() GuildChannel {
	channel := o.channelImpl()
	if channel == nil || channel.GuildChannel() {
		return nil
	}
	return channel
}

// VoiceChannel returns the Option.Value as VoiceChannel
func (o Option) VoiceChannel() VoiceChannel {
	channel := o.channelImpl()
	if channel == nil || channel.VoiceChannel(){
		return nil
	}
	return channel
}

// TextChannel returns the Option.Value as TextChannel
func (o Option) TextChannel() TextChannel {
	channel := o.channelImpl()
	if channel == nil || channel.TextChannel() {
		return nil
	}
	return channel
}

// Category returns the Option.Value as Category
func (o Option) Category() Category {
	channel := o.channelImpl()
	if channel == nil || channel.Category() {
		return nil
	}
	return channel
}

// StoreChannel returns the Option.Value as StoreChannel
func (o Option) StoreChannel() StoreChannel {
	channel := o.channelImpl()
	if channel == nil || channel.StoreChannel() {
		return nil
	}
	return channel
}
