package rate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// TODO: do we need some cleanup task?

//goland:noinspection GoUnusedExportedFunction
func NewRateLimiter(logger log.Logger, config *Config) RateLimiter {
	if logger == nil {
		logger = log.Default()
	}
	if config == nil {
		config = &DefaultConfig
	}
	return &RateLimiterImpl{
		logger:  logger,
		config:  *config,
		hashes:  map[*route.APIRoute]routeHash{},
		buckets: map[hashMajor]*bucket{},
	}
}

//goland:noinspection GoNameStartsWithPackageName
type (
	routeHash string
	hashMajor string

	RateLimiterImpl struct {
		sync.Mutex

		logger log.Logger
		config Config

		// global Rate Limit
		global int64

		// route.APIRoute -> Hash
		hashes map[*route.APIRoute]routeHash
		// Hash + Major Parameter -> bucket
		buckets map[hashMajor]*bucket
	}
)

func (r *RateLimiterImpl) Close(ctx context.Context) {
	// TODO: wait for all buckets to unlock
}

func (r *RateLimiterImpl) Config() Config {
	return r.config
}

func (r *RateLimiterImpl) getRouteHash(route *route.CompiledAPIRoute) hashMajor {
	hash, ok := r.hashes[route.APIRoute]
	if !ok {
		// generate routeHash
		hash = routeHash(route.Method().String() + "+" + route.Path())
		r.hashes[route.APIRoute] = hash
	}
	// return hashMajor
	return hashMajor(string(hash) + "+" + route.MajorParams())
}

func (r *RateLimiterImpl) getBucket(route *route.CompiledAPIRoute, create bool) *bucket {
	hash := r.getRouteHash(route)

	r.Lock()
	defer r.Unlock()
	b, ok := r.buckets[hash]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{
			Remaining: 1,
			Limit:     1,
		}
		r.buckets[hash] = b
	}
	return b
}

func (r *RateLimiterImpl) WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error {
	b := r.getBucket(route, true)
	b.Lock()

	var until time.Time
	now := time.Now()

	if b.Remaining == 0 && b.Reset.After(now) {
		until = b.Reset
	} else {
		until = time.Unix(0, r.global)
	}

	if until.After(now) {
		// TODO: do we want to return early when we know rate limit bigger than ctx deadline?
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return ErrCtxTimeout
		}

		select {
		case <-ctx.Done():
			b.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (r *RateLimiterImpl) UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error {
	b := r.getBucket(route, false)
	if b == nil {
		return nil
	}
	defer b.Unlock()

	bucketID := headers.Get("X-RateLimit-b")

	if bucketID != "" {
		b.ID = bucketID
	}

	global := headers.Get("X-RateLimit-global")

	remaining := headers.Get("X-RateLimit-Remaining")
	reset := headers.Get("X-RateLimit-Reset")
	retryAfter := headers.Get("Retry-After")

	switch {
	case retryAfter != "":
		i, err := strconv.Atoi(retryAfter)
		if err != nil {
			return fmt.Errorf("invalid retryAfter %s: %s", retryAfter, err)
		}

		at := time.Now().Add(time.Duration(i) * time.Second)

		if global != "" {
			r.global = at.UnixNano()
		} else {
			b.Reset = at
		}

	case reset != "":
		unix, err := strconv.ParseInt(reset, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid reset %s: %s", reset, err)
		}

		b.Reset = time.Unix(unix, 0)
	}

	if remaining != "" {
		u, err := strconv.Atoi(remaining)
		if err != nil {
			return fmt.Errorf("invalid remaining %s: %s", remaining, err)
		}

		b.Remaining = u
	} else {
		// Lower remaining one just to be safe
		if b.Remaining > 0 {
			b.Remaining--
		}
	}

	return nil
}

type bucket struct {
	sync.Mutex
	ID        string
	Reset     time.Time
	Remaining int
	Limit     int
}
