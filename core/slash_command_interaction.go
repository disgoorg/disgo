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

type OptionsMap map[string]SlashCommandOption

func (m OptionsMap) Get(name string) *SlashCommandOption {
	if option, ok := m[name]; ok {
		return &option
	}
	return nil
}

func (m OptionsMap) GetAll() []SlashCommandOption {
	options := make([]SlashCommandOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m OptionsMap) GetByType(optionType discord.SlashCommandOptionType) []SlashCommandOption {
	var options []SlashCommandOption
	for _, option := range m {
		if option.Type == optionType {
			options = append(options, option)
		}
	}
	return options
}

func (m OptionsMap) Find(optionFindFunc func(option SlashCommandOption) bool) *SlashCommandOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return &option
		}
	}
	return nil
}

func (m OptionsMap) FindAll(optionFindFunc func(option SlashCommandOption) bool) []SlashCommandOption {
	var options []SlashCommandOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}

// SlashCommandOption holds info about an SlashCommandOption.Value
type SlashCommandOption struct {
	Resolved *Resolved
	Name     string
	Type     discord.SlashCommandOptionType
	Value    interface{}
}

// String returns the SlashCommandOption.Value as string
func (o SlashCommandOption) String() string {
	return o.Value.(string)
}

// Int returns the SlashCommandOption.Value as int
func (o SlashCommandOption) Int() int {
	return o.Value.(int)
}

// Float64 returns the SlashCommandOption.Value as float64
func (o SlashCommandOption) Float64() float64 {
	return o.Value.(float64)
}

// Float32 returns the SlashCommandOption.Value as float32
func (o SlashCommandOption) Float32() float32 {
	return o.Value.(float32)
}

// Bool returns the SlashCommandOption.Value as bool
func (o SlashCommandOption) Bool() bool {
	return o.Value.(bool)
}

// Snowflake returns the SlashCommandOption.Value as Snowflake
func (o SlashCommandOption) Snowflake() discord.Snowflake {
	return discord.Snowflake(o.String())
}

// User returns the SlashCommandOption.Value as User
func (o SlashCommandOption) User() *User {
	return o.Resolved.Users[o.Snowflake()]
}

// Member returns the SlashCommandOption.Value as Member
func (o SlashCommandOption) Member() *Member {
	return o.Resolved.Members[o.Snowflake()]
}

// Role returns the SlashCommandOption.Value as Role
func (o SlashCommandOption) Role() *Role {
	return o.Resolved.Roles[o.Snowflake()]
}

// Channel returns the SlashCommandOption.Value as Channel
func (o SlashCommandOption) Channel() *Channel {
	return o.Resolved.Channels[o.Snowflake()]
}
