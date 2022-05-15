package json

import (
	"bytes"
	"encoding/json"
)

var (
	EmptyBytes = []byte("")
	NullBytes  = []byte("null")
)

func NewPtr[T any](t T) *T {
	return &t
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{
		isNull: true,
	}
}

func OptionalNull[T any]() *Nullable[T] {
	return &Nullable[T]{
		isNull: true,
	}
}

func New[T any](t T) Nullable[T] {
	return Nullable[T]{
		value:  t,
		isNull: false,
	}
}

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
		return NullBytes, nil
	}
	return json.Marshal(n.value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
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
