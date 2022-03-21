package cache

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

// Policy can be used to define your own policy for caching cache
type Policy[T any] func(entity T) bool

// Default discord.Message CachePolicy(s)
//goland:noinspection GoUnusedConst
var (
	MessageCachePolicyNone Policy[discord.Message] = func(_ discord.Message) bool { return false }

	// MessageCachePolicyDuration creates a new CachePolicy which BotCaches discord.Message(s) for the give time.Duration
	MessageCachePolicyDuration = func(duration time.Duration) Policy[discord.Message] {
		return func(message discord.Message) bool {
			return message.CreatedAt.Add(duration).After(time.Now())
		}
	}
	MessageCachePolicyDefault = MessageCachePolicyNone
)

// Default discord.Member CachePolicy(s)
//goland:noinspection GoUnusedGlobalVariable
var (
	MemberCachePolicyNone    Policy[discord.Member] = func(_ discord.Member) bool { return false }
	MemberCachePolicyAll     Policy[discord.Member] = func(_ discord.Member) bool { return true }
	MemberCachePolicyOwner   Policy[discord.Member] = func(member discord.Member) bool { return false /*TODO*/ }
	MemberCachePolicyOnline  Policy[discord.Member] = func(_ discord.Member) bool { return false }
	MemberCachePolicyVoice   Policy[discord.Member] = func(member discord.Member) bool { return false }
	MemberCachePolicyPending Policy[discord.Member] = func(member discord.Member) bool { return member.Pending }
	MemberCachePolicyDefault                        = MemberCachePolicyOwner.Or(MemberCachePolicyVoice)
)

// Or allows you to combine the CachePolicy with another, meaning either of them needs to be true
//goland:noinspection GoUnusedExportedFunction
func (p Policy[T]) Or(policy Policy[T]) Policy[T] {
	return func(entity T) bool {
		return p(entity) || policy(entity)
	}
}

// And allows you to require both CachePolicy(s) to be true for the entity to be cached
//goland:noinspection GoUnusedExportedFunction
func (p Policy[T]) And(policy Policy[T]) Policy[T] {
	return func(entity T) bool {
		return p(entity) && policy(entity)
	}
}

// CachePolicyAny is a shorthand for CachePolicy.Or(CachePolicy).Or(CachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func CachePolicyAny[T any](policies ...Policy[T]) Policy[T] {
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

// CachePolicyAll is a shorthand for CachePolicy.And(CachePolicy).And(CachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func CachePolicyAll[T any](policies ...Policy[T]) Policy[T] {
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
