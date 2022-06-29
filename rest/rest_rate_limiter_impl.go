package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/log"
	"github.com/sasha-s/go-csync"
)

// NewRateLimiter return a new default RateLimiter with the given RateLimiterConfigOpt(s).
func NewRateLimiter(opts ...RateLimiterConfigOpt) RateLimiter {
	config := DefaultRateLimiterConfig()
	config.Apply(opts)

	rateLimiter := &rateLimiterImpl{
		config:  *config,
		hashes:  map[*route.APIRoute]routeHash{},
		buckets: map[hashMajor]*bucket{},
	}

	go rateLimiter.cleanup()

	return rateLimiter
}

type (
	routeHash string
	hashMajor string

	rateLimiterImpl struct {
		config RateLimiterConfig

		// global Rate Limit
		global int64

		// route.APIRoute -> Hash
		hashes   map[*route.APIRoute]routeHash
		hashesMu sync.Mutex
		// Hash + Major Parameter -> bucket
		buckets   map[hashMajor]*bucket
		bucketsMu sync.Mutex
	}
)

func (l *rateLimiterImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *rateLimiterImpl) MaxRetries() int {
	return l.config.MaxRetries
}

func (l *rateLimiterImpl) cleanup() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	for range ticker.C {
		l.doCleanup()
	}
}

func (l *rateLimiterImpl) doCleanup() {
	l.bucketsMu.Lock()
	defer l.bucketsMu.Unlock()
	before := len(l.buckets)
	now := time.Now()
	for hash, b := range l.buckets {
		if !b.mu.TryLock() {
			continue
		}
		if b.Reset.Before(now) {
			l.Logger().Debugf("cleaning up bucket, Hash: %s, ID: %s, Reset: %s", hash, b.ID, b.Reset)
			delete(l.buckets, hash)
		}
		b.mu.Unlock()
	}
	if before != len(l.buckets) {
		l.Logger().Debugf("cleaned up %d rate limit buckets", before-len(l.buckets))
	}
}

func (l *rateLimiterImpl) Close(ctx context.Context) {
	var wg sync.WaitGroup
	for i := range l.buckets {
		wg.Add(1)
		b := l.buckets[i]
		go func() {
			_ = b.mu.CLock(ctx)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (l *rateLimiterImpl) Reset() {
	l.buckets = map[hashMajor]*bucket{}
	l.bucketsMu = sync.Mutex{}
	l.global = 0
	l.hashes = map[*route.APIRoute]routeHash{}
	l.hashesMu = sync.Mutex{}
}

func (l *rateLimiterImpl) getRouteHash(route *route.CompiledAPIRoute) hashMajor {
	l.hashesMu.Lock()
	hash, ok := l.hashes[route.APIRoute]
	if !ok {
		// generate routeHash
		hash = routeHash(route.APIRoute.Method().String() + "+" + route.APIRoute.Path())
		l.hashes[route.APIRoute] = hash
	}
	l.hashesMu.Unlock()
	if route.MajorParams() != "" {
		hash += routeHash("+" + route.MajorParams())
	}
	return hashMajor(hash)
}

func (l *rateLimiterImpl) getBucket(route *route.CompiledAPIRoute, create bool) *bucket {
	hash := l.getRouteHash(route)

	l.Logger().Trace("locking buckets")
	l.bucketsMu.Lock()
	defer func() {
		l.Logger().Trace("unlocking buckets")
		l.bucketsMu.Unlock()
	}()
	b, ok := l.buckets[hash]
	if !ok {
		if !create {
			return nil
		}

		b = &bucket{
			Remaining: 1,
			// we don't know the limit yet
			Limit: -1,
		}
		l.buckets[hash] = b
	}
	return b
}

func (l *rateLimiterImpl) WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error {
	b := l.getBucket(route, true)
	l.Logger().Tracef("locking rest bucket, ID: %s, Limit: %d, Remaining: %d, Reset: %s", b.ID, b.Limit, b.Remaining, b.Reset)
	if err := b.mu.CLock(ctx); err != nil {
		return err
	}

	var until time.Time
	now := time.Now()

	if b.Remaining == 0 && b.Reset.After(now) {
		until = b.Reset
	} else {
		until = time.Unix(0, l.global)
	}

	if until.After(now) {
		// TODO: do we want to return early when we know srate limit bigger than ctx deadline?
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return context.DeadlineExceeded
		}

		select {
		case <-ctx.Done():
			b.mu.Unlock()
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}
	return nil
}

func (l *rateLimiterImpl) UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error {
	b := l.getBucket(route, false)
	if b == nil {
		return nil
	}
	defer func() {
		l.Logger().Tracef("unlocking rest bucket, ID: %s, Limit: %d, Remaining: %d, Reset: %s", b.ID, b.Limit, b.Remaining, b.Reset)
		b.mu.Unlock()
	}()

	// no headers provided means we can't update anything and just unlock it
	if headers == nil {
		return nil
	}
	bucketID := headers.Get("X-RateLimit-Bucket")

	if bucketID != "" {
		b.ID = bucketID
	}

	global := headers.Get("X-RateLimit-Global")
	remaining := headers.Get("X-RateLimit-Remaining")
	limit := headers.Get("X-RateLimit-Limit")
	reset := headers.Get("X-RateLimit-Reset")
	retryAfter := headers.Get("Retry-After")

	l.Logger().Tracef("headers: global %s, remaining: %s, limit: %s, reset: %s, retryAfter: %s", global, remaining, limit, reset, retryAfter)

	switch {
	case retryAfter != "":
		i, err := strconv.Atoi(retryAfter)
		if err != nil {
			return fmt.Errorf("invalid retryAfter %s: %s", retryAfter, err)
		}

		at := time.Now().Add(time.Duration(i) * time.Second)

		if global != "" {
			l.global = at.UnixNano()
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
	mu        csync.Mutex
	ID        string
	Reset     time.Time
	Remaining int
	Limit     int
}
