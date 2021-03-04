package models

type Bit struct {
	Bitfield int64
}

func (b Bit) Add(bits ...interface{}) Bit {
	total := 0
	for bit := range bits {
		total |= bit
	}
	b.Bitfield |= int64(total)
	return b
}

func (b Bit) Remove(bits ...interface{}) Bit {
	total := 0
	for bit := range bits {
		total |= bit
	}
	b.Bitfield &^= int64(total)
	return b
}

func (b Bit) HasAll(bits ...interface{}) bool {
	for bit := range bits {
		if !b.Has(int64(bit)) {
			return false
		}
	}
	return true
}

func (b Bit) Has(bit int64) bool {
	return (b.Bitfield & bit) == bit
}

func (b Bit) MissingAny(bits ...interface{}) bool {
	hasAll := true
	for bit := range bits {
		if !b.Has(int64(bit)) {
			hasAll = false
			break
		}
	}
	return hasAll
}

func (b Bit) Missing(bit int64) bool {
	return !b.Has(bit)
}
