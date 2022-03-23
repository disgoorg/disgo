package main

import (
	"bytes"
	"strings"
	"time"

	"github.com/DisgoOrg/log"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
		_ = event.Respond(discord.InteractionCallbackTypeCreateMessage, discord.MessageCreate{Content: event.Data.Text("test_input")})

	case "test2":
		value := event.Data.Text("test_input")
		_ = event.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, nil)
		go func() {
			time.Sleep(time.Second * 5)
			_, _ = event.Client().Rest().InteractionService().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{Content: &value})
		}()

	case "test3":
		value := event.Data.Text("test_input")
		_ = event.Respond(discord.InteractionCallbackTypeUpdateMessage, discord.MessageUpdate{Content: &value})

	case "test4":
		_ = event.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil)
	}
}

func componentListener(event *events.ComponentInteractionEvent) {
	switch data := event.Data.(type) {
	case discord.ButtonInteractionData:
		ids := strings.Split(data.CustomID().String(), ":")
		switch ids[0] {
		case "modal":
			_ = event.Respond(discord.InteractionCallbackTypeModal, discord.ModalCreate{
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
			_ = event.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, discord.NewMessageCreateBuilder().
				SetContent(data.CustomID().String()).
				Build(),
			)

		case "test2":
			_ = event.Respond(discord.InteractionCallbackTypeDeferredCreateMessage, nil)

		case "test3":
			_ = event.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil)

		case "test4":
			_ = event.Respond(discord.InteractionCallbackTypeUpdateMessage, discord.NewMessageUpdateBuilder().
				SetContent(data.CustomID().String()).
				Build(),
			)
		}

	case discord.SelectMenuInteractionData:
		switch data.CustomID() {
		case "test3":
			if err := event.Respond(discord.InteractionCallbackTypeDeferredUpdateMessage, nil); err != nil {
				log.Errorf("error sending interaction response: %s", err)
			}
			_, _ = event.Client().Rest().InteractionService().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("selected options: %s", data.Values).
				Build(),
			)
		}
	}
}

func applicationCommandListener(event *events.ApplicationCommandInteractionEvent) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "locale":
		err := event.Respond(discord.InteractionCallbackTypeCreateMessage, discord.NewMessageCreateBuilder().
			SetContentf("Guild Locale: %s\nLocale: %s", event.GuildLocale(), event.Locale()).
			Build(),
		)
		if err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}

	case "say":
		_ = event.Respond(discord.InteractionCallbackTypeCreateMessage, discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			ClearAllowedMentions().
			Build(),
		)

	case "test":
		_ = event.Respond(discord.InteractionCallbackTypeCreateMessage, discord.NewMessageCreateBuilder().
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
		_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("gopher").AddFile("gopher.png", bytes.NewBuffer(gopher)).AddFile("gopher.png", bytes.NewBuffer(gopher)).Build())

	case "panic":
		panic("panic in the disco")

	case "party":
		_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().AddStickers("886756806888673321").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "ping":
		_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("pong").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("ping").SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).Build())

	case "test":
		go func() {
			message, err := event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent("test").Build())
			if err != nil {
				log.Errorf("error while sending file: %s", err)
				return
			}
			time.Sleep(time.Second * 2)

			embed := discord.NewEmbedBuilder().SetDescription("edit").Build()
			message, _ = event.Client().Rest().ChannelService().UpdateMessage(event.ChannelID, message.ID, discord.NewMessageUpdateBuilder().SetContent("edit").SetEmbeds(embed, embed).Build())

			time.Sleep(time.Second * 2)

			_, _ = event.Client().Rest().ChannelService().UpdateMessage(event.ChannelID, message.ID, discord.NewMessageUpdateBuilder().SetContent("").SetEmbeds(discord.NewEmbedBuilder().SetDescription("edit2").Build()).Build())
		}()

	case "dm":
		go func() {
			channel, err := event.Client().Rest().UserService().CreateDMChannel(event.Message.Author.ID)
			if err != nil {
				_ = event.Client().Rest().ChannelService().AddReaction(channel.ID(), event.MessageID, "❌")
				return
			}
			_, err = event.Client().Rest().ChannelService().CreateMessage(channel.ID(), discord.NewMessageCreateBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Client().Rest().ChannelService().AddReaction(channel.ID(), event.MessageID, "✅")
			} else {
				_ = event.Client().Rest().ChannelService().AddReaction(channel.ID(), event.MessageID, "❌")
			}
		}()

	case "repeat":
		go func() {
			ch, cls := bot.NewCollector(event.Client(), func(event *events.MessageCreateEvent) bool {
				return !event.Message.Author.BotUser && event.ChannelID == event.ChannelID
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
				_, _ = messageEvent.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContentf("Content: %s, Count: %v", messageEvent.Message.Content, count).SetMessageReferenceByID(event.MessageID).Build())
			}
		}()

	}
}
