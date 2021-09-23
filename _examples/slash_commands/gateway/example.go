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

var (
	token   = os.Getenv("disgo_token")
	guildID = discord.Snowflake(os.Getenv("disgo_guild_id"))

	commands = []discord.ApplicationCommandCreate{
		{
			Type:              discord.ApplicationCommandTypeSlash,
			Name:              "say",
			Description:       "says what you say",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
				{
					Type:        discord.ApplicationCommandOptionTypeString,
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
			},
		},
	}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", info.Version)

	disgo, err := core.NewBotBuilder(token).
		SetGatewayConfig(gateway.Config{
			GatewayIntents: discord.GatewayIntentsNone,
		}).
		SetCacheConfig(core.CacheConfig{
			CacheFlags: core.CacheFlagsDefault,
		}).
		AddEventListeners(&core.ListenerAdapter{
			OnSlashCommand: commandListener,
		}).
		Build()

	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer disgo.Close()

	_, err = disgo.SetGuildCommands(guildID, commands)
	if err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = disgo.ConnectGateway(); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func commandListener(event *core.SlashCommandEvent) {
	if event.CommandName == "say" {
		err := event.Create(core.NewMessageCreateBuilder().
			SetContent(event.Options["message"].String()).
			Build(),
		)
		if err != nil {
			event.Bot().Logger.Error("error on sending response: ", err)
		}
	}
}
