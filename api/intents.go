package api

// Intents is an extension of the Bit structure used when identifying with discord
type Intents int64

// Add allows you to add multiple bits together, producing a new bit
func (p Intents) Add(bits ...Bit) Bit {
	total := Intents(0)
	for _, bit := range bits {
		total |= bit.(Intents)
	}
	p |= total
	return p
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (p Intents) Remove(bits ...Bit) Bit {
	total := Intents(0)
	for _, bit := range bits {
		total |= bit.(Intents)
	}
	p &^= total
	return p
}

// HasAll will ensure that the bit includes all of the bits entered
func (p Intents) HasAll(bits ...Bit) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (p Intents) Has(bit Bit) bool {
	return (p & bit.(Intents)) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (p Intents) MissingAny(bits ...Bit) bool {
	for _, bit := range bits {
		if !p.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (p Intents) Missing(bit Bit) bool {
	return !p.Has(bit)
}

// Constants for the different bit offsets of intents
const (
	IntentsGuilds Intents = 1 << iota
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
	IntentsNone Intents = 0
)
