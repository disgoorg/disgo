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
	logger.Infof("disgo version: %s", info.Version)

	disgo, err := core.NewBuilder(token).
		SetLogger(logger).
		SetHTTPClient(httpClient).
		SetGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentsDefault)).
		SetCacheConfigOpts(core.WithCacheFlags(core.CacheFlagsDefault)).
		Build()

	if err != nil {
		log.Fatalf("error while building disgo: %s", err)
	}

	defer disgo.Close()

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}