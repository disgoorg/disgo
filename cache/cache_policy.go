package cache

import (
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// PolicyNone returns a policy that will never cache anything.
func PolicyNone[T any](_ T) bool { return false }

// PolicyAll returns a policy that will cache all entities.
func PolicyAll[T any](_ T) bool { return true }

// PolicyMembersInclude returns a policy that will only cache members of the given guilds.
func PolicyMembersInclude(guildIDs ...snowflake.ID) Policy[discord.Member] {
	return func(member discord.Member) bool {
		return slices.Contains(guildIDs, member.GuildID)
	}
}

// PolicyMembersPending is a policy that will only cache members that are pending.
func PolicyMembersPending(member discord.Member) bool {
	return member.Pending
}

// PolicyMembersInVoice returns a policy that will only cache members that are connected to an audio channel.
func PolicyMembersInVoice(caches Caches) Policy[discord.Member] {
	return func(member discord.Member) bool {
		_, ok := caches.VoiceState(member.GuildID, member.User.ID)
		return ok
	}
}

// PolicyChannelInclude returns a policy that will only cache channels of the given types.
func PolicyChannelInclude(channelTypes ...discord.ChannelType) Policy[discord.Channel] {
	return func(channel discord.Channel) bool {
		return slices.Contains(channelTypes, channel.Type())
	}
}

// PolicyChannelExclude returns a policy that will not cache channels of the given types.
func PolicyChannelExclude(channelTypes ...discord.ChannelType) Policy[discord.Channel] {
	return func(channel discord.Channel) bool {
		return !slices.Contains(channelTypes, channel.Type())
	}
}

// Policy can be used to define your own policy for when entities should be cached.
type Policy[T any] func(entity T) bool

// Or allows you to combine the CachePolicy with another, meaning either of them needs to be true
func (p Policy[T]) Or(policy Policy[T]) Policy[T] {
	return func(entity T) bool {
		return p(entity) || policy(entity)
	}
}

// And allows you to require both CachePolicy(s) to be true for the entity to be cached
func (p Policy[T]) And(policy Policy[T]) Policy[T] {
	return func(entity T) bool {
		return p(entity) && policy(entity)
	}
}

// AnyPolicy is a shorthand for CachePolicy.Or(CachePolicy).Or(CachePolicy) etc.
func AnyPolicy[T any](policies ...Policy[T]) Policy[T] {
	var policy Policy[T]
	for _, p := range policies {
		if policy == nil {
			policy = p
			continue
		}
		policy = policy.Or(p)
	}
	return policy
}

// AllPolicies is a shorthand for CachePolicy.And(CachePolicy).And(CachePolicy) etc.
func AllPolicies[T any](policies ...Policy[T]) Policy[T] {
	var policy Policy[T]
	for _, p := range policies {
		if policy == nil {
			policy = p
			continue
		}
		policy = policy.And(p)
	}
	return policy
}
