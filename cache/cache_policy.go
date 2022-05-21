package cache

// Policy can be used to define your own policy for when entities should be cached.
type Policy[T any] func(entity T) bool

func PolicyNone[T any](_ T) bool    { return false }
func PolicyAll[T any](_ T) bool     { return true }
func PolicyDefault[T any](t T) bool { return PolicyAll(t) }

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
