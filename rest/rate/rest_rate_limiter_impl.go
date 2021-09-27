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
	"github.com/sasha-s/go-csync"
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
	return &limiterImpl{
		config:  *config,
		hashes:  map[*route.APIRoute]routeHash{},
		buckets: map[hashMajor]*bucket{},
	}
}

//goland:noinspection GoNameStartsWithPackageName
type (
	routeHash string
	hashMajor string

	limiterImpl struct {
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

func (l *limiterImpl) Logger() log.Logger {
	return l.config.Logger
}

func (l *limiterImpl) Close(ctx context.Context) {
	var wg sync.WaitGroup
	for _, b := range l.buckets {
		wg.Add(1)
		b := b
		go func() {
			_ = b.CLock(ctx)
			b.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}

func (l *limiterImpl) Config() Config {
	return l.config
}

func (l *limiterImpl) getRouteHash(route *route.CompiledAPIRoute) hashMajor {
	hash, ok := l.hashes[route.APIRoute]
	if !ok {
		// generate routeHash
		hash = routeHash(route.APIRoute.Method().String() + "+" + route.APIRoute.Path())
		l.hashes[route.APIRoute] = hash
	}
	// return hashMajor
	return hashMajor(string(hash) + "+" + route.MajorParams())
}

func (l *limiterImpl) getBucket(route *route.CompiledAPIRoute, create bool) *bucket {
	hash := l.getRouteHash(route)

	l.Logger().Debug("locking rest rate limiter")
	l.Lock()
	defer func() {
		l.Logger().Debug("unlocking rest rate limiter")
		l.Unlock()
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

func (l *limiterImpl) WaitBucket(ctx context.Context, route *route.CompiledAPIRoute) error {
	b := l.getBucket(route, true)
	l.Logger().Debugf("locking rest bucket: %+v", b)
	if err := b.CLock(ctx); err != nil {
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
		// TODO: do we want to return early when we know rate limit bigger than ctx deadline?
		if deadline, ok := ctx.Deadline(); ok && until.After(deadline) {
			return context.DeadlineExceeded
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

func (l *limiterImpl) UnlockBucket(route *route.CompiledAPIRoute, headers http.Header) error {
	b := l.getBucket(route, false)
	if b == nil {
		return nil
	}
	defer func() {
		l.Logger().Debugf("unlocking rest bucket: %+v", b)
		b.Unlock()
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

	l.Logger().Debugf("headers: global %s, remaining: %s, limit: %s, reset: %s, retryAfter: %s", global, remaining, limit, reset, retryAfter)

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
	csync.Mutex
	ID        string
	Reset     time.Time
	Remaining int
	Limit     int
}
