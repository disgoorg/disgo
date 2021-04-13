package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
	"github.com/DisgoOrg/disgo/api/events"
)

const red = 16711680
const orange = 16562691
const green = 65280

const guildID = "817327181659111454"
const adminRoleID = "817327279583264788"
const testRoleID = "825156597935243304"

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("starting testbot...")

	dgo, err := disgo.NewBuilder(endpoints.Token(os.Getenv("token"))).
		SetLogger(logger).
		SetIntents(api.IntentsGuilds | api.IntentsGuildMessages | api.IntentsGuildMembers).
		SetMemberCachePolicy(api.MemberCachePolicyAll).
		AddEventListeners(&events.ListenerAdapter{
			OnGuildAvailable:       guildAvailListener,
			OnGuildMessageCreate: messageListener,
			OnSlashCommand:         slashCommandListener,
		}).
		Build()
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
		return
	}

	rawCmds := []api.Command{
		{
			Name:              "eval",
			Description:       "runs some go code",
			DefaultPermission: false,
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeString,
					Name:        "code",
					Description: "the code to eval",
					Required:    true,
				},
			},
		},
		{
			Name:              "test",
			Description:       "test test test test test test",
			DefaultPermission: false,
		},
		{
			Name:              "say",
			Description:       "says what you say",
			DefaultPermission: false,
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeString,
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
			},
		},
		{
			Name:              "addrole",
			Description:       "This command adds a role to a member",
			DefaultPermission: false,
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeUser,
					Name:        "member",
					Description: "The member to add a role to",
					Required:    true,
				},
				{
					Type:        api.CommandOptionTypeRole,
					Name:        "role",
					Description: "The role to add to a member",
					Required:    true,
				},
			},
		},
		{
			Name:              "removerole",
			Description:       "This command removes a role from a member",
			DefaultPermission: false,
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeUser,
					Name:        "member",
					Description: "The member to removes a role from",
					Required:    true,
				},
				{
					Type:        api.CommandOptionTypeRole,
					Name:        "role",
					Description: "The role to removes from a member",
					Required:    true,
				},
			},
		},
	}

	// using the api.RestClient directly to avoid the guild needing to be cached
	cmds, err := dgo.RestClient().SetGuildCommands(dgo.ApplicationID(), guildID, rawCmds...)
	if err != nil {
		logger.Errorf("error while registering guild commands: %s", err)
	}

	var cmdsPermissions []api.SetGuildCommandPermissions
	for _, cmd := range cmds {
		var perms api.CommandPermission
		if cmd.Name == "eval" {
			perms = api.CommandPermission{
				ID:         adminRoleID,
				Type:       api.CommandPermissionTypeRole,
				Permission: true,
			}
		} else {
			perms = api.CommandPermission{
				ID:         testRoleID,
				Type:       api.CommandPermissionTypeRole,
				Permission: true,
			}
		}
		cmdsPermissions = append(cmdsPermissions, api.SetGuildCommandPermissions{
			ID:          cmd.ID,
			Permissions: []api.CommandPermission{perms},
		})
	}
	if _, err = dgo.RestClient().SetGuildCommandsPermissions(dgo.ApplicationID(), guildID, cmdsPermissions...); err != nil {
		logger.Errorf("error while setting command permissions: %s", err)
	}

	err = dgo.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to discord: %s", err)
	}

	defer dgo.Close()

	logger.Infof("Bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent) {
	logger.Printf("guild loaded: %s", event.Guild.ID)
}

func slashCommandListener(event *events.SlashCommandEvent) {
	switch event.CommandName {
	case "eval":
		go func() {
			start := time.Now()
			code := event.Option("code").String()
			embed := api.NewEmbedBuilder().
				SetColor(orange).
				AddField("Status", "...", true).
				AddField("Time", "...", true).
				AddField("Code", "```go\n"+code+"\n```", false).
				AddField("Output", "```\n...\n```", false)

			_ = event.Reply(api.NewInteractionResponseBuilder().SetEmbeds(embed.Build()).Build())

			vars := map[string]interface{}{
				"disgo": event.Disgo(),
				"dgo":   event.Disgo(),
				"event": event,
			}
			output, err := gval.Evaluate(code, vars)

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().
					SetEmbeds(embed.
						SetColor(red).
						SetField(0, "Status", "failed", true).
						SetField(3, "Output", "```"+err.Error()+"```", false).
						Build(),
					).
					Build(),
				)
				return
			}
			_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().
				SetEmbeds(embed.
					SetColor(green).
					SetField(0, "Status", "success", true).
					SetField(3, "Output", "```"+fmt.Sprintf("%+v", output)+"```", false).
					Build(),
				).
				Build(),
			)
		}()

	case "say":
		_ = event.Reply(api.NewInteractionResponseBuilder().
			SetContent(event.Option("message").String()).
			SetAllowedMentionsEmpty().
			Build(),
		)

	case "test":
		go func() {
			_ = event.Acknowledge()

			time.Sleep(2 * time.Second)
			_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().
				SetEmbeds(api.NewEmbedBuilder().
					SetDescription("finished with thinking").
					Build(),
				).Build(),
			)

			time.Sleep(1 * time.Second)
			_, _ = event.SendFollowup(api.NewFollowupMessageBuilder().
				SetEmbeds(api.NewEmbedBuilder().
					SetDescription("followup 1").
					Build(),
				).Build(),
			)

			time.Sleep(1 * time.Second)
			_, _ = event.SendFollowup(api.NewFollowupMessageBuilder().
				SetEphemeral(true).
				SetContent("followup 2 only you can see").
				Build(),
			)
		}()

	case "addrole":
		user := event.Option("member").User()
		role := event.Option("role").Role()
		err := event.Disgo().RestClient().AddMemberRole(*event.Interaction.GuildID, user.ID, role.ID)
		if err == nil {
			_ = event.Reply(api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewInteractionResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := event.Option("member").User()
		role := event.Option("role").Role()
		err := event.Disgo().RestClient().RemoveMemberRole(*event.Interaction.GuildID, user.ID, role.ID)
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

func messageListener(event *events.GuildMessageCreateEvent) {
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
