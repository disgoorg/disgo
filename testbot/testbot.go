package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/disgo"
	"github.com/DiscoOrg/disgo/disgo/models"
)

func main() {
	token := os.Getenv("token")
	options := disgo.Options{
		Intents: models.IntentsGuildMessages,
	}
	dgo := disgo.New(token, options)

	e := dgo.Connect()
	if e != nil {
		log.Fatal(e)
	}

	defer dgo.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}