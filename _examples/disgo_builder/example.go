package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token      = os.Getenv("token")
	logger     = log.Default()
	httpClient = http.DefaultClient
)

func main() {
	logger.SetLevel(log.LevelDebug)
	logger.Info("starting example...")
	logger.Info("disgo version: ", info.Version)

	disgo, err := core.NewBotBuilder(token).
		SetLogger(logger).
		SetHTTPClient(httpClient).
		SetGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentsDefault)).
		SetCacheConfigOpts(core.WithCacheFlags(core.CacheFlagsDefault)).
		AddEventListeners(&core.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}).
		Build()
	if err != nil {
		log.Fatal("error while building disgo: ", err)
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
	_, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent(event.Message.Content).Build())
}
