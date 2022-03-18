package main

import (
	"bytes"
	"strings"
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

var listener = &events.ListenerAdapter{
	OnGuildMessageCreate:            messageListener,
	OnApplicationCommandInteraction: applicationCommandListener,
	OnComponentInteraction:          componentListener,
	OnModalSubmit:                   modalListener,
}

func modalListener(event *events.ModalSubmitInteractionEvent) {
	switch event.Data.CustomID {
	case "test1":
		value := *event.Data.Components.Text("test_input")
		_ = event.CreateMessage(discord.MessageCreate{Content: value})

	case "test2":
		value := *event.Data.Components.Text("test_input")
		_ = event.DeferCreateMessage(false)
		go func() {
			time.Sleep(time.Second * 5)
			_, _ = event.UpdateOriginalMessage(discord.MessageUpdate{Content: &value})
		}()

	case "test3":
		value := *event.Data.Components.Text("test_input")
		_ = event.UpdateMessage(discord.MessageUpdate{Content: &value})

	case "test4":
		_ = event.DeferUpdateMessage()
	}
}

func componentListener(event *events.ComponentInteractionEvent) {
	switch data := event.Data.(type) {
	case core.ButtonInteractionData:
		ids := strings.Split(data.CustomID.String(), ":")
		switch ids[0] {
		case "modal":
			_ = event.CreateModal(discord.ModalCreate{
				CustomID: discord.CustomID("test" + ids[1]),
				Title:    "Test" + ids[1] + " Modal",
				Components: []discord.ContainerComponent{
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID:    "test_input",
							Style:       discord.TextInputStyleShort,
							Label:       "qwq",
							Required:    true,
							Placeholder: "test placeholder",
							Value:       "uwu",
						},
					},
				},
			})

		case "test1":
			_ = event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent(data.CustomID.String()).
				Build(),
			)

		case "test2":
			_ = event.DeferCreateMessage(false)

		case "test3":
			_ = event.DeferUpdateMessage()

		case "test4":
			_ = event.UpdateMessage(discord.NewMessageUpdateBuilder().
				SetContent(data.CustomID.String()).
				Build(),
			)
		}

	case core.SelectMenuInteractionData:
		switch data.CustomID {
		case "test3":
			if err := event.DeferUpdateMessage(); err != nil {
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
	case "locale":
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Guild Locale: %s\nLocale: %s", event.GuildLocale, event.Locale).
			Build(),
		)
		if err != nil {
			event.Bot().Logger().Error("error on sending response: ", err)
		}

	case "say":
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(*data.Options.String("message")).
			SetEphemeral(*data.Options.Bool("ephemeral")).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("test").
			AddActionRow(
				discord.NewPrimaryButton("test1", "modal:1"),
				discord.NewPrimaryButton("test2", "modal:2"),
				discord.NewPrimaryButton("test3", "modal:3"),
				discord.NewPrimaryButton("test4", "modal:4"),
			).
			Build(),
		)
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
