package src

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DiscoOrg/disgo/src/endpoints"
	"github.com/DiscoOrg/disgo/src/models"
)

func TestRestClient_Request(t *testing.T) {
	restClient := RestClient{
		Client:    &http.Client{},
		UserAgent: "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)",
	}
	response := &models.GatewayBotRs{}
	err := restClient.Request(endpoints.GatewayBot, nil, response)
	assert.NoError(t, err)
	println(response.URL)
}
