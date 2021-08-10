package gateway

// Intents is an extension of the Bit structure used when identifying with discord
type Intents int64

// Constants for the different bit offsets of Intents
//goland:noinspection GoUnusedConst
const (
	IntentGuilds Intents = 1 << iota
	IntentGuildMembers
	IntentGuildBans
	IntentGuildEmojis
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping

	IntentsNonPrivileged = IntentGuilds |
		IntentGuildBans |
		IntentGuildEmojis |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping

	IntentsPrivileged = IntentGuildMembers |
		IntentGuildPresences

	IntentsAll = IntentsNonPrivileged |
		IntentsPrivileged

	IntentsNone Intents = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (p Intents) Add(bits ...Intents) Intents {
	total := Intents(0)
	for _, bit := range bits {
		total |= bit
	}
	p |= total
	return p
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Intents) Remove(bits ...Intents) Intents {
	total := Intents(0)
	for _, bit := range bits {
		total |= bit
	}
	p &^= total
	return p
}

// HasAll will ensure that the bit includes all of the bits entered
func (p Intents) HasAll(bits ...Intents) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (p Intents) Has(bit Intents) bool {
	return (p & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (p Intents) MissingAny(bits ...Intents) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (p Intents) Missing(bit Intents) bool {
	return !p.Has(bit)
}
