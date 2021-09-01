package oauth2

import "github.com/DisgoOrg/disgo/discord"

var _ EntityBuilder = (*entityBuilderImpl)(nil)

type EntityBuilder interface {
	CreateGuild(guild discord.PartialGuild) *Guild
	CreateUser(user discord.OAuth2User) *User
}

func NewEntityBuilder(client Client) EntityBuilder {
	return &entityBuilderImpl{client: client}
}

type entityBuilderImpl struct {
	client Client
}

func (b *entityBuilderImpl) CreateGuild(guild discord.PartialGuild) *Guild {
	return &Guild{
		PartialGuild: guild,
		Client:       b.client,
	}
}

func (b *entityBuilderImpl) CreateUser(user discord.OAuth2User) *User {
	return &User{
		OAuth2User: user,
		Client:     b.client,
	}
}
