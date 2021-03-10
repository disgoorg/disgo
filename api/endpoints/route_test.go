package endpoints

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIRoute_Compile(t *testing.T) {
	assert.Equal(t, API+ "channels/test1/messages/test2/reactions/test3/@me", PutReaction.Compile("test1", "test2", "test3"))
}

func TestCDNRoute_Compile(t *testing.T) {
	assert.Equal(t, CDN+ "emojis/test1.png", Emote.Compile(PNG, "test1"))
	assert.Equal(t, CDN+ "emojis/test1.gif", Emote.Compile(GIF, "test1"))
}