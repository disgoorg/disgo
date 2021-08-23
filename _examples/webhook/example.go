package webhook

import (
	"context"
	"os"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/info"
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

	client := webhook.New(nil, nil, nil, webhookID, webhookToken)

	ctx, cancel := context.WithCancel(context.Background())

	_, _ = client.CreateMessage(ctx, webhook.NewMessageCreateBuilder().Build())

	cancel()
}
