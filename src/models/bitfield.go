package models

type Bit int64

func (b Bit) Add(bits ...Bit) Bit {
	total := Bit(0)
	for _, bit := range bits {
		total |= bit
	}
	b |= total
	return b
}

func (b Bit) Remove(bits ...Bit) Bit {
	total := Bit(0)
	for _, bit := range bits {
		total |= bit
	}
	b &^= total
	return b
}

func (b Bit) HasAll(bits ...Bit) bool {
	for _, bit := range bits {
		if !b.Has(bit) {
			return false
		}
	}
	return true
}

func (b Bit) Has(bit Bit) bool {
	return (b & bit) == bit
}

func (b Bit) MissingAny(bits ...Bit) bool {
	for _, bit := range bits {
		if !b.Has(bit) {
			return true
		}
	}
	return false
}

func (b Bit) Missing(bit Bit) bool {
	return !b.Has(bit)
}