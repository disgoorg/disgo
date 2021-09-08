package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	VoiceStateFindFunc func(voiceState *VoiceState) bool

	VoiceStateCache interface {
		Get(guildID discord.Snowflake, userID discord.Snowflake) *VoiceState
		GetCopy(guildID discord.Snowflake, userID discord.Snowflake) *VoiceState
		Set(voiceState *VoiceState) *VoiceState
		Remove(guildID discord.Snowflake, userID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*VoiceState
		All() map[discord.Snowflake][]*VoiceState

		GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*VoiceState
		GuildAll(guildID discord.Snowflake) []*VoiceState

		FindFirst(voiceStateFindFunc VoiceStateFindFunc) *VoiceState
		FindAll(voiceStateFindFunc VoiceStateFindFunc) []*VoiceState
	}

	voiceStateCacheImpl struct {
		cacheFlags  CacheFlags
		voiceStates map[discord.Snowflake]map[discord.Snowflake]*VoiceState
	}
)

func NewVoiceStateCache(cacheFlags CacheFlags) VoiceStateCache {
	return &voiceStateCacheImpl{
		cacheFlags:  cacheFlags,
		voiceStates: map[discord.Snowflake]map[discord.Snowflake]*VoiceState{},
	}
}

func (c *voiceStateCacheImpl) Get(guildID discord.Snowflake, userID discord.Snowflake) *VoiceState {
	if _, ok := c.voiceStates[guildID]; !ok {
		return nil
	}
	return c.voiceStates[guildID][userID]
}

func (c *voiceStateCacheImpl) GetCopy(guildID discord.Snowflake, userID discord.Snowflake) *VoiceState {
	return &*c.Get(guildID, userID)
}

func (c *voiceStateCacheImpl) Set(voiceState *VoiceState) *VoiceState {
	if !c.cacheFlags.Missing(CacheFlagVoiceStates) {
		return voiceState
	}
	if _, ok := c.voiceStates[voiceState.GuildID]; !ok {
		c.voiceStates[voiceState.GuildID] = map[discord.Snowflake]*VoiceState{}
	}
	rol, ok := c.voiceStates[voiceState.GuildID][voiceState.UserID]
	if ok {
		*rol = *voiceState
		return rol
	}
	c.voiceStates[voiceState.GuildID][voiceState.UserID] = voiceState

	return voiceState
}

func (c *voiceStateCacheImpl) Remove(guildID discord.Snowflake, userID discord.Snowflake) {
	if _, ok := c.voiceStates[guildID]; !ok {
		return
	}
	delete(c.voiceStates[guildID], userID)
}

func (c *voiceStateCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*VoiceState {
	return c.voiceStates
}

func (c *voiceStateCacheImpl) All() map[discord.Snowflake][]*VoiceState {
	voiceStates := make(map[discord.Snowflake][]*VoiceState, len(c.voiceStates))
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

func (c *voiceStateCacheImpl) GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*VoiceState {
	if _, ok := c.voiceStates[guildID]; !ok {
		return nil
	}
	return c.voiceStates[guildID]
}

func (c *voiceStateCacheImpl) GuildAll(guildID discord.Snowflake) []*VoiceState {
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
