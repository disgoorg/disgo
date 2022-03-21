package discord

import (
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Interaction = (*AutocompleteInteraction)(nil)
)

type AutocompleteInteraction struct {
	baseInteractionImpl
	Data AutocompleteInteractionData `json:"data"`
}

func (AutocompleteInteraction) interaction() {}
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

func (d *AutocompleteInteractionData) MarshalJSON() ([]byte, error) {
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
