package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/bot"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", info.Version)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}),
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
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
	if event.Message.Author.BotUser || event.Message.Author.System {
		return
	}
	if event.Message.Content == "start" {
		go func() {
			ch, cls := bot.NewCollector(event.Client(), func(event *events.MessageCreateEvent) bool {
				return event.ChannelID == event.ChannelID && event.Message.Author.ID == event.Message.Author.ID && event.Message.Content != ""
			})
			i := 1
			str := ">>> "
			ctx, clsCtx := context.WithTimeout(context.Background(), 20*time.Second)
			defer clsCtx()
			for {
				select {
				case <-ctx.Done():
					_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("cancelled").Build())
					return

				case messageEvent := <-ch:
					str += strconv.Itoa(i) + ". " + messageEvent.Message.Content + "\n\n"

					if i == 3 {
						cls()
						_, _ = event.Client().Rest().ChannelService().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetContent(str).Build())
					}
					i++
				}
			}
		}()
	}
}
