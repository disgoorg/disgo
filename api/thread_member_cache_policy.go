package api

// ThreadMemberCachePolicy can be used to define your own policy for caching thread members
type ThreadMemberCachePolicy func(*ThreadMember) bool

// Default member cache policies
var (
	ThreadMemberCachePolicyNone    ThreadMemberCachePolicy = func(_ *ThreadMember) bool { return false }
	ThreadMemberCachePolicyAll     ThreadMemberCachePolicy = func(_ *ThreadMember) bool { return true }
	ThreadMemberCachePolicyActive  ThreadMemberCachePolicy = func(member *ThreadMember) bool { return !member.Thread().ThreadMetadata().Locked }
	ThreadMemberCachePolicyDefault                         = ThreadMemberCachePolicyActive
)

// Or allows you to combine that policy with another, meaning either needs to be true
func (p ThreadMemberCachePolicy) Or(policy ThreadMemberCachePolicy) ThreadMemberCachePolicy {
	return func(member *ThreadMember) bool {
		return p(member) || policy(member)
	}
}

// And allows you to require both policies to be true for the member to be cached
func (p ThreadMemberCachePolicy) And(policy ThreadMemberCachePolicy) ThreadMemberCachePolicy {
	return func(member *ThreadMember) bool {
		return p(member) && policy(member)
	}
}

// ThreadMemberCachePolicyAnyOf is a shorthand for ThreadMemberCachePolicy.Or(ThreadMemberCachePolicy).Or(ThreadMemberCachePolicy) etc.
func ThreadMemberCachePolicyAnyOf(policy ThreadMemberCachePolicy, policies ...ThreadMemberCachePolicy) ThreadMemberCachePolicy {
	for _, p := range policies {
		policy = policy.Or(p)
	}
	return policy
}

// ThreadMemberCachePolicyAllOf is a shorthand for ThreadMemberCachePolicy.And(ThreadMemberCachePolicy).And(ThreadMemberCachePolicy) etc.
func ThreadMemberCachePolicyAllOf(policy ThreadMemberCachePolicy, policies ...ThreadMemberCachePolicy) ThreadMemberCachePolicy {
	for _, p := range policies {
		policy = policy.And(p)
	}
	return policy
}
