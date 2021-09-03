package core

import "github.com/DisgoOrg/disgo/discord"

type ISnowflake interface {
	ID() discord.Snowflake
	Copy() ISnowflake
	Update(entity ISnowflake)
}

type (
	FindFunc    func(entity ISnowflake) bool
	CachePolicy func(entity ISnowflake) bool
)

type GenericCache interface {
	DoCleanUp()

	Get(id discord.Snowflake) ISnowflake
	Set(entity ISnowflake) ISnowflake
	Remove(id discord.Snowflake)

	Cache() map[discord.Snowflake]ISnowflake
	All() []ISnowflake

	FindFirst(findFunc FindFunc) ISnowflake
	FindAll(findFunc FindFunc) []ISnowflake
}

func newGenericCache(cachePolicy CachePolicy) GenericCache {
	return &genericCacheImpl{entities: map[discord.Snowflake]ISnowflake{}, cachePolicy: cachePolicy}
}

type genericCacheImpl struct {
	entities    map[discord.Snowflake]ISnowflake
	cachePolicy func(entity ISnowflake) bool
}

func (c *genericCacheImpl) DoCleanUp() {
	for id, entity := range c.entities {
		if !c.cachePolicy(entity) {
			delete(c.entities, id)
		}
	}
}

func (c *genericCacheImpl) Get(id discord.Snowflake) ISnowflake {
	return c.entities[id]
}

func (c *genericCacheImpl) Set(entity ISnowflake) ISnowflake {
	if !c.cachePolicy(entity) {
		return entity
	}
	e, ok := c.entities[entity.ID()]
	if ok {
		e.Update(entity)
		return e
	}
	c.entities[entity.ID()] = entity
	return entity
}

func (c *genericCacheImpl) Remove(id discord.Snowflake) {
	delete(c.entities, id)
}

func (c *genericCacheImpl) Cache() map[discord.Snowflake]ISnowflake {
	return c.entities
}

func (c *genericCacheImpl) All() []ISnowflake {
	entities := make([]ISnowflake, len(c.entities))
	i := 0
	for _, entity := range c.entities {
		entities[i] = entity
		i++
	}
	return entities
}

func (c *genericCacheImpl) FindFirst(findFunc FindFunc) ISnowflake {
	for _, entity := range c.entities {
		if findFunc(entity) == true {
			return entity
		}
	}
	return nil
}

func (c *genericCacheImpl) FindAll(findFunc FindFunc) []ISnowflake {
	var entities []ISnowflake
	for _, entity := range c.entities {
		if findFunc(entity) == true {
			entities = append(entities, entity)
		}
	}
	return entities
}
