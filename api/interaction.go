package api

import (
	"encoding/json"
	"errors"
)

// InteractionType is the type of Interaction
type InteractionType int

// Supported InteractionType(s)
const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeCommand
	InteractionTypeComponent
)

// Interaction holds the general parameters of each Interaction
type Interaction struct {
	Disgo           Disgo
	ResponseChannel chan InteractionResponse
	Replied         bool
	ID              Snowflake       `json:"id"`
	Type            InteractionType `json:"type"`
	GuildID         *Snowflake      `json:"guild_id,omitempty"`
	ChannelID       *Snowflake      `json:"channel_id,omitempty"`
	Member          *Member         `json:"member,omitempty"`
	User            *User           `json:"User,omitempty"`
	Token           string          `json:"token"`
	Version         int             `json:"version"`
}

// Reply replies to the api.Interaction with the provided api.InteractionResponse
func (i *Interaction) Reply(response InteractionResponse) error {
	if i.Replied {
		return errors.New("you already replied to this interaction")
	}
	i.Replied = true

	if i.FromWebhook() {
		i.ResponseChannel <- response
		return nil
	}

	return i.Disgo.RestClient().SendInteractionResponse(i.ID, i.Token, response)
}

// FromWebhook returns is the Interaction was made via http
func (i *Interaction) FromWebhook() bool {
	return i.ResponseChannel != nil
}

// Guild returns the api.Guild from the api.Cache
func (i *Interaction) Guild() *Guild {
	if i.GuildID == nil {
		return nil
	}
	return i.Disgo.Cache().Guild(*i.GuildID)
}

// DMChannel returns the api.DMChannel from the api.Cache
func (i *Interaction) DMChannel() *DMChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().DMChannel(*i.ChannelID)
}

// MessageChannel returns the api.MessageChannel from the api.Cache
func (i *Interaction) MessageChannel() *MessageChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().MessageChannel(*i.ChannelID)
}

// TextChannel returns the api.TextChannel from the api.Cache
func (i *Interaction) TextChannel() *TextChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().TextChannel(*i.ChannelID)
}

// GuildChannel returns the api.GuildChannel from the api.Cache
func (i *Interaction) GuildChannel() *GuildChannel {
	if i.ChannelID == nil {
		return nil
	}
	return i.Disgo.Cache().GuildChannel(*i.ChannelID)
}

// EditOriginal edits the original api.InteractionResponse
func (i *Interaction) EditOriginal(followupMessage FollowupMessage) (*Message, error) {
	return i.Disgo.RestClient().EditInteractionResponse(i.Disgo.ApplicationID(), i.Token, followupMessage)
}

// DeleteOriginal deletes the original api.InteractionResponse
func (i *Interaction) DeleteOriginal() error {
	return i.Disgo.RestClient().DeleteInteractionResponse(i.Disgo.ApplicationID(), i.Token)
}

// SendFollowup used to send a api.FollowupMessage to an api.Interaction
func (i *Interaction) SendFollowup(followupMessage FollowupMessage) (*Message, error) {
	return i.Disgo.RestClient().SendFollowupMessage(i.Disgo.ApplicationID(), i.Token, followupMessage)
}

// EditFollowup used to edit a api.FollowupMessage from an api.Interaction
func (i *Interaction) EditFollowup(messageID Snowflake, followupMessage FollowupMessage) (*Message, error) {
	return i.Disgo.RestClient().EditFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID, followupMessage)
}

// DeleteFollowup used to delete a api.FollowupMessage from an api.Interaction
func (i *Interaction) DeleteFollowup(messageID Snowflake) error {
	return i.Disgo.RestClient().DeleteFollowupMessage(i.Disgo.ApplicationID(), i.Token, messageID)
}

// FullInteraction is used for easier unmarshalling of different Interaction(s)
type FullInteraction struct {
	ID          Snowflake       `json:"id"`
	Type        InteractionType `json:"type"`
	GuildID     *Snowflake      `json:"guild_id,omitempty"`
	ChannelID   *Snowflake      `json:"channel_id,omitempty"`
	FullMessage *FullMessage    `json:"message,omitempty"`
	Member      *Member         `json:"member,omitempty"`
	User        *User           `json:"User,omitempty"`
	Token       string          `json:"token"`
	Version     int             `json:"version"`
	Data        json.RawMessage `json:"data,omitempty"`
}

// CommandInteraction is a specific Interaction when using Command(s)
type CommandInteraction struct {
	*Interaction
	Data *CommandInteractionData `json:"data,omitempty"`
}

// DeferReply replies to the api.CommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource and shows a loading state
func (i *CommandInteraction) DeferReply(ephemeral bool) error {
	var data InteractionResponseData
	if ephemeral {
		data = InteractionResponseData{Flags: MessageFlagEphemeral}
	}
	return i.Reply(InteractionResponse{Type: InteractionResponseTypeDeferredChannelMessageWithSource, Data: data})
}

// ReplyCreate replies to the api.CommandInteraction with api.InteractionResponseTypeDeferredChannelMessageWithSource & api.InteractionResponseData
func (i *CommandInteraction) ReplyCreate(data InteractionResponseData) error {
	return i.Reply(InteractionResponse{Type: InteractionResponseTypeChannelMessageWithSource, Data: data})
}

