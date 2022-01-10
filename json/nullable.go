package json

import (
	"bytes"
)

var nullBytes = []byte("null")

type Nullable[T any] struct {
	value  T
	isNull bool
}

func (n Nullable[T]) Value() T {
	return n.value
}

func (n Nullable[T]) IsNull() bool {
	return n.isNull
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.isNull {
		return nullBytes, nil
	}
	return Marshal(n.value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		n.isNull = true
		return nil
	}
	n.isNull = false
	return Unmarshal(data, &n.value)
}

func NewNullable[T any](t T) Nullable[T] {
	return Nullable[T]{value: t, isNull: false}
}

func NewNullablePtr[T any](t T) *Nullable[T] {
	n := NewNullable(t)
	return &n
}

func NewNull[T any]() Nullable[T] {
	return Nullable[T]{
		isNull: true,
	}
}

func NewNullPtr[T any]() *Nullable[T] {
	n := NewNull[T]()
	return &n
}

func NewPtr[T any](t T) *T {
	return &t
}
