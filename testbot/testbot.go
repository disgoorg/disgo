package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chebyrash/promise"
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/api/models"
	"github.com/DiscoOrg/disgo/internal"
	"github.com/DiscoOrg/disgo/internal/events"
)

func main() {
	token := os.Getenv("token")
	options := internal.Options{
		Intents: models.IntentsGuildMessages | models.IntentsGuildMembers,
	}
	dgo := disgo.New(token, options)
	dgo.EventManager().AddEventListeners(&events.ListenerAdapter{
		OnGuildMessageReceived: messageHandler,
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

func messageHandler(event events.GuildMessageReceivedEvent) {
	log.Printf("Message received: %s", event.Message.Content)

	switch event.Message.Content {
	case "ping":
		event.TextChannel.SendMessage("pong")
	case "pong":
		event.TextChannel.SendMessage("ping")
	case "dm":
		event.Message.User.OpenDMChannel().Then(func(channel promise.Any) promise.Any {
			return channel.(*models.DMChannel).SendMessage("helo")
		}).Then(func(_ promise.Any) promise.Any {
			return event.Message.AddReaction("✅")
		}).Catch(func(_ error) error {
			event.Message.AddReaction("❌")
			return nil
		})
	}
}