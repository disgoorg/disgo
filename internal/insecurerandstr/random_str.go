// Package insecurerandstr is not secure
package insecurerandstr

import (
	"math/rand/v2"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr returns a random string of the given length.
func RandStr(n int) string {
	b := make([]byte, n)
	for i := range n {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}
