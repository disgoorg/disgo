package gateway

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestIdentifyRateLimiterImpl(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter(WithIdentifyWait(100 * time.Millisecond))

	start := time.Now()

	var wg sync.WaitGroup
	for shardID := range 3 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := r.Wait(context.Background(), id)
			if err != nil {
				t.Errorf("unexpected error for shard %d: %v", id, err)
			}
			r.Unlock(id)
		}(shardID)
	}
	wg.Wait()

	expected := start.Add(200 * time.Millisecond)
	if !time.Now().After(expected.Add(-10*time.Millisecond)) || !time.Now().Before(expected.Add(10*time.Millisecond)) {
		t.Errorf("expected current time to be within 10ms of %v, got %v", expected, time.Now())
	}
}

func TestIdentifyRateLimiterImpl_WithMaxConcurrency(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter(WithIdentifyMaxConcurrency(3), WithIdentifyWait(100*time.Millisecond))

	start := time.Now()

	var wg sync.WaitGroup
	for shardID := range 6 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := r.Wait(context.Background(), id)
			if err != nil {
				t.Errorf("unexpected error for shard %d: %v", id, err)
			}
			r.Unlock(id)
		}(shardID)
	}
	wg.Wait()

	expected := start.Add(100 * time.Millisecond)
	if !time.Now().After(expected.Add(-10*time.Millisecond)) || !time.Now().Before(expected.Add(10*time.Millisecond)) {
		t.Errorf("expected current time to be within 10ms of %v, got %v", expected, time.Now())
	}
}

func TestIdentifyRateLimiterImpl_WaitWithTimeout(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := r.Wait(ctx, 0)
	if err != nil {
		t.Fatalf("unexpected error on first wait: %v", err)
	}

	err = r.Wait(ctx, 0)
	if err == nil {
		t.Errorf("expected error on second wait due to timeout, got nil")
	} else if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded error, got %v", err)
	}

	r.Unlock(0)
}
