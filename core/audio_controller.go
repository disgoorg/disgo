package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// AudioController lets you Connect / Disconnect from a Channel
type AudioController interface {
	// Bot returns the core.Bot instance
	Bot() *Bot

	// Connect sends a discord.GatewayCommand to connect to the specified Channel
	Connect(ctx context.Context, guildID discord.Snowflake, channelID discord.Snowflake) error

	// Disconnect sends a discord.GatewayCommand to disconnect from a Channel
	Disconnect(ctx context.Context, guildID discord.Snowflake) error
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

func (c *audioControllerImpl) Connect(ctx context.Context, guildID discord.Snowflake, channelID discord.Snowflake) error {
	shard, err := c.getShard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodeVoiceStateUpdate, discord.UpdateVoiceStateCommandData{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}

func (c *audioControllerImpl) Disconnect(ctx context.Context, guildID discord.Snowflake) error {
	shard, err := c.getShard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodeVoiceStateUpdate, discord.UpdateVoiceStateCommandData{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *audioControllerImpl) getShard(guildID discord.Snowflake) (gateway.Gateway, error) {
	shard, err := c.Bot().Shard(guildID)
	if err != nil {
		return nil, err
	}
	if !shard.Status().IsConnected() {
		return nil, discord.ErrShardNotConnected
	}
	return shard, nil
}
