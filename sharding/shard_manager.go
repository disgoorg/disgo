package sharding

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"slices"
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

	// ResumeShard resumes a specific shard with the given sessionID and sequence.
	ResumeShard(ctx context.Context, shardID int, state ShardState) error

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

func (m *shardManagerImpl) closeHandler(shard gateway.Gateway, err error, _ bool) {
	var closeError *websocket.CloseError
	if !m.config.AutoScaling || !errors.As(err, &closeError) || gateway.CloseEventCodeByCode(closeError.Code) != gateway.CloseEventCodeShardingRequired {
		return
	}
	m.config.Logger.Debug("shard requires re-sharding", slog.Int("shardID", shard.ShardID()))
	// make sure shard is closed
	shard.Close(context.TODO())

	m.shardsMu.Lock()
	delete(m.shards, shard.ShardID())
	defer m.shardsMu.Unlock()

	oldShardCount := m.config.ShardCount
	newShardCount := shard.ShardCount() * m.config.ShardSplitCount
	if newShardCount > m.config.ShardCount {
		m.config.ShardCount = newShardCount
	}

	newShardID := shard.ShardID()
	var newShardIDs []int
	for range m.config.ShardSplitCount {
		newShardIDs = append(newShardIDs, newShardID)
		newShardID += oldShardCount
	}

	var wg sync.WaitGroup
	for _, shardID := range newShardIDs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := m.openShard(context.Background(), shardID, newShardCount, ShardState{}); err != nil {
				m.config.Logger.Error("failed to re shard", slog.Any("err", err), slog.Int("shard_id", shardID))
			}

			m.config.Logger.Debug("re-sharded shard", slog.Int("shard_id", shardID), slog.Int("shard_count", newShardCount))
		}()
	}
	wg.Wait()
	m.config.Logger.Debug("re-sharded shard", slog.Int("shard_id", shard.ShardID()), slog.String("new_shard_ids", fmt.Sprint(newShardIDs)), slog.Int("new_shard_count", newShardCount))
}

func (m *shardManagerImpl) Open(ctx context.Context) {
	m.config.Logger.Debug("opening shards", slog.String("shard_ids", fmt.Sprint(slices.Collect(maps.Keys(m.config.ShardIDs)))), slog.Int("shard_count", m.config.ShardCount))

	var wg sync.WaitGroup
	for shardID, shardState := range m.config.ShardIDs {
		m.shardsMu.Lock()
		_, ok := m.shards[shardID]
		m.shardsMu.Unlock()
		if ok {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := m.openShard(ctx, shardID, m.config.ShardCount, shardState); err != nil {
				m.config.Logger.Error("failed to open shard", slog.Any("err", err), slog.Int("shard_id", shardID), slog.Int("shard_count", m.config.ShardCount))
			}

			m.config.Logger.Debug("opened shard", slog.Int("shard_id", shardID), slog.Int("shard_count", m.config.ShardCount))
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) Close(ctx context.Context) {
	m.config.Logger.Debug("closing shards", slog.String("shard_ids", fmt.Sprint(slices.Collect(maps.Keys(m.shards)))))
	var wg sync.WaitGroup

	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	for _, shard := range m.shards {
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard.Close(ctx)
		}()
	}
	wg.Wait()
	m.shards = map[int]gateway.Gateway{}
}

func (m *shardManagerImpl) OpenShard(ctx context.Context, shardID int) error {
	if err := m.openShard(ctx, shardID, m.config.ShardCount, ShardState{}); err != nil {
		m.config.Logger.Error("failed to open shard", slog.Any("err", err), slog.Int("shard_id", shardID))
		return err
	}

	m.config.Logger.Debug("opened shard", slog.Int("shard_id", shardID), slog.Int("shard_count", m.config.ShardCount))
	return nil
}

func (m *shardManagerImpl) ResumeShard(ctx context.Context, shardID int, state ShardState) error {
	if err := m.openShard(ctx, shardID, m.config.ShardCount, state); err != nil {
		m.config.Logger.Error("failed to resume shard",
			slog.Any("err", err),
			slog.Int("shard_id", shardID),
			slog.String("session_id", state.SessionID),
			slog.Int("sequence", state.Sequence),
			slog.String("resume_url", state.ResumeURL),
		)
		return err
	}

	m.config.Logger.Debug("resumed shard",
		slog.Int("shard_id", shardID),
		slog.Int("shard_count", m.config.ShardCount),
		slog.String("session_id", state.SessionID),
		slog.Int("sequence", state.Sequence),
		slog.String("resume_url", state.ResumeURL),
	)
	return nil
}

func (m *shardManagerImpl) openShard(ctx context.Context, shardID int, shardCount int, state ShardState) error {
	m.config.Logger.Debug("opening shard",
		slog.Int("shard_id", shardID),
		slog.Int("shard_count", shardCount),
		slog.String("session_id", state.SessionID),
		slog.Int("sequence", state.Sequence),
		slog.String("resume_url", state.ResumeURL),
	)

	opts := append(m.config.GatewayConfigOpts, gateway.WithShardID(shardID), gateway.WithShardCount(shardCount), gateway.WithIdentifyRateLimiter(m.config.IdentifyRateLimiter))
	if state.SessionID != "" {
		opts = append(opts, gateway.WithSessionID(state.SessionID))
	}
	if state.Sequence != 0 {
		opts = append(opts, gateway.WithSequence(state.Sequence))
	}
	if state.ResumeURL != "" {
		opts = append(opts, gateway.WithResumeURL(state.ResumeURL))
	}

	shard := m.config.GatewayCreateFunc(m.token, m.eventHandlerFunc, m.closeHandler, opts...)

	m.shardsMu.Lock()
	m.shards[shardID] = shard
	m.shardsMu.Unlock()

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
	for shard == nil && shardCount != 0 {
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
