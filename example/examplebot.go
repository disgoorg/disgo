package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DisgoOrg/disgo/core"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/log"
	"github.com/PaesslerAG/gval"
)

const red = 16711680
const orange = 16562691
const green = 65280

var token = os.Getenv("disgo_test_token")
var guildID = discord.Snowflake(os.Getenv("guild_id"))
var adminRoleID = discord.Snowflake(os.Getenv("admin_role_id"))
var testRoleID = discord.Snowflake(os.Getenv("test_role_id"))

var client = http.DefaultClient

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting ExampleBot...")
	log.Infof("disgo %s", core.Version)

	dgo, err := core.NewBuilder(token).
		SetRawGatewayEventsEnabled(true).
		SetHTTPClient(client).
		SetGatewayIntents(gateway.GatewayIntentGuilds, gateway.GatewayIntentGuildMessages, gateway.GatewayIntentGuildMembers).
		SetMemberCachePolicy(core.MemberCachePolicyAll).
		AddEventListeners(&events.ListenerAdapter{
			OnRawGateway:         rawGatewayEventListener,
			OnGuildAvailable:     guildAvailListener,
			OnGuildMessageCreate: messageListener,
			OnCommand:            commandListener,
			OnButtonClick:        buttonClickListener,
			OnSelectMenuSubmit:   selectMenuSubmitListener,
		}).
		Build()
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
		return
	}

	rawCmds := []discord.CommandCreate{
		{
			Name:              "eval",
			Description:       "runs some go code",
			DefaultPermission: true,
			Options: []discord.CommandOption{
				{
					Type:        discord.CommandOptionTypeString,
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
			Options: []discord.CommandOption{
				{
					Type:        discord.CommandOptionTypeString,
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
			Options: []discord.CommandOption{
				{
					Type:        discord.CommandOptionTypeUser,
					Name:        "member",
					Description: "The member to add a role to",
					Required:    true,
				},
				{
					Type:        discord.CommandOptionTypeRole,
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
			Options: []discord.CommandOption{
				{
					Type:        discord.CommandOptionTypeUser,
					Name:        "member",
					Description: "The member to removes a role from",
					Required:    true,
				},
				{
					Type:        discord.CommandOptionTypeRole,
					Name:        "role",
					Description: "The role to removes from a member",
					Required:    true,
				},
			},
		},
	}

	// using the discord.RestClient directly to avoid the guild needing to be cached
	cmds, err := dgo.RestServices().ApplicationService().SetGuildCommands(dgo.ApplicationID(), guildID, rawCmds...)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	var cmdsPermissions []discord.SetGuildCommandPermissions
	for _, cmd := range cmds {
		var perms discord.CommandPermission
		if cmd.Name == "eval" {
			perms = discord.CommandPermission{
				ID:         adminRoleID,
				Type:       discord.CommandPermissionTypeRole,
				Permission: true,
			}
		} else {
			perms = discord.CommandPermission{
				ID:         testRoleID,
				Type:       discord.CommandPermissionTypeRole,
				Permission: true,
			}
			cmdsPermissions = append(cmdsPermissions, discord.SetGuildCommandPermissions{
				ID:          cmd.ID,
				Permissions: []discord.CommandPermission{perms},
			})
		}
		cmdsPermissions = append(cmdsPermissions, discord.SetGuildCommandPermissions{
			ID:          cmd.ID,
			Permissions: []discord.CommandPermission{perms},
		})
	}
	if _, err = dgo.RestServices().SetGuildCommandsPermissions(dgo.ApplicationID(), guildID, cmdsPermissions...); err != nil {
		log.Errorf("error while setting command permissions: %s", err)
	}

	err = dgo.Connect()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	defer dgo.Close()

	log.Infof("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent) {
	log.Infof("guild loaded: %s", event.Guild.ID)
}

func rawGatewayEventListener(event *events.RawEvent) {
	if event.Type == gateway.GatewayEventInteractionCreate {
		println(string(event.RawPayload))
	}
}

func buttonClickListener(event *events.ButtonClickEvent) {
	switch event.CustomID() {
	case "test1":
		_ = event.Respond(discord.InteractionResponseTypeChannelMessageWithSource,
			discord.NewMessageCreateBuilder().
				SetContent(event.CustomID()).
				Build(),
		)

	case "test2":
		_ = event.Respond(discord.InteractionResponseTypeDeferredChannelMessageWithSource, nil)

	case "test3":
		_ = event.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil)

	case "test4":
		_ = event.Respond(discord.InteractionResponseTypeUpdateMessage,
			discord.NewMessageCreateBuilder().
				SetContent(event.CustomID()).
				Build(),
		)
	}
}

func selectMenuSubmitListener(event *events.SelectMenuSubmitEvent) {
	switch event.CustomID() {
	case "test3":
		if err := event.DeferEdit(); err != nil {
			log.Errorf("error sending interaction response: %s", err)
		}
		_, _ = event.SendFollowup(discord.NewMessageCreateBuilder().
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
			embed := core.NewEmbedBuilder().
				SetColor(orange).
				AddField("Status", "...", true).
				AddField("Time", "...", true).
				AddField("Code", "```go\n"+code+"\n```", false).
				AddField("Output", "```\n...\n```", false)
			_ = event.Reply(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

			start := time.Now()
			output, err := gval.Evaluate(code, map[string]interface{}{
				"disgo": event.Disgo(),
				"dgo":   event.Disgo(),
				"event": event,
			})

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, err = event.Interaction.EditOriginal(discord.NewMessageUpdateBuilder().
					SetEmbeds(embed.
						SetColor(red).
						SetField(0, "Status", "Failed", true).
						SetField(3, "Output", "```"+err.Error()+"```", false).
						Build(),
					).
					Build(),
				)
				if err != nil {
					log.Errorf("error sending interaction response: %s", err)
				}
				return
			}
			_, err = event.Interaction.EditOriginal(discord.NewMessageUpdateBuilder().
				SetEmbeds(embed.
					SetColor(green).
					SetField(0, "Status", "Success", true).
					SetField(3, "Output", "```"+fmt.Sprintf("%+v", output)+"```", false).
					Build(),
				).
				Build(),
			)
			if err != nil {
				log.Errorf("error sending interaction response: %s", err)
			}
		}()

	case "say":
		_ = event.Reply(discord.NewMessageCreateBuilder().
			SetContent(event.Option("message").String()).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		reader, _ := os.Open("gopher.png")
		if err := event.Reply(discord.NewMessageCreateBuilder().
			SetContent("test message").
			AddFile("gopher.png", reader).
			AddActionRow(
				discord.NewPrimaryButton("test1", "test1", nil),
				discord.NewPrimaryButton("test2", "test2", nil),
				discord.NewPrimaryButton("test3", "test3", nil),
				discord.NewPrimaryButton("test4", "test4", nil),
			).
			AddActionRow(
				discord.NewSelectMenu("test3", "test", 1, 1,
					discord.NewSelectOption("test1", "1"),
					discord.NewSelectOption("test2", "2"),
					discord.NewSelectOption("test3", "3"),
				),
			).
			Build(),
		); err != nil {
			log.Errorf("error sending interaction response: %s", err)
		}

	case "addrole":
		user := event.Option("member").User()
		role := event.Option("role").Role()

		if err := event.Disgo().RestServices().AddMemberRole(*event.Interaction.GuildID, user.ID, role.ID); err == nil {
			_ = event.Reply(discord.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(discord.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := event.Option("member").User()
		role := event.Option("role").Role()

		if err := event.Disgo().RestServices().RemoveMemberRole(*event.Interaction.GuildID, user.ID, role.ID); err == nil {
			_ = event.Reply(discord.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(discord.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(16711680).SetDescriptionf("Failed to remove %s from %s", role, user).Build(),
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
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.MessageChannel().SendMessage(discord.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				log.Errorf("error while sending file: %s", err)
				return
			}
			time.Sleep(time.Second * 2)

			embed := core.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = message.Update(discord.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = message.Update(discord.NewMessageUpdateBuilder().SetContent("").SetEmbeds(core.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.SendMessage(discord.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()
	}
}
