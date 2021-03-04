package disgo

import (
	"net/http"

	"github.com/Alex-R-31/disgo/src/models"
)

type Disgo struct {
	Token   string
	Intents models.Intent
	Client  *http.Client
}
