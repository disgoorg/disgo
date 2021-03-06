package src

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	token := os.Getenv("token")
	client := New(token)

	e := client.Connect()
	if e != nil {
		log.Fatal(e)
	}

	code := <- client.SigChannel
	log.Infof("exiting with code: %d", code)
}