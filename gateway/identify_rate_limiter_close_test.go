package gateway

import (
	"context"
	"testing"
	"time"
)

// TestIdentifyRateLimiterImpl_CloseThenWait ensures Close releases its internal
// lock so the limiter stays usable after being closed.
func TestIdentifyRateLimiterImpl_CloseThenWait(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter(WithIdentifyWait(100 * time.Millisecond))

	if err := r.Wait(context.Background(), 0); err != nil {
		t.Fatalf("unexpected error on first wait: %v", err)
	}
	r.Unlock(0)

	r.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- r.Wait(ctx, 0)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error waiting after close: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Wait deadlocked after Close: r.mu left locked")
	}

	r.Unlock(0)
}

// TestIdentifyRateLimiterImpl_CloseKeepsLockedBucketLocked ensures Close does
// not unlock buckets it failed to acquire before its context expired.
func TestIdentifyRateLimiterImpl_CloseKeepsLockedBucketLocked(t *testing.T) {
	t.Parallel()

	r := NewIdentifyRateLimiter(WithIdentifyWait(10 * time.Second))

	if err := r.Wait(context.Background(), 0); err != nil {
		t.Fatalf("unexpected error on first wait: %v", err)
	}
	// bucket 0 stays locked by the caller

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		r.Close(ctx)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Close deadlocked")
	}

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer waitCancel()

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- r.Wait(waitCtx, 0)
	}()

	select {
	case err := <-waitDone:
		if err == nil {
			t.Fatal("bucket was spuriously unlocked during Close")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait deadlocked after Close: r.mu left locked")
	}
}
