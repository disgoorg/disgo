package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullBool(t *testing.T) {
	type v struct {
		NullBool NullBool `json:"null_bool"`
	}
	cases := []struct {
		input NullBool
		data  string
	}{
		{*NewNullBool(), `{"null_bool":null}`},
		{*NewBool(true), `{"null_bool":true}`},
		{*NewBool(false), `{"null_bool":false}`},
	}

	for _, c := range cases {
		data, err := Marshal(v{
			NullBool: c.input,
		})
		assert.NoError(t, err)
		assert.Equal(t, c.data, string(data))
	}
}

func TestNullBoolPtr(t *testing.T) {
	type v struct {
		NullBool *NullBool `json:"null_bool,omitempty"`
	}
	cases := []struct {
		input *NullBool
		data  string
	}{
		{nil, `{}`},
		{NewNullBool(), `{"null_bool":null}`},
		{NewBool(true), `{"null_bool":true}`},
		{NewBool(false), `{"null_bool":false}`},
	}

	for _, c := range cases {
		data, err := Marshal(v{
			NullBool: c.input,
		})
		assert.NoError(t, err)
		assert.Equal(t, c.data, string(data))
	}
}
