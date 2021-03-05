package models

type Intent int
type GatewayIntent struct {
	Bit
}

func newGatewayIntent(bitfield int64) GatewayIntent {
	return GatewayIntent{Bit{Bitfield: 1 << bitfield}}
}

var (
	GatewayIntentsGuilds                 = newGatewayIntent(0)
	GatewayGatewayIntentsGuildMembers    = newGatewayIntent(1)
	GatewayGatewayIntentsGuildBans       = newGatewayIntent(2)
	GatewayIntentsGuildEmojis            = newGatewayIntent(3)
	GatewayIntentsGuildIntegrations      = newGatewayIntent(4)
	GatewayIntentsGuildWebhooks          = newGatewayIntent(5)
	GatewayIntentsGuildInvites           = newGatewayIntent(6)
	GatewayIntentsGuildVoiceStates       = newGatewayIntent(7)
	GatewayIntentsGuildPresences         = newGatewayIntent(8)
	GatewayIntentsGuildMessages          = newGatewayIntent(9)
	GatewayIntentsGuildMessageReactions  = newGatewayIntent(10)
	GatewayIntentsGuildMessageTyping     = newGatewayIntent(11)
	GatewayIntentsDirectMessages         = newGatewayIntent(12)
	GatewayIntentsDirectMessageReactions = newGatewayIntent(13)
	GatewayIntentsDirectMessageTyping    = newGatewayIntent(14)

	GatewayIntentsNone GatewayIntent = GatewayIntent{Bit{Bitfield: 0}}
)

const (
	IntentsGuilds Intent = 1 << iota
	IntentsGuildMembers
	IntentsGuildBans
	IntentsGuildEmojis
	IntentsGuildIntegrations
	IntentsGuildWebhooks
	IntentsGuildInvites
	IntentsGuildVoiceStates
	IntentsGuildPresences
	IntentsGuildMessages
	IntentsGuildMessageReactions
	IntentsGuildMessageTyping
	IntentsDirectMessages
	IntentsDirectMessageReactions
	IntentsDirectMessageTyping

	IntentsAllWithoutPrivileged = IntentsGuilds |
		IntentsGuildBans |
		IntentsGuildEmojis |
		IntentsGuildIntegrations |
		IntentsGuildWebhooks |
		IntentsGuildInvites |
		IntentsGuildVoiceStates |
		IntentsGuildMessages |
		IntentsGuildMessageReactions |
		IntentsGuildMessageTyping |
		IntentsDirectMessages |
		IntentsDirectMessageReactions |
		IntentsDirectMessageTyping
	IntentsAll = IntentsAllWithoutPrivileged |
		IntentsGuildMembers |
		IntentsGuildPresences
	IntentsNone Intent = 0
)
