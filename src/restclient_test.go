package src

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DiscoOrg/disgo/src/endpoints"
)

func TestRestClient_Request(t *testing.T) {
	restClient := RestClient{
		Client: &http.Client{},
	}
	response := &struct {
		Url    string `json:"url"`
		Shards int
	}{}
	err := restClient.Request(endpoints.GatewayRoute, nil, response)
	assert.NoError(t, err)
	println(response.Url)
}
