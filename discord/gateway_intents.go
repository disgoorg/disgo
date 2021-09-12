package discord

// GatewayIntents is an extension of the Bit structure used when identifying with discord
type GatewayIntents int64

// Constants for the different bit offsets of GatewayIntents
//goland:noinspection GoUnusedConst
const (
	GatewayIntentGuilds GatewayIntents = 1 << iota
	GatewayIntentGuildMembers
	GatewayIntentGuildBans
	GatewayIntentGuildEmojis
	GatewayIntentGuildIntegrations
	GatewayIntentGuildWebhooks
	GatewayIntentGuildInvites
	GatewayIntentGuildVoiceStates
	GatewayIntentGuildPresences
	GatewayIntentGuildMessages
	GatewayIntentGuildMessageReactions
	GatewayIntentGuildMessageTyping
	GatewayIntentDirectMessages
	GatewayIntentDirectMessageReactions
	GatewayIntentDirectMessageTyping

	GatewayIntentsNonPrivileged = GatewayIntentGuilds |
		GatewayIntentGuildBans |
		GatewayIntentGuildEmojis |
		GatewayIntentGuildIntegrations |
		GatewayIntentGuildWebhooks |
		GatewayIntentGuildInvites |
		GatewayIntentGuildVoiceStates |
		GatewayIntentGuildMessages |
		GatewayIntentGuildMessageReactions |
		GatewayIntentGuildMessageTyping |
		GatewayIntentDirectMessages |
		GatewayIntentDirectMessageReactions |
		GatewayIntentDirectMessageTyping

	GatewayIntentsPrivileged = GatewayIntentGuildMembers |
		GatewayIntentGuildPresences

	GatewayIntentsAll = GatewayIntentsNonPrivileged |
		GatewayIntentsPrivileged

	GatewayIntentsDefault = GatewayIntentsNone

	GatewayIntentsNone GatewayIntents = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (i GatewayIntents) Add(bits ...GatewayIntents) GatewayIntents {
	for _, bit := range bits {
		i |= bit
	}
	return i
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (i GatewayIntents) Remove(bits ...GatewayIntents) GatewayIntents {
	for _, bit := range bits {
		i &^= bit
	}
	return i
}

// Has will ensure that the bit includes all the bits entered
func (i GatewayIntents) Has(bits ...GatewayIntents) bool {
	for _, bit := range bits {
		if (i & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (i GatewayIntents) Missing(bits ...GatewayIntents) bool {
	for _, bit := range bits {
		if (i & bit) != bit {
			return true
		}
	}
	return false
}
