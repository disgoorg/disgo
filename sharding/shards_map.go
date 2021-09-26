package sharding

import (
	"sync"

	"github.com/DisgoOrg/disgo/gateway"
)

func NewShardsMap() *ShardsMap {
	return &ShardsMap{Shards: map[int]gateway.Gateway{}}
}

type ShardsMap struct {
	sync.RWMutex
	Shards map[int]gateway.Gateway
}

func (m *ShardsMap) Get(shardId int) gateway.Gateway {
	m.RLock()
	defer m.RUnlock()
	return m.Shards[shardId]
}

func (m *ShardsMap) Set(shardId int, shard gateway.Gateway) {
	m.Lock()
	defer m.Unlock()
	m.Shards[shardId] = shard
}

func (m *ShardsMap) Has(shardId int) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.Shards[shardId]
	return ok
}

func (m *ShardsMap) Delete(shardId int) gateway.Gateway {
	m.RLock()
	shard, ok := m.Shards[shardId]
	m.RUnlock()
	if ok {
		m.Lock()
		delete(m.Shards, shardId)
		m.Unlock()
	}
	return shard
}
