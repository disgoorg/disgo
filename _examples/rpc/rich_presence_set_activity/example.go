package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rpc"
)

var (
	clientID = snowflake.GetEnv("disgo_client_id")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Info("example is starting...")

	client, err := rpc.New(clientID)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Not allowed to set CreatedAt, maybe only parameter when receiving presence?
	if _, err := client.SetActivity(os.Getpid(), discord.Activity{
		Type:    discord.ActivityTypeGame,
		Details: json.Ptr("Lorem Ipsum"),
		State:   json.Ptr("dolor sit amet"),
		Timestamps: &discord.ActivityTimestamps{
			Start: json.Ptr(time.Now()),
		},
	}); err != nil {
		log.Fatal(err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
