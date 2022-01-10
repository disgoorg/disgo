package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type SlashCommandOption interface {
	discord.SlashCommandOption
}

type SlashCommandOptionSubCommand struct {
	discord.SlashCommandOptionSubCommand
	Options []SlashCommandOption
}

type SlashCommandOptionSubCommandGroup struct {
	discord.SlashCommandOptionSubCommandGroup
	Options []SlashCommandOptionSubCommand
}

type SlashCommandOptionString struct {
	discord.SlashCommandOptionString
	Resolved *SlashCommandResolved
}

func (o *SlashCommandOptionString) MentionedUsers() []*User {
	matches := discord.MentionTypeUser.FindAllStringSubmatch(o.Value, -1)
	users := make([]*User, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		users[i] = o.Resolved.Users[discord.Snowflake(matches[i][1])]
	}
	return users
}

func (o *SlashCommandOptionString) MentionedMembers() []*Member {
	matches := discord.MentionTypeUser.FindAllStringSubmatch(o.Value, -1)
	members := make([]*Member, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		if member, ok := o.Resolved.Members[discord.Snowflake(matches[i][1])]; ok {
			members[i] = member
		}
	}
	return members
}

func (o *SlashCommandOptionString) MentionedChannels() []Channel {
	matches := discord.MentionTypeChannel.FindAllStringSubmatch(o.Value, -1)
	channels := make([]Channel, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		if channel, ok := o.Resolved.Channels[discord.Snowflake(matches[i][1])]; ok {
			channels[i] = channel
		}
	}
	return channels
}

func (o *SlashCommandOptionString) MentionedRoles() []*Role {
	matches := discord.MentionTypeRole.FindAllStringSubmatch(o.Value, -1)
	roles := make([]*Role, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		if role, ok := o.Resolved.Roles[discord.Snowflake(matches[i][1])]; ok {
			roles[i] = role
		}
	}
	return roles
}

type SlashCommandOptionInt struct {
	discord.SlashCommandOptionInt
}

type SlashCommandOptionBool struct {
	discord.SlashCommandOptionBool
}

type SlashCommandOptionUser struct {
	discord.SlashCommandOptionUser
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionUser) User() *User {
	return o.Resolved.Users[o.Value]
}

func (o SlashCommandOptionUser) Member() *Member {
	return o.Resolved.Members[o.Value]
}

type SlashCommandOptionChannel struct {
	discord.SlashCommandOptionChannel
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionChannel) Channel() Channel {
	return o.Resolved.Channels[o.Value]
}

type SlashCommandOptionRole struct {
	discord.SlashCommandOptionRole
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionRole) Role() *Role {
	return o.Resolved.Roles[o.Value]
}

type SlashCommandOptionMentionable struct {
	discord.SlashCommandOptionMentionable
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionMentionable) User() *User {
	return o.Resolved.Users[o.Value]
}

func (o SlashCommandOptionMentionable) Member() *Member {
	return o.Resolved.Members[o.Value]
}

func (o SlashCommandOptionMentionable) Role() *Role {
	return o.Resolved.Roles[o.Value]
}

type SlashCommandOptionFloat struct {
	discord.SlashCommandOptionFloat
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

func (m SlashCommandOptionsMap) Channel(name string) Channel {
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
