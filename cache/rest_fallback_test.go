package cache_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

type mockRest struct {
	rest.Rest
	getMemberFunc  func(guildID snowflake.ID, userID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error)
	getChannelFunc func(channelID snowflake.ID, opts ...rest.RequestOpt) (discord.Channel, error)
	getRolesFunc   func(guildID snowflake.ID, opts ...rest.RequestOpt) ([]discord.Role, error)
}

func (m *mockRest) GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
	if m.getMemberFunc != nil {
		return m.getMemberFunc(guildID, userID, opts...)
	}
	return nil, errors.New("not found")
}

func (m *mockRest) GetChannel(channelID snowflake.ID, opts ...rest.RequestOpt) (discord.Channel, error) {
	if m.getChannelFunc != nil {
		return m.getChannelFunc(channelID, opts...)
	}
	return nil, errors.New("not found")
}

func (m *mockRest) GetRoles(guildID snowflake.ID, opts ...rest.RequestOpt) ([]discord.Role, error) {
	if m.getRolesFunc != nil {
		return m.getRolesFunc(guildID, opts...)
	}
	return nil, errors.New("not found")
}

func TestRestFallbackCaches_Member(t *testing.T) {
	caches := cache.New(cache.WithCaches(cache.FlagsAll))
	guildID := snowflake.ID(123)
	userID := snowflake.ID(456)

	restCalls := 0
	mock := &mockRest{
		getMemberFunc: func(gID snowflake.ID, uID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
			restCalls++
			return &discord.Member{
				GuildID: gID,
				User:    discord.User{ID: uID, Username: "testuser"},
			}, nil
		},
	}

	// 1. Without fallback enabled
	disabledWrapper := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagsNone)
	if _, ok := disabledWrapper.Member(guildID, userID); ok {
		t.Fatalf("expected not ok when disabled")
	}
	if restCalls != 0 {
		t.Fatalf("expected 0 rest calls when disabled, got %d", restCalls)
	}

	// 2. With fallback enabled
	enabledWrapper := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagMembers)
	member, ok := enabledWrapper.Member(guildID, userID)
	if !ok || member.User.Username != "testuser" {
		t.Fatalf("expected to fetch member via rest, got ok=%v, member=%+v", ok, member)
	}
	if restCalls != 1 {
		t.Fatalf("expected 1 rest call, got %d", restCalls)
	}

	// 3. Second call should hit the cache without another REST call
	member2, ok := enabledWrapper.Member(guildID, userID)
	if !ok || member2.User.Username != "testuser" {
		t.Fatalf("expected cache hit on second call")
	}
	if restCalls != 1 {
		t.Fatalf("expected rest calls to remain 1, got %d", restCalls)
	}
}

func TestRestFallbackCaches_ConcurrentMissDeduplication(t *testing.T) {
	caches := cache.New(cache.WithCaches(cache.FlagsAll))
	guildID := snowflake.ID(123)
	userID := snowflake.ID(789)

	var restCalls atomic.Int32
	mock := &mockRest{
		getMemberFunc: func(gID snowflake.ID, uID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
			restCalls.Add(1)
			// Simulate network latency so concurrent callers overlap
			time.Sleep(20 * time.Millisecond)
			return &discord.Member{
				GuildID: gID,
				User:    discord.User{ID: uID, Username: "concurrent_user"},
			}, nil
		},
	}

	wrapper := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagMembers)

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			member, ok := wrapper.Member(guildID, userID)
			if !ok || member.User.Username != "concurrent_user" {
				t.Errorf("failed to get member: ok=%v, member=%+v", ok, member)
			}
		}()
	}
	wg.Wait()

	// In-flight singleflight request coalescing ensures only 1 REST call was made
	if calls := restCalls.Load(); calls != 1 {
		t.Fatalf("expected exactly 1 REST call due to request deduplication, got %d", calls)
	}
}

func TestRestFallbackCaches_NonGuildChannel(t *testing.T) {
	caches := cache.New(cache.WithCaches(cache.FlagsAll))
	channelID := snowflake.ID(999)

	mock := &mockRest{
		getChannelFunc: func(cID snowflake.ID, opts ...rest.RequestOpt) (discord.Channel, error) {
			// Return a direct message channel, which does not implement discord.GuildChannel
			return &discord.DMChannel{}, nil
		},
	}

	wrapper := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagChannels)
	ch, ok := wrapper.Channel(channelID)
	if ok || ch != nil {
		t.Fatalf("expected non-guild channel to return (nil, false), got ok=%v, ch=%+v", ok, ch)
	}
}

func TestRestFallbackCaches_PreventDoubleWrapping(t *testing.T) {
	caches := cache.New(cache.WithCaches(cache.FlagsAll))
	mock := &mockRest{}

	wrapper1 := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagMembers)
	wrapper2 := cache.NewRestFallbackCaches(wrapper1, mock, cache.RestFallbackFlagChannels)

	// Verify that wrapper2 is the exact same instance as wrapper1, not a nested wrapper
	if wrapper1 != wrapper2 {
		t.Fatalf("expected NewRestFallbackCaches to reuse existing wrapper instance to avoid double-wrapping")
	}
}

func TestRestFallbackCaches_RestError(t *testing.T) {
	caches := cache.New(cache.WithCaches(cache.FlagsAll))
	mock := &mockRest{
		getMemberFunc: func(guildID snowflake.ID, userID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
			return nil, errors.New("discord 404 not found")
		},
	}

	wrapper := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagMembers)
	_, ok := wrapper.Member(1, 2)
	if ok {
		t.Fatalf("expected error from REST to return ok=false")
	}

	// Verify it was not stored in the underlying cache
	if _, ok := caches.Member(1, 2); ok {
		t.Fatalf("expected failed request not to be cached")
	}
}

func TestRestFallbackCaches_Role(t *testing.T) {
	caches := cache.New(cache.WithCaches(cache.FlagsAll))
	guildID := snowflake.ID(10)
	roleID := snowflake.ID(20)

	restCalls := 0
	mock := &mockRest{
		getRolesFunc: func(gID snowflake.ID, opts ...rest.RequestOpt) ([]discord.Role, error) {
			restCalls++
			return []discord.Role{
				{ID: roleID, GuildID: gID, Name: "Admin"},
			}, nil
		},
	}

	wrapper := cache.NewRestFallbackCaches(caches, mock, cache.RestFallbackFlagRoles)
	role, ok := wrapper.Role(guildID, roleID)
	if !ok || role.Name != "Admin" {
		t.Fatalf("expected role Admin, got ok=%v, role=%+v", ok, role)
	}
	if restCalls != 1 {
		t.Fatalf("expected 1 rest call, got %d", restCalls)
	}

	// Cache hit
	role2, ok := wrapper.Role(guildID, roleID)
	if !ok || role2.Name != "Admin" {
		t.Fatalf("expected cache hit for role")
	}
	if restCalls != 1 {
		t.Fatalf("expected rest calls to remain 1, got %d", restCalls)
	}
}
