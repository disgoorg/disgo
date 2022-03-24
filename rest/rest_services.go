package rest

import (
	"context"
)

var _ Rest = (*servicesImpl)(nil)

// NewRest returns a new default Rest
func NewRest(restClient Client) Rest {
	return &servicesImpl{
		restClient:           restClient,
		applications:         NewApplications(restClient),
		oauth2:               NewOAuth2(restClient),
		gateway:              NewGateway(restClient),
		guilds:               NewGuilds(restClient),
		members:              NewMembers(restClient),
		channels:             NewChannels(restClient),
		threads:              NewThreads(restClient),
		interactions:         NewInteractions(restClient),
		invites:              NewInvites(restClient),
		guildTemplates:       NewGuildTemplates(restClient),
		users:                NewUsers(restClient),
		voice:                NewVoice(restClient),
		webhooks:             NewWebhooks(restClient),
		stageInstances:       NewStageInstances(restClient),
		emojis:               NewEmojis(restClient),
		stickers:             NewStickers(restClient),
		guildScheduledEvents: NewGuildScheduledEvents(restClient),
	}
}

// Rest is a manager for all of disgo's HTTP requests
type Rest interface {
	RestClient() Client
	Close(ctx context.Context)

	Applications() Applications
	OAuth2() OAuth2
	Gateway() Gateway
	Guilds() Guilds
	Members() Members
	Channels() Channels
	Threads() Threads
	Interactions() Interactions
	Invites() Invites
	GuildTemplates() GuildTemplates
	Users() Users
	Voice() Voice
	Webhooks() Webhooks
	StageInstances() StageInstances
	Emojis() Emojis
	Stickers() Stickers
	GuildScheduledEvents() GuildScheduledEvents
}

type servicesImpl struct {
	restClient Client

	applications         Applications
	oauth2               OAuth2
	gateway              Gateway
	guilds               Guilds
	members              Members
	channels             Channels
	threads              Threads
	interactions         Interactions
	invites              Invites
	guildTemplates       GuildTemplates
	users                Users
	voice                Voice
	webhooks             Webhooks
	stageInstances       StageInstances
	emojis               Emojis
	stickers             Stickers
	guildScheduledEvents GuildScheduledEvents
}

func (s *servicesImpl) RestClient() Client {
	return s.restClient
}

func (s *servicesImpl) Close(ctx context.Context) {
	s.restClient.Close(ctx)
}

func (s *servicesImpl) Applications() Applications {
	return s.applications
}

func (s *servicesImpl) OAuth2() OAuth2 {
	return s.oauth2
}

func (s *servicesImpl) Gateway() Gateway {
	return s.gateway
}

func (s *servicesImpl) Guilds() Guilds {
	return s.guilds
}

func (s *servicesImpl) Members() Members {
	return s.members
}

func (s *servicesImpl) Channels() Channels {
	return s.channels
}

func (s *servicesImpl) Threads() Threads {
	return s.threads
}

func (s *servicesImpl) Interactions() Interactions {
	return s.interactions
}

func (s *servicesImpl) Invites() Invites {
	return s.invites
}

func (s *servicesImpl) GuildTemplates() GuildTemplates {
	return s.guildTemplates
}

func (s *servicesImpl) Users() Users {
	return s.users
}

func (s *servicesImpl) Voice() Voice {
	return s.voice
}

func (s *servicesImpl) Webhooks() Webhooks {
	return s.webhooks
}

func (s *servicesImpl) StageInstances() StageInstances {
	return s.stageInstances
}

func (s *servicesImpl) Emojis() Emojis {
	return s.emojis
}

func (s *servicesImpl) Stickers() Stickers {
	return s.stickers
}

func (s *servicesImpl) GuildScheduledEvents() GuildScheduledEvents {
	return s.guildScheduledEvents
}
