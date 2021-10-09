package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DisgoOrg/disgo/bot"

	"github.com/DisgoOrg/disgo/events"

	"github.com/DisgoOrg/disgo/core"
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
	if event.Message.Author.IsBot || event.Message.Author.IsSystem {
		return
	}
	if event.Message.Content == "start" {
		go func() {
			ch, cls := event.Bot().Collectors.NewMessageCollector(func(message *core.Message) bool {
				return message.ChannelID == event.ChannelID && message.Author.ID == event.Message.Author.ID && message.Content != ""
			})
			i := 1
			str := ">>> "
			for message := range ch {
				str += strconv.Itoa(i) + ". " + message.Content + "\n\n"

				if i == 3 {
					cls()
					_, _ = message.Channel().CreateMessage(core.NewMessageCreateBuilder().SetContent(str).Build())
				}
				i++
			}
		}()

	}
}
