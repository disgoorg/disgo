package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/json"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Infof("disgo version: %s", disgo.Version)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentAutoModerationConfiguration, gateway.IntentAutoModerationExecution)),
		bot.WithEventListenerFunc(func(event *events.Ready) {
			go showCaseAutoMod(event.Client())
		}),
		bot.WithEventListenerFunc(func(event *events.AutoModerationRuleCreate) {
			fmt.Printf("rule created: %#v\n", event.AutoModerationRule)
		}),
		bot.WithEventListenerFunc(func(event *events.AutoModerationRuleUpdate) {
			fmt.Printf("rule updated: %#v\n", event.AutoModerationRule)
		}),
		bot.WithEventListenerFunc(func(event *events.AutoModerationRuleDelete) {
			fmt.Printf("rule deleted: %#v\n", event.AutoModerationRule)
		}),
		bot.WithEventListenerFunc(func(event *events.AutoModerationActionExecution) {
			fmt.Printf("action executed: %#v\n", event.EventAutoModerationActionExecution)
		}),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func showCaseAutoMod(client bot.Client) {
	rule, err := client.Rest().CreateAutoModerationRule(guildID, discord.AutoModerationRuleCreate{
		Name:        "test-rule",
		EventType:   discord.AutoModerationEventTypeMessageSend,
		TriggerType: discord.AutoModerationTriggerTypeKeyword,
		TriggerMetadata: &discord.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*test*"},
		},
		Actions: []discord.AutoModerationAction{
			{
				Type: discord.AutoModerationActionTypeSendAlertMessage,
				Metadata: &discord.AutoModerationActionMetadata{
					ChannelID: channelID,
				},
			},
			{
				Type: discord.AutoModerationActionTypeBlockMessage,
			},
		},
		Enabled: json.Ptr(true),
	})
	if err != nil {
		log.Error("error while creating rule: ", err)
		return
	}

	time.Sleep(time.Second * 10)

	rule, err = client.Rest().UpdateAutoModerationRule(guildID, rule.ID, discord.AutoModerationRuleUpdate{
		Name: json.Ptr("test-rule-updated"),
		TriggerMetadata: &discord.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*test2*"},
		},
		Actions: &[]discord.AutoModerationAction{
			{
				Type: discord.AutoModerationActionTypeSendAlertMessage,
				Metadata: &discord.AutoModerationActionMetadata{
					ChannelID: channelID,
				},
			},
		},
	})
	if err != nil {
		log.Error("error while updating rule: ", err)
		return
	}

	time.Sleep(time.Second * 10)

	err = client.Rest().DeleteAutoModerationRule(guildID, rule.ID)
	if err != nil {
		log.Error("error while deleting rule: ", err)
		return
	}

}
