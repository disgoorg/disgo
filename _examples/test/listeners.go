package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/DisgoOrg/disgo/core/events"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
	"github.com/PaesslerAG/gval"
)

var listener = &events.ListenerAdapter{
	OnGuildMessageCreate:            messageListener,
	OnApplicationCommandInteraction: applicationCommandListener,
	OnComponentInteraction:          componentListener,
}

func componentListener(event *events.ComponentInteractionEvent) {
	switch data := event.Data.(type) {
	case *core.ButtonInteractionData:
		switch data.CustomID {
		case "test1":
			_ = event.Create(discord.NewMessageCreateBuilder().
				SetContent(data.CustomID.String()).
				Build(),
			)

		case "test2":
			_ = event.DeferCreate(false)

		case "test3":
			_ = event.DeferUpdate()

		case "test4":
			_ = event.Update(discord.NewMessageUpdateBuilder().
				SetContent(data.CustomID.String()).
				Build(),
			)
		}

	case *core.SelectMenuInteractionData:
		switch data.CustomID {
		case "test3":
			if err := event.DeferUpdate(); err != nil {
				log.Errorf("error sending interaction response: %s", err)
			}
			_, _ = event.CreateFollowupMessage(discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("selected options: %s", data.Values).
				Build(),
			)
		}
	}
}

func applicationCommandListener(event *events.ApplicationCommandInteractionEvent) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName {
	case "eval":
		go func() {
			code := *data.Options.String("code")
			embed := discord.NewEmbedBuilder().
				SetColor(orange).
				AddField("Status", "...", true).
				AddField("Time", "...", true).
				AddField("Code", "```go\n"+code+"\n```", false).
				AddField("Output", "```\n...\n```", false)
			_ = event.Create(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

			start := time.Now()
			output, err := gval.Evaluate(code, map[string]interface{}{
				"bot":   event.Bot(),
				"event": event,
			})

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, err = event.UpdateResponse(discord.NewMessageUpdateBuilder().
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
			_, err = event.UpdateResponse(discord.NewMessageUpdateBuilder().
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
		_ = event.Create(discord.NewMessageCreateBuilder().
			SetContent(*data.Options.String("message")).
			SetEphemeral(*data.Options.Bool("ephemeral")).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		go func() {
			_ = event.DeferCreate(true)
			members, err := event.Guild().RequestMembersWithQuery("", 0)
			if err != nil {
				_, _ = event.UpdateResponse(discord.NewMessageUpdateBuilder().SetContentf("failed to load members. error: %s", err).Build())
				return
			}
			_, _ = event.UpdateResponse(discord.NewMessageUpdateBuilder().
				SetContentf("loaded %d members", len(members)).
				Build(),
			)
		}()

	case "addrole":
		user := data.Options.User("member")
		role := data.Options.Role("role")

		if err := event.Bot().RestServices.GuildService().AddMemberRole(*event.GuildID, user.ID, role.ID); err == nil {
			_ = event.Create(discord.NewMessageCreateBuilder().AddEmbeds(
				discord.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Create(discord.NewMessageCreateBuilder().AddEmbeds(
				discord.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := data.Options.User("member")
		role := data.Options.Role("role")

		if err := event.Bot().RestServices.GuildService().RemoveMemberRole(*event.GuildID, user.ID, role.ID); err == nil {
			_ = event.Create(discord.NewMessageCreateBuilder().AddEmbeds(
				discord.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Create(discord.NewMessageCreateBuilder().AddEmbeds(
				discord.NewEmbedBuilder().SetColor(16711680).SetDescriptionf("Failed to remove %s from %s", role, user).Build(),
			).Build())
		}
	}
}

func messageListener(event *events.GuildMessageCreateEvent) {
	if event.Message.Author.BotUser {
		return
	}

	switch event.Message.Content {
	case "gopher":
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent("gopher").AddFile("gopher.png", bytes.NewBuffer(gopher)).AddFile("gopher.png", bytes.NewBuffer(gopher)).Build())

	case "panic":
		panic("panic in the disco")

	case "party":
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().AddStickers("886756806888673321").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "ping":
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.Channel().CreateMessage(discord.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				log.Errorf("error while sending file: %s", err)
				return
			}
			time.Sleep(time.Second * 2)

			embed := discord.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = message.Update(discord.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = message.Update(discord.NewMessageUpdateBuilder().SetContent("").SetEmbeds(discord.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.CreateMessage(discord.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()

	case "repeat":
		go func() {
			ch, cls := event.Bot().Collectors.NewMessageCollector(func(m *core.Message) bool {
				return !m.Author.BotUser && m.ChannelID == event.ChannelID
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

				_, _ = msg.Reply(discord.NewMessageCreateBuilder().SetContentf("Content: %s, Count: %v", msg.Content, count).Build())
			}
		}()

	}
}
