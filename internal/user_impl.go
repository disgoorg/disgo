package internal

import (
	"fmt"
	"strings"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
)

var _ api.User = (UserImpl)(nil)
var _ api.Mentionable = (UserImpl)(nil)

type UserImpl struct {
	disgo         api.Disgo
	id            api.Snowflake `json:"id"`
	username      string        `json:"username"`
	discriminator int           `json:"discriminator"`
	avatar        *string       `json:"avatar"`
	bot           bool          `json:"bot"`
	flags         api.UserFlags `json:"public_flags"`
}

func (u UserImpl) Disgo() api.Disgo {
	return u.disgo
}

func (u UserImpl) ID() api.Snowflake {
	return u.id
}

func (u UserImpl) Username() string {
	return u.username
}

func (u UserImpl) Discriminator() int {
	return u.discriminator
}

// Tag returns the user's Username#Discriminator
func (u UserImpl) Tag() string {
	return fmt.Sprintf("%s#%d", u.username, u.discriminator)
}

func (u UserImpl) AvatarURL() *string {
	if u.avatar == nil {
		return nil
	}
	animated := strings.HasPrefix(*u.avatar, "a_")
	format := endpoints.PNG
	if animated {
		format = endpoints.GIF
	}
	a := endpoints.UserAvatar.Compile(format, u.id, *u.avatar).Route()
	return &a
}

func (u UserImpl) EffectiveAvatarURL() string {
	a := u.AvatarURL()
	if a != nil {
		return *a
	}
	return endpoints.DefaultUserAvatar.Compile(endpoints.PNG, u.discriminator%5).Route()
}

func (u UserImpl) Bot() bool {
	return u.bot
}

func (u UserImpl) Flags() api.UserFlags {
	return u.flags
}

// Mention returns the user as a mention
func (u UserImpl) Mention() string {
	return "<@" + u.id.String() + ">"
}

func (u UserImpl) String() string {
	return u.Mention()
}

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u UserImpl) OpenDMChannel() (*api.DMChannel, error) {
	return u.Disgo().RestClient().OpenDMChannel(u.ID())
}
