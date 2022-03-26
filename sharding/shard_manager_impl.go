package sharding

import (
	"context"
	"sync"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding/srate"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake"
)

var _ ShardManager = (*shardManagerImpl)(nil)

func New(token string, eventHandlerFunc gateway.EventHandlerFunc, opts ...ConfigOpt) ShardManager {
	config := DefaultConfig()
	config.Apply(opts)

	return &shardManagerImpl{
		shards:           NewShardsMap(),
		token:            token,
		eventHandlerFunc: eventHandlerFunc,
		config:           *config,
	}
}

type shardManagerImpl struct {
	shards *ShardsMap

	token            string
	eventHandlerFunc gateway.EventHandlerFunc
	config           Config
}

func (m *shardManagerImpl) Logger() log.Logger {
	return m.config.Logger
}

func (m *shardManagerImpl) RateLimiter() srate.Limiter {
	return m.config.RateLimiter
}

func (m *shardManagerImpl) Open(ctx context.Context) {
	m.Logger().Infof("opening %s shards...", m.config.Shards)
	var wg sync.WaitGroup

	for _, shardInt := range m.config.Shards.Values() {
		shardID := shardInt
		if m.shards.Has(shardID) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer m.RateLimiter().UnlockBucket(shardID)
			if err := m.RateLimiter().WaitBucket(ctx, shardID); err != nil {
				m.Logger().Errorf("failed to wait shard bucket %d: %s", shardID, err)
				return
			}

			shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(m.config.ShardCount))...)
			m.shards.Set(shardID, shard)
			if err := shard.Open(ctx); err != nil {
				m.Logger().Errorf("failed to open shard %d: %s", shardID, err)
			}
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) ReOpen(ctx context.Context) {
	m.Logger().Infof("reopening %s shards...", m.config.Shards)
	var wg sync.WaitGroup

	for shardID := range m.shards.AllIDs() {
		wg.Add(1)
		shard := m.shards.Get(shardID)
		go func() {
			defer wg.Done()
			if shard != nil {
				shard.Close(ctx)
			}
			if err := shard.ReOpen(ctx, time.Second); err != nil {
				m.Logger().Errorf("failed to reopen shard %d: %s", shard.ShardID(), err)
			}
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) Close(ctx context.Context) {
	m.Logger().Infof("closing %v shards...", m.config.Shards)
	var wg sync.WaitGroup

	for shardID := range m.shards.AllIDs() {
		m.config.Shards.Delete(shardID)
		shard := m.shards.Delete(shardID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard.Close(ctx)
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) OpenShard(ctx context.Context, shardID int) error {
	m.Logger().Infof("opening shard %d...", shardID)
	shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(m.config.ShardCount))...)
	m.config.Shards.Add(shardID)
	m.shards.Set(shardID, shard)
	return shard.Open(ctx)
}

func (m *shardManagerImpl) ReOpenShard(ctx context.Context, shardID int) error {
	m.Logger().Infof("reopening shard %d...", shardID)
	shard := m.shards.Get(shardID)
	if shard != nil {
		shard.Close(ctx)
	}
	return shard.Open(ctx)
}

func (m *shardManagerImpl) CloseShard(ctx context.Context, shardID int) {
	m.Logger().Infof("closing shard %d...", shardID)
	m.config.Shards.Delete(shardID)
	shard := m.shards.Delete(shardID)
	if shard != nil {
		shard.Close(ctx)
	}
}

func (m *shardManagerImpl) GetGuildShard(guildId snowflake.Snowflake) gateway.Gateway {
	return m.Shard(ShardIDByGuild(guildId, m.config.ShardCount))
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	return m.shards.Get(shardID)
}

func (m *shardManagerImpl) Shards() *ShardsMap {
	return m.shards
}
