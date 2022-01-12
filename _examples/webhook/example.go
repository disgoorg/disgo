package main

import (
	"context"
	"os"
	"sync"
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
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting webhook example...")
	log.Info("disgo version: ", info.Version)

	// construct new webhook client
	client := webhook.NewClient(webhookID, webhookToken)
	defer client.Close(context.TODO())

	// new sync.WaitGroup to await all messages to be sent before shutting down
	var wg sync.WaitGroup

	// send 10 messages with the webhook
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go send(&wg, client, i)
	}

	// wait for all messages to be sent
	wg.Wait()
	log.Info("existing webhook example...")
}

// send(s) a message to the webhook
func send(wg *sync.WaitGroup, client webhook.Client, i int) {
	defer wg.Done()

	if _, err := client.CreateMessage(discord.NewWebhookMessageCreateBuilder().
		SetContentf("test %d", i).
		Build(),
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
	); err != nil {
		log.Errorf("error sending message %d: %s", i, err)
	}
}
