package models

type Intent Bit

func (i Intent) Bit() Bit {
	return Bit(i)
}

func (i Intent) Add(bits ...Bit) Intent {
	return Intent(i.Bit().Add(bits...))
}

func (i Intent) Remove(bits ...Bit) Intent {
	return Intent(i.Bit().Remove(bits...))
}

func (i Intent) HasAll(bits ...Bit) bool {
	return i.Bit().HasAll(bits...)
}

func (i Intent) Has(bit Bit) bool {
	return i.Bit().Has(bit)
}

func (i Intent) MissingAny(bits ...Bit) bool {
	return i.Bit().MissingAny(bits...)
}

func (i Intent) Missing(bits Bit) bool {
	return i.Bit().Missing(bits)
}

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
