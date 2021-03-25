package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

func main() {
	token := os.Getenv("token")

	dgo, err := disgo.NewBuilder(token).
		SetLogLevel(log.InfoLevel).
		SetIntents(api.IntentsGuilds | api.IntentsGuildMessages | api.IntentsGuildMembers).
		SetMemberCachePolicy(api.MemberCachePolicyAll).
		AddEventListeners(&events.ListenerAdapter{
			OnGuildAvailable: guildAvailListener,
			OnGuildMessageReceived: messageListener,
			OnSlashCommand: slashCommandListener,
		}).
		Build()
	if err != nil {
		return
	}


	err = dgo.Connect()
	if err != nil {
		log.Fatal(err)
	}
	//_, _ = dgo.CreateCommand("test", "this is a test command").Create()

	defer dgo.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent){
	log.Printf("guild loaded: %s", event.GuildID)
}

func slashCommandListener(event *events.SlashCommandEvent){
	if event.Name == "test" {
		go func() {
			_ = event.Reply(*api.NewInteractionResponseBuilder().
				SetContent("test").
				SetEphemeral(true).
				AddEmbeds(
					api.NewEmbedBuilder().SetDescription("test1").Build(),
					api.NewEmbedBuilder().SetDescription("test2").Build(),
				).
				Build(),
			)
		}()
	}
}

func messageListener(event *events.GuildMessageReceivedEvent) {
	if event.Message.Author.IsBot {
		return
	}
	if event.Message.Content == nil {
		return
	}

	switch *event.Message.Content {
	case "ping":
		_, _ = event.MessageChannel().SendMessage("pong")
		
	case "pong":
		_, _ = event.MessageChannel().SendMessage("ping")

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.SendMessage("helo")
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()
	}
}
