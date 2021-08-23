package main

import (
	"context"
	"net/http"
	"os"
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

	client := webhook.New(webhookID, webhookToken, webhook.WithLogger(logger), webhook.WithRestClientConfigOpts(rest.WithHTTPClient(httpClient)))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	_, _ = client.CreateMessage(webhook.NewMessageCreateBuilder().SetContent("test").Build(), rest.WithCtx(ctx), rest.WithReason("this adds a reason header"))

	cancel()
}
