package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type AutocompleteInteraction struct {
	*InteractionFields
	CommandID           discord.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             AutocompleteOptionsMap
}

func (i *AutocompleteInteraction) InteractionType() discord.InteractionType {
	return discord.InteractionTypeAutocomplete
}

func (i *AutocompleteInteraction) Respond(callbackType discord.InteractionCallbackType, callbackData discord.InteractionCallbackData, opts ...rest.RequestOpt) error {
	return respond(i.InteractionFields, callbackType, callbackData, opts...)
}

func (i *AutocompleteInteraction) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return result(i.InteractionFields, choices, opts...)
}

func (i *AutocompleteInteraction) ResultMapString(resultMap map[string]string, opts ...rest.RequestOpt) error {
	return resultMapString(i.InteractionFields, resultMap, opts...)
}

func (i *AutocompleteInteraction) ResultMapInt(resultMap map[string]int, opts ...rest.RequestOpt) error {
	return resultMapInt(i.InteractionFields, resultMap, opts...)
}

func (i *AutocompleteInteraction) ResultMapFloat(resultMap map[string]float64, opts ...rest.RequestOpt) error {
	return resultMapFloat(i.InteractionFields, resultMap, opts...)
}

// CommandPath returns the ApplicationCommand path
func (i *AutocompleteInteraction) CommandPath() string {
	path := i.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

// Guild returns the Guild from the Caches
func (i *AutocompleteInteraction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Bot.Caches.GuildCache().Get(*i.GuildID)
}

// Channel returns the Channel from the Caches
func (i *AutocompleteInteraction) Channel() *Channel {
	return i.Bot.Caches.ChannelCache().Get(i.ChannelID)
}

type AutocompleteOptionsMap map[string]discord.AutocompleteOption

func (m AutocompleteOptionsMap) Get(name string) discord.AutocompleteOption {
	if option, ok := m[name]; ok {
		return option
	}
	return nil
}

func (m AutocompleteOptionsMap) StringOption(name string) *SlashCommandOptionString {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionString); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) String(name string) *string {
	option := m.StringOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m AutocompleteOptionsMap) IntOption(name string) *SlashCommandOptionInt {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionInt); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Int(name string) *int {
	option := m.IntOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m AutocompleteOptionsMap) BoolOption(name string) *SlashCommandOptionBool {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionBool); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Bool(name string) *bool {
	option := m.BoolOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m AutocompleteOptionsMap) UserOption(name string) *SlashCommandOptionUser {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionUser); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) User(name string) *User {
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

func (m AutocompleteOptionsMap) Member(name string) *Member {
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

func (m AutocompleteOptionsMap) ChannelOption(name string) *SlashCommandOptionChannel {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionChannel); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Channel(name string) *Channel {
	option := m.ChannelOption(name)
	if option == nil {
		return nil
	}
	return option.Channel()
}

func (m AutocompleteOptionsMap) RoleOption(name string) *SlashCommandOptionRole {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionRole); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Role(name string) *Role {
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

func (m AutocompleteOptionsMap) MentionableOption(name string) *SlashCommandOptionMentionable {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionMentionable); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Snowflake(name string) *discord.Snowflake {
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

func (m AutocompleteOptionsMap) FloatOption(name string) *SlashCommandOptionFloat {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(SlashCommandOptionFloat); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Float(name string) *float64 {
	option := m.FloatOption(name)
	if option == nil {
		return nil
	}
	return &option.Value
}

func (m AutocompleteOptionsMap) GetAll() []discord.AutocompleteOption {
	options := make([]discord.AutocompleteOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m AutocompleteOptionsMap) GetByType(optionType discord.ApplicationCommandOptionType) []discord.AutocompleteOption {
	return m.FindAll(func(option discord.AutocompleteOption) bool {
		return option.Type() == optionType
	})
}

func (m AutocompleteOptionsMap) Find(optionFindFunc func(option discord.AutocompleteOption) bool) discord.AutocompleteOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return option
		}
	}
	return nil
}

func (m AutocompleteOptionsMap) FindAll(optionFindFunc func(option discord.AutocompleteOption) bool) []discord.AutocompleteOption {
	var options []discord.AutocompleteOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}

func (m AutocompleteOptionsMap) Focused(name string) bool {
	option := m.Get(name)
	switch o := option.(type) {
	case discord.AutocompleteOptionString:
		return o.Focused
	case discord.AutocompleteOptionInt:
		return o.Focused
	case discord.AutocompleteOptionBool:
		return o.Focused
	case discord.AutocompleteOptionUser:
		return o.Focused
	case discord.AutocompleteOptionChannel:
		return o.Focused
	case discord.AutocompleteOptionRole:
		return o.Focused
	case discord.AutocompleteOptionMentionable:
		return o.Focused
	case discord.AutocompleteOptionFloat:
		return o.Focused
	default:
		return false
	}
}

func (m AutocompleteOptionsMap) FocusedOption() discord.AutocompleteOption {
	return m.Find(func(option discord.AutocompleteOption) bool {
		switch o := option.(type) {
		case discord.AutocompleteOptionString:
			return o.Focused
		case discord.AutocompleteOptionInt:
			return o.Focused
		case discord.AutocompleteOptionBool:
			return o.Focused
		case discord.AutocompleteOptionUser:
			return o.Focused
		case discord.AutocompleteOptionChannel:
			return o.Focused
		case discord.AutocompleteOptionRole:
			return o.Focused
		case discord.AutocompleteOptionMentionable:
			return o.Focused
		case discord.AutocompleteOptionFloat:
			return o.Focused
		default:
			return false
		}
	})
}
