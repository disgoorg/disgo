package core

import "github.com/DisgoOrg/snowflake"

type (
	VoiceStateFindFunc func(voiceState *VoiceState) bool

	VoiceStateCache interface {
		Get(guildID snowflake.Snowflake, userID snowflake.Snowflake) *VoiceState
		GetCopy(guildID snowflake.Snowflake, userID snowflake.Snowflake) *VoiceState
		Set(voiceState *VoiceState) *VoiceState
		Remove(guildID snowflake.Snowflake, userID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*VoiceState
		All() map[snowflake.Snowflake][]*VoiceState

		GuildCache(guildID snowflake.Snowflake) map[snowflake.Snowflake]*VoiceState
		GuildAll(guildID snowflake.Snowflake) []*VoiceState

		FindFirst(voiceStateFindFunc VoiceStateFindFunc) *VoiceState
		FindAll(voiceStateFindFunc VoiceStateFindFunc) []*VoiceState
	}

	voiceStateCacheImpl struct {
		cacheFlags  CacheFlags
		voiceStates map[snowflake.Snowflake]map[snowflake.Snowflake]*VoiceState
	}
)

func NewVoiceStateCache(cacheFlags CacheFlags) VoiceStateCache {
	return &voiceStateCacheImpl{
		cacheFlags:  cacheFlags,
		voiceStates: map[snowflake.Snowflake]map[snowflake.Snowflake]*VoiceState{},
	}
}

func (c *voiceStateCacheImpl) Get(guildID snowflake.Snowflake, userID snowflake.Snowflake) *VoiceState {
	if _, ok := c.voiceStates[guildID]; !ok {
		return nil
	}
	return c.voiceStates[guildID][userID]
}

func (c *voiceStateCacheImpl) GetCopy(guildID snowflake.Snowflake, userID snowflake.Snowflake) *VoiceState {
	if voiceState := c.Get(guildID, userID); voiceState != nil {
		vs := *voiceState
		return &vs
	}
	return nil
}

func (c *voiceStateCacheImpl) Set(voiceState *VoiceState) *VoiceState {
	if c.cacheFlags.Missing(CacheFlagVoiceStates) {
		return voiceState
	}
	if _, ok := c.voiceStates[voiceState.GuildID]; !ok {
		c.voiceStates[voiceState.GuildID] = map[snowflake.Snowflake]*VoiceState{}
	}
	rol, ok := c.voiceStates[voiceState.GuildID][voiceState.UserID]
	if ok {
		*rol = *voiceState
		return rol
	}
	c.voiceStates[voiceState.GuildID][voiceState.UserID] = voiceState

	return voiceState
}

func (c *voiceStateCacheImpl) Remove(guildID snowflake.Snowflake, userID snowflake.Snowflake) {
	if _, ok := c.voiceStates[guildID]; !ok {
		return
	}
	delete(c.voiceStates[guildID], userID)
}

func (c *voiceStateCacheImpl) Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*VoiceState {
	return c.voiceStates
}

func (c *voiceStateCacheImpl) All() map[snowflake.Snowflake][]*VoiceState {
	voiceStates := make(map[snowflake.Snowflake][]*VoiceState, len(c.voiceStates))
	for guildID, guildVoiceStates := range c.voiceStates {
		voiceStates[guildID] = make([]*VoiceState, len(guildVoiceStates))
		i := 0
		for _, voiceStateVoiceState := range guildVoiceStates {
			voiceStates[guildID] = append(voiceStates[guildID], voiceStateVoiceState)
		}
		i++
	}
	return voiceStates
}

func (c *voiceStateCacheImpl) GuildCache(guildID snowflake.Snowflake) map[snowflake.Snowflake]*VoiceState {
	if _, ok := c.voiceStates[guildID]; !ok {
		return nil
	}
	return c.voiceStates[guildID]
}

func (c *voiceStateCacheImpl) GuildAll(guildID snowflake.Snowflake) []*VoiceState {
	if _, ok := c.voiceStates[guildID]; !ok {
		return nil
	}
	voiceStates := make([]*VoiceState, len(c.voiceStates[guildID]))
	i := 0
	for _, voiceState := range c.voiceStates[guildID] {
		voiceStates = append(voiceStates, voiceState)
		i++
	}
	return voiceStates
}

func (c *voiceStateCacheImpl) FindFirst(voiceStateFindFunc VoiceStateFindFunc) *VoiceState {
	for _, guildVoiceStates := range c.voiceStates {
		for _, voiceState := range guildVoiceStates {
			if voiceStateFindFunc(voiceState) {
				return voiceState
			}
		}
	}
	return nil
}

func (c *voiceStateCacheImpl) FindAll(voiceStateFindFunc VoiceStateFindFunc) []*VoiceState {
	var voiceStates []*VoiceState
	for _, guildVoiceStates := range c.voiceStates {
		for _, voiceState := range guildVoiceStates {
			if voiceStateFindFunc(voiceState) {
				voiceStates = append(voiceStates, voiceState)
			}
		}
	}
	return voiceStates
}
