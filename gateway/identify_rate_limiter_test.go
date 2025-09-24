package gateway

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIdentifyRateLimiterImpl(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter(WithIdentifyWait(100 * time.Millisecond))

	start := time.Now()

	var wg sync.WaitGroup
	for shardID := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.Wait(context.Background(), shardID)
			assert.NoError(t, err)
			r.Unlock(shardID)
		}()
	}
	wg.Wait()

	expected := start.Add(200 * time.Millisecond)
	assert.WithinDuration(t, expected, time.Now(), 10*time.Millisecond)
}

func TestIdentifyRateLimiterImpl_WithMaxConcurrency(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter(WithIdentifyMaxConcurrency(3), WithIdentifyWait(100*time.Millisecond))

	start := time.Now()

	var wg sync.WaitGroup
	for shardID := range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.Wait(context.Background(), shardID)
			assert.NoError(t, err)
			r.Unlock(shardID)
		}()
	}
	wg.Wait()

	expected := start.Add(100 * time.Millisecond)
	assert.WithinDuration(t, expected, time.Now(), 10*time.Millisecond)
}

func TestIdentifyRateLimiterImpl_WaitWithTimeout(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := r.Wait(ctx, 0)
	assert.NoError(t, err)

	err = r.Wait(ctx, 0)
	if assert.Error(t, err) {
		assert.Equal(t, context.DeadlineExceeded, err)
	}

	r.Unlock(0)
}
