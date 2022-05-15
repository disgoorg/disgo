package sharding

import (
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding/srate"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
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
	shards         *ShardsMap
	guildsPerShard map[int]int

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

func (m *shardManagerImpl) onEvent(gatewayEventType discord.GatewayEventType, sequenceNumber discord.GatewaySequence, shardID int, payload io.Reader) {
	if m.config.AutoScaling {
		var bufferRead bytes.Buffer
		payloadReader := io.TeeReader(payload, &bufferRead)
		var err error
		switch gatewayEventType {
		case discord.GatewayEventTypeReady:
			var event discord.GatewayEventReady
			if err = json.NewDecoder(payloadReader).Decode(&event); err != nil {
				m.Logger().Error("failed to decode ready event", "error", err)
				break
			}
			m.guildsPerShard[shardID] = len(event.Guilds)

		case discord.GatewayEventTypeGuildCreate:
			var event discord.GatewayGuild
			if err = json.NewDecoder(payloadReader).Decode(&event); err != nil {
				m.Logger().Error("failed to decode guild create event", "error", err)
				break
			}
			m.guildsPerShard[shardID] += 1

		case discord.GatewayEventTypeGuildDelete:
			var event discord.UnavailableGuild
			if err = json.NewDecoder(payloadReader).Decode(&event); err != nil {
				m.Logger().Error("failed to decode guild delete event", "error", err)
				break
			}
			if !event.Unavailable {
				m.guildsPerShard[shardID] -= 1
			}
		}

	}
	m.eventHandlerFunc(gatewayEventType, sequenceNumber, shardID, payload)
}

func (m *shardManagerImpl) Open(ctx context.Context) {
	m.Logger().Debugf("opening %s shards...", m.config.Shards)
	var wg sync.WaitGroup

	for shardInt := range m.config.ShardIDs {
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
	m.Logger().Debugf("reopening %s shards...", m.config.Shards)
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
	m.Logger().Debugf("closing %v shards...", m.config.Shards)
	var wg sync.WaitGroup

	for shardID := range m.shards.AllIDs() {
		shard := m.shards.Delete(shardID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard.Close(ctx)
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) OpenShard(ctx context.Context, shardID int, shardCount int) error {
	m.Logger().Debugf("opening shard %d...", shardID)
	shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(m.config.ShardCount))...)
	m.config.Shards.Add(shardID)
	m.shards.Set(shardID, shard)
	return shard.Open(ctx)
}

func (m *shardManagerImpl) ReOpenShard(ctx context.Context, shardID int) error {
	m.Logger().Debugf("reopening shard %d...", shardID)
	shard := m.shards.Get(shardID)
	if shard != nil {
		shard.Close(ctx)
	}
	return shard.Open(ctx)
}

func (m *shardManagerImpl) CloseShard(ctx context.Context, shardID int) {
	m.Logger().Debugf("closing shard %d...", shardID)
	shard := m.shards.Delete(shardID)
	if shard != nil {
		shard.Close(ctx)
	}
}

func (m *shardManagerImpl) ShardByGuildID(guildId snowflake.ID) gateway.Gateway {
	var shardCount int
	m.shards.For(func(shardID int, shard gateway.Gateway) {
		if shard.ShardCount() > shardCount {
			shardCount = shard.ShardCount()
		}
	})

	var shard gateway.Gateway
	for shard == nil || shardCount != 0 {
		shard = m.Shard(ShardIDByGuild(guildId, shardCount))
		shardCount /= 2
	}
	return shard
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	return m.shards.Get(shardID)
}

func (m *shardManagerImpl) Shards() *ShardsMap {
	return m.shards
}
