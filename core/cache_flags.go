package core

// CacheFlags are used to enable/disable certain internal caches
type CacheFlags int

// values for CacheFlags
//goland:noinspection GoUnusedConst
const (
	CacheFlagsNone  CacheFlags = 0
	CacheFlagGuilds CacheFlags = 1 << iota

	CacheFlagTextChannels
	CacheFlagDMChannels
	CacheFlagVoiceChannels
	CacheFlagCategories
	CacheFlagNewsChannels
	CacheFlagStoreChannels
	CacheFlagStageChannels

	CacheFlagRoles
	CacheFlagRoleTags

	CacheFlagEmojis
	CacheFlagStickers

	CacheFlagVoiceStates

	CacheFlagStageInstances

	CacheFlagsAllChannels = CacheFlagTextChannels |
		CacheFlagDMChannels |
		CacheFlagVoiceChannels |
		CacheFlagCategories |
		CacheFlagNewsChannels |
		CacheFlagStoreChannels |
		CacheFlagStageChannels

	CacheFlagsDefault = CacheFlagGuilds |
		CacheFlagsAllChannels |
		CacheFlagRoles |
		CacheFlagEmojis |
		CacheFlagStickers |
		CacheFlagVoiceStates

	CacheFlagsFullRoles = CacheFlagRoles |
		CacheFlagRoleTags

	CacheFlagsAll = CacheFlagsAllChannels |
		CacheFlagsFullRoles |
		CacheFlagEmojis |
		CacheFlagStickers |
		CacheFlagVoiceStates |
		CacheFlagStageInstances
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
