package cache

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
		FlagStageInstances
)

// Add allows you to add multiple bits together, producing a new bit
func (f Flags) Add(bits ...Flags) Flags {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f Flags) Remove(bits ...Flags) Flags {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func (f Flags) Has(bits ...Flags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (f Flags) Missing(bits ...Flags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}
