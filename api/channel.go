package api

import (
	"time"

	"github.com/DisgoOrg/restclient"
)

var _ Channel = (*ChannelImpl)(nil)
var _ MessageChannel = (*ChannelImpl)(nil)
var _ GuildChannel = (*ChannelImpl)(nil)
var _ VoiceChannel = (*ChannelImpl)(nil)
var _ DMChannel = (*ChannelImpl)(nil)
var _ TextChannel = (*ChannelImpl)(nil)
var _ Category = (*ChannelImpl)(nil)
var _ StoreChannel = (*ChannelImpl)(nil)

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
const (
	ChannelTypeText ChannelType = iota
	ChannelTypeDM
	ChannelTypeVoice
	ChannelTypeGroupDM
	ChannelTypeCategory
	ChannelTypeNews
	ChannelTypeStore
	ChannelTypeNewsThread
	ChannelTypePublicThread
	ChannelTypePrivateThread
	ChannelTypeStage
)

// ChannelImpl is a generic discord channel object
type ChannelImpl struct {
	Disgo_            Disgo
	ID_               Snowflake       `json:"id"`
	Name_             *string         `json:"name,omitempty"`
	Type_             ChannelType     `json:"type"`
	LastMessageID_    *Snowflake      `json:"last_message_id,omitempty"`
	GuildID_          *Snowflake      `json:"guild_id,omitempty"`
	Position_         *int            `json:"position,omitempty"`
	Topic_            *string         `json:"topic,omitempty"`
	NSFW_             *bool           `json:"nsfw,omitempty"`
	Bitrate_          *int            `json:"bitrate,omitempty"`
	UserLimit_        *int            `json:"user_limit,omitempty"`
	RateLimitPerUser_ *int            `json:"rate_limit_per_user,omitempty"`
	Recipients_       []*User         `json:"recipients,omitempty"`
	Icon_             *string         `json:"icon,omitempty"`
	OwnerID_          *Snowflake      `json:"owner_id,omitempty"`
	ApplicationID_    *Snowflake      `json:"application_id,omitempty"`
	ParentID_         *Snowflake      `json:"parent_id,omitempty"`
	Permissions_      *Permissions    `json:"permissions,omitempty"`
	LastPinTimestamp_ *time.Time      `json:"last_pin_timestamp,omitempty"`
	MessageCount_     int             `json:"message_count"`
	MemberCount_      int             `json:"member_count"`
	ThreadMetadata_   *ThreadMetadata `json:"thread_metadata,omitempty"`
}

type Channel interface {
	Disgo() Disgo
	ID() Snowflake
	Name() string
	Type() ChannelType

	MessageChannel() bool
	GuildChannel() bool
	Thread() bool
	TextChannel() bool
	VoiceChannel() bool
	DMChannel() bool
	Category() bool
	NewsChannel() bool
	StoreChannel() bool
	NewsThread() bool
	PublicThread() bool
	PrivateThread() bool
}

func (c *ChannelImpl) Disgo() Disgo {
	return c.Disgo_
}

func (c *ChannelImpl) ID() Snowflake {
	return c.ID_
}

func (c *ChannelImpl) Name() string {
	return *c.Name_
}

func (c *ChannelImpl) Type() ChannelType {
	return c.Type_
}

func (c *ChannelImpl) MessageChannel() bool {
	return c.TextChannel() || c.NewsChannel() || c.Thread() || c.DMChannel()
}

func (c *ChannelImpl) GuildChannel() bool {
	return c.Category() || c.NewsChannel() || c.TextChannel() || c.VoiceChannel() || c.Thread()
}

func (c *ChannelImpl) Thread() bool {
	return c.NewsThread() || c.PublicThread() || c.PrivateThread()
}

func (c *ChannelImpl) DMChannel() bool {
	return c.Type() != ChannelTypeDM
}

func (c *ChannelImpl) TextChannel() bool {
	return c.Type() != ChannelTypeText
}

func (c *ChannelImpl) VoiceChannel() bool {
	return c.Type() != ChannelTypeVoice
}

func (c *ChannelImpl) Category() bool {
	return c.Type() != ChannelTypeCategory
}

func (c *ChannelImpl) NewsChannel() bool {
	return c.Type() != ChannelTypeNews
}

func (c *ChannelImpl) StoreChannel() bool {
	return c.Type() != ChannelTypeStore
}

func (c *ChannelImpl) NewsThread() bool {
	return c.Type() != ChannelTypeNewsThread
}

func (c *ChannelImpl) PublicThread() bool {
	return c.Type() != ChannelTypePublicThread
}

func (c *ChannelImpl) PrivateThread() bool {
	return c.Type() != ChannelTypePrivateThread
}

func (c *ChannelImpl) Stage() bool {
	return c.Type() != ChannelTypeStage
}

