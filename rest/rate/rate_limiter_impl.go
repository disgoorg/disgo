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
func NewLimiter(config *Config) Limiter {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	return &LimiterImpl{
		config:  *config,
		hashes:  map[*route.APIRoute]routeHash{},
		buckets: map[hashMajor]*bucket{},
	}
}

//goland:noinspection GoNameStartsWithPackageName
type (
	routeHash string
	hashMajor string

	LimiterImpl struct {
		sync.Mutex

		config Config

		// global Rate Limit
		global int64

		// route.APIRoute -> Hash
		hashes map[*route.APIRoute]routeHash
		// Hash + Major Parameter -> bucket
		buckets map[hashMajor]*bucket
	}
)

func (r *LimiterImpl) Logger() log.Logger {
	return r.config.Logger
}

func (r *LimiterImpl) Close(ctx context.Context) {
	// TODO: wait for all buckets to unlock
}

func (r *LimiterImpl) Config() Config {
	return r.config
}

func (r *LimiterImpl) getRouteHash(route *route.CompiledAPIRoute) hashMajor {
	hash, ok := r.hashes[route.APIRoute]
	if !ok {
		// generate routeHash
		hash = routeHash(route.Method().String() + "+" + route.Path())
		r.hashes[route.APIRoute] = hash
	}
	// return hashMajor
	return hashMajor(string(hash) + "+" + route.MajorParams())
}

func (r *LimiterImpl) getBucket(route *route.CompiledAPIRoute, create bool) *bucket {
	hash := r.getRouteHash(route)

	r.Logger().Debug("locking rate limiter")
	r.Lock()
	defer func() {
		r.Logger().Debug("unlocking rate limiter")
		r.Unlock()
	}()
	b, ok := r.buckets[hash]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{
			Remaining: 1,
			// we don't know the limit yet
			Limit: -1,
		}
		r.buckets[hash] = b
	}
	return b
}

func (r *LimiterImpl) WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error {
	b := r.getBucket(route, true)
	r.Logger().Debugf("locking bucket: %+v", b)
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

func (r *LimiterImpl) UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error {
	b := r.getBucket(route, false)
	if b == nil {
		return nil
	}
	defer func() {
		r.Logger().Debugf("unlocking bucket: %+v", b)
		b.Unlock()
	}()

	bucketID := headers.Get("X-RateLimit-Bucket")

	if bucketID != "" {
		b.ID = bucketID
	}

	global := headers.Get("X-RateLimit-Global")
	remaining := headers.Get("X-RateLimit-Remaining")
	limit := headers.Get("X-RateLimit-Limit")
	reset := headers.Get("X-RateLimit-Reset")
	retryAfter := headers.Get("Retry-After")

	r.Logger().Debugf("headers: global %s, remaining: %s, limit: %s, reset: %s, retryAfter: %s", global, remaining, limit, reset, retryAfter)

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
		unix, err := strconv.ParseFloat(reset, 64)
		if err != nil {
			return fmt.Errorf("invalid reset %s: %s", reset, err)
		}

		sec := int64(unix)
		b.Reset = time.Unix(sec, int64((unix-float64(sec))*float64(time.Second)))
	}

	if limit != "" {
		u, err := strconv.Atoi(limit)
		if err != nil {
			return fmt.Errorf("invalid limit %s: %s", limit, err)
		}

		b.Limit = u
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
