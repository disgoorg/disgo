package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

const (
	red    = 16711680
	orange = 16562691
	green  = 65280
)

var (
	token       = os.Getenv("disgo_token")
	guildID     = discord.Snowflake(os.Getenv("disgo_guild_id"))
	adminRoleID = discord.Snowflake(os.Getenv("disgo_admin_role_id"))
	testRoleID  = discord.Snowflake(os.Getenv("disgo_test_role_id"))
)

func main() {
	log.SetDefault(log.New(log.LstdFlags | log.Lshortfile))
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("bot version: %s", info.Version)

	bot, err := core.NewBotBuilder(token).
		SetRawEventsEnabled(true).
		SetGatewayConfig(gateway.Config{
			GatewayIntents: discord.GatewayIntentsAll,
			Compress:       true,
		}).
		SetCacheConfig(core.CacheConfig{
			CacheFlags:        core.CacheFlagsDefault,
			MemberCachePolicy: core.MemberCachePolicyAll,
		}).
		AddEventListeners(getListener()).
		Build()
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	registerCommands(bot)

	err = bot.Connect()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	defer bot.Close()

	log.Infof("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
