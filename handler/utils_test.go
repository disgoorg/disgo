package handler

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindHandler(t *testing.T) {
	handlers := []*handlerHolder[CommandHandler]{
		{path: "/test"},
		{path: "/test/{id}"},
		{path: "/test/{id}/test"},
	}

	testCases := []struct {
		path   string
		ok     bool
		values map[string]string
	}{
		{path: "/test", ok: true, values: map[string]string{}},
		{path: "/test/123", ok: true, values: map[string]string{"id": "123"}},
		{path: "/test/123/test", ok: true, values: map[string]string{"id": "123"}},
		{path: "/test/123/test/123", ok: false, values: nil},
		{path: "/bla", ok: false, values: nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.path[1:], func(t *testing.T) {
			_, values, ok := findHandler(handlers, testCase.path, CommandDelimiter)

			assert.Equal(t, testCase.ok, ok)
			assert.True(t, reflect.DeepEqual(values, testCase.values))
		})
	}
}
