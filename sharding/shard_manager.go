package sharding

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"sync"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/gateway"
)

// DefaultShardSplitCount is the default count a shard should be split into when it needs re-sharding.
const DefaultShardSplitCount = 2

// ShardManager manages multiple gateway.Gateway connections.
// For more information on sharding see: https://discord.com/developers/docs/topics/gateway#sharding
type ShardManager interface {
	// Open opens all configured shards.
	Open(ctx context.Context)

	// Close closes all shards.
	Close(ctx context.Context)

	// OpenShard opens a specific shard.
	OpenShard(ctx context.Context, shardID int) error

	// CloseShard closes a specific shard.
	CloseShard(ctx context.Context, shardID int)

	// ShardByGuildID returns the gateway.Gateway for the shard that contains the given guild.
	ShardByGuildID(guildId snowflake.ID) gateway.Gateway

	// Shard returns the gateway.Gateway for the given shard ID.
	Shard(shardID int) gateway.Gateway

	// Shards returns all shards. This function is thread-safe.
	Shards() iter.Seq[gateway.Gateway]
}

// ShardIDByGuild returns the shard ID for the given guildID and shardCount.
func ShardIDByGuild(guildID snowflake.ID, shardCount int) int {
	return int((uint64(guildID) >> 22) % uint64(shardCount))
}

var _ ShardManager = (*shardManagerImpl)(nil)

// New creates a new default ShardManager with the given token, eventHandlerFunc and ConfigOpt(s).
func New(token string, eventHandlerFunc gateway.EventHandlerFunc, opts ...ConfigOpt) ShardManager {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &shardManagerImpl{
		shards:           map[int]gateway.Gateway{},
		token:            token,
		eventHandlerFunc: eventHandlerFunc,
		config:           cfg,
	}
}

type shardManagerImpl struct {
	shards   map[int]gateway.Gateway
	shardsMu sync.Mutex

	token            string
	eventHandlerFunc gateway.EventHandlerFunc
	config           config
}

func (m *shardManagerImpl) closeHandler(shard gateway.Gateway, err error) {
	var closeError *websocket.CloseError
	if !m.config.AutoScaling || !errors.As(err, &closeError) || gateway.CloseEventCodeByCode(closeError.Code) != gateway.CloseEventCodeShardingRequired {
		return
	}
	m.config.Logger.Debug("shard requires re-sharding", slog.Int("shardID", shard.ShardID()))
	// make sure shard is closed
	shard.Close(context.TODO())

	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()

	delete(m.shards, shard.ShardID())
	delete(m.config.ShardIDs, shard.ShardID())

	newShardCount := shard.ShardCount() * m.config.ShardSplitCount
	if newShardCount > m.config.ShardCount {
		m.config.ShardCount = newShardCount
	}

	newShardID := shard.ShardID()
	var newShardIDs []int
	for len(newShardIDs) < m.config.ShardSplitCount {
		newShardIDs = append(newShardIDs, newShardID)
		newShardID += m.config.ShardSplitCount
	}

	var wg sync.WaitGroup
	for i := range newShardIDs {
		shardID := newShardIDs[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.config.RateLimiter.WaitBucket(context.TODO(), shardID); err != nil {
				m.config.Logger.Error("failed to wait shard bucket", slog.Any("err", err), slog.Int("shard_id", shardID))
				return
			}
			defer m.config.RateLimiter.UnlockBucket(shardID)

			newShard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.closeHandler, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(newShardCount))...)
			m.shards[shardID] = newShard
			if err := newShard.Open(context.TODO()); err != nil {
				m.config.Logger.Error("failed to re shard", slog.Any("err", err), slog.Int("shard_id", shardID))
			}
		}()
	}
	wg.Wait()
	m.config.Logger.Debug("re-sharded shard", slog.Int("shard_id", shard.ShardID()), slog.String("new_shard_ids", fmt.Sprint(newShardIDs)), slog.Int("new_shard_count", newShardCount))
}

func (m *shardManagerImpl) Open(ctx context.Context) {
	m.config.Logger.Debug("opening shards", slog.String("shard_ids", fmt.Sprint(m.config.ShardIDs)))
	var wg sync.WaitGroup

	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	for shardID := range m.config.ShardIDs {
		if _, ok := m.shards[shardID]; ok {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.config.RateLimiter.WaitBucket(ctx, shardID); err != nil {
				m.config.Logger.Error("failed to wait shard bucket", slog.Any("err", err), slog.Int("shard_id", shardID))
				return
			}
			defer m.config.RateLimiter.UnlockBucket(shardID)

			shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.closeHandler, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(m.config.ShardCount))...)
			m.shards[shardID] = shard
			if err := shard.Open(ctx); err != nil {
				m.config.Logger.Error("failed to open shard", slog.Any("err", err), slog.Int("shard_id", shardID))
			}
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) Close(ctx context.Context) {
	m.config.Logger.Debug("closing shards", slog.String("shard_ids", fmt.Sprint(m.config.ShardIDs)))
	var wg sync.WaitGroup

	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	for shardID := range m.shards {
		shard := m.shards[shardID]
		delete(m.shards, shardID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard.Close(ctx)
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) OpenShard(ctx context.Context, shardID int) error {
	return m.openShard(ctx, shardID, m.config.ShardCount)
}

func (m *shardManagerImpl) openShard(ctx context.Context, shardID int, shardCount int) error {
	m.config.Logger.Debug("opening shard", slog.Int("shard_id", shardID))

	if err := m.config.RateLimiter.WaitBucket(ctx, shardID); err != nil {
		return err
	}
	defer m.config.RateLimiter.UnlockBucket(shardID)
	shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.closeHandler, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(shardCount))...)

	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	m.config.ShardIDs[shardID] = struct{}{}
	m.shards[shardID] = shard
	return shard.Open(ctx)
}

func (m *shardManagerImpl) CloseShard(ctx context.Context, shardID int) {
	m.config.Logger.Debug("closing shard", slog.Int("shard_id", shardID))
	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	shard, ok := m.shards[shardID]
	if ok {
		shard.Close(ctx)
		delete(m.shards, shardID)
	}
}

func (m *shardManagerImpl) ShardByGuildID(guildId snowflake.ID) gateway.Gateway {
	shardCount := m.config.ShardCount
	var shard gateway.Gateway
	for shard == nil || shardCount != 0 {
		shard = m.Shard(ShardIDByGuild(guildId, shardCount))
		shardCount /= m.config.ShardSplitCount
	}
	return shard
}

func (m *shardManagerImpl) Shard(shardID int) gateway.Gateway {
	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	return m.shards[shardID]
}

func (m *shardManagerImpl) Shards() iter.Seq[gateway.Gateway] {
	return func(yield func(gateway.Gateway) bool) {
		m.shardsMu.Lock()
		defer m.shardsMu.Unlock()
		for _, shard := range m.shards {
			if !yield(shard) {
				return
			}
		}
	}
}
