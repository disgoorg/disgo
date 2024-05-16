package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
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
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(commandListener),
	)
	if err != nil {
		panic("error while building disgo instance: " + err.Error())
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		panic("error while registering commands: " + err.Error())
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		panic("error while connecting to gateway: " + err.Error())
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
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
			event.Client().Logger().Error("error on sending response", slog.Any("err", err))
		}
	}
}
