package models

// Bit is a utility for interacting with bitfields
type Bit int64

// Add allows you to add multiple bits together, producing a new bit
func (b Bit) Add(bits ...Bit) Bit {
	total := Bit(0)
	for _, bit := range bits {
		total |= bit
	}
	b |= total
	return b
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (b Bit) Remove(bits ...Bit) Bit {
	total := Bit(0)
	for _, bit := range bits {
		total |= bit
	}
	b &^= total
	return b
}

// HasAll will ensure that the bit includes all of the bits entered
func (b Bit) HasAll(bits ...Bit) bool {
	for _, bit := range bits {
		if !b.Has(bit) {
			return false
		}
	}
	return true
}

// Has will check whether the Bit contains another bit
func (b Bit) Has(bit Bit) bool {
	return (b & bit) == bit
}

// MissingAny will check whether the bit is missing any one of the bits
func (b Bit) MissingAny(bits ...Bit) bool {
	for _, bit := range bits {
		if !b.Has(bit) {
			return true
		}
	}
	return false
}

// Missing will do the inverse of Bit.Has
func (b Bit) Missing(bit Bit) bool {
	return !b.Has(bit)
}