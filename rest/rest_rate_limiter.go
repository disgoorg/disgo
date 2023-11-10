package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sasha-s/go-csync"
)

const (
	// MaxRetries is the maximum number of retries the client should do
	MaxRetries = 10
	// CleanupInterval is the interval at which the rate limiter cleans up old buckets
	CleanupInterval = time.Second * 10
)

// RateLimiter can be used to supply your own rate limit implementation
type RateLimiter interface {
	// MaxRetries returns the maximum number of retries the client should do
	MaxRetries() int

	// Close gracefully closes the RateLimiter.
	// If the context deadline is exceeded, the RateLimiter will be closed immediately.
	Close(ctx context.Context)

	// Reset resets the rate limiter to its initial state
	Reset()

	// WaitBucket waits for the given bucket to be available for new requests & locks it
	WaitBucket(ctx context.Context, endpoint *CompiledEndpoint) error

	// UnlockBucket unlocks the given bucket and calculates the rate limit for the next request
	UnlockBucket(endpoint *CompiledEndpoint, rs *http.Response) error
}

// NewRateLimiter return a new default RateLimiter with the given RateLimiterConfigOpt(s).
func NewRateLimiter(opts ...RateLimiterConfigOpt) RateLimiter {
	config := DefaultRateLimiterConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "rest_rate_limiter"))

	rateLimiter := &rateLimiterImpl{
		config:  *config,
		hashes:  map[*Endpoint]string{},
		buckets: map[string]*bucket{},
	}

	go rateLimiter.cleanup()

	return rateLimiter
}

type (
	rateLimiterImpl struct {
		config RateLimiterConfig

		// global Rate Limit
		global time.Time

		// APIRoute -> Hash
		hashes   map[*Endpoint]string
		hashesMu sync.Mutex
		// Hash + Major Parameter -> bucket
		buckets   map[string]*bucket
		bucketsMu sync.Mutex
	}
)

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
			l.config.Logger.Debug("cleaning up bucket", slog.String("hash", hash), slog.String("id", b.ID), slog.Time("reset", b.Reset))
			delete(l.buckets, hash)
		}
		b.mu.Unlock()
	}
	if before != len(l.buckets) {
		l.config.Logger.Debug("cleaned up rate limit buckets", slog.Int("before", before), slog.Int("after", len(l.buckets)), slog.Int("removed", before-len(l.buckets)))
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
	l.buckets = map[string]*bucket{}
	l.bucketsMu = sync.Mutex{}
	l.global = time.Time{}
	l.hashes = map[*Endpoint]string{}
	l.hashesMu = sync.Mutex{}
}

func (l *rateLimiterImpl) getRouteHash(endpoint *CompiledEndpoint) string {
	l.hashesMu.Lock()
	hash, ok := l.hashes[endpoint.Endpoint]
	if !ok {
		// generate routeHash
		hash = endpoint.Endpoint.Method + "+" + endpoint.Endpoint.Route
		l.hashes[endpoint.Endpoint] = hash
	}
	l.hashesMu.Unlock()
	if endpoint.MajorParams != "" {
		hash += "+" + endpoint.MajorParams
	}
	return hash
}

