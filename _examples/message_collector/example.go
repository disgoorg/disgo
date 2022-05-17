package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", disgo.Version)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
	if event.Message.Author.Bot || event.Message.Author.System {
		return
	}
	if event.Message.Content == "start" {
		go func() {
			ch, cls := bot.NewEventCollector(event.Client(), func(event *events.MessageCreateEvent) bool {
				return event.ChannelID == event.ChannelID && event.Message.Author.ID == event.Message.Author.ID && event.Message.Content != ""
			})
			i := 1
			str := ">>> "
			ctx, clsCtx := context.WithTimeout(context.Background(), 20*time.Second)
			defer clsCtx()
			for {
				select {
				case <-ctx.Done():
					_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("cancelled").Build())
					return

				case messageEvent := <-ch:
					str += strconv.Itoa(i) + ". " + messageEvent.Message.Content + "\n\n"

					if i == 3 {
						cls()
						_, _ = event.Client().Rest().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetContent(str).Build())
					}
					i++
				}
			}
		}()
	}
}
