package api

// InteractionType is the type of Interaction
type InteractionType int

// Constants for InteractionType
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
)

// An Interaction is the slash command object you receive when a user uses one of your commands
type Interaction struct {
	ID        Snowflake        `json:"id"`
	Type      InteractionType  `json:"type"`
	Data      *InteractionData `json:"data,omitempty"`
	GuildID   *Snowflake       `json:"guild_id,omitempty"`
	ChannelID *Snowflake       `json:"channel_id,omitempty"`
	Member    *Member          `json:"member,omitempty"`
	User      *User            `json:"User,omitempty"`
	Token     string           `json:"token"`
	Version   int              `json:"version"`
}

// InteractionData is the command data payload
type InteractionData struct {
	ID       Snowflake     `json:"id"`
	Name     string        `json:"name"`
	Resolved *Resolved     `json:"resolved"`
	Options  []*OptionData `json:"options,omitempty"`
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
	Options []*OptionData     `json:"options,omitempty"`
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

// Channel returns the Option.Value as ChannelImpl
func (o Option) Channel() Channel {
	return o.Resolved.Channels[o.Snowflake()]
}

// MessageChannel returns the Option.Value as MessageChannel
func (o Option) MessageChannel() MessageChannel {
	channel := o.Channel()
	if channel == nil || !channel.MessageChannel() {
		return nil
	}
	return channel.(MessageChannel)
}

// GuildChannel returns the Option.Value as GuildChannel
func (o Option) GuildChannel() GuildChannel {
	channel := o.Channel()
	if channel == nil || !channel.GuildChannel() {
		return nil
	}
	return channel.(GuildChannel)
}

// VoiceChannel returns the Option.Value as VoiceChannel
func (o Option) VoiceChannel() VoiceChannel {
	channel := o.Channel()
	if channel == nil || !channel.VoiceChannel() {
		return nil
	}
	return channel.(VoiceChannel)
}

// TextChannel returns the Option.Value as TextChannel
func (o Option) TextChannel() TextChannel {
	channel := o.Channel()
	if channel == nil || !channel.TextChannel() {
		return nil
	}
	return channel.(TextChannel)
}

// Thread returns the Option.Value as Thread
func (o Option) Thread() Thread {
	channel := o.Channel()
	if channel == nil || !channel.Thread() {
		return nil
	}
	return channel.(Thread)
}

// Category returns the Option.Value as Category
func (o Option) Category() Category {
	channel := o.Channel()
	if channel == nil || !channel.Category() {
		return nil
	}
	return channel.(Category)
}

// StoreChannel returns the Option.Value as StoreChannel
func (o Option) StoreChannel() StoreChannel {
	channel := o.Channel()
	if channel == nil || !channel.StoreChannel() {
		return nil
	}
	return channel.(StoreChannel)
}
