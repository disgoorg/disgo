package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

	disgo, err := core.NewBot(token,
		core.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		core.WithEventListeners(&core.ListenerAdapter{
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

func onMessageCreate(event *core.MessageCreateEvent) {
	if event.Message.Author.IsBot || event.Message.Author.IsSystem {
		return
	}
	if event.Message.Content == "start" {
		go func() {
			ch, cls := event.Channel().CollectMessages(func(message *core.Message) bool {
				return message.ChannelID == event.ChannelID && message.Author.ID == event.Message.Author.ID && message.Content != ""
			})
			i := 1
			str := ">>> "
			for message := range ch {
				if i > 3 {
					cls()
					_, _ = message.Channel().CreateMessage(core.NewMessageCreateBuilder().SetContent(str).Build())
				}
				str += strconv.Itoa(i) + ". " + message.Content + "\n\n"
				i++
			}
		}()

	}
}
