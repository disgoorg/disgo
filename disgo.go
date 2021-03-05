package disgo

import (
	"github.com/Alex-R-31/disgo/src"
	"github.com/Alex-R-31/disgo/src/models"
)

type Disgo struct {
	Token      string
	Intents    models.Intent
	RestClient src.RestClient
}
