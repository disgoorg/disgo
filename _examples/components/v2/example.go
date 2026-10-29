package main

import (
	"bytes"
	"context"
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	//go:embed thumbnail.jpg
	thumbnail []byte

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "test",
			Description: "test",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "ephemeral",
					Description: "if the message should be ephemeral",
					Required:    false,
				},
			},
			IntegrationTypes: []discord.ApplicationIntegrationType{
				discord.ApplicationIntegrationTypeGuildInstall,
				discord.ApplicationIntegrationTypeUserInstall,
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
	slog.SetLogLoggerLevel(slog.LevelDebug)

	r := handler.New()
	r.SlashCommand("/test", onTest)
	r.SlashCommand("/modal", onModal)
	r.Modal("/modal", onModalSubmit)

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListeners(r),
	)
	if err != nil {
		slog.Error("error while building bot", slog.Any("err", err))
		return
	}
	defer client.Close(context.TODO())

	if err = handler.SyncCommands(client, commands, []snowflake.ID{guildID}); err != nil {
		slog.Error("error while syncing commands", slog.Any("err", err))
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
		return
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onTest(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	flags := discord.MessageFlagIsComponentsV2
	if ephemeral, ok := data.OptBool("ephemeral"); !ok || ephemeral {
		flags = flags.Add(discord.MessageFlagEphemeral)
	}

	return e.CreateMessage(discord.MessageCreate{
		Flags: flags,
		Components: []discord.LayoutComponent{
			discord.NewContainer(
				discord.NewSection(
					discord.NewTextDisplay("**Name: [Seeing Red](https://open.spotify.com/track/65qBr6ToDUjTD1RiE1H4Gl)**"),
					discord.NewTextDisplay("**Artist: [Architects](https://open.spotify.com/artist/3ZztVuWxHzNpl0THurTFCv)**"),
					discord.NewTextDisplay("**Album:  [The Sky, The Earth & All Between](https://open.spotify.com/album/2W82VyyIFAXigJEiLm5TT1)**"),
				).WithAccessory(discord.NewThumbnail("attachment://thumbnail.png")),
				discord.NewTextDisplay("`0:08`/`3:40`"),
				discord.NewTextDisplay("[🔘▬▬▬▬▬▬▬▬▬]"),
				discord.NewSmallSeparator(),
				discord.NewActionRow(
					discord.NewPrimaryButton("", "/player/previous").WithEmoji(discord.ComponentEmoji{Name: "⏮"}),
					discord.NewPrimaryButton("", "/player/pause_play").WithEmoji(discord.ComponentEmoji{Name: "⏯"}),
					discord.NewPrimaryButton("", "/player/next").WithEmoji(discord.ComponentEmoji{Name: "⏭"}),
					discord.NewDangerButton("", "/player/stop").WithEmoji(discord.ComponentEmoji{Name: "⏹"}),
					discord.NewPrimaryButton("", "/player/like").WithEmoji(discord.ComponentEmoji{Name: "❤️"}),
				),
			).WithAccentColor(0x5c5fea),
		},
		Files: []*discord.File{
			discord.NewFile("thumbnail.png", "", bytes.NewReader(thumbnail)),
		},
	})
}

func onModal(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.Modal(discord.ModalCreate{
		CustomID: "/modal",
		Title:    "Test Modal",
		Components: []discord.LayoutComponent{
			discord.NewTextDisplay("This is a modal"),
			discord.NewLabel("Test Input", discord.NewShortTextInput("test_input")),
			discord.NewLabel("Test Select", discord.NewUserSelectMenu("test_user", "Select a user")),
		},
	})
}

func onModalSubmit(e *handler.ModalEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "You submitted the modal!",
		Flags:   discord.MessageFlagEphemeral,
	})
}
