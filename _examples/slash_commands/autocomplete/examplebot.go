package main

import (
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/log"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = discord.Snowflake(os.Getenv("disgo_guild_id"))

	//go:embed data.json
	rawData []byte

	autocompleteData   map[string][]string
	autocompleteGroups []string
)

func main() {
	log.SetDefault(log.New(log.LstdFlags | log.Lshortfile))
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("bot version: %s", info.Version)

	err := json.Unmarshal(rawData, &autocompleteData)
	if err != nil {
		log.Error("failed to parse rawData: ", err)
	}
	autocompleteGroups = make([]string, len(autocompleteData))
	i := 0
	for group := range autocompleteData {
		autocompleteGroups[i] = group
		i++
	}

	disgo, err := bot.New(token,
		bot.WithRawEventsEnabled(),
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentsNone),
			gateway.WithCompress(true),
		),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnSlashCommand:                   slashCommandListener,
			OnApplicationCommandAutocomplete: applicationCommandAutocompleteListener,
		}),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	registerCommands(disgo)

	err = disgo.ConnectGateway()
	if err != nil {
		log.Fatalf("error while connecting to discord: %s", err)
	}

	defer disgo.Close()

	log.Infof("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func applicationCommandAutocompleteListener(event *events.ApplicationCommandAutocompleteEvent) {
	switch event.CommandName {
	case "autocomplete":
		go func() {
			focused := event.FocusedOption()
			var targets []string
			if focused.Name == "group" {
				targets = autocompleteGroups
			} else {
				targets = autocompleteData[event.Options["group"].String()]
			}
			result := fuzzy.FindFold(focused.String(), targets)

			if focused.Name != "group" && focused.String() != "" {
				if len(result) > 24 {
					result = result[:24]
				}
				result = append([]string{focused.String()}, result...)
			} else if len(result) > 25 {
				result = result[:25]
			}

			choices := make([]discord.ApplicationCommandOptionChoice, len(result))
			for i, value := range result {
				choices[i] = discord.ApplicationCommandOptionChoice{
					Name:  value,
					Value: value,
				}
			}
			if err := event.Result(choices); err != nil {
				event.Bot().Logger.Error("failed to return autocomplete choices: ", err)
			}
		}()
	}
}

func slashCommandListener(event *events.SlashCommandEvent) {
	switch event.CommandName {
	case "autocomplete":
		_ = event.Create(core.NewMessageCreateBuilder().SetContentf("you selected `%v` of group `%v`", event.Options["value"].String(), event.Options["group"].String()).Build())
	}
}

func registerCommands(bot *core.Bot) {
	if _, err := bot.SetGuildCommands(guildID, []discord.ApplicationCommandCreate{
		{
			Type:              discord.ApplicationCommandTypeSlash,
			Name:              "autocomplete",
			Description:       "autocomplete",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
				{
					Type:         discord.ApplicationCommandOptionTypeString,
					Name:         "group",
					Description:  "group",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discord.ApplicationCommandOptionTypeString,
					Name:         "value",
					Description:  "value",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}); err != nil {
		log.Fatalf("error while registering guild commands: %s", err)
	}
}
