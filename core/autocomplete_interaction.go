package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/snowflake"
)

type AutocompleteInteractionFilter func(autocompleteInteraction *AutocompleteInteraction) bool

var _ Interaction = (*AutocompleteInteraction)(nil)

type AutocompleteInteraction struct {
	*BaseInteraction
	Data AutocompleteInteractionData
}

func (i AutocompleteInteraction) interaction() {}
func (i AutocompleteInteraction) Type() discord.InteractionType {
	return discord.InteractionTypeAutocomplete
}

func (i AutocompleteInteraction) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

func (i AutocompleteInteraction) ResultMapString(resultMap map[string]string, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceString{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return i.Result(choices, opts...)
}

func (i AutocompleteInteraction) ResultMapInt(resultMap map[string]int, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceInt{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return i.Result(choices, opts...)
}

func (i AutocompleteInteraction) ResultMapFloat(resultMap map[string]float64, opts ...rest.RequestOpt) error {
	choices := make([]discord.AutocompleteChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.AutocompleteChoiceFloat{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return i.Result(choices, opts...)
}

type AutocompleteInteractionData struct {
	discord.AutocompleteInteractionData
	SubCommandName      *string
	SubCommandGroupName *string
	Options             AutocompleteOptionsMap
}

// CommandPath returns the ApplicationCommand path
func (i *AutocompleteInteractionData) CommandPath() string {
	path := i.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

type AutocompleteOptionsMap map[string]discord.AutocompleteOption

func (m AutocompleteOptionsMap) Get(name string) discord.AutocompleteOption {
	if option, ok := m[name]; ok {
		return option
	}
	return nil
}

func (m AutocompleteOptionsMap) StringOption(name string) *discord.AutocompleteOptionString {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionString); ok {
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

func (m AutocompleteOptionsMap) IntOption(name string) *discord.AutocompleteOptionInt {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionInt); ok {
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

func (m AutocompleteOptionsMap) BoolOption(name string) *discord.AutocompleteOptionBool {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionBool); ok {
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

func (m AutocompleteOptionsMap) UserOption(name string) *discord.AutocompleteOptionUser {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionUser); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) ChannelOption(name string) *discord.AutocompleteOptionChannel {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionChannel); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) RoleOption(name string) *discord.AutocompleteOptionRole {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionRole); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) MentionableOption(name string) *discord.AutocompleteOptionMentionable {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionMentionable); ok {
		return &opt
	}
	return nil
}

func (m AutocompleteOptionsMap) Snowflake(name string) *snowflake.Snowflake {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	switch opt := option.(type) {
	case discord.AutocompleteOptionChannel:
		return &opt.Value

	case discord.AutocompleteOptionRole:
		return &opt.Value

	case discord.AutocompleteOptionUser:
		return &opt.Value

	case discord.AutocompleteOptionMentionable:
		return &opt.Value

	default:
		return nil
	}
}

func (m AutocompleteOptionsMap) FloatOption(name string) *discord.AutocompleteOptionFloat {
	option := m.Get(name)
	if option == nil {
		return nil
	}
	if opt, ok := option.(discord.AutocompleteOptionFloat); ok {
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
	return m.Get(name).Focused()
}

func (m AutocompleteOptionsMap) FocusedOption() discord.AutocompleteOption {
	return m.Find(func(option discord.AutocompleteOption) bool {
		return option.Focused()
	})
}
