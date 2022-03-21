package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Interaction = (*ApplicationCommandInteraction)(nil)
)

type ApplicationCommandInteraction struct {
	baseInteractionImpl
	Data ApplicationCommandInteractionData `json:"data"`
}

func (i *ApplicationCommandInteraction) UnmarshalJSON(data []byte) error {
	type applicationCommandInteraction ApplicationCommandInteraction
	var vInteraction struct {
		Data json.RawMessage `json:"data"`
		applicationCommandInteraction
	}

	if err := json.Unmarshal(data, &vInteraction); err != nil {
		return err
	}
	var cType struct {
		Type ApplicationCommandType `json:"type"`
	}
	if err := json.Unmarshal(vInteraction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ApplicationCommandInteractionData
		err             error
	)

	switch cType.Type {
	case ApplicationCommandTypeSlash:
		v := SlashCommandInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeUser:
		v := UserCommandInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeMessage:
		v := MessageCommandInteractionData{}
		err = json.Unmarshal(vInteraction.Data, &v)
		interactionData = v

	default:
		return fmt.Errorf("unkown application rawInteraction data with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	*i = ApplicationCommandInteraction(vInteraction.applicationCommandInteraction)

	i.Data = interactionData
	return nil
}

func (ApplicationCommandInteraction) Type() InteractionType {
	return InteractionTypeApplicationCommand
}

func (i ApplicationCommandInteraction) SlashCommandInteractionData() SlashCommandInteractionData {
	return i.Data.(SlashCommandInteractionData)
}

func (i ApplicationCommandInteraction) UserCommandInteractionData() UserCommandInteractionData {
	return i.Data.(UserCommandInteractionData)
}

func (i ApplicationCommandInteraction) MessageCommandInteractionData() MessageCommandInteractionData {
	return i.Data.(MessageCommandInteractionData)
}

func (ApplicationCommandInteraction) interaction() {}

type ApplicationCommandInteractionData interface {
	Type() ApplicationCommandType
	CommandID() snowflake.Snowflake
	CommandName() string

	applicationCommandInteractionData()
}

type rawSlashCommandInteractionData struct {
	ID       snowflake.Snowflake  `json:"id"`
	Name     string               `json:"name"`
	Resolved SlashCommandResolved `json:"resolved"`
	Options  []SlashCommandOption `json:"options"`
}

func (d *rawSlashCommandInteractionData) UnmarshalJSON(data []byte) error {
	type alias rawSlashCommandInteractionData
	var v struct {
		Options []UnmarshalSlashCommandOption `json:"options"`
		alias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*d = rawSlashCommandInteractionData(v.alias)
	d.Options = make([]SlashCommandOption, len(v.Options))
	for i := range v.Options {
		d.Options[i] = v.Options[i].SlashCommandOption
	}

	return nil
}

type SlashCommandInteractionData struct {
	id                  snowflake.Snowflake
	name                string
	SubCommandName      *string
	SubCommandGroupName *string
	Resolved            SlashCommandResolved
	Options             map[string]SlashCommandOption
}

func (d *SlashCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawSlashCommandInteractionData

	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.Resolved = iData.Resolved

	d.Options = make(map[string]SlashCommandOption)
	if len(iData.Options) > 0 {
		flattenedOptions := iData.Options

		unmarshalOption := flattenedOptions[0]
		if option, ok := unmarshalOption.(SlashCommandOptionSubCommandGroup); ok {
			d.SubCommandGroupName = &option.OptionName
			flattenedOptions = make([]SlashCommandOption, len(option.Options))
			for ii := range option.Options {
				flattenedOptions[ii] = option.Options[ii]
			}
			unmarshalOption = option.Options[0]
		}
		if option, ok := unmarshalOption.(SlashCommandOptionSubCommand); ok {
			d.SubCommandName = &option.OptionName
			flattenedOptions = option.Options
		}

		for _, option := range flattenedOptions {
			d.Options[option.Name()] = option
		}
	}
	return nil
}

func (d SlashCommandInteractionData) MarshalJSON() ([]byte, error) {
	options := make([]SlashCommandOption, len(d.Options))
	i := 0
	for _, option := range d.Options {
		options[i] = option
		i++
	}

	if d.SubCommandName != nil {
		options = []SlashCommandOption{
			SlashCommandOptionSubCommand{
				OptionName: *d.SubCommandName,
				Options:    options,
			},
		}
	}
	if d.SubCommandGroupName != nil {
		subCommandOptions := make([]SlashCommandOptionSubCommand, len(options))
		for ii := range options {
			subCommandOptions[ii] = options[ii].(SlashCommandOptionSubCommand)
		}
		options = []SlashCommandOption{
			SlashCommandOptionSubCommandGroup{
				OptionName: *d.SubCommandGroupName,
				Options:    subCommandOptions,
			},
		}
	}

	return json.Marshal(rawSlashCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Resolved: d.Resolved,
		Options:  options,
	})
}

func (SlashCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (d SlashCommandInteractionData) CommandID() snowflake.Snowflake {
	return d.id
}

func (d SlashCommandInteractionData) CommandName() string {
	return d.name
}

func (d SlashCommandInteractionData) Option(name string) (SlashCommandOption, bool) {
	option, ok := d.Options[name]
	return option, ok
}

func (d SlashCommandInteractionData) StringOption(name string) (SlashCommandOptionString, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionString)
		return opt, ok
	}
	return SlashCommandOptionString{}, false
}

func (d SlashCommandInteractionData) OptString(name string) (string, bool) {
	if option, ok := d.StringOption(name); ok {
		return option.Value, true
	}
	return "", false
}

func (d SlashCommandInteractionData) String(name string) string {
	if option, ok := d.OptString(name); ok {
		return option
	}
	return ""
}

func (d SlashCommandInteractionData) IntOption(name string) (SlashCommandOptionInt, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionInt)
		return opt, ok
	}
	return SlashCommandOptionInt{}, false
}

