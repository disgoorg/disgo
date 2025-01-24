package main

import (
	"bytes"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var listener = &events.ListenerAdapter{
	OnGuildMessageCreate:            messageListener,
	OnApplicationCommandInteraction: applicationCommandListener,
	OnComponentInteraction:          componentListener,
	OnModalSubmit:                   modalListener,
	OnAutocompleteInteraction:       autocompleteListener,
}

func modalListener(event *events.ModalSubmitInteractionCreate) {
	switch event.Data.CustomID {
	case "test1":
		_ = event.CreateMessage(discord.MessageCreate{Content: event.Data.Text("test_input")})

	case "test2":
		value := event.Data.Text("test_input")
		_ = event.DeferCreateMessage(false)
		go func() {
			time.Sleep(time.Second * 5)
			_, _ = event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{Content: &value})
		}()

	case "test3":
		value := event.Data.Text("test_input")
		_ = event.UpdateMessage(discord.MessageUpdate{Content: &value})

	case "test4":
		_ = event.DeferUpdateMessage()
	}
}

func componentListener(event *events.ComponentInteractionCreate) {
	switch data := event.Data.(type) {
	case discord.ButtonInteractionData:
		ids := strings.Split(data.CustomID(), ":")
		switch ids[0] {
		case "modal":
			_ = event.Modal(discord.ModalCreate{
				CustomID: "test" + ids[1],
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
				SetContent(data.CustomID()).
				Build(),
			)

		case "test2":
			_ = event.DeferCreateMessage(false)

		case "test3":
			_ = event.DeferUpdateMessage()

		case "test4":
			_ = event.UpdateMessage(discord.NewMessageUpdateBuilder().
				SetContent(data.CustomID()).
				Build(),
			)
		}

	case discord.StringSelectMenuInteractionData:
		switch data.CustomID() {
		case "test3":
			if err := event.DeferUpdateMessage(); err != nil {
				slog.Error("error sending interaction response", slog.Any("err", err))
			}
			_, _ = event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("selected options: %s", data.Values).
				Build(),
			)
		}

	case discord.MentionableSelectMenuInteractionData:
		switch data.CustomID() {
		case "test4":
			if err := event.DeferUpdateMessage(); err != nil {
				slog.Error("error sending interaction response", slog.Any("err", err))
			}
			_, _ = event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("selected mentionable options: %s", data.Values).
				Build(),
			)
		}
	}
}

func applicationCommandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "locale":
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Guild Locale: %s\nLocale: %s", event.GuildLocale(), event.Locale()).
			Build(),
		)
		if err != nil {
			slog.Error("error on sending response", slog.Any("err", err))
		}

	case "say":
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
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

	case "test2":
		selectMenu := discord.NewMentionableSelectMenu("test4", "select users/members/roles")
		selectMenu.MaxValues = 3

		_ = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("test2").
			AddActionRow(selectMenu).
			Build(),
		)
	}
}

func autocompleteListener(event *events.AutocompleteInteractionCreate) {
	switch event.Data.CommandName {
	case "test2":
		if err := event.AutocompleteResult([]discord.AutocompleteChoice{
			discord.AutocompleteChoiceInt{
				Name:  "test1",
				Value: 1,
			},
			discord.AutocompleteChoiceInt{
				Name:  "test2",
				Value: 2,
			},
		}); err != nil {
			slog.Error("error on sending response", slog.Any("err", err))
		}
	}
}

func messageListener(event *events.GuildMessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	switch event.Message.Content {
	case "avatar":
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
			SetContentf("Avatar: %s", event.Message.Member.EffectiveAvatarURL()).
			Build(),
		)
	case "channel":
		ch, _ := event.Channel()
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
			SetContentf("channel:\n```\n%#v\n```", ch).
			Build(),
		)
	case "gopher":
		message, err := event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
			SetContent("gopher").
			AddFile("gopher.png", "this is a gopher", bytes.NewBuffer(gopher)).
			AddFile("gopher.png", "", bytes.NewBuffer(gopher)).
			Build(),
		)
		if err != nil {
			slog.Error("error on sending response", slog.Any("err", err))
		}
		time.Sleep(1 * time.Second)
		_, err = event.Client().Rest().UpdateMessage(event.ChannelID, message.ID, discord.NewMessageUpdateBuilder().
			SetContent("edited gopher").
			RetainAttachments(message.Attachments[0]).
			Build(),
		)
		if err != nil {
			slog.Error("error on updating response", slog.Any("err", err))
		}

	case "panic":
		panic("panic in the disco")

	case "party":
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().AddStickers(886756806888673321).SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "ping":
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				slog.Error("error while sending file", slog.Any("err", err))
				return
			}
			time.Sleep(time.Second * 2)

			embed := discord.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = event.Client().Rest().UpdateMessage(event.ChannelID, message.ID, discord.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = event.Client().Rest().UpdateMessage(event.ChannelID, message.ID, discord.NewMessageUpdateBuilder().SetContent("").SetEmbeds(discord.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Client().Rest().CreateDMChannel(event.Message.Author.ID)
			if err != nil {
				_ = event.Client().Rest().AddReaction(channel.ID(), event.MessageID, "❌")
				return
			}
			_, err = event.Client().Rest().CreateMessage(channel.ID(), discord.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Client().Rest().AddReaction(channel.ID(), event.MessageID, "✅")
			} else {
				_ = event.Client().Rest().AddReaction(channel.ID(), event.MessageID, "❌")
			}
		}()

	case "repeat":
		go func() {
			ch, cls := bot.NewEventCollector(event.Client(), func(event *events.MessageCreate) bool {
				return !event.Message.Author.Bot && event.ChannelID == event.ChannelID
			})

			var count = 0
			for {
				count++
				if count >= 10 {
					cls()
					return
				}
				messageEvent, ok := <-ch

				if !ok {
					return
				}
				_, _ = messageEvent.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContentf("Content: %s, Count: %v", messageEvent.Message.Content, count).SetMessageReferenceByID(messageEvent.MessageID).Build())
			}
		}()

	}
}
