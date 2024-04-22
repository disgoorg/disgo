package gateway

import "github.com/disgoorg/disgo/internal/flags"

// Intents is an extension of the Bit structure used when identifying with discord
type Intents int64

// Constants for the different bit offsets of Intents
const (
	IntentGuilds Intents = 1 << iota
	IntentGuildMembers
	IntentGuildModeration
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
	_
	_
	IntentGuildMessagePolls
	IntentDirectMessagePolls

	IntentsGuild = IntentGuilds |
		IntentGuildMembers |
		IntentGuildModeration |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildPresences |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentGuildScheduledEvents |
		IntentGuildMessagePolls

	IntentsDirectMessage = IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentDirectMessagePolls

	IntentsMessagePolls = IntentGuildMessagePolls |
		IntentDirectMessagePolls

	IntentsNonPrivileged = IntentGuilds |
		IntentGuildModeration |
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
		IntentAutoModerationExecution |
		IntentGuildMessagePolls |
		IntentDirectMessagePolls

	IntentsPrivileged = IntentGuildMembers |
		IntentGuildPresences | IntentMessageContent

	IntentsAll = IntentsNonPrivileged |
		IntentsPrivileged

	IntentsDefault = IntentsNone

	IntentsNone Intents = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (i Intents) Add(bits ...Intents) Intents {
	return flags.Add(i, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (i Intents) Remove(bits ...Intents) Intents {
	return flags.Remove(i, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (i Intents) Has(bits ...Intents) bool {
	return flags.Has(i, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (i Intents) Missing(bits ...Intents) bool {
	return flags.Missing(i, bits...)
}