func (d SlashCommandInteractionData) OptInt(name string) (int, bool) {
	if option, ok := d.IntOption(name); ok {
		return option.Value, true
	}
	return 0, false
}

func (d SlashCommandInteractionData) Int(name string) int {
	if option, ok := d.OptInt(name); ok {
		return option
	}
	return 0
}

func (d SlashCommandInteractionData) BoolOption(name string) (SlashCommandOptionBool, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionBool)
		return opt, ok
	}
	return SlashCommandOptionBool{}, false
}

func (d SlashCommandInteractionData) OptBool(name string) (bool, bool) {
	if option, ok := d.BoolOption(name); ok {
		return option.Value, true
	}
	return false, false
}

func (d SlashCommandInteractionData) Bool(name string) bool {
	if option, ok := d.OptBool(name); ok {
		return option
	}
	return false
}

func (d SlashCommandInteractionData) UserOption(name string) (SlashCommandOptionUser, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionUser)
		return opt, ok
	}
	return SlashCommandOptionUser{}, false
}

func (d SlashCommandInteractionData) OptUser(name string) (User, bool) {
	if option, ok := d.Option(name); ok {
		var (
			userID   snowflake.Snowflake
			userIDOk bool
		)
		switch opt := option.(type) {
		case SlashCommandOptionUser:
			opt, ok := option.(SlashCommandOptionUser)
			userID = opt.Value
			userIDOk = ok
		case SlashCommandOptionMentionable:
			opt, ok := option.(SlashCommandOptionMentionable)
			userID = opt.Value
			userIDOk = ok
		}
		if userIDOk {
			user, userOk := d.Resolved.Users[userID]
			return user, userOk
		}
	}
	return User{}, false
}

func (d SlashCommandInteractionData) User(name string) User {
	if user, ok := d.OptUser(name); ok {
		return user
	}
	return User{}
}

func (d SlashCommandInteractionData) OptMember(name string) (ResolvedMember, bool) {
	if option, ok := d.Option(name); ok {
		var (
			userID   snowflake.Snowflake
			userIDOk bool
		)
		switch opt := option.(type) {
		case SlashCommandOptionUser:
			opt, ok := option.(SlashCommandOptionUser)
			userID = opt.Value
			userIDOk = ok
		case SlashCommandOptionMentionable:
			opt, ok := option.(SlashCommandOptionMentionable)
			userID = opt.Value
			userIDOk = ok
		}
		if userIDOk {
			member, memberOk := d.Resolved.Members[userID]
			return member, memberOk
		}
	}
	return ResolvedMember{}, false
}

