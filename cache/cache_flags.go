package cache

// Types are used to enable/disable certain internal caches
type Types int

// values for CacheTypes
const (
	TypeGuilds Types = 1 << iota
	TypeGuildScheduledEvents
	TypeMembers
	TypeThreadMembers
	TypeMessages
	TypePresences
	TypeChannels
	TypeRoles
	TypeEmojis
	TypeStickers
	TypeVoiceStates
	TypeStageInstances
	TypesNone Types = 0

	TypesAll = TypeGuilds |
		TypeGuildScheduledEvents |
		TypeMembers |
		TypeThreadMembers |
		TypeMessages |
		TypePresences |
		TypeChannels |
		TypeRoles |
		TypeEmojis |
		TypeStickers |
		TypeVoiceStates |
		TypeStageInstances
)

// Add allows you to add multiple bits together, producing a new bit
func (f Types) Add(bits ...Types) Types {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f Types) Remove(bits ...Types) Types {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func (f Types) Has(bits ...Types) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func (f Types) Missing(bits ...Types) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}
