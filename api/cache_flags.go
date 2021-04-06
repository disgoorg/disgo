package api

// CacheFlags are used to enable/disable certain internal caches
type CacheFlags int

// values for CacheFlags
const (
	CacheFlagDMChannels CacheFlags = 0 << iota
	CacheFlagCategories
	CacheFlagTextChannels
	CacheFlagVoiceChannels
	CacheFlagStoreChannels
	CacheFlagRoles
	CacheFlagEmotes
	CacheFlagVoiceState

	CacheFlagsDefault = CacheFlagDMChannels |
		CacheFlagCategories |
		CacheFlagTextChannels |
		CacheFlagVoiceChannels |
		CacheFlagStoreChannels |
		CacheFlagRoles |
		CacheFlagEmotes
)

// Add allows you to add multiple bits together, producing a new bit
func (c CacheFlags) Add(bits ...Bit) Bit {
	total := CacheFlags(0)
	for _, bit := range bits {
		total |= bit.(CacheFlags)
	}
	c |= total
	return c
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (c CacheFlags) Remove(bits ...Bit) Bit {
	total := CacheFlags(0)
	for _, bit := range bits {
		total |= bit.(CacheFlags)
	}
	c &^= total
	return c
}

// HasAll will ensure that the bit includes all of the bits entered
func (c CacheFlags) HasAll(bits ...Bit) bool {
	for _, bit := range bits {
		if !c.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (c CacheFlags) Has(bit Bit) bool {
	return (c & bit.(CacheFlags)) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (c CacheFlags) MissingAny(bits ...Bit) bool {
	for _, bit := range bits {
		if !c.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (c CacheFlags) Missing(bit Bit) bool {
	return !c.Has(bit)
}
