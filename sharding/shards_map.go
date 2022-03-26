package sharding

import (
	"sync"

	"github.com/disgoorg/disgo/gateway"
)

func NewShardsMap() *ShardsMap {
	return &ShardsMap{shards: map[int]gateway.Gateway{}}
}

type ShardsMap struct {
	mu     sync.RWMutex
	shards map[int]gateway.Gateway
}

func (m *ShardsMap) AllIDs() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]int, len(m.shards))
	i := 0
	for id := range m.shards {
		ids[i] = id
		i++
	}
	return ids
}

func (m *ShardsMap) Get(shardId int) gateway.Gateway {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.shards[shardId]
}

func (m *ShardsMap) Set(shardId int, shard gateway.Gateway) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shards[shardId] = shard
}

func (m *ShardsMap) Has(shardId int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.shards[shardId]
	return ok
}

func (m *ShardsMap) Delete(shardId int) gateway.Gateway {
	m.mu.RLock()
	shard, ok := m.shards[shardId]
	m.mu.RUnlock()
	if ok {
		m.mu.Lock()
		delete(m.shards, shardId)
		m.mu.Unlock()
	}
	return shard
}

func (m *ShardsMap) For(forFunc func(shardID int, shard gateway.Gateway)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for shardID, shard := range m.shards {
		forFunc(shardID, shard)
	}
}
