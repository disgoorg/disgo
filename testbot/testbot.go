package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/models"
	"github.com/DiscoOrg/disgo/models/events/guild"
)

func main() {
	token := os.Getenv("token")
	options := disgo.Options{
		Intents: models.IntentsGuildMessages,
	}
	dgo := disgo.New(token, options)
	dgo.AddEventHandlers(onGuildReady)

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

func onGuildReady(event guild.GuildAvailableEvent){
	log.Infof("lol it wörks BRUH guild: %#v", event.Guild())
}