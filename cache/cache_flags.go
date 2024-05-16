package cache

import "github.com/disgoorg/disgo/internal/flags"

// Flags are used to enable/disable certain internal caches
type Flags int

// values for CacheFlags
const (
	FlagGuilds Flags = 1 << iota
	FlagGuildScheduledEvents
	FlagMembers
	FlagThreadMembers
	FlagMessages
	FlagPresences
	FlagChannels
	FlagRoles
	FlagEmojis
	FlagStickers
	FlagVoiceStates
	FlagStageInstances
	FlagGuildSoundboardSounds

	FlagsNone Flags = 0
	FlagsAll        = FlagGuilds |
		FlagGuildScheduledEvents |
		FlagMembers |
		FlagThreadMembers |
		FlagMessages |
		FlagPresences |
		FlagChannels |
		FlagRoles |
		FlagEmojis |
		FlagStickers |
		FlagVoiceStates |
		FlagStageInstances |
		FlagGuildSoundboardSounds
)

// Add allows you to add multiple bits together, producing a new bit
func (f Flags) Add(bits ...Flags) Flags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f Flags) Remove(bits ...Flags) Flags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f Flags) Has(bits ...Flags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f Flags) Missing(bits ...Flags) bool {
	return flags.Missing(f, bits...)
}
