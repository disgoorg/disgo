package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DisgoOrg/disgo/collectors"
	"github.com/DisgoOrg/disgo/events"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
	"github.com/PaesslerAG/gval"
)

var listener = &events.ListenerAdapter{
	OnGuildMessageCreate: messageListener,
	OnSlashCommand:       slashCommandListener,
	OnButtonClick:        buttonClickListener,
	OnSelectMenuSubmit:   selectMenuSubmitListener,
}

func buttonClickListener(event *events.ButtonClickEvent) {
	switch event.CustomID {
	case "test1":
		_ = event.Respond(discord.InteractionCallbackTypeChannelMessageWithSource,
			core.NewMessageCreateBuilder().
				SetContent(event.CustomID).
				Build(),
		)

	case "test2":
		_ = event.Respond(discord.InteractionCallbackTypeDeferredChannelMessageWithSource, nil)

	case "test3":
		_ = event.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil)

	case "test4":
		_ = event.Respond(discord.InteractionCallbackTypeUpdateMessage,
			core.NewMessageCreateBuilder().
				SetContent(event.CustomID).
				Build(),
		)
	}
}

func selectMenuSubmitListener(event *events.SelectMenuSubmitEvent) {
	switch event.CustomID {
	case "test3":
		if err := event.DeferUpdate(); err != nil {
			log.Errorf("error sending interaction response: %s", err)
		}
		_, _ = event.CreateFollowup(core.NewMessageCreateBuilder().
			SetEphemeral(true).
			SetContentf("selected options: %s", event.Values).
			Build(),
		)
	}
}

func slashCommandListener(event *events.SlashCommandEvent) {
	switch event.CommandName {
	case "eval":
		go func() {
			code := event.Options["code"].String()
			embed := core.NewEmbedBuilder().
				SetColor(orange).
				AddField("Status", "...", true).
				AddField("Time", "...", true).
				AddField("Code", "```go\n"+code+"\n```", false).
				AddField("Output", "```\n...\n```", false)
			_ = event.Create(core.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

			start := time.Now()
			output, err := gval.Evaluate(code, map[string]interface{}{
				"bot":   event.Bot(),
				"event": event,
			})

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, err = event.UpdateOriginal(core.NewMessageUpdateBuilder().
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
			_, err = event.UpdateOriginal(core.NewMessageUpdateBuilder().
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
		_ = event.Create(core.NewMessageCreateBuilder().
			SetContent(event.Options["message"].String()).
			SetEphemeral(event.Options["ephemeral"].Bool()).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		go func() {
			_ = event.DeferCreate(true)
			members, err := event.Guild().LoadAllMembers()
			if err != nil {
				_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().SetContentf("failed to load members. error: %s", err).Build())
			}
			_, _ = event.UpdateOriginal(core.NewMessageUpdateBuilder().
				SetContentf("loaded %d members", len(members)).
				Build(),
			)
		}()

	case "addrole":
		user := event.Options["member"].User()
		role := event.Options["role"].Role()

		if err := event.Bot().RestServices.GuildService().AddMemberRole(*event.GuildID, user.ID, role.ID); err == nil {
			_ = event.Create(core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Create(core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := event.Options["member"].User()
		role := event.Options["role"].Role()

		if err := event.Bot().RestServices.GuildService().RemoveMemberRole(*event.GuildID, user.ID, role.ID); err == nil {
			_ = event.Create(core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Create(core.NewMessageCreateBuilder().AddEmbeds(
				core.NewEmbedBuilder().SetColor(16711680).SetDescriptionf("Failed to remove %s from %s", role, user).Build(),
			).Build())
		}
	}
}

func messageListener(event *events.GuildMessageCreateEvent) {
	if event.Message.Author.IsBot {
		return
	}

	switch event.Message.Content {
	case "panic":
		panic("panic in the disco")

	case "party":
		_, _ = event.Message.Reply(core.NewMessageCreateBuilder().AddStickers("886756806888673321").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "ping":
		_, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.Channel().CreateMessage(core.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				log.Errorf("error while sending file: %s", err)
				return
			}
			time.Sleep(time.Second * 2)

			embed := core.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = message.Update(core.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = message.Update(core.NewMessageUpdateBuilder().SetContent("").SetEmbeds(core.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.CreateMessage(core.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()

	case "repeat":
		go func() {
			ch, cls := collectors.NewMessageCollectorByChannel(event.Channel(), func(m *core.Message) bool {
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

				_, _ = msg.Reply(core.NewMessageCreateBuilder().SetContentf("Content: %s, Count: %v", msg.Content, count).Build())
			}
		}()

	}
}
