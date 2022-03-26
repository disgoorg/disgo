package json

import "encoding/json"

var (
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	Indent        = json.Indent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

type (
	RawMessage  = json.RawMessage
	Unmarshaler = json.Unmarshaler
	Marshaler   = json.Marshaler
)
