package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/cache"
	"github.com/disgoorg/DisGo"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetSnowflakeEnv("disgo_guild_id")

	//go:embed gopher.png
	gopher []byte
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("bot version: %s", DisGo.Version)

	client, err := disgo.New(token,
		//bot.WithRawEventsEnabled(),
		bot.WithGatewayConfigOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentsNonPrivileged),
			gateway.WithPresence(discord.NewListeningPresence("your bullshit", discord.OnlineStatusOnline, false)),
		),
		bot.WithCacheConfigOpts(
			cache.WithCacheFlags(cache.FlagsAll),
			cache.WithMemberCachePolicy(cache.MemberCachePolicyAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterNone),
		bot.WithEventListeners(listener),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	registerCommands(client)

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	defer client.Close(context.TODO())

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
