package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

var (
	token     = os.Getenv("disgo_token")
	publicKey = os.Getenv("disgo_public_key")
	guildID   = snowflake.GetSnowflakeEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			CommandName:       "say",
			Description:       "says what you say",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "ephemeral",
					Description: "If the response should only be visible to you",
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

	// use custom ed25519 verify implementation
	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}

	client, err := disgo.New(token,
		bot.WithHTTPServerConfigOpts(
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

	defer client.Close(context.TODO())

	if _, err = client.Rest().ApplicationService().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = client.StartHTTPServer(); err != nil {
		log.Fatal("error while starting http server: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionEvent) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		err := event.Respond(discord.InteractionCallbackTypeCreateMessage, discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			Build(),
		)
		if err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}
	}
}
