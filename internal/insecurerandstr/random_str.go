// Package insecurerandstr is not secure
package insecurerandstr

import (
	"math/rand/v2"
	"strings"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr returns a random string of the given length.
func RandStr(n int) string {
	var b strings.Builder
	b.Grow(n)
	for range n {
		b.WriteByte(letters[rand.IntN(len(letters))])
	}
	return b.String()
}
