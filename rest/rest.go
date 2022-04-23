package rest

import (
	"context"
)

var _ Rest = (*servicesImpl)(nil)

// NewRest returns a new default Rest
func NewRest(client Client) Rest {
	return &servicesImpl{
		Client:               client,
		Applications:         NewApplications(client),
		OAuth2:               NewOAuth2(client),
		Gateway:              NewGateway(client),
		Guilds:               NewGuilds(client),
		Members:              NewMembers(client),
		Channels:             NewChannels(client),
		Threads:              NewThreads(client),
		Interactions:         NewInteractions(client),
		Invites:              NewInvites(client),
		GuildTemplates:       NewGuildTemplates(client),
		Users:                NewUsers(client),
		Voice:                NewVoice(client),
		Webhooks:             NewWebhooks(client),
		StageInstances:       NewStageInstances(client),
		Emojis:               NewEmojis(client),
		Stickers:             NewStickers(client),
		GuildScheduledEvents: NewGuildScheduledEvents(client),
	}
}

// Rest is a manager for all of disgo's HTTP requests
type Rest interface {
	Client
	Close(ctx context.Context)

	Applications
	OAuth2
	Gateway
	Guilds
	Members
	Channels
	Threads
	Interactions
	Invites
	GuildTemplates
	Users
	Voice
	Webhooks
	StageInstances
	Emojis
	Stickers
	GuildScheduledEvents
}

type servicesImpl struct {
	Client

	Applications
	OAuth2
	Gateway
	Guilds
	Members
	Channels
	Threads
	Interactions
	Invites
	GuildTemplates
	Users
	Voice
	Webhooks
	StageInstances
	Emojis
	Stickers
	GuildScheduledEvents
}

func (s *servicesImpl) Close(ctx context.Context) {
	s.Client.Close(ctx)
}
