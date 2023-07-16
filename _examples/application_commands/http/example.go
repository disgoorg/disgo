package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/httpserver"
)

var (
	token     = os.Getenv("disgo_token")
	publicKey = os.Getenv("disgo_public_key")
	guildID   = snowflake.GetEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "say",
			Description: "says what you say",
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
	log.Info("disgo version: ", disgo.Version)

	// use custom ed25519 verify implementation
	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}

	client, err := disgo.New(token,
		bot.WithHTTPServerConfigOpts(publicKey,
			httpserver.WithURL("/interactions/callback"),
			httpserver.WithAddress(":80"),
		),
		bot.WithEventListenerFunc(commandListener),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = client.OpenHTTPServer(); err != nil {
		log.Fatal("error while starting http server: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		if err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			Build(),
		); err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}
	}
}
