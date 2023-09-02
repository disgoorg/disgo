package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type baseInteraction struct {
	id             snowflake.ID
	applicationID  snowflake.ID
	token          string
	version        int
	guildID        *snowflake.ID
	channelID      snowflake.ID
	channel        InteractionChannel
	locale         Locale
	guildLocale    *Locale
	member         *ResolvedMember
	user           *User
	appPermissions *Permissions
}

func (i baseInteraction) ID() snowflake.ID {
	return i.id
}
func (i baseInteraction) ApplicationID() snowflake.ID {
	return i.applicationID
}
func (i baseInteraction) Token() string {
	return i.token
}
func (i baseInteraction) Version() int {
	return i.version
}
func (i baseInteraction) GuildID() *snowflake.ID {
	return i.guildID
}

// Deprecated: Use Channel() instead
func (i baseInteraction) ChannelID() snowflake.ID {
	return i.channelID
}
func (i baseInteraction) Channel() InteractionChannel {
	return i.channel
}
func (i baseInteraction) Locale() Locale {
	return i.locale
}
func (i baseInteraction) GuildLocale() *Locale {
	return i.guildLocale
}
func (i baseInteraction) Member() *ResolvedMember {
	return i.member
}
func (i baseInteraction) User() User {
	if i.user != nil {
		return *i.user
	}
	return i.member.User
}

func (i baseInteraction) AppPermissions() *Permissions {
	return i.appPermissions
}

func (i baseInteraction) CreatedAt() time.Time {
	return i.id.Time()
}
