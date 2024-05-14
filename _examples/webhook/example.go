package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

var (
	webhookID    = snowflake.GetEnv("webhook_id")
	webhookToken = os.Getenv("webhook_token")
)

func main() {
	slog.Info("starting webhook example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	// construct new webhook client
	client := webhook.New(webhookID, webhookToken)
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
	slog.Info("exiting webhook example...")
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
		slog.Error("error sending message to webhook", slog.Any("error", err), slog.Int("i", i))
	}
}
