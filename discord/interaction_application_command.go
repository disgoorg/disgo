package discord

import (
	"fmt"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

var (
	_ Interaction = (*ApplicationCommandInteraction)(nil)
)

type ApplicationCommandInteraction struct {
	baseInteraction
	Data ApplicationCommandInteractionData `json:"data"`
}

func (i *ApplicationCommandInteraction) UnmarshalJSON(data []byte) error {
	var interaction struct {
		rawInteraction
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}

	var cType struct {
		Type ApplicationCommandType `json:"type"`
	}
	if err := json.Unmarshal(interaction.Data, &cType); err != nil {
		return err
	}

	var (
		interactionData ApplicationCommandInteractionData
		err             error
	)
	switch cType.Type {
	case ApplicationCommandTypeSlash:
		v := SlashCommandInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeUser:
		v := UserCommandInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v

	case ApplicationCommandTypeMessage:
		v := MessageCommandInteractionData{}
		err = json.Unmarshal(interaction.Data, &v)
		interactionData = v
		if interaction.GuildID != nil {
			for id := range v.Resolved.Messages {
				msg := v.Resolved.Messages[id]
				msg.GuildID = interaction.GuildID
				v.Resolved.Messages[id] = msg
			}
		}

	default:
		return fmt.Errorf("unknown application rawInteraction data with type %d received", cType.Type)
	}
	if err != nil {
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

	i.Data = interactionData
	return nil
}

func (i ApplicationCommandInteraction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		rawInteraction
		Data ApplicationCommandInteractionData `json:"data"`
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
	CommandID() snowflake.ID
	CommandName() string
	GuildID() *snowflake.ID

	applicationCommandInteractionData()
}

type rawSlashCommandInteractionData struct {
	ID       snowflake.ID                 `json:"id"`
	Name     string                       `json:"name"`
	Type     ApplicationCommandType       `json:"type"`
	GuildID  *snowflake.ID                `json:"guild_id,omitempty"`
	Resolved SlashCommandResolved         `json:"resolved"`
	Options  []internalSlashCommandOption `json:"options"`
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
	if len(v.Options) > 0 {
		d.Options = make([]internalSlashCommandOption, len(v.Options))
		for i := range v.Options {
			d.Options[i] = v.Options[i].internalSlashCommandOption
		}
	}

	return nil
}

type SlashCommandInteractionData struct {
	id                  snowflake.ID
	name                string
	guildID             *snowflake.ID
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
	d.guildID = iData.GuildID
	d.Resolved = iData.Resolved

	d.Options = make(map[string]SlashCommandOption)
	if len(iData.Options) > 0 {
		flattenedOptions := iData.Options

		unmarshalOption := flattenedOptions[0]
		if option, ok := unmarshalOption.(SlashCommandOptionSubCommandGroup); ok {
			d.SubCommandGroupName = &option.Name
			flattenedOptions = make([]internalSlashCommandOption, len(option.Options))
			for ii := range option.Options {
				flattenedOptions[ii] = option.Options[ii]
			}
			unmarshalOption = option.Options[0]
		}
		if option, ok := unmarshalOption.(SlashCommandOptionSubCommand); ok {
			d.SubCommandName = &option.Name

			flattenedOptions = make([]internalSlashCommandOption, len(option.Options))
			for i := range option.Options {
				flattenedOptions[i] = option.Options[i]
			}
		}

		for _, option := range flattenedOptions {
			d.Options[option.name()] = option.(SlashCommandOption)
		}
	}
	return nil
}

func (d SlashCommandInteractionData) MarshalJSON() ([]byte, error) {
	options := make([]internalSlashCommandOption, 0, len(d.Options))
	for _, option := range d.Options {
		options = append(options, option)
	}

	if d.SubCommandName != nil {
		subCmd := SlashCommandOptionSubCommand{
			Name:    *d.SubCommandName,
			Options: make([]SlashCommandOption, 0, len(options)),
			Type:    ApplicationCommandOptionTypeSubCommand,
		}
		for _, option := range options {
			subCmd.Options = append(subCmd.Options, option.(SlashCommandOption))
		}
		options = []internalSlashCommandOption{subCmd}
	}

	if d.SubCommandGroupName != nil {
		groupCmd := SlashCommandOptionSubCommandGroup{
			Name:    *d.SubCommandGroupName,
			Options: make([]SlashCommandOptionSubCommand, 0, len(options)),
			Type:    ApplicationCommandOptionTypeSubCommandGroup,
		}
		for _, option := range options {
			groupCmd.Options = append(groupCmd.Options, option.(SlashCommandOptionSubCommand))
		}
		options = []internalSlashCommandOption{groupCmd}
	}

	return json.Marshal(rawSlashCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Type:     d.Type(),
		GuildID:  d.guildID,
		Resolved: d.Resolved,
		Options:  options,
	})
}

func (SlashCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeSlash
}

func (d SlashCommandInteractionData) CommandID() snowflake.ID {
	return d.id
}

func (d SlashCommandInteractionData) CommandName() string {
	return d.name
}

func (d SlashCommandInteractionData) CommandPath() string {
	path := "/" + d.name
	if d.SubCommandGroupName != nil {
		path += "/" + *d.SubCommandGroupName
	}
	if d.SubCommandName != nil {
		path += "/" + *d.SubCommandName
	}
	return path
}

func (d SlashCommandInteractionData) GuildID() *snowflake.ID {
	return d.guildID
}

func (d SlashCommandInteractionData) Option(name string) (SlashCommandOption, bool) {
	option, ok := d.Options[name]
	return option, ok
}

func (d SlashCommandInteractionData) OptString(name string) (string, bool) {
	if option, ok := d.Option(name); ok {
		var v string
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return "", false
}

func (d SlashCommandInteractionData) String(name string) string {
	if option, ok := d.OptString(name); ok {
		return option
	}
	return ""
}

func (d SlashCommandInteractionData) OptInt(name string) (int, bool) {
	if option, ok := d.Option(name); ok {
		var v int
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return 0, false
}

func (d SlashCommandInteractionData) Int(name string) int {
	if option, ok := d.OptInt(name); ok {
		return option
	}
	return 0
}

func (d SlashCommandInteractionData) OptBool(name string) (bool, bool) {
	if option, ok := d.Option(name); ok {
		var v bool
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return false, false
}

func (d SlashCommandInteractionData) Bool(name string) bool {
	if option, ok := d.OptBool(name); ok {
		return option
	}
	return false
}

func (d SlashCommandInteractionData) OptUser(name string) (User, bool) {
	if option, ok := d.Option(name); ok {
		var userID snowflake.ID
		if err := json.Unmarshal(option.Value, &userID); err == nil {
			user, ok := d.Resolved.Users[userID]
			return user, ok
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
		var userID snowflake.ID
		if err := json.Unmarshal(option.Value, &userID); err == nil {
			user, ok := d.Resolved.Members[userID]
			return user, ok
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

func (d SlashCommandInteractionData) OptChannel(name string) (ResolvedChannel, bool) {
	if option, ok := d.Option(name); ok {
		var channelID snowflake.ID
		if err := json.Unmarshal(option.Value, &channelID); err == nil {
			channel, ok := d.Resolved.Channels[channelID]
			return channel, ok
		}
	}
	return ResolvedChannel{}, false
}

func (d SlashCommandInteractionData) Channel(name string) ResolvedChannel {
	if channel, ok := d.OptChannel(name); ok {
		return channel
	}
	return ResolvedChannel{}
}

func (d SlashCommandInteractionData) OptRole(name string) (Role, bool) {
	if option, ok := d.Option(name); ok {
		var roleID snowflake.ID
		if err := json.Unmarshal(option.Value, &roleID); err == nil {
			role, ok := d.Resolved.Roles[roleID]
			return role, ok
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

func (d SlashCommandInteractionData) OptSnowflake(name string) (snowflake.ID, bool) {
	if option, ok := d.Option(name); ok {
		var id snowflake.ID
		if err := json.Unmarshal(option.Value, &id); err == nil {
			return id, ok
		}
	}
	return 0, false
}

func (d SlashCommandInteractionData) Snowflake(name string) snowflake.ID {
	if id, ok := d.OptSnowflake(name); ok {
		return id
	}
	return 0
}

func (d SlashCommandInteractionData) OptFloat(name string) (float64, bool) {
	if option, ok := d.Option(name); ok {
		var v float64
		if err := json.Unmarshal(option.Value, &v); err == nil {
			return v, true
		}
	}
	return 0, false
}

func (d SlashCommandInteractionData) Float(name string) float64 {
	if value, ok := d.OptFloat(name); ok {
		return value
	}
	return 0
}

func (d SlashCommandInteractionData) OptAttachment(name string) (Attachment, bool) {
	if option, ok := d.Option(name); ok {
		var v snowflake.ID
		if err := json.Unmarshal(option.Value, &v); err == nil {
			attachment, ok := d.Resolved.Attachments[v]
			return attachment, ok
		}
	}
	return Attachment{}, false
}

func (d SlashCommandInteractionData) Attachment(name string) Attachment {
	if attachment, ok := d.OptAttachment(name); ok {
		return attachment
	}
	return Attachment{}
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
		return option.Type == optionType
	})
}

func (d SlashCommandInteractionData) Find(optionFindFunc func(option SlashCommandOption) bool) (SlashCommandOption, bool) {
	for _, option := range d.Options {
		if optionFindFunc(option) {
			return option, true
		}
	}
	return SlashCommandOption{}, false
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
	Users       map[snowflake.ID]User            `json:"users,omitempty"`
	Members     map[snowflake.ID]ResolvedMember  `json:"members,omitempty"`
	Roles       map[snowflake.ID]Role            `json:"roles,omitempty"`
	Channels    map[snowflake.ID]ResolvedChannel `json:"channels,omitempty"`
	Attachments map[snowflake.ID]Attachment      `json:"attachments,omitempty"`
}

func (r *SlashCommandResolved) UnmarshalJSON(data []byte) error {
	type slashCommandResolved SlashCommandResolved
	var v slashCommandResolved
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = SlashCommandResolved(v)
	for id, member := range r.Members {
		if user, ok := r.Users[id]; ok {
			member.User = user
			r.Members[id] = member
		}
	}
	return nil
}

type ContextCommandInteractionData interface {
	ApplicationCommandInteractionData
	TargetID() snowflake.ID

	contextCommandInteractionData()
}

var (
	_ ApplicationCommandInteractionData = (*UserCommandInteractionData)(nil)
	_ ContextCommandInteractionData     = (*UserCommandInteractionData)(nil)
)

type rawUserCommandInteractionData struct {
	ID       snowflake.ID           `json:"id"`
	Name     string                 `json:"name"`
	Type     ApplicationCommandType `json:"type"`
	GuildID  *snowflake.ID          `json:"guild_id,omitempty"`
	Resolved UserCommandResolved    `json:"resolved"`
	TargetID snowflake.ID           `json:"target_id"`
}

type UserCommandInteractionData struct {
	id       snowflake.ID
	name     string
	guildID  *snowflake.ID
	Resolved UserCommandResolved `json:"resolved"`
	targetID snowflake.ID
}

func (d *UserCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawUserCommandInteractionData
	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.guildID = iData.GuildID
	d.Resolved = iData.Resolved
	d.targetID = iData.TargetID
	return nil
}

func (d *UserCommandInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawUserCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Type:     d.Type(),
		GuildID:  d.guildID,
		Resolved: d.Resolved,
		TargetID: d.targetID,
	})
}

func (UserCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeUser
}

func (d UserCommandInteractionData) CommandID() snowflake.ID {
	return d.id
}

func (d UserCommandInteractionData) CommandName() string {
	return d.name
}

func (d UserCommandInteractionData) GuildID() *snowflake.ID {
	return d.guildID
}

func (d UserCommandInteractionData) TargetID() snowflake.ID {
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
	Users   map[snowflake.ID]User           `json:"users,omitempty"`
	Members map[snowflake.ID]ResolvedMember `json:"members,omitempty"`
}

func (r *UserCommandResolved) UnmarshalJSON(data []byte) error {
	type userCommandResolved UserCommandResolved
	var v userCommandResolved
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = UserCommandResolved(v)
	for id, member := range r.Members {
		if user, ok := r.Users[id]; ok {
			member.User = user
			r.Members[id] = member
		}
	}
	return nil
}

var (
	_ ApplicationCommandInteractionData = (*MessageCommandInteractionData)(nil)
	_ ContextCommandInteractionData     = (*MessageCommandInteractionData)(nil)
)

type rawMessageCommandInteractionData struct {
	ID       snowflake.ID           `json:"id"`
	Name     string                 `json:"name"`
	Type     ApplicationCommandType `json:"type"`
	GuildID  *snowflake.ID          `json:"guild_id,omitempty"`
	Resolved MessageCommandResolved `json:"resolved"`
	TargetID snowflake.ID           `json:"target_id"`
}

type MessageCommandInteractionData struct {
	id       snowflake.ID
	name     string
	guildID  *snowflake.ID
	Resolved MessageCommandResolved `json:"resolved"`
	targetID snowflake.ID
}

func (d *MessageCommandInteractionData) UnmarshalJSON(data []byte) error {
	var iData rawMessageCommandInteractionData
	if err := json.Unmarshal(data, &iData); err != nil {
		return err
	}
	d.id = iData.ID
	d.name = iData.Name
	d.guildID = iData.GuildID
	d.Resolved = iData.Resolved
	d.targetID = iData.TargetID
	return nil
}

func (d *MessageCommandInteractionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawMessageCommandInteractionData{
		ID:       d.id,
		Name:     d.name,
		Type:     d.Type(),
		GuildID:  d.guildID,
		Resolved: d.Resolved,
		TargetID: d.targetID,
	})
}

func (MessageCommandInteractionData) Type() ApplicationCommandType {
	return ApplicationCommandTypeMessage
}

func (d MessageCommandInteractionData) CommandID() snowflake.ID {
	return d.id
}

func (d MessageCommandInteractionData) CommandName() string {
	return d.name
}

func (d MessageCommandInteractionData) GuildID() *snowflake.ID {
	return d.guildID
}

func (d MessageCommandInteractionData) TargetID() snowflake.ID {
	return d.targetID
}

func (d MessageCommandInteractionData) TargetMessage() Message {
	return d.Resolved.Messages[d.targetID]
}

func (MessageCommandInteractionData) applicationCommandInteractionData() {}
func (MessageCommandInteractionData) contextCommandInteractionData()     {}

type MessageCommandResolved struct {
	Messages map[snowflake.ID]Message `json:"messages,omitempty"`
}
