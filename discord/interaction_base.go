package discord

import "github.com/DisgoOrg/snowflake"

type BaseInteraction interface {
	ID() snowflake.Snowflake
	ApplicationID() snowflake.Snowflake
	Token() string
	Version() int
	GuildID() *snowflake.Snowflake
	ChannelID() snowflake.Snowflake
	Locale() Locale
	GuildLocale() *Locale
	Member() *ResolvedMember
	User() User
}

type baseInteractionImpl struct {
	id            snowflake.Snowflake
	applicationID snowflake.Snowflake
	token         string
	version       int
	guildID       *snowflake.Snowflake
	channelID     snowflake.Snowflake
	locale        Locale
	guildLocale   *Locale
	member        *ResolvedMember
	user          *User
}

func (i baseInteractionImpl) ID() snowflake.Snowflake {
	return i.id
}
func (i baseInteractionImpl) ApplicationID() snowflake.Snowflake {
	return i.applicationID
}
func (i baseInteractionImpl) Token() string {
	return i.token
}
func (i baseInteractionImpl) Version() int {
	return i.version
}
func (i baseInteractionImpl) GuildID() *snowflake.Snowflake {
	return i.guildID
}
func (i baseInteractionImpl) ChannelID() snowflake.Snowflake {
	return i.channelID
}
func (i baseInteractionImpl) Locale() Locale {
	return i.locale
}
func (i baseInteractionImpl) GuildLocale() *Locale {
	return i.guildLocale
}
func (i baseInteractionImpl) Member() *ResolvedMember {
	return i.member
}
func (i baseInteractionImpl) User() User {
	if i.user != nil {
		return *i.user
	}
	return i.member.User
}
func (baseInteractionImpl) interaction() {}
