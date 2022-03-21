package bot

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/snowflake"
)

// AudioController lets you Connect / Disconnect from a Channel
type AudioController interface {
	// Client returns the bot.Client instance
	Client() Client

	// Connect sends a discord.GatewayCommand to connect to the specified Channel
	Connect(ctx context.Context, guildID snowflake.Snowflake, channelID snowflake.Snowflake) error

	// Disconnect sends a discord.GatewayCommand to disconnect from a Channel
	Disconnect(ctx context.Context, guildID snowflake.Snowflake) error
}

func NewAudioController(client Client) AudioController {
	return &audioControllerImpl{bot: client}
}

type audioControllerImpl struct {
	bot Client
}

func (c *audioControllerImpl) Client() Client {
	return c.bot
}

func (c *audioControllerImpl) Connect(ctx context.Context, guildID snowflake.Snowflake, channelID snowflake.Snowflake) error {
	shard, err := c.getShard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodeVoiceStateUpdate, discord.UpdateVoiceStateCommandData{
		GuildID:   guildID,
		ChannelID: &channelID,
	}))
}

func (c *audioControllerImpl) Disconnect(ctx context.Context, guildID snowflake.Snowflake) error {
	shard, err := c.getShard(guildID)
	if err != nil {
		return err
	}
	return shard.Send(ctx, discord.NewGatewayCommand(discord.GatewayOpcodeVoiceStateUpdate, discord.UpdateVoiceStateCommandData{
		GuildID:   guildID,
		ChannelID: nil,
	}))
}

func (c *audioControllerImpl) getShard(guildID snowflake.Snowflake) (gateway.Gateway, error) {
	shard, err := c.Client().Shard(guildID)
	if err != nil {
		return nil, err
	}
	if shard.Status() != gateway.StatusReady {
		return nil, discord.ErrShardNotConnected
	}
	return shard, nil
}
