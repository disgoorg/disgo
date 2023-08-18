package sharding

import (
	"context"
	"errors"
	"sync"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/gateway"
)

var _ ShardManager = (*shardManagerImpl)(nil)

// New creates a new default ShardManager with the given token, eventHandlerFunc and ConfigOpt(s).
func New(token string, eventHandlerFunc gateway.EventHandlerFunc, opts ...ConfigOpt) ShardManager {
	config := DefaultConfig()
	config.Apply(opts)

	return &shardManagerImpl{
		shards:           map[int]gateway.Gateway{},
		token:            token,
		eventHandlerFunc: eventHandlerFunc,
		config:           *config,
	}
}

type shardManagerImpl struct {
	shards   map[int]gateway.Gateway
	shardsMu sync.Mutex

	token            string
	eventHandlerFunc gateway.EventHandlerFunc
	config           Config
}

func (m *shardManagerImpl) closeHandler(shard gateway.Gateway, err error) {
	var closeError *websocket.CloseError
	if !m.config.AutoScaling || !errors.As(err, &closeError) || gateway.CloseEventCodeByCode(closeError.Code) != gateway.CloseEventCodeShardingRequired {
		return
	}
	m.config.Logger.Debugf("shard %d requires re-sharding", shard.ShardID())
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
				m.config.Logger.Errorf("failed to wait shard bucket %d: %s", shardID, err)
				return
			}
			defer m.config.RateLimiter.UnlockBucket(shardID)

			newShard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.closeHandler, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(newShardCount))...)
			m.shards[shardID] = newShard
			if err := newShard.Open(context.TODO()); err != nil {
				m.config.Logger.Errorf("failed to re shard %d, error: %s", shardID, err)
			}
		}()
	}
	wg.Wait()
	m.config.Logger.Debugf("re-sharded shard %d into newShards: %d, newShardCount: %d", shard.ShardID(), newShardIDs, newShardCount)
}

func (m *shardManagerImpl) Open(ctx context.Context) {
	m.config.Logger.Debugf("opening %+v shards...", m.config.ShardIDs)
	var wg sync.WaitGroup

	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	for shardInt := range m.config.ShardIDs {
		shardID := shardInt
		if _, ok := m.shards[shardID]; ok {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.config.RateLimiter.WaitBucket(ctx, shardID); err != nil {
				m.config.Logger.Errorf("failed to wait shard bucket %d: %s", shardID, err)
				return
			}
			defer m.config.RateLimiter.UnlockBucket(shardID)

			shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.closeHandler, append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(m.config.ShardCount))...)
			m.shards[shardID] = shard
			if err := shard.Open(ctx); err != nil {
				m.config.Logger.Errorf("failed to open shard %d: %s", shardID, err)
			}
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) Close(ctx context.Context) {
	m.config.Logger.Debugf("closing %v shards...", m.config.ShardIDs)
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
	m.config.Logger.Debugf("opening shard %d...", shardID)

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
	m.config.Logger.Debugf("closing shard %d...", shardID)
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

func (m *shardManagerImpl) Shards() map[int]gateway.Gateway {
	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	shards := make(map[int]gateway.Gateway, len(m.shards))
	for shardID, shard := range m.shards {
		shards[shardID] = shard
	}
	return m.shards
}
