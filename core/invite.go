package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Invite struct {
	discord.Invite
	Bot        Bot
	Inviter    *User
	TargetUser *User
}

// URL returns the invite URL in format like https://discord.gg/{code}
func (i *Invite) URL() string {
	compiledRoute, _ := route.InviteURL.Compile(nil, i.Code)
	return compiledRoute.URL()
}

// Delete deletes this Invite
func (i *Invite) Delete(opts ...rest.RequestOpt) (*Invite, error) {
	invite, err := i.Bot.RestServices().InviteService().DeleteInvite(i.Code, opts...)
	if err != nil {
		return nil, err
	}
	return i.Bot.EntityBuilder().CreateInvite(*invite, CacheStrategyNo), nil
}
