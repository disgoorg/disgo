package gateway

// Intents is an extension of the Bit structure used when identifying with discord
type Intents int64

// Constants for the different bit offsets of Intents
const (
	IntentGuilds Intents = 1 << iota
	IntentGuildMembers
	IntentGuildBans
	IntentGuildEmojisAndStickers
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
	IntentMessageContent
	IntentGuildScheduledEvents
	_
	_
	_
	IntentAutoModerationConfiguration
	IntentAutoModerationExecution

	IntentsGuild = IntentGuilds |
		IntentGuildMembers |
		IntentGuildBans |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildPresences |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentGuildScheduledEvents

	IntentsDirectMessage = IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping

	IntentsNonPrivileged = IntentGuilds |
		IntentGuildBans |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution

	IntentsPrivileged = IntentGuildMembers |
		IntentGuildPresences | IntentMessageContent

	IntentsAll = IntentsNonPrivileged |
		IntentsPrivileged

	IntentsDefault = IntentsNone

	IntentsNone Intents = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (i Intents) Add(bits ...Intents) Intents {
	for _, bit := range bits {
		i |= bit
	}
	return i
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (i Intents) Remove(bits ...Intents) Intents {
	for _, bit := range bits {
		i &^= bit
	}
	return i
}

// Has will ensure that the bit includes all the bits entered
func (i Intents) Has(bits ...Intents) bool {
	for _, bit := range bits {
		if (i & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (i Intents) Missing(bits ...Intents) bool {
	for _, bit := range bits {
		if (i & bit) != bit {
			return true
		}
	}
	return false
}
