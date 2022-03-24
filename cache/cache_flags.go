package cache

// Flags are used to enable/disable certain internal caches
type Flags int

// values for CacheFlags
const (
	FlagGuilds Flags = 1 << iota
	FlagGuildScheduledEvents
	FlagUsers
	FlagMembers
	FlagThreadMembers
	FlagMessages
	FlagPresences

	FlagGuildTextChannels
	FlagDMChannels
	FlagGuildVoiceChannels
	FlagGroupDMChannels
	FlagGuildCategories
	FlagGuildNewsChannels
	FlagGuildNewsThreads
	FlagGuildPublicThreads
	FlagGuildPrivateThreads
	FlagGuildStageVoiceChannels

	FlagRoles
	FlagRoleTags

	FlagEmojis
	FlagStickers

	FlagVoiceStates

	FlagStageInstances

	FlagsNone Flags = 0

	FlagsAllChannels = FlagGuildTextChannels |
		FlagDMChannels |
		FlagGuildVoiceChannels |
		FlagGroupDMChannels |
		FlagGuildCategories |
		FlagGuildNewsChannels |
		FlagGuildNewsThreads |
		FlagGuildPublicThreads |
		FlagGuildPrivateThreads |
		FlagGuildStageVoiceChannels

	FlagsAllThreads = FlagGuildNewsThreads |
		FlagGuildPublicThreads |
		FlagGuildPrivateThreads

	FlagsDefault = FlagGuilds |
		FlagsAllChannels |
		FlagRoles |
		FlagEmojis |
		FlagStickers |
		FlagVoiceStates

	FlagsFullRoles = FlagRoles |
		FlagRoleTags

	FlagsAll = FlagGuilds |
		FlagGuildScheduledEvents |
		FlagsAllChannels |
		FlagsFullRoles |
		FlagEmojis |
		FlagStickers |
		FlagVoiceStates |
		FlagStageInstances |
		FlagPresences
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
