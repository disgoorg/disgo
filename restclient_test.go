package disgo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DiscoOrg/disgo/endpoints"
	"github.com/DiscoOrg/disgo/models"
)

func TestRestClient_Request(t *testing.T) {
	token := os.Getenv("token")
	dgo := New(token, Options{})

	restClient := dgo.RestClient()
	response := &models.GatewayBotRs{}
	err := restClient.Request(endpoints.GatewayBot, nil, response)
	assert.NoError(t, err)
	println(response.URL)
}
