package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/util"

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
	log.Infof("disgo %s", util.Version)

	disgo, err := core.NewBuilder(token).
		SetRawEventsEnabled(true).
		SetHTTPClient(client).
		SetGatewayConfig(gateway.Config{
			GatewayIntents: gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMembers,
		}).
		SetHTTPServerConfig(httpserver.Config{
			URL:       "/interactions/callback",
			Port:      ":80",
			PublicKey: "your public key from the developer dashboard",
		}).
		SetCacheConfig(core.CacheConfig{
			CacheFlags:        core.CacheFlagsDefault,
			MemberCachePolicy: core.MemberCachePolicyAll,
		}).
		AddEventListeners(&events.ListenerAdapter{
			OnRawGateway:         rawGatewayEventListener,
			OnGuildAvailable:     guildAvailListener,
			OnGuildMessageCreate: messageListener,
			OnSlashCommand:       commandListener,
			OnButtonClick:        buttonClickListener,
			OnSelectMenuSubmit:   selectMenuSubmitListener,
		}).
		Build()

	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	rawCmds := []discord.ApplicationCommandCreate{
		{
			Name:              "eval",
			Description:       "runs some go code",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
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
			Options: []discord.ApplicationCommandOption{
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
			Options: []discord.ApplicationCommandOption{
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
			Options: []discord.ApplicationCommandOption{
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
	cmds, err := disgo.RestServices().ApplicationService().SetGuildCommands(context.Background(), disgo.ApplicationID(), guildID, rawCmds...)
	if err != nil {
		log.Errorf("error while registering guild commands: %s", err)
	}

	var cmdsPermissions []discord.GuildCommandPermissionsSet
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
			cmdsPermissions = append(cmdsPermissions, discord.GuildCommandPermissionsSet{
				ID:          cmd.ID,
				Permissions: []discord.CommandPermission{perms},
			})
		}
		cmdsPermissions = append(cmdsPermissions, discord.GuildCommandPermissionsSet{
			ID:          cmd.ID,
			Permissions: []discord.CommandPermission{perms},
		})
	}
	if _, err = disgo.SetGuildCommandsPermissions(context.Background(), guildID, cmdsPermissions...); err != nil {
		log.Errorf("error while setting command permissions: %s", err)
	}

	err = disgo.Connect()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	defer disgo.Close()

	log.Infof("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent) {
	log.Infof("guild loaded: %s", event.Guild.ID)
}

func rawGatewayEventListener(event *events.RawEvent) {
	if event.Type == string(gateway.EventTypeInteractionCreate) {
		println(string(event.RawPayload))
	}
}

func buttonClickListener(event *events.ButtonClickEvent) {
	switch event.CustomID() {
	case "test1":
		_ = event.Respond(context.Background(), discord.InteractionResponseTypeChannelMessageWithSource,
			core.NewMessageCreateBuilder().
				SetContent(event.CustomID()).
				Build(),
		)

	case "test2":
		_ = event.Respond(context.Background(), discord.InteractionResponseTypeDeferredChannelMessageWithSource, nil)

	case "test3":
		_ = event.Respond(context.Background(), discord.InteractionResponseTypeDeferredUpdateMessage, nil)

	case "test4":
		_ = event.Respond(context.Background(), discord.InteractionResponseTypeUpdateMessage,
			core.NewMessageCreateBuilder().
				SetContent(event.CustomID()).
				Build(),
		)
	}
}

func selectMenuSubmitListener(event *events.SelectMenuSubmitEvent) {
	switch event.CustomID() {
	case "test3":
		if err := event.DeferUpdate(); err != nil {
			log.Errorf("error sending interaction response: %s", err)
		}
		_, _ = event.CreateFollowup(context.Background(), core.NewMessageCreateBuilder().
			SetEphemeral(true).
			SetContentf("selected options: %s", event.Values()).
			Build(),
		)
	}
}

func commandListener(event *events.SlashCommandEvent) {
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
			_ = event.Reply(context.Background(), core.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

			start := time.Now()
			output, err := gval.Evaluate(code, map[string]interface{}{
				"disgo": event.Disgo(),
				"dgo":   event.Disgo(),
				"event": event,
			})

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, err = event.Interaction.UpdateOriginal(context.Background(), core.NewMessageUpdateBuilder().
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
			_, err = event.Interaction.UpdateOriginal(context.Background(), core.NewMessageUpdateBuilder().
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
		_ = event.Reply(context.Background(), core.NewMessageCreateBuilder().
			SetContent(event.Option("message").String()).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		reader, _ := os.Open("gopher.png")
		if err := event.Reply(context.Background(), core.NewMessageCreateBuilder().
			SetContent("test message").
			AddFile("gopher.png", reader).
			AddActionRow(
				core.NewPrimaryButton("test1", "test1", nil),
				core.NewPrimaryButton("test2", "test2", nil),
				core.NewPrimaryButton("test3", "test3", nil),
				core.NewPrimaryButton("test4", "test4", nil),
			).
			AddActionRow(
				core.NewSelectMenu("test3", "test", 1, 1,
					core.NewSelectOption("test1", "1"),
					core.NewSelectOption("test2", "2"),
					core.NewSelectOption("test3", "3"),
				),
			).
			Build(),
		); err != nil {
			log.Errorf("error sending interaction response: %s", err)
		}

	case "addrole":
		user := event.Option("member").User()
		role := event.Option("role").Role()

		if err := event.Disgo().RestServices().GuildService().AddMemberRole(context.Background(), *event.Interaction.GuildID, user.ID, role.ID); err == nil {
			_ = event.Reply(context.Background(), core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(context.Background(), core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := event.Option("member").User()
		role := event.Option("role").Role()

		if err := event.Disgo().RestServices().GuildService().RemoveMemberRole(context.Background(), *event.Interaction.GuildID, user.ID, role.ID); err == nil {
			_ = event.Reply(context.Background(), core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(context.Background(), core.NewMessageCreateBuilder().AddEmbeds(
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
		_, _ = event.Message.Reply(context.Background(), core.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(context.Background(), core.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.MessageChannel().CreateMessage(context.Background(), core.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				log.Errorf("error while sending file: %s", err)
				return
			}
			time.Sleep(time.Second * 2)

			embed := core.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = message.Update(context.Background(), core.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = message.Update(context.Background(), core.NewMessageUpdateBuilder().SetContent("").SetEmbeds(core.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel(context.Background())
			if err != nil {
				_ = event.Message.AddReaction(context.Background(), "❌")
				return
			}
			_, err = channel.CreateMessage(context.Background(), core.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction(context.Background(), "✅")
			} else {
				_ = event.Message.AddReaction(context.Background(), "❌")
			}
		}()

	case "repeat":
		go func() {
			ch, cls := event.MessageChannel().CollectMessages(func(m *core.Message) bool {
				return !m.Author.IsBot && m.ChannelID == event.ChannelID
			})

			var count = 0
			for {
				count++
				if count >= 10 {
					cls()
					return
				}

				msg, ok := <-ch

				if !ok {
					return
				}

				_, _ = msg.Reply(context.Background(), core.NewMessageCreateBuilder().SetContentf("Content: %s, Count: %v", *msg.Content, count).Build())
			}
		}()

	}
}
