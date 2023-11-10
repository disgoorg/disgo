// Package insecurerandstr is not secure
package insecurerandstr

import (
	"math/rand"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randStr = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandStr returns a random string of the given length.
func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[randStr.Intn(len(letters))]
	}
	return string(b)
}
