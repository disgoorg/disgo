package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type SelfUser struct {
	discord.OAuth2User
	Bot  *Bot
	User *User
}

// Update updates the SelfUser with the given payload
func (u *SelfUser) Update(updateSelfUser discord.SelfUserUpdate, opts ...rest.RequestOpt) (*SelfUser, error) {
	selfUser, err := u.Bot.RestServices.UserService().UpdateSelfUser(updateSelfUser, opts...)
	if err != nil {
		return nil, err
	}
	return u.Bot.EntityBuilder.CreateSelfUser(*selfUser, CacheStrategyNoWs), nil
}

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u *SelfUser) OpenDMChannel(_ ...rest.RequestOpt) (*Channel, error) {
	return nil, discord.ErrSelfDM
}
