package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name: "say",
			NameLocalizations: map[discord.Locale]string{
				discord.LocaleEnglishGB: "say",
				discord.LocaleGerman:    "sagen",
			},
			Description: "says what you say",
			DescriptionLocalizations: map[discord.Locale]string{
				discord.LocaleEnglishGB: "says what you say",
				discord.LocaleGerman:    "sagt was du sagst",
			},
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name: "message",
					NameLocalizations: map[discord.Locale]string{
						discord.LocaleEnglishGB: "message",
						discord.LocaleGerman:    "nachricht",
					},
					Description: "What to say",
					DescriptionLocalizations: map[discord.Locale]string{
						discord.LocaleEnglishGB: "What to say",
						discord.LocaleGerman:    "Was soll ich sagen?",
					},
					Required: true,
				},
				discord.ApplicationCommandOptionBool{
					Name: "ephemeral",
					NameLocalizations: map[discord.Locale]string{
						discord.LocaleEnglishGB: "ephemeral",
						discord.LocaleGerman:    "kurzlebig",
					},
					Description: "If the response should only be visible to you",
					DescriptionLocalizations: map[discord.Locale]string{
						discord.LocaleEnglishGB: "If the response should only be visible to you",
						discord.LocaleGerman:    "Wenn die Antwort nur dir sichtbar sein soll",
					},
					Required: true,
				},
			},
		},
	}
)

func main() {
	log.SetLevel(log.LevelTrace)
	log.Info("starting example...")
	log.Infof("disgo version: %s", disgo.Version)

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
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

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			Build(),
		)
		if err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}
	}
}
