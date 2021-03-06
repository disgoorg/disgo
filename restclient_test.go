package disgo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DiscoOrg/disgo/disgo/endpoints"
	"github.com/DiscoOrg/disgo/disgo/models"
)

func TestRestClient_Request(t *testing.T) {
	restClient := RestClientImpl{
		client:    &http.Client{},
	}
	response := &models.GatewayBotRs{}
	err := restClient.Request(endpoints.GatewayBot, nil, response)
	assert.NoError(t, err)
	println(response.URL)
}
