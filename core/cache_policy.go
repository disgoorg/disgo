package core

var (
	MemberCachePolicyNone    CachePolicy[Member] = func(_ Member) bool { return false }
	MemberCachePolicyAll     CachePolicy[Member] = func(_ Member) bool { return true }
	MemberCachePolicyOwner   CachePolicy[Member] = func(member Member) bool { return member.IsOwner() }
	MemberCachePolicyOnline  CachePolicy[Member] = func(_ Member) bool { return false }
	MemberCachePolicyVoice   CachePolicy[Member] = func(member Member) bool { return false }
	MemberCachePolicyPending CachePolicy[Member] = func(member Member) bool { return member.Pending }
	MemberCachePolicyDefault                     = MemberCachePolicyOwner.Or(MemberCachePolicyVoice)
)

var (
	MessageCachePolicyNone    CachePolicy[Message] = func(_ Message) bool { return false }
	MessageCachePolicyDefault                      = MessageCachePolicyNone
)

type CachePolicy[T any] func(T) bool

// Or allows you to combine that policy with another, meaning either needs to be true
func (p CachePolicy[T]) Or(policy CachePolicy[T]) CachePolicy[T] {
	return func(t T) bool {
		return p(t) || policy(t)
	}
}

// And allows you to require both policies to be true for the member to be cached
func (p CachePolicy[T]) And(policy CachePolicy[T]) CachePolicy[T] {
	return func(t T) bool {
		return p(t) && policy(t)
	}
}

// CachePolicyAnyOf is a shorthand for CachePolicy.Or(CachePolicy).Or(CachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func CachePolicyAnyOf[T any](policy CachePolicy[T], policies ...CachePolicy[T]) CachePolicy[T] {
	for i := range policies {
		policy = policy.Or(policies[i])
	}
	return policy
}

// CachePolicyAllOf is a shorthand for CachePolicy.And(CachePolicy).And(CachePolicy) etc.
//goland:noinspection GoUnusedExportedFunction
func CachePolicyAllOf[T any](policy CachePolicy[T], policies ...CachePolicy[T]) CachePolicy[T] {
	for i := range policies {
		policy = policy.And(policies[i])
	}
	return policy
}
