package core

import (
	"time"
)

// MessageCachePolicy can be used to define your own policy for caching messages
type MessageCachePolicy func(*Message) bool

// Default member caches policies
//goland:noinspection GoUnusedConst
var (
	MessageCachePolicyNone    MessageCachePolicy = func(_ *Message) bool { return false }
	MessageCachePolicyDefault                    = MessageCachePolicyNone
)

// Or allows you to combine that policy with another, meaning either needs to be true
//goland:noinspection GoUnusedExportedFunction
func (p MessageCachePolicy) Or(policy MessageCachePolicy) MessageCachePolicy {
	return func(message *Message) bool {
		return p(message) || policy(message)
	}
}

// And allows you to require both policies to be true for the member to be cached
//goland:noinspection GoUnusedExportedFunction
func (p MessageCachePolicy) And(policy MessageCachePolicy) MessageCachePolicy {
	return func(message *Message) bool {
		return p(message) && policy(message)
	}
}

// MessageCachePolicyDuration creates a new MessageCachePolicy which caches messages for the give duration
//goland:noinspection GoUnusedExportedFunction
func MessageCachePolicyDuration(duration time.Duration) MessageCachePolicy {
	return func(message *Message) bool {
		return message.CreatedAt.Add(duration).After(time.Now())
	}
}

// MessageCachePolicyAny is a shorthand for MessageCachePolicy.Or(MessageCachePolicy).Or(MessageCachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func MessageCachePolicyAny(policy MessageCachePolicy, policies ...MessageCachePolicy) MessageCachePolicy {
	for _, p := range policies {
		policy = policy.Or(p)
	}
	return policy
}

// MessageCachePolicyAll is a shorthand for MessageCachePolicy.And(MessageCachePolicy).And(MessageCachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func MessageCachePolicyAll(policy MessageCachePolicy, policies ...MessageCachePolicy) MessageCachePolicy {
	for _, p := range policies {
		policy = policy.And(p)
	}
	return policy
}