func (d SlashCommandInteractionData) Member(name string) ResolvedMember {
	if member, ok := d.OptMember(name); ok {
		return member
	}
	return ResolvedMember{}
}

func (d SlashCommandInteractionData) ChannelOption(name string) (SlashCommandOptionChannel, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionChannel)
		return opt, ok
	}
	return SlashCommandOptionChannel{}, false
}

func (d SlashCommandInteractionData) OptChannel(name string) (ResolvedChannel, bool) {
	if option, ok := d.ChannelOption(name); ok {
		channel, ok := d.Resolved.Channels[option.Value]
		return channel, ok
	}
	return ResolvedChannel{}, false
}

func (d SlashCommandInteractionData) Channel(name string) ResolvedChannel {
	if channel, ok := d.OptChannel(name); ok {
		return channel
	}
	return ResolvedChannel{}
}

func (d SlashCommandInteractionData) RoleOption(name string) (SlashCommandOptionRole, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionRole)
		return opt, ok
	}
	return SlashCommandOptionRole{}, false
}

func (d SlashCommandInteractionData) OptRole(name string) (Role, bool) {
	if option, ok := d.Option(name); ok {
		var (
			roleID   snowflake.Snowflake
			roleIDOk bool
		)
		switch opt := option.(type) {
		case SlashCommandOptionRole:
			opt, ok := option.(SlashCommandOptionRole)
			roleID = opt.Value
			roleIDOk = ok
		case SlashCommandOptionMentionable:
			opt, ok := option.(SlashCommandOptionMentionable)
			roleID = opt.Value
			roleIDOk = ok
		}
		if roleIDOk {
			role, roleOk := d.Resolved.Roles[roleID]
			return role, roleOk
		}
	}
	return Role{}, false
}

func (d SlashCommandInteractionData) Role(name string) Role {
	if role, ok := d.OptRole(name); ok {
		return role
	}
	return Role{}
}

func (d SlashCommandInteractionData) MentionableOption(name string) (SlashCommandOptionMentionable, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionMentionable)
		return opt, ok
	}
	return SlashCommandOptionMentionable{}, false
}

func (d SlashCommandInteractionData) OptSnowflake(name string) (snowflake.Snowflake, bool) {
	if option, ok := d.Option(name); ok {
		switch opt := option.(type) {
		case SlashCommandOptionChannel:
			return opt.Value, true
		case SlashCommandOptionRole:
			return opt.Value, true
		case SlashCommandOptionUser:
			return opt.Value, true
		case SlashCommandOptionMentionable:
			return opt.Value, true
		}
	}
	return "", false
}

func (d SlashCommandInteractionData) Snowflake(name string) snowflake.Snowflake {
	if id, ok := d.OptSnowflake(name); ok {
		return id
	}
	return ""
}

func (d SlashCommandInteractionData) FloatOption(name string) (SlashCommandOptionFloat, bool) {
	if option, ok := d.Option(name); ok {
		opt, ok := option.(SlashCommandOptionFloat)
		return opt, ok
	}
	return SlashCommandOptionFloat{}, false
}

func (d SlashCommandInteractionData) OptFloat(name string) (float64, bool) {
	if option, ok := d.FloatOption(name); ok {
		return option.Value, true
	}
	return 0, false
}

func (d SlashCommandInteractionData) Float(name string) float64 {
	if option, ok := d.FloatOption(name); ok {
		return option.Value
	}
	return 0
}

func (d SlashCommandInteractionData) All() []SlashCommandOption {
	options := make([]SlashCommandOption, len(d.Options))
	i := 0
	for _, option := range d.Options {
		options[i] = option
		i++
	}
	return options
}

func (d SlashCommandInteractionData) GetByType(optionType ApplicationCommandOptionType) []SlashCommandOption {
	return d.FindAll(func(option SlashCommandOption) bool {
		return option.Type() == optionType
	})
}

func (d SlashCommandInteractionData) Find(optionFindFunc func(option SlashCommandOption) bool) (SlashCommandOption, bool) {
	for _, option := range d.Options {
		if optionFindFunc(option) {
			return option, true
		}
	}
	return nil, false
}

