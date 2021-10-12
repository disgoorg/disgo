package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/events"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = discord.Snowflake(os.Getenv("disgo_guild_id"))

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:              "root-command",
			Description:       "root command",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "sub-command-group",
					Description: "sub command group",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "sub-command",
							Description: "sub command",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "string",
									Description: "string option",
									Required:    false,
									Choices: []discord.ApplicationCommandOptionChoiceString{
										{
											Name:  "0",
											Value: "0",
										},
										{
											Name:  "1",
											Value: "1",
										},
									},
									Autocomplete: false,
								},
								discord.ApplicationCommandOptionString{
									Name:         "string",
									Description:  "string option",
									Required:     false,
									Autocomplete: true,
								},
								discord.ApplicationCommandOptionInt{
									Name:        "int",
									Description: "int option",
									Required:    false,
									Choices: []discord.ApplicationCommandOptionChoiceInt{
										{
											Name:  "0",
											Value: 0,
										},
										{
											Name:  "1",
											Value: 0,
										},
									},
									Autocomplete: false,
								},
								discord.ApplicationCommandOptionInt{
									Name:         "int",
									Description:  "int option",
									Required:     false,
									Autocomplete: true,
								},
								discord.ApplicationCommandOptionBool{
									Name:        "bool",
									Description: "bool option",
									Required:    false,
								},
								discord.ApplicationCommandOptionUser{
									Name:        "user",
									Description: "user option",
									Required:    false,
								},
								discord.ApplicationCommandOptionChannel{
									Name:         "channel",
									Description:  "channel option",
									Required:     false,
									ChannelTypes: []discord.ChannelType{discord.ChannelTypeText},
								},
								discord.ApplicationCommandOptionRole{
									Name:        "role",
									Description: "role option",
									Required:    false,
								},
								discord.ApplicationCommandOptionMentionable{
									Name:        "mentionable",
									Description: "mentionable option",
									Required:    false,
								},
								discord.ApplicationCommandOptionFloat{
									Name:        "float",
									Description: "float option",
									Required:    false,
									Choices: []discord.ApplicationCommandOptionChoiceFloat{
										{
											Name:  "1.1",
											Value: 1.1,
										},
										{
											Name:  "2.1",
											Value: 2.1,
										},
									},
									Autocomplete: false,
								},
								discord.ApplicationCommandOptionFloat{
									Name:         "float",
									Description:  "float option",
									Required:     false,
									Autocomplete: true,
								},
							},
						},
					},
				},
			},
		},
		discord.UserCommandCreate{
			Name:              "test",
			DefaultPermission: true,
		},
		discord.MessageCommandCreate{
			Name:              "test",
			DefaultPermission: true,
		},
	}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", info.Version)

	disgo, err := bot.New(token,
		bot.WithGatewayOpts(gateway.WithGatewayIntents(discord.GatewayIntentsNone)),
		bot.WithCacheOpts(core.WithCacheFlags(core.CacheFlagsDefault)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnSlashCommand: commandListener,
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer disgo.Close()

	_, err = disgo.SetGuildCommands(guildID, commands)
	if err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = disgo.ConnectGateway(); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func commandListener(event *events.SlashCommandEvent) {
	if event.CommandName == "say" {
		err := event.Create(core.NewMessageCreateBuilder().
			SetContent(event.Options["message"].String()).
			Build(),
		)
		if err != nil {
			event.Bot().Logger.Error("error on sending response: ", err)
		}
	}
}
