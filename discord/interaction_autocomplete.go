package discord

import (
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

var (
	_ Interaction = (*AutocompleteInteraction)(nil)
)

type AutocompleteInteraction struct {
	baseInteraction
	Data AutocompleteInteractionData `json:"data"`
}

func (i *AutocompleteInteraction) UnmarshalJSON(data []byte) error {
	var interaction struct {
		rawInteraction
		Data AutocompleteInteractionData `json:"data"`
	}
	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}

	i.baseInteraction.id = interaction.ID
	i.baseInteraction.applicationID = interaction.ApplicationID
	i.baseInteraction.token = interaction.Token
	i.baseInteraction.version = interaction.Version
	i.baseInteraction.guildID = interaction.GuildID
	i.baseInteraction.channelID = interaction.ChannelID
	i.baseInteraction.channel = interaction.Channel
	i.baseInteraction.locale = interaction.Locale
	i.baseInteraction.guildLocale = interaction.GuildLocale
	i.baseInteraction.member = interaction.Member
	i.baseInteraction.user = interaction.User
	i.baseInteraction.appPermissions = interaction.AppPermissions

	i.Data = interaction.Data
	return nil
}

func (i AutocompleteInteraction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		rawInteraction
		Data AutocompleteInteractionData `json:"data"`
	}{
		rawInteraction: rawInteraction{
			ID:             i.id,
			Type:           i.Type(),
			ApplicationID:  i.applicationID,
			Token:          i.token,
			Version:        i.version,
			GuildID:        i.guildID,
			ChannelID:      i.channelID,
			Channel:        i.channel,
			Locale:         i.locale,
			GuildLocale:    i.guildLocale,
			Member:         i.member,
			User:           i.user,
			AppPermissions: i.appPermissions,
		},
		Data: i.Data,
	})
}

func (AutocompleteInteraction) Type() InteractionType {
	return InteractionTypeAutocomplete
}

func (AutocompleteInteraction) interaction() {}

type rawAutocompleteInteractionData struct {
	ID      snowflake.ID                 `json:"id"`
	Name    string                       `json:"name"`
	GuildID *snowflake.ID                `json:"guild_id"`
	Options []internalAutocompleteOption `json:"options"`
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
	if len(v.Options) > 0 {
		d.Options = make([]internalAutocompleteOption, len(v.Options))
		for i := range v.Options {
			d.Options[i] = v.Options[i].internalAutocompleteOption
		}
	}

	return nil
}

type AutocompleteInteractionData struct {
	CommandID           snowflake.ID
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	GuildID             *snowflake.ID
	Options             map[string]AutocompleteOption
}

func (d *AutocompleteInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawAutocompleteInteractionData

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.CommandID = iData.ID
	d.CommandName = iData.Name
	d.GuildID = iData.GuildID

	d.Options = make(map[string]AutocompleteOption)
	if len(iData.Options) > 0 {
		flattenedOptions := iData.Options

		unmarshalOption := flattenedOptions[0]
		if option, ok := unmarshalOption.(AutocompleteOptionSubCommandGroup); ok {
			d.SubCommandGroupName = &option.Name
			flattenedOptions = make([]internalAutocompleteOption, len(option.Options))
			for ii := range option.Options {
				flattenedOptions[ii] = option.Options[ii]
			}
			unmarshalOption = option.Options[0]
		}
		if option, ok := unmarshalOption.(AutocompleteOptionSubCommand); ok {
			d.SubCommandName = &option.Name

			flattenedOptions = make([]internalAutocompleteOption, len(option.Options))
			for i := range option.Options {
				flattenedOptions[i] = option.Options[i]
			}
		}

		for _, option := range flattenedOptions {
			d.Options[option.name()] = option.(AutocompleteOption)
		}
	}
	return nil
}

func (d AutocompleteInteractionData) MarshalJSON() ([]byte, error) {
	options := make([]internalAutocompleteOption, 0, len(d.Options))
	for _, option := range d.Options {
		options = append(options, option)
	}

	if d.SubCommandName != nil {
		subCmd := AutocompleteOptionSubCommand{
			Name:    *d.SubCommandName,
			Options: make([]AutocompleteOption, 0, len(options)),
			Type:    ApplicationCommandOptionTypeSubCommand,
		}
		for _, option := range options {
			subCmd.Options = append(subCmd.Options, option.(AutocompleteOption))
		}
		options = []internalAutocompleteOption{subCmd}
	}

	if d.SubCommandGroupName != nil {
		groupCmd := AutocompleteOptionSubCommandGroup{
			Name:    *d.SubCommandGroupName,
			Options: make([]AutocompleteOptionSubCommand, 0, len(options)),
			Type:    ApplicationCommandOptionTypeSubCommandGroup,
		}
		for _, option := range options {
			groupCmd.Options = append(groupCmd.Options, option.(AutocompleteOptionSubCommand))
		}
		options = []internalAutocompleteOption{groupCmd}
	}

	return json.Marshal(rawAutocompleteInteractionData{
		ID:      d.CommandID,
		Name:    d.CommandName,
		GuildID: d.GuildID,
		Options: options,
	})
}

func (d AutocompleteInteractionData) CommandPath() string {
	path := "/" + d.CommandName
	if d.SubCommandGroupName != nil {
		path += "/" + *d.SubCommandGroupName
	}
	if d.SubCommandName != nil {
		path += "/" + *d.SubCommandName
	}
	return path
}

func (d AutocompleteInteractionData) Option(name string) (AutocompleteOption, bool) {
	option, ok := d.Options[name]
	return option, ok
}

func (d AutocompleteInteractionData) OptString(name string) (string, bool) {
	if option, ok := d.Option(name); ok {
		var v string
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return "", false
}

func (d AutocompleteInteractionData) String(name string) string {
	if option, ok := d.OptString(name); ok {
		return option
	}
	return ""
}

func (d AutocompleteInteractionData) OptInt(name string) (int, bool) {
	if option, ok := d.Option(name); ok {
		var v int
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return 0, false
}

func (d AutocompleteInteractionData) Int(name string) int {
	if option, ok := d.OptInt(name); ok {
		return option
	}
	return 0
}

func (d AutocompleteInteractionData) OptBool(name string) (bool, bool) {
	if option, ok := d.Option(name); ok {
		var v bool
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return false, false
}

func (d AutocompleteInteractionData) Bool(name string) bool {
	if option, ok := d.OptBool(name); ok {
		return option
	}
	return false
}

func (d AutocompleteInteractionData) OptSnowflake(name string) (snowflake.ID, bool) {
	if option, ok := d.Option(name); ok {
		var v snowflake.ID
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return 0, false
}

func (d AutocompleteInteractionData) Snowflake(name string) snowflake.ID {
	if id, ok := d.OptSnowflake(name); ok {
		return id
	}
	return 0
}

func (d AutocompleteInteractionData) OptFloat(name string) (float64, bool) {
	if option, ok := d.Option(name); ok {
		var v float64
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return 0, false
}

func (d AutocompleteInteractionData) Float(name string) float64 {
	if float, ok := d.OptFloat(name); ok {
		return float
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
		return option.Type == optionType
	})
}

func (d AutocompleteInteractionData) Find(optionFindFunc func(option AutocompleteOption) bool) (AutocompleteOption, bool) {
	for _, option := range d.Options {
		if optionFindFunc(option) {
			return option, true
		}
	}
	return AutocompleteOption{}, false
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
