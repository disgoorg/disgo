package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type SelfUser struct {
	discord.SelfUser
	Disgo Disgo
	User  *User
}

// Update updates the SelfUser with the given payload
func (u *SelfUser) Update(updateSelfUser discord.UpdateSelfUser) (*SelfUser, rest.Error) {
	selfUser, err := u.Disgo.RestServices().UserService().UpdateSelfUser(updateSelfUser)
	if err != nil {
		return nil, err
	}
	return u.Disgo.EntityBuilder().CreateSelfUser(*selfUser, CacheStrategyNoWs), nil
}

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u *SelfUser) OpenDMChannel(opts ...rest.RequestOpt) (DMChannel, rest.Error) {
	return nil, rest.NewError(nil, discord.ErrSelfDM)
}
