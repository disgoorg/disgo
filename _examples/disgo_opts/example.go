package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
)

func main() {
	disgo, err := core.NewBot(os.Getenv("token"),
		core.WithGatewayConfigOpts(
			gateway.WithGatewayIntents(
				discord.GatewayIntentGuilds,
				discord.GatewayIntentGuildMessages,
				discord.GatewayIntentDirectMessages,
			),
		),
		core.WithCacheConfigOpts(core.WithCacheFlags(core.CacheFlagsDefault)),
		core.WithEventListeners(&core.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	defer disgo.Close()

	if err = disgo.ConnectGateway(); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *core.MessageCreateEvent) {
	_, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent(event.Message.Content).Build())
}