// MessageChannel is used for sending Message(s) to User(s)
type MessageChannel interface {
	Channel
	LastMessageID() *Snowflake
	LastPinTimestamp() *time.Time
	CreateMessage(messageCreate MessageCreate) (*Message, restclient.RestError)
	UpdateMessage(messageID Snowflake, messageUpdate MessageUpdate) (*Message, restclient.RestError)
	DeleteMessage(messageID Snowflake) restclient.RestError
	BulkDeleteMessages(messageIDs ...Snowflake) restclient.RestError
}

func (c *ChannelImpl) LastMessageID() *Snowflake {
	return c.LastMessageID_
}

func (c *ChannelImpl) LastPinTimestamp() *time.Time {
	return c.LastPinTimestamp_
}

// CreateMessage sends a Message to a TextChannel
func (c ChannelImpl) CreateMessage(messageCreate MessageCreate) (*Message, restclient.RestError) {
	// Todo: attachments
	return c.Disgo().RestClient().CreateMessage(c.ID(), messageCreate)
}

// UpdateMessage edits a Message in this TextChannel
func (c ChannelImpl) UpdateMessage(messageID Snowflake, messageUpdate MessageUpdate) (*Message, restclient.RestError) {
	return c.Disgo().RestClient().UpdateMessage(c.ID(), messageID, messageUpdate)
}

// DeleteMessage allows you to edit an existing Message sent by you
func (c ChannelImpl) DeleteMessage(messageID Snowflake) restclient.RestError {
	return c.Disgo().RestClient().DeleteMessage(c.ID(), messageID)
}

// BulkDeleteMessages allows you bulk delete Message(s)
func (c ChannelImpl) BulkDeleteMessages(messageIDs ...Snowflake) restclient.RestError {
	return c.Disgo().RestClient().BulkDeleteMessages(c.ID(), messageIDs...)
}

// DMChannel is used for interacting in private Message(s) with users
type DMChannel interface {
	MessageChannel
}

// GuildChannel is a generic type for all server channels
type GuildChannel interface {
	Channel
	Guild() *Guild
	GuildID() Snowflake
	Permissions() Permissions
	ParentID() *Snowflake
	Parent() Category
	Position() int
}

// Guild returns the channel's Guild
func (c *ChannelImpl) Guild() *Guild {
	return c.Disgo().Cache().Guild(c.GuildID())
}

// GuildID returns the channel's Guild ID
func (c *ChannelImpl) GuildID() Snowflake {
	if !c.GuildChannel() || c.GuildID_ == nil {
		panic("unsupported operation")
	}
	return *c.GuildID_
}

func (c *ChannelImpl) Permissions() Permissions {
	if !c.GuildChannel() || c.Permissions_ == nil {
		panic("unsupported operation")
	}
	return *c.Permissions_
}

func (c *ChannelImpl) ParentID() *Snowflake {
	if !c.GuildChannel() || c.ParentID_ == nil {
		panic("unsupported operation")
	}
	return c.ParentID_
}

func (c *ChannelImpl) Parent() Category {
	if c.ParentID() == nil {
		return nil
	}
	return c.Disgo().Cache().Category(*c.ParentID())
}

func (c *ChannelImpl) Position() int {
	if !c.GuildChannel() || c.Permissions_ == nil {
		panic("unsupported operation")
	}
	return *c.Position_
}

// Category groups text & voice channels in servers together
type Category interface {
	GuildChannel
}

// VoiceChannel adds methods specifically for interacting with discord's voice
type VoiceChannel interface {
	GuildChannel
	Connect() error
	Bitrate() int
}

// Connect sends a api.GatewayCommand to connect to this VoiceChannel
func (c *ChannelImpl) Connect() error {
	return c.Disgo().AudioController().Connect(c.GuildID(), c.ID())
}

func (c *ChannelImpl) Bitrate() int {
	if c.Bitrate_ == nil {
		panic("unsupported operation")
	}
	return *c.Bitrate_
}

// TextChannel allows you to interact with discord's text channels
type TextChannel interface {
	GuildChannel
	MessageChannel
	NSFW() bool
	Topic() *string
}

func (c *ChannelImpl) NSFW() bool {
	if c.NSFW_ == nil {
		panic("unsupported operation")
	}
	return *c.NSFW_
}

func (c *ChannelImpl) Topic() *string {
	return c.Topic_
}

// NewsChannel allows you to interact with discord's text channels
type NewsChannel interface {
	TextChannel
	CrosspostMessage(messageID Snowflake) (*Message, restclient.RestError)
}

// CrosspostMessage crossposts an existing Message
func (c ChannelImpl) CrosspostMessage(messageID Snowflake) (*Message, restclient.RestError) {
	if c.Type() != ChannelTypeNews {
		panic("channel type is not NEWS")
	}
	return c.Disgo().RestClient().CrosspostMessage(c.ID(), messageID)
}

// StoreChannel allows you to interact with discord's store channels
type StoreChannel interface {
	GuildChannel
}

type ChannelCreate struct {
}
type ChannelUpdate struct {
}
