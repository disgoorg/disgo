package core

//goland:noinspection GoUnusedGlobalVariable
var DefaultCacheConfig = CacheConfig{
	CacheFlags:         CacheFlagsDefault,
	MemberCachePolicy:  MemberCachePolicyDefault,
	MessageCachePolicy: MessageCachePolicyDefault,
}

type CacheConfig struct {
	CacheFlags         CacheFlags
	MemberCachePolicy  CachePolicy[Member]
	MessageCachePolicy CachePolicy[Message]
}

type CacheConfigOpt func(config *CacheConfig)

func (c *CacheConfig) Apply(opts []CacheConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCacheFlags(cacheFlags ...CacheFlags) CacheConfigOpt {
	return func(config *CacheConfig) {
		var flags CacheFlags
		for _, flag := range cacheFlags {
			flags = flags.Add(flag)
		}
		config.CacheFlags = flags
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberCachePolicy(memberCachePolicy CachePolicy[Member]) CacheConfigOpt {
	return func(config *CacheConfig) {
		config.MemberCachePolicy = memberCachePolicy
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMessageCachePolicy(messageCachePolicy CachePolicy[Message]) CacheConfigOpt {
	return func(config *CacheConfig) {
		config.MessageCachePolicy = messageCachePolicy
	}
}