func (l *rateLimiterImpl) getBucket(endpoint *CompiledEndpoint, create bool) *bucket {
	hash := l.getRouteHash(endpoint)

	l.config.Logger.Debug("locking buckets")
	l.bucketsMu.Lock()
	defer func() {
		l.config.Logger.Debug("unlocking buckets")
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

func (l *rateLimiterImpl) WaitBucket(ctx context.Context, endpoint *CompiledEndpoint) error {
	b := l.getBucket(endpoint, true)
	l.config.Logger.Debug("locking rest bucket", slog.String("id", b.ID), slog.Int("limit", b.Limit), slog.Int("remaining", b.Remaining), slog.Time("reset", b.Reset))
	if err := b.mu.CLock(ctx); err != nil {
		return err
	}

	var until time.Time
	now := time.Now()

	if b.Remaining == 0 && b.Reset.After(now) {
		until = b.Reset
	} else {
		until = l.global
	}

	if until.After(now) {
		// TODO: do we want to return early when we know the rate limit bigger than ctx deadline?
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

func (l *rateLimiterImpl) UnlockBucket(endpoint *CompiledEndpoint, rs *http.Response) error {
	b := l.getBucket(endpoint, false)
	if b == nil {
		return nil
	}
	defer func() {
		l.config.Logger.Debug("unlocking rest bucket", slog.String("id", b.ID), slog.Int("limit", b.Limit), slog.Int("remaining", b.Remaining), slog.Time("reset", b.Reset))
		b.mu.Unlock()
	}()

	// no response provided means we can't update anything and just unlock it
	if rs == nil || rs.Header == nil {
		return nil
	}
	bucketHeader := rs.Header.Get("X-RateLimit-Bucket")

	// if we don't have a bucket header, we can't update anything
	if bucketHeader == "" {
		return nil
	}

	b.ID = bucketHeader

	global := rs.Header.Get("X-RateLimit-Global") != ""
	cloudflare := rs.Header.Get("via") == ""
	remainingHeader := rs.Header.Get("X-RateLimit-Remaining")
	limitHeader := rs.Header.Get("X-RateLimit-Limit")
	resetHeader := rs.Header.Get("X-RateLimit-Reset")
	resetAfterHeader := rs.Header.Get("X-RateLimit-Reset-After")
	retryAfterHeader := rs.Header.Get("Retry-After")

	l.config.Logger.Debug("ratelimit response headers", slog.Int("code", rs.StatusCode), slog.Bool("global", global), slog.Bool("cloudflare", cloudflare), slog.String("remaining", remainingHeader), slog.String("limit", limitHeader), slog.String("reset", resetHeader), slog.String("reset_after", resetAfterHeader), slog.String("retry_after", retryAfterHeader))

	// we hit a rate limit. let's see if it was global cloudflare or a route specific one
	if rs.StatusCode == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(retryAfterHeader)
		if err != nil {
			return fmt.Errorf("invalid retryAfter %s: %w", retryAfterHeader, err)
		}
		reset := time.Now().Add(time.Second * time.Duration(retryAfter))
		if global {
			l.global = reset
			l.config.Logger.Warn("global rate limit exceeded", slog.Int("retry_after", retryAfter))
		} else if cloudflare {
			l.global = reset
			l.config.Logger.Warn("cloudflare rate limit exceeded", slog.Int("retry_after", retryAfter))
		} else {
			b.Remaining = 0
			b.Reset = reset
			l.config.Logger.Warn("rate limit exceeded", slog.String("endpoint", endpoint.URL), slog.Int("retry_after", retryAfter))
		}
		return nil
	}

	if limitHeader != "" {
		limit, err := strconv.Atoi(limitHeader)
		if err != nil {
			return fmt.Errorf("invalid limit %s: %w", limitHeader, err)
		}
		b.Limit = limit
	}

	if remainingHeader != "" {
		remaining, err := strconv.Atoi(remainingHeader)
		if err != nil {
			return fmt.Errorf("invalid remaining %s: %w", remainingHeader, err)
		}
		b.Remaining = remaining
	}

	// we prioritize the reset after header over the reset header as it's more accurate due to clock differences
	if resetAfterHeader != "" {
		resetAfter, err := strconv.ParseFloat(resetAfterHeader, 64)
		if err != nil {
			return fmt.Errorf("invalid reset after %s: %w", resetAfterHeader, err)
		}

		b.Reset = time.Now().Add(time.Duration(resetAfter) * time.Second)
	} else if resetHeader != "" {
		reset, err := strconv.ParseFloat(resetHeader, 64)
		if err != nil {
			return fmt.Errorf("invalid reset %s: %w", resetHeader, err)
		}

		sec := int64(reset)
		b.Reset = time.Unix(sec, int64((reset-float64(sec))*float64(time.Second)))
	} else {
		return fmt.Errorf("no reset or reset after header found in response")
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
