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
func NewRateLimiter(logger log.Logger) RateLimiter {
	if logger == nil {
		logger = log.Default()
	}
	return &RateLimiterImpl{
		logger:  logger,
		Hashes:  map[*route.APIRoute]routeHash{},
		Buckets: map[hashMajor]*bucket{},
	}
}

//goland:noinspection GoNameStartsWithPackageName
type (
	routeHash string
	hashMajor string

	RateLimiterImpl struct {
		sync.Mutex

		logger log.Logger

		// Global Rate Limit
		Global int64

		// route.APIRoute -> Hash
		Hashes map[*route.APIRoute]routeHash
		// Hash + Major Parameter -> bucket
		Buckets map[hashMajor]*bucket
	}
)

func (r *RateLimiterImpl) Close(ctx context.Context) {
	// TODO: wait for all buckets to unlock
}

func (r *RateLimiterImpl) getRouteHash(route *route.CompiledAPIRoute) hashMajor {
	hash, ok := r.Hashes[route.APIRoute]
	if !ok {
		// generate routeHash
		hash = routeHash(route.Method().String() + "+" + route.Path())
		r.Hashes[route.APIRoute] = hash
	}
	// return hashMajor
	return hashMajor(string(hash) + "+" + route.MajorParams())
}

func (r *RateLimiterImpl) getBucket(route *route.CompiledAPIRoute, create bool) *bucket {
	hash := r.getRouteHash(route)

	r.Lock()
	defer r.Unlock()
	b, ok := r.Buckets[hash]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{
			Remaining: 1,
			Limit:     1,
		}
		r.Buckets[hash] = b
	}
	return b
}

func (r *RateLimiterImpl) WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error {
	bucket := r.getBucket(route, true)
	bucket.Lock()

	var until time.Time
	now := time.Now()

	if bucket.Remaining == 0 && bucket.Reset.After(now) {
		until = bucket.Reset
	} else {
		until = time.Unix(0, r.Global)
	}

	if until.After(now) {
		// TODO: do we want to return early when we know rate limit bigger than ctx deadline?
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return ErrCtxTimeout
		}

		select {
		case <-ctx.Done():
			bucket.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (r *RateLimiterImpl) UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error {
	bucket := r.getBucket(route, false)
	if bucket == nil {
		return nil
	}
	defer bucket.Unlock()

	bucketID := headers.Get("X-RateLimit-bucket")

	if bucketID != "" {
		bucket.ID = bucketID
	}

	global := headers.Get("X-RateLimit-Global")

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
			r.Global = at.UnixNano()
		} else {
			bucket.Reset = at
		}

	case reset != "":
		unix, err := strconv.ParseInt(reset, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid reset %s: %s", reset, err)
		}

		bucket.Reset = time.Unix(unix, 0)
	}

	if remaining != "" {
		u, err := strconv.Atoi(remaining)
		if err != nil {
			return fmt.Errorf("invalid remaining %s: %s", remaining, err)
		}

		bucket.Remaining = u
	} else {
		// Lower remaining one just to be safe
		if bucket.Remaining > 0 {
			bucket.Remaining--
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
