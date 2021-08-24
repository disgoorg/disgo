package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Integration struct {
	discord.Integration
	Disgo       Disgo
	GuildID     discord.Snowflake
	User        *User
	Application *IntegrationApplication
}

// Guild returns the Guild the Integration belongs to
func (i *Integration) Guild() *Guild {
	return i.Disgo.Cache().GuildCache().Get(i.GuildID)
}

// Member returns the Member the Integration uses
func (i *Integration) Member() *Member {
	if i.User == nil {
		return nil
	}
	return i.Disgo.Cache().MemberCache().Get(i.GuildID, i.User.ID)
}

// Role returns the Subscriber Role the Integration uses
func (i *Integration) Role() *Role {
	if i.RoleID == nil {
		return nil
	}
	return i.Disgo.Cache().RoleCache().Get(*i.RoleID)
}

// Delete deletes the Integration from the Guild
func (i *Integration) Delete(opts ...rest.RequestOpt) rest.Error {
	return i.Disgo.RestServices().GuildService().DeleteIntegration(i.GuildID, i.ID, opts...)
}

type IntegrationApplication struct {
	discord.IntegrationApplication
	Bot *User
}
