package discord

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/json"
)

var (
	_ Interaction = (*AutocompleteInteraction)(nil)
)

type AutocompleteInteraction struct {
	baseInteractionImpl
	Data AutocompleteInteractionData `json:"data"`
}

func (AutocompleteInteraction) Type() InteractionType {
	return InteractionTypeAutocomplete
}

type rawAutocompleteInteractionData struct {
	ID      snowflake.Snowflake  `json:"id"`
	Name    string               `json:"name"`
	Options []AutocompleteOption `json:"options"`
}

func (d *rawAutocompleteInteractionData) UnmarshalJSON(data []byte) error {
	type alias rawAutocompleteInteractionData
	var v struct {
		Options []UnmarshalAutocompleteOption `json:"options"`
		alias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*d = rawAutocompleteInteractionData(v.alias)
	d.Options = make([]AutocompleteOption, len(v.Options))
	for i := range v.Options {
		d.Options[i] = v.Options[i].AutocompleteOption
	}

	return nil
}

type AutocompleteInteractionData struct {
	ID                  snowflake.Snowflake
	Name                string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             map[string]AutocompleteOption
}

func (d *AutocompleteInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawAutocompleteInteractionData

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.ID = iData.ID
	d.Name = iData.Name

	d.Options = make(map[string]AutocompleteOption)
	if len(iData.Options) > 0 {
		flattenedOptions := iData.Options

		unmarshalOption := flattenedOptions[0]
		if option, ok := unmarshalOption.(AutocompleteOptionSubCommandGroup); ok {
			d.SubCommandGroupName = &option.OptionName
			flattenedOptions = make([]AutocompleteOption, len(option.Options))
			for ii := range option.Options {
				flattenedOptions[ii] = option.Options[ii]
			}
			unmarshalOption = option.Options[0]
		}
		if option, ok := unmarshalOption.(AutocompleteOptionSubCommand); ok {
			d.SubCommandName = &option.OptionName
			flattenedOptions = option.Options
		}

		for _, option := range flattenedOptions {
			d.Options[option.Name()] = option
		}
	}
	return nil
}

func (d AutocompleteInteractionData) MarshalJSON() ([]byte, error) {
	options := make([]AutocompleteOption, len(d.Options))
	i := 0
	for _, option := range d.Options {
		options[i] = option
		i++
	}

	if d.SubCommandName != nil {
		options = []AutocompleteOption{
			AutocompleteOptionSubCommand{
				OptionName: *d.SubCommandName,
				Options:    options,
			},
		}
	}
	if d.SubCommandGroupName != nil {
		subCommandOptions := make([]AutocompleteOptionSubCommand, len(options))
		for ii := range options {
			subCommandOptions[ii] = options[ii].(AutocompleteOptionSubCommand)
		}
		options = []AutocompleteOption{
			AutocompleteOptionSubCommandGroup{
				OptionName: *d.SubCommandGroupName,
				Options:    subCommandOptions,
			},
		}
	}

	return json.Marshal(rawAutocompleteInteractionData{
		ID:      d.ID,
		Name:    d.Name,
		Options: options,
	})
}

func (d AutocompleteInteractionData) Option(name string) (AutocompleteOption, bool) {
	option, ok := d.Options[name]
	return option, ok
}

func (d AutocompleteInteractionData) StringOption(name string) (AutocompleteOptionString, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionString)
		return opt, ok
	}
	return AutocompleteOptionString{}, false
}

func (d AutocompleteInteractionData) OptString(name string) (string, bool) {
	if option, ok := d.StringOption(name); ok {
		return option.Value, true
	}
	return "", false
}

func (d AutocompleteInteractionData) String(name string) string {
	if option, ok := d.OptString(name); ok {
		return option
	}
	return ""
}

func (d AutocompleteInteractionData) IntOption(name string) (AutocompleteOptionInt, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionInt)
		return opt, ok
	}
	return AutocompleteOptionInt{}, false
}

func (d AutocompleteInteractionData) OptInt(name string) (int, bool) {
	if option, ok := d.IntOption(name); ok {
		return option.Value, true
	}
	return 0, false
}

func (d AutocompleteInteractionData) Int(name string) int {
	if option, ok := d.OptInt(name); ok {
		return option
	}
	return 0
}

func (d AutocompleteInteractionData) BoolOption(name string) (AutocompleteOptionBool, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionBool)
		return opt, ok
	}
	return AutocompleteOptionBool{}, false
}

func (d AutocompleteInteractionData) OptBool(name string) (bool, bool) {
	if option, ok := d.BoolOption(name); ok {
		return option.Value, true
	}
	return false, false
}

func (d AutocompleteInteractionData) Bool(name string) bool {
	if option, ok := d.OptBool(name); ok {
		return option
	}
	return false
}

func (d AutocompleteInteractionData) UserOption(name string) (AutocompleteOptionUser, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionUser)
		return opt, ok
	}
	return AutocompleteOptionUser{}, false
}

func (d AutocompleteInteractionData) ChannelOption(name string) (AutocompleteOptionChannel, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionChannel)
		return opt, ok
	}
	return AutocompleteOptionChannel{}, false
}

func (d AutocompleteInteractionData) RoleOption(name string) (AutocompleteOptionRole, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionRole)
		return opt, ok
	}
	return AutocompleteOptionRole{}, false
}

func (d AutocompleteInteractionData) MentionableOption(name string) (AutocompleteOptionMentionable, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionMentionable)
		return opt, ok
	}
	return AutocompleteOptionMentionable{}, false
}

func (d AutocompleteInteractionData) OptSnowflake(name string) (snowflake.Snowflake, bool) {
	if option, ok := d.Option(name); ok {
		switch opt := option.(type) {
		case AutocompleteOptionChannel:
			return opt.Value, true
		case AutocompleteOptionRole:
			return opt.Value, true
		case AutocompleteOptionUser:
			return opt.Value, true
		case AutocompleteOptionMentionable:
			return opt.Value, true
		}
	}
	return "", false
}

func (d AutocompleteInteractionData) Snowflake(name string) snowflake.Snowflake {
	if id, ok := d.OptSnowflake(name); ok {
		return id
	}
	return ""
}

func (d AutocompleteInteractionData) FloatOption(name string) (AutocompleteOptionFloat, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(AutocompleteOptionFloat)
		return opt, ok
	}
	return AutocompleteOptionFloat{}, false
}

func (d AutocompleteInteractionData) OptFloat(name string) (float64, bool) {
	if option, ok := d.FloatOption(name); ok {
		return option.Value, true
	}
	return 0, false
}

func (d AutocompleteInteractionData) Float(name string) float64 {
	if option, ok := d.FloatOption(name); ok {
		return option.Value
	}
	return 0
}

func (d AutocompleteInteractionData) All() []AutocompleteOption {
	options := make([]AutocompleteOption, len(d.Options))
	i := 0
	for _, option := range d.Options {
		options[i] = option
		i++
	}
	return options
}

func (d AutocompleteInteractionData) GetByType(optionType ApplicationCommandOptionType) []AutocompleteOption {
	return d.FindAll(func(option AutocompleteOption) bool {
		return option.Type() == optionType
	})
}

func (d AutocompleteInteractionData) Find(optionFindFunc func(option AutocompleteOption) bool) (AutocompleteOption, bool) {
	for _, option := range d.Options {
		if optionFindFunc(option) {
			return option, true
		}
	}
	return nil, false
}

func (d AutocompleteInteractionData) FindAll(optionFindFunc func(option AutocompleteOption) bool) []AutocompleteOption {
	var options []AutocompleteOption
	for _, option := range d.Options {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}