func (d SlashCommandInteractionData) FindAll(optionFindFunc func(option SlashCommandOption) bool) []SlashCommandOption {
	var options []SlashCommandOption
	for _, option := range d.Options {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}

func (SlashCommandInteractionData) applicationCommandInteractionData() {}

type SlashCommandResolved struct {
	Users    map[snowflake.Snowflake]User            `json:"users,omitempty"`
	Members  map[snowflake.Snowflake]ResolvedMember  `json:"members,omitempty"`
	Roles    map[snowflake.Snowflake]Role            `json:"roles,omitempty"`
	Channels map[snowflake.Snowflake]ResolvedChannel `json:"channels,omitempty"`
}

type ContextCommandInteractionData interface {
	ApplicationCommandInteractionData
	TargetID() snowflake.Snowflake

	contextCommandInteractionData()
}

var (
	_ ApplicationCommandInteractionData = (*UserCommandInteractionData)(nil)
	_ ContextCommandInteractionData     = (*UserCommandInteractionData)(nil)
)

type rawUserCommandInteractionData struct {
	ID       snowflake.Snowflake `json:"id"`
	Name     string              `json:"name"`
	Resolved UserCommandResolved `json:"resolved"`
	TargetID snowflake.Snowflake `json:"target_id"`
}

type UserCommandInteractionData struct {
	id       snowflake.Snowflake
	name     string
	Resolved UserCommandResolved `json:"resolved"`
	targetID snowflake.Snowflake
}

func (d *UserCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawUserCommandInteractionData
	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.Resolved = iData.Resolved
	d.targetID = iData.TargetID
	return nil
}

func (d *UserCommandInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawUserCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Resolved: d.Resolved,
		TargetID: d.targetID,
	})
}

func (UserCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}
func (d UserCommandInteractionData) CommandID() snowflake.Snowflake {
	return d.id
}
func (d UserCommandInteractionData) CommandName() string {
	return d.name
}
func (d UserCommandInteractionData) TargetID() snowflake.Snowflake {
	return d.targetID
}
func (d UserCommandInteractionData) TargetUser() User {
	return d.Resolved.Users[d.targetID]
}
func (d UserCommandInteractionData) TargetMember() ResolvedMember {
	return d.Resolved.Members[d.targetID]
}

func (UserCommandInteractionData) applicationCommandInteractionData() {}
func (UserCommandInteractionData) contextCommandInteractionData()     {}

type UserCommandResolved struct {
	Users   map[snowflake.Snowflake]User           `json:"users,omitempty"`
	Members map[snowflake.Snowflake]ResolvedMember `json:"members,omitempty"`
}

var (
	_ ApplicationCommandInteractionData = (*MessageCommandInteractionData)(nil)
	_ ContextCommandInteractionData     = (*MessageCommandInteractionData)(nil)
)

type rawMessageCommandInteractionData struct {
	ID       snowflake.Snowflake    `json:"id"`
	Name     string                 `json:"name"`
	Resolved MessageCommandResolved `json:"resolved"`
	TargetID snowflake.Snowflake    `json:"target_id"`
}

type MessageCommandInteractionData struct {
	id       snowflake.Snowflake
	name     string
	Resolved MessageCommandResolved `json:"resolved"`
	targetID snowflake.Snowflake
}

func (d *MessageCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawMessageCommandInteractionData
	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.Resolved = iData.Resolved
	d.targetID = iData.TargetID
	return nil
}

func (d *MessageCommandInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawMessageCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Resolved: d.Resolved,
		TargetID: d.targetID,
	})
}

func (MessageCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}
func (d MessageCommandInteractionData) CommandID() snowflake.Snowflake {
	return d.id
}
func (d MessageCommandInteractionData) CommandName() string {
	return d.name
}
func (d MessageCommandInteractionData) TargetID() snowflake.Snowflake {
	return d.targetID
}
func (d MessageCommandInteractionData) TargetMessage() Message {
	return d.Resolved.Messages[d.targetID]
}

func (MessageCommandInteractionData) applicationCommandInteractionData() {}
func (MessageCommandInteractionData) contextCommandInteractionData()     {}

type MessageCommandResolved struct {
	Messages map[snowflake.Snowflake]Message `json:"messages,omitempty"`
}
