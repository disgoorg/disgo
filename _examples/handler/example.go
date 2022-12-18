package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Replies with pong",
		},
		discord.SlashCommandCreate{
			Name:        "ping2",
			Description: "Replies with pong2",
		},
		discord.SlashCommandCreate{
			Name:        "ping3",
			Description: "Replies with pong3",
		},
	}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", disgo.Version)

	mux := handler.New()
	mux.HandleCommand("/ping", handlePingCommand)
	mux.HandleCommand("/{name}", handleCommand)
	mux.HandleComponent("button1:{data}", handleComponent)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentDirectMessages)),
		bot.WithEventListeners(mux),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	// register commands
	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		log.Fatal("error while setting global commands: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func handlePingCommand(client bot.Client, event *handler.CommandEvent) error {
	return event.CreateMessage(discord.MessageCreate{
		Content: "pong",
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.NewPrimaryButton("button1", "button1/testData"),
			},
		},
	})
}

func handleCommand(client bot.Client, event *handler.CommandEvent) error {
	commandName := event.Variables["name"]
	return event.CreateMessage(discord.MessageCreate{Content: "commandName: " + commandName})
}

func handleComponent(client bot.Client, event *handler.ComponentEvent) error {
	data := event.Variables["data"]
	return event.CreateMessage(discord.MessageCreate{Content: "component: " + data})
}
