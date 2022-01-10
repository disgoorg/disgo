package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullable_MarshalJSON(t *testing.T) {
	cases := []struct {
		input Nullable[bool]
		data  string
	}{
		{NewNull[bool](), `{"nullable":null,"test":false}`},
		{NewNullable(true), `{"nullable":true,"test":false}`},
		{NewNullable(false), `{"nullable":false,"test":false}`},
	}

	for _, c := range cases {
		type v struct {
			Nullable Nullable[bool] `json:"nullable"`
			Test     bool           `json:"test"`
		}
		data, err := Marshal(v{
			Nullable: c.input,
		})
		assert.NoError(t, err)
		assert.Equal(t, c.data, string(data))
	}
}

func TestOptional_MarshalJSON(t *testing.T) {
	cases := []struct {
		input *bool
		data  string
	}{
		{nil, `{"test":false}`},
		{NewPtr(true), `{"optional":true,"test":false}`},
		{NewPtr(false), `{"optional":false,"test":false}`},
	}

	for _, c := range cases {
		type v struct {
			Optional *bool `json:"optional,omitempty"`
			Test     bool  `json:"test"`
		}
		data, err := Marshal(v{
			Optional: c.input,
		})
		assert.NoError(t, err)
		assert.Equal(t, c.data, string(data))
	}
}

func TestOptionalNullable_MarshalJSON(t *testing.T) {
	cases := []struct {
		input *Nullable[bool]
		data  string
	}{
		{nil, `{"test":false}`},
		{NewNullPtr[bool](), `{"optional_nullable":null,"test":false}`},
		{NewNullablePtr(true), `{"optional_nullable":true,"test":false}`},
		{NewNullablePtr(false), `{"optional_nullable":false,"test":false}`},
	}

	for _, c := range cases {
		type v struct {
			Optional *Nullable[bool] `json:"optional_nullable,omitempty"`
			Test     bool            `json:"test"`
		}
		data, err := Marshal(v{
			Optional: c.input,
		})
		assert.NoError(t, err)
		assert.Equal(t, c.data, string(data))
	}
}

func TestNullable_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		expected Nullable[bool]
		input    string
	}{
		{NewNull[bool](), `{"nullable":null,"test":false}`},
		{NewNullable(true), `{"nullable":true,"test":false}`},
		{NewNullable(false), `{"nullable":false,"test":false}`},
	}

	for _, c := range cases {
		var v struct {
			Nullable Nullable[bool] `json:"nullable"`
		}
		err := Unmarshal([]byte(c.input), &v)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, v.Nullable)
	}
}

func TestOptional_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		expected *bool
		input    string
	}{
		{nil, `{"nullable":null,"test":false}`},
		{NewPtr(true), `{"nullable":true,"test":false}`},
		{NewPtr(false), `{"nullable":false,"test":false}`},
	}

	for _, c := range cases {
		var v struct {
			Nullable *bool `json:"nullable"`
		}
		err := Unmarshal([]byte(c.input), &v)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, v.Nullable)
	}
}
