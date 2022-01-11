package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token     = os.Getenv("disgo_token")
	publicKey = os.Getenv("disgo_public_key")
	guildID   = discord.Snowflake(os.Getenv("disgo_guild_id"))

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:              "say",
			Description:       "says what you say",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
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
	log.Info("disgo version: ", info.Version)

	disgo, err := bot.New(token,
		bot.WithHTTPServerOpts(
			httpserver.WithURL("/interactions/callback"),
			httpserver.WithPort(":80"),
			httpserver.WithPublicKey(publicKey),
		),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: commandListener,
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer disgo.Close(context.TODO())

	if _, err = disgo.SetGuildCommands(guildID, commands); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = disgo.StartHTTPServer(); err != nil {
		log.Fatal("error while starting http server: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionEvent) {
	data := event.SlashCommandInteractionData()
	if data.CommandName == "say" {
		err := event.Create(discord.NewMessageCreateBuilder().
			SetContent(*data.Options.String("message")).
			Build(),
		)
		if err != nil {
			event.Bot().Logger().Error("error on sending response: ", err)
		}
	}
}
