package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetSnowflakeEnv("disgo_test_guild_id")

	//go:embed gopher.png
	gopher []byte
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("bot version: %s", info.Version)

	disgo, err := bot.New(token,
		//bot.WithRawEventsEnabled(),
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentsNonPrivileged),
			gateway.WithPresence(core.NewListeningPresence("your bullshit", discord.OnlineStatusOnline, false)),
		),
		bot.WithCacheOpts(
			core.WithCacheFlags(core.CacheFlagsAll),
			core.WithMemberCachePolicy(core.MemberCachePolicyAll),
		),
		bot.WithMemberChunkingFilter(core.MemberChunkingFilterNone),
		bot.WithEventListeners(listener),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	registerCommands(disgo)

	if err = disgo.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	defer disgo.Close(context.TODO())

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
