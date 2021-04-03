package endpoints

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIRoute_Compile(t *testing.T) {
	assert.Equal(t, API+"/channels/test1/messages/test2/reactions/test3/@me", AddReaction.Compile("test1", "test2", "test3").Route())
}

func TestCDNRoute_Compile(t *testing.T) {
	assert.Equal(t, CDN+"/emojis/test1.png", Emote.Compile(PNG, "test1").Route())
	assert.Equal(t, CDN+"/emojis/test1.gif", Emote.Compile(GIF, "test1").Route())
}

func TestCustomRoute_Compile(t *testing.T) {
	testAPI := NewCustomRoute(GET, "https://test.de/{test}")

	assert.Equal(t, "https://test.de/test", testAPI.Compile("test").Route())
}