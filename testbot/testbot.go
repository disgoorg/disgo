package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/models"
)

func main() {
	token := "ODE3NDAzMTgyNTI2MzY1NzA2.YEJACQ.43XOAVFCmCtWzLYAHIX6oh13hE0"
	options := disgo.Options{
		Intents: models.IntentsGuildMessages | models.IntentsGuildMembers,
	}
	dgo := disgo.New(token, options)
	dgo.EventManager().AddEventListeners(&disgo.ListenerAdapter{
		OnGuildMessageReceived: func(event disgo.GuildMessageReceivedEvent) {
			log.Printf("Message received: %s", event.Message.Content)
		},
	})

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
