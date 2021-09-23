package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// AudioController lets you Connect / Disconnect from a Channel
type AudioController interface {
	// Bot returns the core.Bot instance
	Bot() *Bot

	// Connect sends a discord.GatewayCommand to connect to a Channel
	Connect(guildID discord.Snowflake, channelID discord.Snowflake) error

	// Disconnect sends a discord.GatewayCommand to disconnect from a Channel
	Disconnect(guildID discord.Snowflake) error
}

func NewAudioController(bot *Bot) AudioController {
	return &audioControllerImpl{bot: bot}
}

type audioControllerImpl struct {
	bot *Bot
}

func (c *audioControllerImpl) Bot() *Bot {
	return c.bot
}

func (c *audioControllerImpl) Connect(guildID discord.Snowflake, channelID discord.Snowflake) error {
	shard, err := c.getShard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(discord.NewGatewayCommand(discord.OpVoiceStateUpdate, discord.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}

func (c *audioControllerImpl) Disconnect(guildID discord.Snowflake) error {
	shard, err := c.getShard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(discord.NewGatewayCommand(discord.OpVoiceStateUpdate, discord.UpdateVoiceStateCommand{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *audioControllerImpl) getShard(guildID discord.Snowflake) (gateway.Gateway, error) {
	var shard gateway.Gateway
	if c.Bot().HasGateway() {
		shard = c.Bot().Gateway
	} else if c.Bot().HasShardManager() {
		shard = c.Bot().ShardManager.GetGuildShard(guildID)
	} else {
		return nil, discord.ErrNoGatewayOrShardManager
	}
	if shard == nil {
		return nil, discord.ErrShardNotFound
	}
	if !shard.Status().IsConnected() {
		return nil, discord.ErrShardNotConnected
	}
	return shard, nil
}
