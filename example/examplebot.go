package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

const red = 0xdb1222
const orange = 0xff7700
const green = 0x00fc00

var token = os.Getenv("token")
var guildID = api.Snowflake(os.Getenv("guild_id"))
var adminRoleID = api.Snowflake(os.Getenv("admin_role_id"))
var testRoleID = api.Snowflake(os.Getenv("test_role_id"))
var emoteID = api.Snowflake(os.Getenv("test_emote_id"))

var logger = logrus.New()
var client = http.DefaultClient

func main() {
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("starting ExampleBot...")
	logger.Infof("disgo %s", api.Version)

	dgo, err := disgo.NewBuilder(token).
		SetLogger(logger).
		SetRawGatewayEventsEnabled(true).
		SetHTTPClient(client).
		SetGatewayIntents(api.GatewayIntentsNonPrivileged | api.GatewayIntentGuildMembers).
		SetMemberCachePolicy(api.MemberCachePolicyAll).
		AddEventListeners(events.ListenerAdapter{
			OnRawGateway:         rawGatewayEventListener,
			OnGuildAvailable:     guildAvailListener,
			OnGuildMessageCreate: messageListener,
			OnCommand:            commandListener,
			OnButtonClick:        buttonClickListener,
			OnSelectMenuSubmit:   selectMenuSubmitListener,
		}).
		Build()
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
		return
	}

	rawCmds := []api.CommandCreate{
		{
			Name:              "eval",
			Description:       "runs some go code",
			DefaultPermission: true,
			Options: []api.CommandOption{
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
			DefaultPermission: true,
		},
		{
			Name:              "say",
			Description:       "says what you say",
			DefaultPermission: true,
			Options: []api.CommandOption{
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
			DefaultPermission: true,
			Options: []api.CommandOption{
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
			DefaultPermission: true,
			Options: []api.CommandOption{
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
			cmdsPermissions = append(cmdsPermissions, api.SetGuildCommandPermissions{
				ID:          cmd.ID,
				Permissions: []api.CommandPermission{perms},
			})
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

	logger.Infof("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent) {
	logger.Printf("guild loaded: %s", event.Guild.ID)
}

func rawGatewayEventListener(event *events.RawGatewayEvent) {
	if event.Type == api.GatewayEventInteractionCreate {
		println(string(event.RawPayload))
	}
}

func buttonClickListener(event *events.ButtonClickEvent) {
	switch event.CustomID() {
	case "test1":
		_ = event.Respond(api.InteractionResponseTypeChannelMessageWithSource,
			api.NewMessageCreateBuilder().
				SetContent(event.CustomID()).
				Build(),
		)

	case "test2":
		_ = event.Respond(api.InteractionResponseTypeDeferredChannelMessageWithSource, nil)

	case "test3":
		_ = event.Respond(api.InteractionResponseTypeDeferredUpdateMessage, nil)

	case "test4":
		_ = event.Respond(api.InteractionResponseTypeUpdateMessage,
			api.NewMessageCreateBuilder().
				SetContent(event.CustomID()).
				Build(),
		)
	}
}

func selectMenuSubmitListener(event *events.SelectMenuSubmitEvent) {
	switch event.CustomID() {
	case "test3":
		if err := event.DeferEdit(); err != nil {
			logger.Errorf("error sending interaction response: %s", err)
		}
		_, _ = event.SendFollowup(api.NewMessageCreateBuilder().
			SetEphemeral(true).
			SetContentf("selected options: %s", event.Values()).
			Build(),
		)
	}
}

func commandListener(event *events.CommandEvent) {
	switch event.CommandName() {
	case "eval":
		go func() {
			code := event.Option("code").String()
			embed := api.NewEmbedBuilder().
				SetColor(orange).
				AddField("Status", "...", true).
				AddField("Time", "...", true).
				AddField("Code", "```go\n"+code+"\n```", false).
				AddField("Output", "```\n...\n```", false)
			_ = event.Reply(api.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

			start := time.Now()
			output, err := gval.Evaluate(code, map[string]interface{}{
				"disgo": event.Disgo(),
				"dgo":   event.Disgo(),
				"event": event,
			})

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, err = event.Interaction.EditOriginal(api.NewMessageUpdateBuilder().
					SetEmbeds(embed.
						SetColor(red).
						SetField(0, "Status", "Failed", true).
						SetField(3, "Output", "```"+err.Error()+"```", false).
						Build(),
					).
					Build(),
				)
				if err != nil {
					logger.Errorf("error sending interaction response: %s", err)
				}
				return
			}
			_, err = event.Interaction.EditOriginal(api.NewMessageUpdateBuilder().
				SetEmbeds(embed.
					SetColor(green).
					SetField(0, "Status", "Success", true).
					SetField(3, "Output", "```"+fmt.Sprintf("%+v", output)+"```", false).
					Build(),
				).
				Build(),
			)
			if err != nil {
				logger.Errorf("error sending interaction response: %s", err)
			}
		}()

	case "say":
		_ = event.Reply(api.NewMessageCreateBuilder().
			SetContent(event.Option("message").String()).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		reader, _ := os.Open("gopher.png")
		if err := event.Reply(api.NewMessageCreateBuilder().
			SetContent("test message").
			AddFile("gopher.png", reader).
			AddActionRow(
				api.NewPrimaryButton("test1", "test1", nil),
				api.NewPrimaryButton("test2", "test2", nil),
				api.NewPrimaryButton("test3", "test3", nil),
				api.NewPrimaryButton("test4", "test4", nil),
			).
			AddActionRow(
				api.NewSelectMenu("test3", "test", 1, 1,
					api.NewSelectOption("test1", "1"),
					api.NewSelectOption("test2", "2"),
					api.NewSelectOption("test3", "3"),
				),
			).
			Build(),
		); err != nil {
			logger.Errorf("error sending interaction response: %s", err)
		}

	case "addrole":
		user := event.Option("member").User()
		role := event.Option("role").Role()

		if err := event.Disgo().RestClient().AddMemberRole(*event.Interaction.GuildID, user.ID, role.ID); err == nil {
			_ = event.Reply(api.NewMessageCreateBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewMessageCreateBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := event.Option("member").User()
		role := event.Option("role").Role()

		if err := event.Disgo().RestClient().RemoveMemberRole(*event.Interaction.GuildID, user.ID, role.ID); err == nil {
			_ = event.Reply(api.NewMessageCreateBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewMessageCreateBuilder().AddEmbeds(
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
		_, _ = event.Message.Reply(api.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&api.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(api.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&api.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.MessageChannel().CreateMessage(api.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				logger.Errorf("error while sending file: %s", err)
				return
			}
			time.Sleep(time.Second * 2)

			embed := api.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = message.Update(api.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = message.Update(api.NewMessageUpdateBuilder().SetContent("").SetEmbeds(api.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Message.Author.CreateDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.CreateMessage(api.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()
	}
}
