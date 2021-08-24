package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/DisgoOrg/log"
)

var (
	webhookID    = discord.Snowflake(os.Getenv("webhook_id"))
	webhookToken = os.Getenv("webhook_token")
	logger       = log.Default()
	httpClient   = http.DefaultClient
)

func main() {
	logger.SetLevel(log.LevelDebug)
	logger.Info("starting ExampleBot...")
	logger.Infof("disgo %s", info.Version)

	client := webhook.New(webhookID, webhookToken,
		webhook.WithLogger(logger),
		webhook.WithRestClientConfigOpts(
			rest.WithHTTPClient(httpClient),
		),
	)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)

	for i := 1; i <= 10; i++ {
		go send(ctx, client, i)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func send(ctx context.Context, client webhook.Client, i int) {
	_, _ = client.CreateMessage(webhook.NewMessageCreateBuilder().SetContentf("test %d", i).Build(), rest.WithCtx(ctx), rest.WithReason("this adds a reason header"))
}
