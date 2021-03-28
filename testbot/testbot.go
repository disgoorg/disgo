package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

func main() {
	log.Infof("starting testbot...")
	token := os.Getenv("token")
	publicKey := os.Getenv("public-key")

	dgo, err := disgo.NewBuilder(token).
		SetLogLevel(log.InfoLevel).
		SetIntents(api.IntentsGuilds|api.IntentsGuildMessages|api.IntentsGuildMembers).
		SetMemberCachePolicy(api.MemberCachePolicyAll).
		SetWebhookServerProperties("/webhooks/interactions/callback", 80, publicKey).
		AddEventListeners(&events.ListenerAdapter{
			OnGuildAvailable:       guildAvailListener,
			OnGuildMessageReceived: messageListener,
			OnSlashCommand:         slashCommandListener,
		}).
		Build()
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
		return
	}

	_, err = dgo.RestClient().SetGuildCommands(dgo.ApplicationID(), "817327181659111454",
		api.SlashCommand{
			Name:        "test",
			Description: "test test test test test test",
		},
		api.SlashCommand{
			Name:        "say",
			Description: "says what you say",
			Options: []*api.CommandOption{
				{
					Type:        api.OptionTypeString,
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
			},
		},
		api.SlashCommand{
			Name:        "addrole",
			Description: "This command adds a role to a member",
			Options: []*api.CommandOption{
				{
					Type:        api.OptionTypeUser,
					Name:        "member",
					Description: "The member to add a role to",
					Required:    true,
				},
				{
					Type:        api.OptionTypeRole,
					Name:        "role",
					Description: "The role to add to a member",
					Required:    true,
				},
			},
		},
		api.SlashCommand{
			Name:        "removerole",
			Description: "This command removes a role from a member",
			Options: []*api.CommandOption{
				{
					Type:        api.OptionTypeUser,
					Name:        "member",
					Description: "The member to removes a role from",
					Required:    true,
				},
				{
					Type:        api.OptionTypeRole,
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

	err = dgo.Start()
	if err != nil {
		log.Fatalf("error while starting webhookserver: %s", err)
	}

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

func guildAvailListener(event *events.GuildAvailableEvent) {
	log.Printf("guild loaded: %s", event.GuildID)
}

func slashCommandListener(event *events.SlashCommandEvent) {
	switch event.Interaction.Data.Name {
	case "say":
		_ = event.Reply(api.NewInteractionResponseBuilder().
			SetContent(event.OptionByName("message").String()).
			SetAllowedMentionsEmpty().
			Build(),
		)
	case "test":
		_ = event.Reply(api.NewInteractionResponseBuilder().
			SetContent("test").
			SetEphemeral(true).
			AddEmbeds(
				api.NewEmbedBuilder().SetDescription("test1").Build(),
				api.NewEmbedBuilder().SetDescription("test2").Build(),
			).
			Build(),
		)
	case "addrole":
		user := event.OptionByName("member").User()
		role := event.OptionByName("role").Role()
		err := event.Disgo.RestClient().AddMemberRole(*event.Interaction.GuildID, user.ID, role.ID)
		if err == nil {
			_ = event.Reply(api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(16711680).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}
	case "removerole":
		user := event.OptionByName("member").User()
		role := event.OptionByName("role").Role()
		err := event.Disgo.RestClient().RemoveMemberRole(*event.Interaction.GuildID, user.ID, role.ID)
		if err == nil {
			_ = event.Reply(api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(16711680).SetDescriptionf("Failed to remove %s from %s", role, user).Build(),
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
		_, _ = event.Message.Reply(api.NewMessageBuilder().SetContent("pong").SetAllowedMentions(&api.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(api.NewMessageBuilder().SetContent("ping").SetAllowedMentions(&api.AllowedMentions{RepliedUser: false}).Build())

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.SendMessage(api.NewMessageBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()
	}
}
