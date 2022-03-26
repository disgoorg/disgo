package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	APITestRoute = NewAPIRoute(PUT, "/channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me", "wait")
	CDNTestRoute = NewCDNRoute("/emojis/{emote.id}", PNG, GIF)
)

func TestAPIRoute_Compile(t *testing.T) {
	queryParams := map[string]any{
		"wait": true,
	}
	compiledRoute, err := APITestRoute.Compile(queryParams, "test1", "test2", "test3")
	assert.NoError(t, err)
	assert.Equal(t, API+"/channels/test1/messages/test2/reactions/test3/@me?wait=true", compiledRoute.URL())
}

func TestCDNRoute_Compile(t *testing.T) {
	compiledRoute, err := CDNTestRoute.Compile(nil, PNG, 256, "test1")
	assert.NoError(t, err)
	assert.Equal(t, CDN+"/emojis/test1.png?size=256", compiledRoute.URL())

	compiledRoute, err = CDNTestRoute.Compile(nil, GIF, 512, "test1")
	assert.NoError(t, err)
	assert.Equal(t, CDN+"/emojis/test1.gif?size=512", compiledRoute.URL())
}

func TestCustomRoute_Compile(t *testing.T) {
	testAPI := NewCustomAPIRoute(GET, "https://test.de/", "test/{test}")

	compiledRoute, err := testAPI.Compile(nil, "test")
	assert.NoError(t, err)
	assert.Equal(t, "https://test.de/test/test", compiledRoute.URL())
}
