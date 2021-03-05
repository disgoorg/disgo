package models

// Intent is an extension of the Bit structure used when identifying with discord
type Intent Bit

// Bit returns the Bit of the Intent
func (i Intent) Bit() Bit {
	return Bit(i)
}

// Add calls the Bit Add method
func (i Intent) Add(bits ...Bit) Intent {
	return Intent(i.Bit().Add(bits...))
}

// Remove calls the Bit Remove method
func (i Intent) Remove(bits ...Bit) Intent {
	return Intent(i.Bit().Remove(bits...))
}

// HasAll calls the Bit HasAll method
func (i Intent) HasAll(bits ...Bit) bool {
	return i.Bit().HasAll(bits...)
}

// Has calls the Bit Has method
func (i Intent) Has(bit Bit) bool {
	return i.Bit().Has(bit)
}

// MissingAny calls the Bit MissingAny method
func (i Intent) MissingAny(bits ...Bit) bool {
	return i.Bit().MissingAny(bits...)
}

// Missing calls the Bit Missing method
func (i Intent) Missing(bits Bit) bool {
	return i.Bit().Missing(bits)
}

// Constants for the different bit offsets of intents
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
