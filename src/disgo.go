package src

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/src/models"
)

// Disgo is the main discord client
type Disgo struct {
	Token      string
	Intents    models.Intent
	RestClient RestClient
	Gateway    Gateway
	SigChannel chan int
}

func (d Disgo) Connect() error {
	err := d.Gateway.Open()
	if err != nil {
		log.Errorf("Unable to connect to gateway. error: %s", err)
		return err
	}
	return nil
}

func New(token string) Disgo {
	disgo := Disgo{
		Token:   token,
		Intents: 0,
		SigChannel: make(chan int),
	}

	disgo.RestClient = RestClient{
		Disgo:     &disgo,
		Client:    &http.Client{},
		UserAgent: "DiscordBot (https://github.com/disgoorg/disgo, 0.0.1)",
	}

	disgo.Gateway = Gateway{
		Disgo: &disgo,
	}
	return disgo
}
