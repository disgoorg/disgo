package src

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DiscoOrg/disgo/src/endpoints"
)

func TestRestClient_Request(t *testing.T) {
	restClient := RestClient{
		Client:    &http.Client{},
		Token:     "",
		UserAgent: "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)",
	}
	response := &struct {
		Url    string `json:"url"`
		Shards int    `json:"shards"`
	}{}
	err := restClient.Request(endpoints.User, nil, response, "312617227490951168")
	assert.NoError(t, err)
	println(response.Url)
}
