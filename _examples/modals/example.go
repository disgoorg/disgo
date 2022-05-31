package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
)

var (
	token = os.Getenv("disgo_token")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "modal",
			Description: "brings up a modal",
			IntegrationTypes: []discord.ApplicationIntegrationType{
				discord.ApplicationIntegrationTypeUserInstall,
				discord.ApplicationIntegrationTypeGuildInstall,
			},
			Contexts: []discord.InteractionContextType{
				discord.InteractionContextTypeGuild,
				discord.InteractionContextTypeBotDM,
				discord.InteractionContextTypePrivateChannel,
			},
		},
	}
)

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNone)),
		bot.WithEventListenerFunc(commandListener),
		bot.WithEventListenerFunc(modalListener),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if err = handler.SyncCommands(client, commands, nil); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "modal" {
		if err := event.Modal(discord.NewModalCreateBuilder().
			SetTitle("Modal Title").
			SetCustomID("modal-id").
			AddLabel("short text", discord.NewShortTextInput("short-text-input")).
			AddLabel("paragraph text", discord.NewParagraphTextInput("paragraph-text-input")).
			AddLabel("select menu", discord.NewStringSelectMenu("select-menu", "select something idiot",
				discord.NewStringSelectMenuOption("helo", "helo"),
				discord.NewStringSelectMenuOption("uwu", "uwu"),
				discord.NewStringSelectMenuOption("owo", "owo"),
			).
				WithMinValues(0).
				WithMaxValues(2),
			).
			Build(),
		); err != nil {
			event.Client().Logger.Error("error creating modal", slog.Any("err", err))
		}
	}
}

func modalListener(event *events.ModalSubmitInteractionCreate) {
	var content string
	for component := range event.Data.AllComponents() {
		switch c := component.(type) {
		case discord.TextInputComponent:
			content += c.CustomID + ": " + c.Value + "\n"
		case discord.StringSelectMenuComponent:
			content += c.CustomID + ": " + strings.Join(c.Values, ", ") + "\n"
		}
	}

	if err := event.CreateMessage(discord.MessageCreate{
		Content: content,
	}); err != nil {
		event.Client().Logger.Error("error creating modal", slog.Any("err", err))
	}
}
