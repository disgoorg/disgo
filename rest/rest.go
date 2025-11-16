package rest

// Rest is a manager for all of disgo's HTTP requests
type Rest interface {
	Client

	Applications
	OAuth2
	Gateway
	Guilds
	AutoModeration
	Members
	Channels
	Threads
	Interactions
	Invites
	GuildTemplates
	Users
	Voice
	Webhooks
	SoundboardSounds
	StageInstances
	Emojis
	Stickers
	SKUs
	GuildScheduledEvents
}

var _ Rest = (*restImpl)(nil)

// New returns a new default Rest
func New(client Client, opts ...ConfigOpt) Rest {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &restImpl{
		Client:               client,
		Applications:         NewApplications(client),
		OAuth2:               NewOAuth2(client),
		Gateway:              NewGateway(client),
		Guilds:               NewGuilds(client),
		AutoModeration:       NewAutoModeration(client),
		Members:              NewMembers(client),
		Channels:             NewChannels(client, cfg.DefaultAllowedMentions),
		Threads:              NewThreads(client),
		Interactions:         NewInteractions(client, cfg.DefaultAllowedMentions),
		Invites:              NewInvites(client),
		GuildTemplates:       NewGuildTemplates(client),
		Users:                NewUsers(client),
		Voice:                NewVoice(client),
		Webhooks:             NewWebhooks(client, cfg.DefaultAllowedMentions),
		SoundboardSounds:     NewSoundboardSounds(client),
		StageInstances:       NewStageInstances(client),
		Emojis:               NewEmojis(client),
		Stickers:             NewStickers(client),
		SKUs:                 NewSKUs(client),
		GuildScheduledEvents: NewGuildScheduledEvents(client),
	}
}

type restImpl struct {
	Client

	Applications
	OAuth2
	Gateway
	Guilds
	AutoModeration
	Members
	Channels
	Threads
	Interactions
	Invites
	GuildTemplates
	Users
	Voice
	Webhooks
	SoundboardSounds
	StageInstances
	Emojis
	Stickers
	SKUs
	GuildScheduledEvents
}