// ButtonInteraction is a specific Interaction when CLicked on Button(s)
type ButtonInteraction struct {
	*Interaction
	Message *Message               `json:"message,omitempty"`
	Data    *ButtonInteractionData `json:"data,omitempty"`
}

// DeferEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeDeferredUpdateMessage and cancels the loading state
func (i *ButtonInteraction) DeferEdit() error {
	return i.Reply(InteractionResponse{Type: InteractionResponseTypeDeferredUpdateMessage})
}

// ReplyEdit replies to the api.ButtonInteraction with api.InteractionResponseTypeUpdateMessage & api.InteractionResponseData which edits the original api.Message
func (i *ButtonInteraction) ReplyEdit(data InteractionResponseData) error {
	return i.Reply(InteractionResponse{Type: InteractionResponseTypeUpdateMessage, Data: data})
}

// CommandInteractionData is the command data payload
type CommandInteractionData struct {
	ID       Snowflake     `json:"id"`
	Name     string        `json:"name"`
	Resolved *Resolved     `json:"resolved,omitempty"`
	Options  []*OptionData `json:"options,omitempty"`
}

// ButtonInteractionData is the command data payload
type ButtonInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
}

// Resolved contains resolved mention data
type Resolved struct {
	Users    map[Snowflake]*User    `json:"users,omitempty"`
	Members  map[Snowflake]*Member  `json:"members,omitempty"`
	Roles    map[Snowflake]*Role    `json:"roles,omitempty"`
	Channels map[Snowflake]*Channel `json:"channels,omitempty"`
}

// OptionData is used for options or subcommands in your slash commands
type OptionData struct {
	Name    string            `json:"name"`
	Type    CommandOptionType `json:"type"`
	Value   interface{}       `json:"value,omitempty"`
	Options []*OptionData     `json:"options,omitempty"`
}

// Option holds info about an Option.Value
type Option struct {
	Resolved *Resolved
	Name     string
	Type     CommandOptionType
	Value    interface{}
}

// String returns the Option.Value as string
func (o Option) String() string {
	return o.Value.(string)
}

// Integer returns the Option.Value as int
func (o Option) Integer() int {
	return o.Value.(int)
}

// Bool returns the Option.Value as bool
func (o Option) Bool() bool {
	return o.Value.(bool)
}

// Snowflake returns the Option.Value as Snowflake
func (o Option) Snowflake() Snowflake {
	return Snowflake(o.String())
}

// User returns the Option.Value as User
func (o Option) User() *User {
	return o.Resolved.Users[o.Snowflake()]
}

// Member returns the Option.Value as Member
func (o Option) Member() *Member {
	return o.Resolved.Members[o.Snowflake()]
}

// Role returns the Option.Value as Role
func (o Option) Role() *Role {
	return o.Resolved.Roles[o.Snowflake()]
}

// Channel returns the Option.Value as Channel
func (o Option) Channel() *Channel {
	return o.Resolved.Channels[o.Snowflake()]
}

// MessageChannel returns the Option.Value as MessageChannel
func (o Option) MessageChannel() *MessageChannel {
	channel := o.Channel()
	if channel == nil || (channel.Type != ChannelTypeText && channel.Type != ChannelTypeNews) {
		return nil
	}
	return &MessageChannel{Channel: *channel}
}

// GuildChannel returns the Option.Value as GuildChannel
func (o Option) GuildChannel() *GuildChannel {
	channel := o.Channel()
	if channel == nil || (channel.Type != ChannelTypeText && channel.Type != ChannelTypeNews && channel.Type != ChannelTypeCategory && channel.Type != ChannelTypeStore && channel.Type != ChannelTypeVoice) {
		return nil
	}
	return &GuildChannel{Channel: *channel}
}

// VoiceChannel returns the Option.Value as VoiceChannel
func (o Option) VoiceChannel() *VoiceChannel {
	channel := o.Channel()
	if channel == nil || channel.Type != ChannelTypeVoice {
		return nil
	}
	return &VoiceChannel{GuildChannel: GuildChannel{Channel: *channel}}
}

// TextChannel returns the Option.Value as TextChannel
func (o Option) TextChannel() *TextChannel {
	channel := o.Channel()
	if channel == nil || (channel.Type != ChannelTypeText && channel.Type != ChannelTypeNews) {
		return nil
	}
	return &TextChannel{GuildChannel: GuildChannel{Channel: *channel}, MessageChannel: MessageChannel{Channel: *channel}}
}

// Category returns the Option.Value as Category
func (o Option) Category() *Category {
	channel := o.Channel()
	if channel == nil || channel.Type != ChannelTypeCategory {
		return nil
	}
	return &Category{GuildChannel: GuildChannel{Channel: *channel}}
}

// StoreChannel returns the Option.Value as StoreChannel
func (o Option) StoreChannel() *StoreChannel {
	channel := o.Channel()
	if channel == nil || channel.Type != ChannelTypeStore {
		return nil
	}
	return &StoreChannel{GuildChannel: GuildChannel{Channel: *channel}}
}
