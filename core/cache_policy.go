package core

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

// CachePolicy can be used to define your own policy for caching cache
type CachePolicy[T any] func(entity T) bool

// Default discord.Message CachePolicy(s)
//goland:noinspection GoUnusedConst
var (
	MessageCachePolicyNone CachePolicy[discord.Message] = func(_ discord.Message) bool { return false }

	// MessageCachePolicyDuration creates a new CachePolicy which BotCaches discord.Message(s) for the give time.Duration
	MessageCachePolicyDuration = func(duration time.Duration) CachePolicy[discord.Message] {
		return func(message discord.Message) bool {
			return message.CreatedAt.Add(duration).After(time.Now())
		}
	}
	MessageCachePolicyDefault = MessageCachePolicyNone
)

// Default discord.Member CachePolicy(s)
//goland:noinspection GoUnusedGlobalVariable
var (
	MemberCachePolicyNone    CachePolicy[discord.Member] = func(_ discord.Member) bool { return false }
	MemberCachePolicyAll     CachePolicy[discord.Member] = func(_ discord.Member) bool { return true }
	MemberCachePolicyOwner   CachePolicy[discord.Member] = func(member discord.Member) bool { return false /*TODO*/ }
	MemberCachePolicyOnline  CachePolicy[discord.Member] = func(_ discord.Member) bool { return false }
	MemberCachePolicyVoice   CachePolicy[discord.Member] = func(member discord.Member) bool { return false }
	MemberCachePolicyPending CachePolicy[discord.Member] = func(member discord.Member) bool { return member.Pending }
	MemberCachePolicyDefault                             = MemberCachePolicyOwner.Or(MemberCachePolicyVoice)
)

// Or allows you to combine the CachePolicy with another, meaning either of them needs to be true
//goland:noinspection GoUnusedExportedFunction
func (p CachePolicy[T]) Or(policy CachePolicy[T]) CachePolicy[T] {
	return func(entity T) bool {
		return p(entity) || policy(entity)
	}
}

// And allows you to require both CachePolicy(s) to be true for the entity to be cached
//goland:noinspection GoUnusedExportedFunction
func (p CachePolicy[T]) And(policy CachePolicy[T]) CachePolicy[T] {
	return func(entity T) bool {
		return p(entity) && policy(entity)
	}
}

// CachePolicyAny is a shorthand for CachePolicy.Or(CachePolicy).Or(CachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func CachePolicyAny[T any](policies ...CachePolicy[T]) CachePolicy[T] {
	var policy CachePolicy[T]
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
func CachePolicyAll[T any](policies ...CachePolicy[T]) CachePolicy[T] {
	var policy CachePolicy[T]
	for _, p := range policies {
		if policy == nil {
			policy = p
			continue
		}
		policy = policy.And(p)
	}
	return policy
}
