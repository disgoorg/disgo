package endpoints

import (
	"strings"
)

type Token string

func (t Token) MarshalJSON() ([]byte, error) {
	return []byte("\""+t+"\""), nil
}

func (t *Token) UnmarshalJSON(raw []byte) error {
	*t = Token(strings.ReplaceAll(string(raw), "\"", ""))
	return nil
}

func (t Token) String() string {
	return strings.Repeat("*", len(t))
}

func (t Token) GoString() string {
	return t.String()
}
