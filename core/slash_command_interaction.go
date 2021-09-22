package core

import "github.com/DisgoOrg/disgo/discord"

type SlashCommandInteraction struct {
	*ApplicationCommandInteraction
	SlashCommandInteractionData
}

// CommandPath returns the ApplicationCommand path
func (i *SlashCommandInteraction) CommandPath() string {
	path := i.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

type SlashCommandInteractionData struct {
	SubCommandName      *string
	SubCommandGroupName *string
	Options             OptionsMap
}

type OptionsMap map[string]ApplicationCommandOption

func (m OptionsMap) Get(name string) *ApplicationCommandOption {
	if option, ok := m[name]; ok {
		return &option
	}
	return nil
}

func (m OptionsMap) GetAll() []ApplicationCommandOption {
	options := make([]ApplicationCommandOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m OptionsMap) GetByType(optionType discord.ApplicationCommandOptionType) []ApplicationCommandOption {
	return m.FindAll(func(option ApplicationCommandOption) bool {
		return option.Type == optionType
	})
}

func (m OptionsMap) Find(optionFindFunc func(option ApplicationCommandOption) bool) *ApplicationCommandOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return &option
		}
	}
	return nil
}

func (m OptionsMap) FindAll(optionFindFunc func(option ApplicationCommandOption) bool) []ApplicationCommandOption {
	var options []ApplicationCommandOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}

// ApplicationCommandOption holds info about an ApplicationCommandOption.Value
type ApplicationCommandOption struct {
	Resolved *Resolved
	Name     string
	Type     discord.ApplicationCommandOptionType
	Value    interface{}
	Focused  bool
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
func (o ApplicationCommandOption) Channel() *Channel {
	return o.Resolved.Channels[o.Snowflake()]
}
