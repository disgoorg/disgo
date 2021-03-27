package api

import "time"

// MemberCachePolicy can be used to define your own policy for caching members
type MemberCachePolicy func(*Member) bool

// Default member cache policies
var (
	MemberCachePolicyNone    MemberCachePolicy = func(_ *Member) bool { return false }
	MemberCachePolicyAll     MemberCachePolicy = func(_ *Member) bool { return true }
	MemberCachePolicyOwner   MemberCachePolicy = func(member *Member) bool { return member.IsOwner() }
	MemberCachePolicyOnline  MemberCachePolicy = func(_ *Member) bool { return false }
	MemberCachePolicyVoice   MemberCachePolicy = func(member *Member) bool { return false }
	MemberCachePolicyPending MemberCachePolicy = func(member *Member) bool { return member.Pending }
	MemberCachePolicyDefault                   = MemberCachePolicyOwner.Or(MemberCachePolicyVoice)
)

// Or allows you to combine that policy with another, meaning either needs to be true
func (p MemberCachePolicy) Or(policy MemberCachePolicy) MemberCachePolicy {
	return func(member *Member) bool {
		return p(member) || policy(member)
	}
}

// And allows you to require both policies to be true for the member to be cached
func (p MemberCachePolicy) And(policy MemberCachePolicy) MemberCachePolicy {
	return func(member *Member) bool {
		return p(member) && policy(member)
	}
}

// MemberCachePolicyAnyOf is a shorthand for MemberCachePolicy.Or(MemberCachePolicy).Or(MemberCachePolicy) etc.
func MemberCachePolicyAnyOf(policy MemberCachePolicy, policies ...MemberCachePolicy) MemberCachePolicy {
	for _, p := range policies {
		policy = policy.Or(p)
	}
	return policy
}

// MemberCachePolicyAllOf is a shorthand for MemberCachePolicy.And(MemberCachePolicy).And(MemberCachePolicy) etc.
func MemberCachePolicyAllOf(policy MemberCachePolicy, policies ...MemberCachePolicy) MemberCachePolicy {
	for _, p := range policies {
		policy = policy.And(p)
	}
	return policy
}

// MessageCachePolicy can be used to define your own policy for caching messages
type MessageCachePolicy func(*Message) bool

// Default member cache policies
var (
	MessageCachePolicyNone    MessageCachePolicy = func(_ *Message) bool { return false }
	MessageCachePolicyDefault                    = MessageCachePolicyNone
)

// Or allows you to combine that policy with another, meaning either needs to be true
func (p MessageCachePolicy) Or(policy MessageCachePolicy) MessageCachePolicy {
	return func(message *Message) bool {
		return p(message) || policy(message)
	}
}

// And allows you to require both policies to be true for the member to be cached
func (p MessageCachePolicy) And(policy MessageCachePolicy) MessageCachePolicy {
	return func(message *Message) bool {
		return p(message) && policy(message)
	}
}

// MessageCachePolicyDuration creates a new MessageCachePolicy which caches messages for the give duration
func MessageCachePolicyDuration(duration time.Duration) MessageCachePolicy {
	return func(message *Message) bool {
		return message.CreatedAt.Add(duration).After(time.Now())
	}
}

// MessageCachePolicyAny is a shorthand for MessageCachePolicy.Or(MessageCachePolicy).Or(MessageCachePolicy) etc.
func MessageCachePolicyAny(policy MessageCachePolicy, policies ...MessageCachePolicy) MessageCachePolicy {
	for _, p := range policies {
		policy = policy.Or(p)
	}
	return policy
}

// MessageCachePolicyAll is a shorthand for MessageCachePolicy.And(MessageCachePolicy).And(MessageCachePolicy) etc.
func MessageCachePolicyAll(policy MessageCachePolicy, policies ...MessageCachePolicy) MessageCachePolicy {
	for _, p := range policies {
		policy = policy.And(p)
	}
	return policy
}
