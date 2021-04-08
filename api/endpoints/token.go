package endpoints

import (
	"strings"
)

// Token holds a discord token and is used to keep your logs clean from critical information
type Token string

// MarshalJSON makes sure we don#t send ********* to discords as tokens
func (t Token) MarshalJSON() ([]byte, error) {
	return []byte("\""+t+"\""), nil
}

// UnmarshalJSON makes sure we parse tokens from discord correctly
func (t *Token) UnmarshalJSON(raw []byte) error {
	*t = Token(strings.ReplaceAll(string(raw), "\"", ""))
	return nil
}

// String masks the token
func (t Token) String() string {
	return strings.Repeat("*", len(t))
}

// GoString masks the token
func (t Token) GoString() string {
	return t.String()
}
