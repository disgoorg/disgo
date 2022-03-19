package json

import (
	"bytes"
	"encoding/json"
)

var nullBytes = []byte("null")

//goland:noinspection GoUnusedExportedFunction
func NewPtr[T any](t T) *T {
	v := t
	return &v
}

//goland:noinspection GoUnusedExportedFunction
func Null[T any]() Nullable[T] {
	return Nullable[T]{
		isNull: true,
	}
}

//goland:noinspection GoUnusedExportedFunction
func OptionalNull[T any]() *Nullable[T] {
	return &Nullable[T]{
		isNull: true,
	}
}

//goland:noinspection GoUnusedExportedFunction
func New[T any](t T) Nullable[T] {
	return Nullable[T]{
		value:  t,
		isNull: false,
	}
}

//goland:noinspection GoUnusedExportedFunction
func NewOptional[T any](t T) *Nullable[T] {
	return &Nullable[T]{
		value:  t,
		isNull: false,
	}
}

type Nullable[T any] struct {
	value  T
	isNull bool
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.isNull {
		return nullBytes, nil
	}
	return json.Marshal(n.value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		n.isNull = true
		return nil
	}
	return json.Unmarshal(data, &n.value)
}

func (n Nullable[T]) Value() T {
	return n.value
}

func (n Nullable[T]) IsNull() bool {
	return n.isNull
}
