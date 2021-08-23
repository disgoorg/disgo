package webhook

import (
	"context"
	"os"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/webhook"
	"github.com/DisgoOrg/log"
)

var (
	webhookID    = discord.Snowflake(os.Getenv("webhook_id"))
	webhookToken = os.Getenv("webhook_token")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting ExampleBot...")
	log.Infof("disgo %s", info.Version)

	client := webhook.New(nil, nil, webhookID, webhookToken)

	ctx, cancel := context.WithCancel(context.Background())

	_, _ = client.CreateMessage(webhook.NewMessageCreateBuilder().Build(), rest.WithCtx(ctx), rest.WithReason("this adds a reason header"))

	cancel()
}
