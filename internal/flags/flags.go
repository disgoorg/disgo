package flags

import "golang.org/x/exp/constraints"

// Add allows you to add multiple bits together, producing a new bit
func Add[T constraints.Integer](f T, bits ...T) T {
	for _, bit := range bits {
		f |= bit
	}
	return f
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func Remove[T constraints.Integer](f T, bits ...T) T {
	for _, bit := range bits {
		f &^= bit
	}
	return f
}

// Has will ensure that the bit includes all the bits entered
func Has[T constraints.Integer](f T, bits ...T) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return false
		}
	}
	return true
}

// Missing will check whether the bit is missing any one of the bits
func Missing[T constraints.Integer](f T, bits ...T) bool {
	for _, bit := range bits {
		if (f & bit) != bit {
			return true
		}
	}
	return false
}
