package rest

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRateLimiterUnlockRetryAfter(t *testing.T) {
	for _, test := range []struct {
		name       string
		bucket     bool
		global     bool
		cloudflare bool
	}{
		{name: "route integer", bucket: true},
		{name: "global integer", bucket: true, global: true},
		{name: "cloudflare integer", bucket: true, cloudflare: true},
		{name: "global without bucket", global: true},
		{name: "cloudflare without bucket", cloudflare: true},
		{name: "route without bucket"},
	} {
		t.Run(test.name, func(t *testing.T) {
			limiter := &rateLimiterImpl{config: defaultRateLimiterConfig(), buckets: map[string]*bucket{}}
			endpoint := GetChannel.Compile(nil, 123)
			if err := limiter.Wait(context.Background(), endpoint); err != nil {
				t.Fatal(err)
			}
			header := http.Header{}
			header.Set("Retry-After", "2")
			if test.bucket {
				header.Set("X-RateLimit-Bucket", "test-bucket")
			}
			if test.global {
				header.Set("X-RateLimit-Global", "true")
			}
			if !test.cloudflare {
				header.Set("Via", "1.1 google")
			}
			before := time.Now()
			if err := limiter.Unlock(endpoint, &http.Response{StatusCode: http.StatusTooManyRequests, Header: header}); err != nil {
				t.Fatal(err)
			}
			after := time.Now()
			delay := 2 * time.Second
			reset := limiter.getBucket(endpoint, false).Reset
			waitEndpoint := endpoint
			if test.global || test.cloudflare {
				reset = limiter.global
				waitEndpoint = GetChannel.Compile(nil, 456)
			}
			if reset.Before(before.Add(delay)) || reset.After(after.Add(delay)) {
				t.Fatalf("reset = %v, want between %v and %v", reset, before.Add(delay), after.Add(delay))
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			if err := limiter.Wait(ctx, waitEndpoint); !errors.Is(err, context.DeadlineExceeded) {
				if err == nil {
					_ = limiter.Unlock(waitEndpoint, nil)
				}
				t.Fatalf("Wait = %v, want context.DeadlineExceeded", err)
			}
		})
	}
}

func TestRateLimiterUnlockResetAfter(t *testing.T) {
	for _, test := range []struct {
		value string
		delay time.Duration
	}{
		{value: "0.25", delay: 250 * time.Millisecond},
		{value: "1.25", delay: 1250 * time.Millisecond},
		{value: "2", delay: 2 * time.Second},
	} {
		t.Run(test.value, func(t *testing.T) {
			limiter := &rateLimiterImpl{config: defaultRateLimiterConfig(), buckets: map[string]*bucket{}}
			endpoint := GetChannel.Compile(nil, 123)
			if err := limiter.Wait(context.Background(), endpoint); err != nil {
				t.Fatal(err)
			}
			header := http.Header{}
			header.Set("X-RateLimit-Bucket", "test-bucket")
			header.Set("X-RateLimit-Limit", "5")
			header.Set("X-RateLimit-Remaining", "0")
			header.Set("X-RateLimit-Reset-After", test.value)
			before := time.Now()
			if err := limiter.Unlock(endpoint, &http.Response{StatusCode: http.StatusOK, Header: header}); err != nil {
				t.Fatal(err)
			}
			after := time.Now()
			reset := limiter.getBucket(endpoint, false).Reset
			if reset.Before(before.Add(test.delay)) || reset.After(after.Add(test.delay)) {
				t.Fatalf("reset = %v, want between %v and %v", reset, before.Add(test.delay), after.Add(test.delay))
			}
		})
	}
}
