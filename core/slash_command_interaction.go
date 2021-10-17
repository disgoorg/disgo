package core

import "github.com/DisgoOrg/disgo/discord"

type SlashCommandInteractionFilter func(slashCommandInteraction *SlashCommandInteraction) bool

type SlashCommandInteraction struct {
	discord.SlashCommandInteraction
	CreateInteraction
	FollowupInteraction
	CommandID           discord.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Resolved            *SlashCommandResolved
	Options             SlashCommandOptionsMap
}

// CommandPath returns the ApplicationCommand path
func (i *SlashCommandInteraction) CommandPath() string {
	path := i.Data.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

// Guild returns the Guild from the Caches
func (i *SlashCommandInteraction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.GuildCache().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *SlashCommandInteraction) Channel() *Channel {
	return i.Bot.Caches.ChannelCache().Get(i.ChannelID)
}

// SlashCommandResolved contains resolved mention data for SlashCommand(s)
type SlashCommandResolved struct {
	Users    map[discord.Snowflake]*User
	Members  map[discord.Snowflake]*Member
	Roles    map[discord.Snowflake]*Role
	Channels map[discord.Snowflake]*Channel
}

type SlashCommandOptionsMap map[string]SlashCommandOption

func (m SlashCommandOptionsMap) Get(name string) SlashCommandOption {
	if option, ok := m[name]; ok {
		return option
	}
	return nil
}

func (m SlashCommandOptionsMap) StringOption(name string) *SlashCommandOptionString {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionString); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) String(name string) *string {
	option := m.StringOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m SlashCommandOptionsMap) IntOption(name string) *SlashCommandOptionInt {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionInt); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) Int(name string) *int {
	option := m.IntOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m SlashCommandOptionsMap) BoolOption(name string) *SlashCommandOptionBool {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionBool); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) Bool(name string) *bool {
	option := m.BoolOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m SlashCommandOptionsMap) UserOption(name string) *SlashCommandOptionUser {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionUser); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) User(name string) *User {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	switch opt := option.(type) {
	case SlashCommandOptionUser:
		return opt.User()
	case SlashCommandOptionMentionable:
		return opt.User()
	default:
		return nil
	}
}

func (m SlashCommandOptionsMap) Member(name string) *Member {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	switch opt := option.(type) {
	case SlashCommandOptionUser:
		return opt.Member()
	case SlashCommandOptionMentionable:
		return opt.Member()
	default:
		return nil
	}
}

func (m SlashCommandOptionsMap) ChannelOption(name string) *SlashCommandOptionChannel {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionChannel); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) Channel(name string) *Channel {
	option := m.ChannelOption(name)
	if option == nil {
		return nil
	}
	return option.Channel()
}

func (m SlashCommandOptionsMap) RoleOption(name string) *SlashCommandOptionRole {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionRole); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) Role(name string) *Role {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	switch opt := option.(type) {
	case SlashCommandOptionRole:
		return opt.Role()
	case SlashCommandOptionMentionable:
		return opt.Role()
	default:
		return nil
	}
}

func (m SlashCommandOptionsMap) MentionableOption(name string) *SlashCommandOptionMentionable {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionMentionable); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) Snowflake(name string) *discord.Snowflake {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	switch opt := option.(type) {
	case SlashCommandOptionChannel:
		return &opt.Value

	case SlashCommandOptionRole:
		return &opt.Value

	case SlashCommandOptionUser:
		return &opt.Value

	case SlashCommandOptionMentionable:
		return &opt.Value

	default:
		return nil
	}
}

func (m SlashCommandOptionsMap) FloatOption(name string) *SlashCommandOptionFloat {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionFloat); ok {
		return &opt
	}
	return nil
}

func (m SlashCommandOptionsMap) Float(name string) *float64 {
	option := m.FloatOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m SlashCommandOptionsMap) GetAll() []SlashCommandOption {
	options := make([]SlashCommandOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m SlashCommandOptionsMap) GetByType(optionType discord.ApplicationCommandOptionType) []SlashCommandOption {
	return m.FindAll(func(option SlashCommandOption) bool {
		return option.Type() == optionType
	})
}

func (m SlashCommandOptionsMap) Find(optionFindFunc func(option SlashCommandOption) bool) SlashCommandOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return option
		}
	}
	return nil
}

func (m SlashCommandOptionsMap) FindAll(optionFindFunc func(option SlashCommandOption) bool) []SlashCommandOption {
	var options []SlashCommandOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}
