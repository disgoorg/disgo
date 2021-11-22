package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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

	disgo, err := bot.New(token,
		bot.WithGatewayOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer disgo.Close()

	if err = disgo.ConnectGateway(); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
	if event.Message.Author.Bot || event.Message.Author.System {
		return
	}
	if event.Message.Content == "start" {
		go func() {
			ch, cls := event.Bot().Collectors.NewMessageCollector(func(message *core.Message) bool {
				return message.ChannelID == event.ChannelID && message.Author.ID == event.Message.Author.ID && message.Content != ""
			})
			i := 1
			str := ">>> "
			ctx, clsCtx := context.WithTimeout(context.Background(), 20*time.Second)
			defer clsCtx()
			for {
				select {
				case <-ctx.Done():
					_, _ = event.Channel().CreateMessage(discord.NewMessageCreateBuilder().SetContent("cancelled").Build())
					return

				case message := <-ch:
					str += strconv.Itoa(i) + ". " + message.Content + "\n\n"

					if i == 3 {
						cls()
						_, _ = message.Channel().CreateMessage(discord.NewMessageCreateBuilder().SetContent(str).Build())
					}
					i++
				}
			}
		}()
	}
}
