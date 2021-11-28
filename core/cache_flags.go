package core

// CacheFlags are used to enable/disable certain internal caches
type CacheFlags int

// values for CacheFlags
//goland:noinspection GoUnusedConst
const (
	CacheFlagGuilds CacheFlags = 1 << iota

	CacheFlagPresences

	CacheFlagGuildTextChannels
	CacheFlagDMChannels
	CacheFlagGuildVoiceChannels
	CacheFlagGroupDMChannels
	CacheFlagGuildCategories
	CacheFlagGuildNewsChannels
	CacheFlagGuildStoreChannels
	CacheFlagGuildNewsThreads
	CacheFlagGuildPublicThreads
	CacheFlagGuildPrivateThreads
	CacheFlagGuildStageVoiceChannels

	CacheFlagRoles
	CacheFlagRoleTags

	CacheFlagEmojis
	CacheFlagStickers

	CacheFlagVoiceStates

	CacheFlagStageInstances

	CacheFlagsNone CacheFlags = 0

	CacheFlagsAllChannels = CacheFlagGuildTextChannels |
		CacheFlagDMChannels |
		CacheFlagGuildVoiceChannels |
		CacheFlagGroupDMChannels |
		CacheFlagGuildCategories |
		CacheFlagGuildNewsChannels |
		CacheFlagGuildStoreChannels |
		CacheFlagGuildNewsThreads |
		CacheFlagGuildPublicThreads |
		CacheFlagGuildPrivateThreads |
		CacheFlagGuildStageVoiceChannels

	CacheFlagsAllThreads = CacheFlagGuildNewsThreads |
		CacheFlagGuildPublicThreads |
		CacheFlagGuildPrivateThreads

	CacheFlagsDefault = CacheFlagGuilds |
		CacheFlagsAllChannels |
		CacheFlagRoles |
		CacheFlagEmojis |
		CacheFlagStickers |
		CacheFlagVoiceStates

	CacheFlagsFullRoles = CacheFlagRoles |
		CacheFlagRoleTags

	CacheFlagsAll = CacheFlagGuilds |
		CacheFlagsAllChannels |
		CacheFlagsFullRoles |
		CacheFlagEmojis |
		CacheFlagStickers |
		CacheFlagVoiceStates |
		CacheFlagStageInstances |
		CacheFlagPresences
)

// Add allows you to add multiple bits together, producing a new bit
func (f CacheFlags) Add(bits ...CacheFlags) CacheFlags {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f CacheFlags) Remove(bits ...CacheFlags) CacheFlags {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func (f CacheFlags) Has(bits ...CacheFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (f CacheFlags) Missing(bits ...CacheFlags) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}
