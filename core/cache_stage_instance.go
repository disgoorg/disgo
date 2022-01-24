package core

import "github.com/DisgoOrg/snowflake"

type (
	StageInstanceFindFunc func(stageInstance *StageInstance) bool

	StageInstanceCache interface {
		Get(stageInstanceID snowflake.Snowflake) *StageInstance
		GetCopy(stageInstanceID snowflake.Snowflake) *StageInstance
		Set(stageInstance *StageInstance) *StageInstance
		Remove(stageInstanceID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]*StageInstance
		All() []*StageInstance

		FindFirst(stageInstanceFindFunc StageInstanceFindFunc) *StageInstance
		FindAll(stageInstanceFindFunc StageInstanceFindFunc) []*StageInstance
	}

	stageInstanceCacheImpl struct {
		cacheFlags     CacheFlags
		stageInstances map[snowflake.Snowflake]*StageInstance
	}
)

func NewStageInstanceCache(cacheFlags CacheFlags) StageInstanceCache {
	return &stageInstanceCacheImpl{
		cacheFlags:     cacheFlags,
		stageInstances: map[snowflake.Snowflake]*StageInstance{},
	}
}

func (c *stageInstanceCacheImpl) Get(stageInstanceID snowflake.Snowflake) *StageInstance {
	return c.stageInstances[stageInstanceID]
}

func (c *stageInstanceCacheImpl) GetCopy(stageInstanceID snowflake.Snowflake) *StageInstance {
	if stageInstance := c.Get(stageInstanceID); stageInstance != nil {
		st := *stageInstance
		return &st
	}
	return nil
}

func (c *stageInstanceCacheImpl) Set(stageInstance *StageInstance) *StageInstance {
	if c.cacheFlags.Missing(CacheFlagStageInstances) {
		return stageInstance
	}
	stI, ok := c.stageInstances[stageInstance.ID]
	if ok {
		*stI = *stageInstance
		return stI
	}
	c.stageInstances[stageInstance.ID] = stageInstance
	return stageInstance
}

func (c *stageInstanceCacheImpl) Remove(id snowflake.Snowflake) {
	delete(c.stageInstances, id)
}

func (c *stageInstanceCacheImpl) Cache() map[snowflake.Snowflake]*StageInstance {
	return c.stageInstances
}

func (c *stageInstanceCacheImpl) All() []*StageInstance {
	stageInstances := make([]*StageInstance, len(c.stageInstances))
	i := 0
	for _, stageInstance := range c.stageInstances {
		stageInstances[i] = stageInstance
		i++
	}
	return stageInstances
}

func (c *stageInstanceCacheImpl) FindFirst(stageInstanceFindFunc StageInstanceFindFunc) *StageInstance {
	for _, stI := range c.stageInstances {
		if stageInstanceFindFunc(stI) {
			return stI
		}
	}
	return nil
}

func (c *stageInstanceCacheImpl) FindAll(stageInstanceFindFunc StageInstanceFindFunc) []*StageInstance {
	var stageInstances []*StageInstance
	for _, stI := range c.stageInstances {
		if stageInstanceFindFunc(stI) {
			stageInstances = append(stageInstances, stI)
		}
	}
	return stageInstances
}
