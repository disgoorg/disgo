package api

// GatewayIntents is an extension of the Bit structure used when identifying with discord
type GatewayIntents int64

// Constants for the different bit offsets of GatewayIntents
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

	GatewayIntentsNone GatewayIntents = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (p GatewayIntents) Add(bits ...GatewayIntents) GatewayIntents {
	total := GatewayIntents(0)
	for _, bit := range bits {
		total |= bit
	}
	p |= total
	return p
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p GatewayIntents) Remove(bits ...GatewayIntents) GatewayIntents {
	total := GatewayIntents(0)
	for _, bit := range bits {
		total |= bit
	}
	p &^= total
	return p
}

// HasAll will ensure that the bit includes all of the bits entered
func (p GatewayIntents) HasAll(bits ...GatewayIntents) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (p GatewayIntents) Has(bit GatewayIntents) bool {
	return (p & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (p GatewayIntents) MissingAny(bits ...GatewayIntents) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (p GatewayIntents) Missing(bit GatewayIntents) bool {
	return !p.Has(bit)
}
