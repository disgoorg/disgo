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
	token := os.Getenv("token")
	options := disgo.Options{
		Intents: models.IntentsGuildMessages | models.IntentsGuildMembers,
	}
	dgo := disgo.New(token, options)
	dgo.EventManager().AddEventListeners(&listener{})

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

type listener struct {
	disgo.ListenerAdapter
}

func (l listener) OnGenericEvent(event disgo.GenericEvent) {
	println("overridden OnGenericEvent")
}


func (l listener) OnGuildMessageReceived(event disgo.GuildMessageReceivedEvent) {
	println("hm")
	log.Infof("received message in guild: %s, channel: %s, content: %s", event.Guild.ID, event.TextChannel.ID, event.Message.Content)
}