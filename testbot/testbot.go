package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

func main() {
	log.Infof("starting testbot...")
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
		log.Fatalf("error while building disgo instance: %s", err)
		return
	}

	commands, err := dgo.RestClient().SetGuildApplicationCommands(dgo.ApplicationID(), "817327181659111454",
		api.ApplicationCommand{
			Name:        "addrole",
			Description: "This command adds a role to a member",
			Options: []api.ApplicationCommandOption{
				{
					Type:        api.ApplicationCommandOptionTypeUser,
					Name:        "member",
					Description: "The member to add a role to",
					Required:    true,
				},
				{
					Type:        api.ApplicationCommandOptionTypeRole,
					Name:        "role",
					Description: "The role to add to a member",
					Required:    true,
				},
			},
		},
		api.ApplicationCommand{
			Name:        "removerole",
			Description: "This command removes a role from a member",
			Options: []api.ApplicationCommandOption{
				{
					Type:        api.ApplicationCommandOptionTypeUser,
					Name:        "member",
					Description: "The member to removes a role from",
					Required:    true,
				},
				{
					Type:        api.ApplicationCommandOptionTypeRole,
					Name:        "role",
					Description: "The role to removes from a member",
					Required:    true,
				},
			},
		},
	)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	log.Printf("%#v", commands)

	err = dgo.Connect()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	defer dgo.Close()

	log.Infof("Bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent){
	log.Printf("guild loaded: %s", event.GuildID)
}

func slashCommandListener(event *events.SlashCommandEvent){
	switch event.Interaction.Data.Name {
	case "test":
		_ = event.Reply(*api.NewInteractionResponseBuilder().
			SetContent("test").
			SetEphemeral(true).
			AddEmbeds(
				api.NewEmbedBuilder().SetDescription("test1").Build(),
				api.NewEmbedBuilder().SetDescription("test2").Build(),
			).
			Build(),
		)
	case "addrole":
		userID := event.Interaction.Data.Options[0].Snowflake()
		roleID := event.Interaction.Data.Options[1].Snowflake()
		err := event.Disgo.RestClient().AddMemberRole(*event.GuildID, userID, roleID)
		if err == nil {
			_ = event.Reply(*api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(intPtr(65280)).SetDescriptionf("Added <@&%v> to <@%v>", roleID, userID).Build(),
			).Build())
		} else {
			_ = event.Reply(*api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(intPtr(16711680)).SetDescriptionf("Failed to add <@&%v> to <@%v>", roleID, userID).Build(),
			).Build())
		}
	case "removerole":
		userID := event.Interaction.Data.Options[0].Snowflake()
		roleID := event.Interaction.Data.Options[1].Snowflake()
		err := event.Disgo.RestClient().RemoveMemberRole(*event.GuildID, userID, roleID)
		if err == nil {
			_ = event.Reply(*api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(intPtr(65280)).SetDescriptionf("Removed <@&%v> from <@%v>", roleID, userID).Build(),
			).Build())
		} else {
			_ = event.Reply(*api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(intPtr(16711680)).SetDescriptionf("Failed to remove <@&%v> from <@%v>", roleID, userID).Build(),
			).Build())
		}
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
		_, _ = event.Message.Reply(api.Message{Content: strPtr("ping")})
		
	case "pong":
		_, _ = event.Message.Reply(api.Message{Content: strPtr("pong")})

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.SendMessage(api.Message{Content: strPtr("helo")})
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()
	}
}

func strPtr(str string) *string {
	return &str
}

func intPtr(int int) *int {
	return &int
}