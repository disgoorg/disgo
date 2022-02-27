package sharding

import (
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

var _ ShardManager = (*shardManagerImpl)(nil)

func newGateway(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
	newConfig := *config
	newConfig.ShardID = shardID
	newConfig.ShardCount = shardCount
	return gateway.New(token, url, eventHandlerFunc, &newConfig)
}

func New(token string, gatewayURL string, eventHandlerFunc gateway.EventHandlerFunc, config *Config) ShardManager {
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.ShardIDs == nil || len(config.ShardIDs) == 0 {
		config.ShardIDs = map[int]struct{}{0: {}}
	}
	if config.ShardCount == 0 {
		config.ShardCount = len(config.ShardIDs)
	}
	if config.GatewayConfig == nil {
		config.GatewayConfig = &gateway.DefaultConfig
	}
	if config.GatewayCreateFunc == nil {
		config.GatewayCreateFunc = newGateway
	}
	if config.RateLimiter == nil {
		config.RateLimiter = srate.NewLimiter(&srate.DefaultConfig)
	}
	return &shardManagerImpl{
		shards:           NewShardsMap(),
		token:            token,
		gatewayURL:       gatewayURL,
		eventHandlerFunc: eventHandlerFunc,
		config:           *config,
	}
}

type shardManagerImpl struct {
	shards         *ShardsMap
	guildsPerShard map[int]int

	token            string
	gatewayURL       string
	eventHandlerFunc gateway.EventHandlerFunc
	config           Config
}

func (m *shardManagerImpl) Logger() log.Logger {
	return m.config.Logger
}

func (m *shardManagerImpl) Config() Config {
	return m.config
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
	m.Logger().Infof("opening %s shards...", m.config.ShardIDs)
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

			shard := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, m.config.ShardCount, m.onEvent, m.config.GatewayConfig)
			m.shards.Set(shardID, shard)
			if err := shard.Open(ctx); err != nil {
				m.Logger().Errorf("failed to open shard %d: %s", shardID, err)
			}
		}()
	}
	wg.Wait()
}

func (m *shardManagerImpl) ReOpen(ctx context.Context) {
	m.Logger().Infof("reopening %s shards...", m.config.ShardIDs)
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
	m.Logger().Infof("closing %v shards...", m.config.ShardIDs)
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
	m.Logger().Infof("opening shard %d...", shardID)
	shard := m.config.GatewayCreateFunc(m.token, m.gatewayURL, shardID, shardCount, m.eventHandlerFunc, m.config.GatewayConfig)
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
	shard := m.shards.Delete(shardID)
	if shard != nil {
		shard.Close(ctx)
	}
}

func (m *shardManagerImpl) ShardByGuildID(guildId snowflake.Snowflake) gateway.Gateway {
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
