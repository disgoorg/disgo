package disgo

import "github.com/DiscoOrg/disgo/models"

// Options is the configuration used when creating the client
type Options struct {
	Intents     models.Intent
	RestTimeout int
}
