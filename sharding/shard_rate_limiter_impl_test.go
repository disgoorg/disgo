package sharding

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShardRateLimiterImpl(t *testing.T) {
	t.Parallel()

	r := NewRateLimiter(WithIdentifyWait(100 * time.Millisecond))

	start := time.Now()

	var wg sync.WaitGroup
	for shardID := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.WaitBucket(context.Background(), shardID)
			assert.NoError(t, err)
			r.UnlockBucket(shardID)
		}()
	}
	wg.Wait()

	expected := start.Add(200 * time.Millisecond)
	assert.WithinDuration(t, expected, time.Now(), 10*time.Millisecond)
}

func TestShardRateLimiterImpl_WithMaxConcurrency(t *testing.T) {
	t.Parallel()

	r := NewRateLimiter(WithMaxConcurrency(3), WithIdentifyWait(100*time.Millisecond))

	start := time.Now()

	var wg sync.WaitGroup
	for shardID := range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.WaitBucket(context.Background(), shardID)
			assert.NoError(t, err)
			r.UnlockBucket(shardID)
		}()
	}
	wg.Wait()

	expected := start.Add(100 * time.Millisecond)
	assert.WithinDuration(t, expected, time.Now(), 10*time.Millisecond)
}

func TestShardRateLimiterImpl_WaitBucketWithTimeout(t *testing.T) {
	t.Parallel()

	r := NewRateLimiter()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := r.WaitBucket(ctx, 0)
	assert.NoError(t, err)

	err = r.WaitBucket(ctx, 0)
	if assert.Error(t, err) {
		assert.Equal(t, context.DeadlineExceeded, err)
	}

	r.UnlockBucket(0)
}
