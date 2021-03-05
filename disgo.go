package disgo

import (
	"github.com/DiscoOrg/disgo/src"
	"github.com/DiscoOrg/disgo/src/models"
)

// Disgo is the main discord client
type Disgo struct {
	Token      string
	Intents    models.Intent
	RestClient src.RestClient
}
