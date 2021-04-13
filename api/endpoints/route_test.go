package endpoints

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIRoute_Compile(t *testing.T) {
	compiledRoute, err := AddReaction.Compile("test1", "test2", "test3")
	assert.NoError(t, err)
	assert.Equal(t, API+"/channels/test1/messages/test2/reactions/test3/@me", compiledRoute.Route())
}

func TestCDNRoute_Compile(t *testing.T) {
	compiledRoute, err := Emote.Compile(PNG, "test1")
	assert.NoError(t, err)
	assert.Equal(t, CDN+"/emojis/test1.png", compiledRoute.Route())

	compiledRoute, err = Emote.Compile(GIF, "test1")
	assert.NoError(t, err)
	assert.Equal(t, CDN+"/emojis/test1.gif", compiledRoute.Route())
}

func TestCustomRoute_Compile(t *testing.T) {
	testAPI := NewCustomRoute(GET, "https://test.de/{test}")

	compiledRoute, err := testAPI.Compile("test")
	assert.NoError(t, err)
	assert.Equal(t, "https://test.de/test", compiledRoute.Route())
}
