package src

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DiscoOrg/disgo/src/endpoints"
)

func TestRestClient_Request(t *testing.T) {
	restClient := RestClient{
		Client:    &http.Client{},
		Token:     os.Getenv("token"),
		UserAgent: "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)",
	}
	response := &struct {
		URL    string `json:"url"`
		Shards int    `json:"shards"`
	}{}
	err := restClient.Request(endpoints.GatewayBot, nil, response)
	assert.NoError(t, err)
	println(response.URL)
}
