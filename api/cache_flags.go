package api

// CacheFlags are used to enable/disable certain internal caches
type CacheFlags int

// values for CacheFlags
const (
	CacheFlagsNone      CacheFlags = 0
	CacheFlagDMChannels CacheFlags = 1 << iota
	CacheFlagCategories
	CacheFlagTextChannels
	CacheFlagThreads
	CacheFlagVoiceChannels
	CacheFlagStoreChannels
	CacheFlagRoles
	CacheFlagRoleTags
	CacheFlagEmotes
	CacheFlagVoiceState
	CacheFlagCommands
	CacheFlagCommandPermissions

	CacheFlagsChannels = CacheFlagDMChannels |
		CacheFlagCategories |
		CacheFlagTextChannels |
		CacheFlagThreads |
		CacheFlagVoiceChannels |
		CacheFlagStoreChannels

	CacheFlagsDefault = CacheFlagsChannels |
		CacheFlagRoles |
		CacheFlagEmotes

	CacheFlagsFullRoles = CacheFlagRoles |
		CacheFlagRoleTags

	CacheFlagsFullCommands = CacheFlagCommands |
		CacheFlagCommandPermissions

	CacheFlagsAll = CacheFlagsChannels |
		CacheFlagsFullRoles |
		CacheFlagEmotes |
		CacheFlagVoiceState |
		CacheFlagsFullCommands
)

// Add allows you to add multiple bits together, producing a new bit
func (c CacheFlags) Add(bits ...CacheFlags) CacheFlags {
	total := CacheFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	c |= total
	return c
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (c CacheFlags) Remove(bits ...CacheFlags) CacheFlags {
	total := CacheFlags(0)
	for _, bit := range bits {
		total |= bit
	}
	c &^= total
	return c
}

// HasAll will ensure that the bit includes all of the bits entered
func (c CacheFlags) HasAll(bits ...CacheFlags) bool {
	for _, bit := range bits {
		if !c.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (c CacheFlags) Has(bit CacheFlags) bool {
	return (c & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (c CacheFlags) MissingAny(bits ...CacheFlags) bool {
	for _, bit := range bits {
		if !c.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (c CacheFlags) Missing(bit CacheFlags) bool {
	return !c.Has(bit)
}
