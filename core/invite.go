package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Invite struct {
	discord.Invite
	Disgo      Disgo
	Inviter    *User
	TargetUser *User
}

// URL returns the invite URL in format like https://discord.gg/{code}
func (i *Invite) URL() string {
	compiledRoute, _ := route.InviteURL.Compile(nil, i.Code)
	return compiledRoute.URL()
}

func (i *Invite) Delete(opts ...rest.RequestOpt) (*Invite, rest.Error) {
	invite, err := i.Disgo.RestServices().InviteService().DeleteInvite(i.Code, opts...)
	if err != nil {
		return nil, err
	}
	return i.Disgo.EntityBuilder().CreateInvite(*invite, CacheStrategyNo), nil
}
