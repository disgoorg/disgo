package api

// Bit is a utility for interacting with bitfields
type Bit interface {
	Add(bits ...Bit) Bit
	Remove(bits ...Bit) Bit
	HasAll(bits ...Bit) bool
	Has(bit Bit) bool
	MissingAny(bits ...Bit) bool
	Missing(bit Bit) bool
}
