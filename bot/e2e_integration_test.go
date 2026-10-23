package bot_test

import (
	"errors"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

type fullMockRest struct {
	rest.Rest
	getMemberCount int
	getChannelCount int
	getRoleCount int
}

func (m *fullMockRest) GetMember(guildID snowflake.ID, userID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
	m.getMemberCount++
	if userID == 999 {
		return nil, errors.New("user not found")
	}
	return &discord.Member{
		GuildID: guildID,
		User:    discord.User{ID: userID, Username: "FetchedUser"},
	}, nil
}

func (m *fullMockRest) GetChannel(channelID snowflake.ID, opts ...rest.RequestOpt) (discord.Channel, error) {
	m.getChannelCount++
	var uc discord.UnmarshalChannel
	if err := uc.UnmarshalJSON([]byte(`{"id":"888","type":0,"guild_id":"500","name":"general"}`)); err != nil {
		return nil, err
	}
	return uc.Channel, nil
}

func (m *fullMockRest) GetRoles(guildID snowflake.ID, opts ...rest.RequestOpt) ([]discord.Role, error) {
	m.getRoleCount++
	return []discord.Role{
		{ID: 101, GuildID: guildID, Name: "AdminRole"},
		{ID: 102, GuildID: guildID, Name: "ModRole"},
	}, nil
}

func TestEndToEnd_LRU_TTL_And_RestFallback(t *testing.T) {
	guildID := snowflake.ID(500)
	user1 := snowflake.ID(1)
	user2 := snowflake.ID(2)
	user3 := snowflake.ID(3)

	// 1. Construct a Caches instance using our GroupedLRUCache for members with capacity 2 per guild
	memberLRU := cache.NewGroupedLRUCache[discord.Member](
		cache.FlagsAll,
		cache.FlagMembers,
		nil,
		2, // only keep 2 members per guild!
	)

	channelLRU := cache.NewLRUCache[discord.GuildChannel](
		cache.FlagsAll,
		cache.FlagChannels,
		nil,
		5,
	)

	caches := cache.New(
		cache.WithCaches(cache.FlagsAll),
		cache.WithMemberCache(cache.NewMemberCache(memberLRU)),
		cache.WithChannelCache(cache.NewChannelCache(channelLRU)),
	)

	mock := &fullMockRest{}

	// 2. Wrap with REST fallback enabled for Members, Channels, Roles
	fallbackCaches := cache.NewRestFallbackCaches(
		caches,
		mock,
		cache.RestFallbackFlagMembers|cache.RestFallbackFlagChannels|cache.RestFallbackFlagRoles,
	)

	// 3. Request user1 -> cache miss -> REST call -> cached
	m1, ok := fallbackCaches.Member(guildID, user1)
	if !ok || m1.User.Username != "FetchedUser" {
		t.Fatalf("expected user1 fetched, got ok=%v", ok)
	}
	if mock.getMemberCount != 1 {
		t.Fatalf("expected 1 rest call, got %d", mock.getMemberCount)
	}

	// 4. Request user1 again -> cache hit -> 0 extra REST calls
	m1Again, ok := fallbackCaches.Member(guildID, user1)
	if !ok || m1Again.User.Username != "FetchedUser" {
		t.Fatalf("expected user1 cached")
	}
	if mock.getMemberCount != 1 {
		t.Fatalf("expected still 1 rest call, got %d", mock.getMemberCount)
	}

	// 5. Request user2 -> fetched and cached (cache now has user1, user2)
	_, _ = fallbackCaches.Member(guildID, user2)
	if mock.getMemberCount != 2 {
		t.Fatalf("expected 2 rest calls, got %d", mock.getMemberCount)
	}

	// 6. Request user3 -> capacity is 2! user1 was LRU -> user1 is evicted!
	_, _ = fallbackCaches.Member(guildID, user3)
	if mock.getMemberCount != 3 {
		t.Fatalf("expected 3 rest calls, got %d", mock.getMemberCount)
	}

	// 7. Request user2 -> hit
	_, _ = fallbackCaches.Member(guildID, user2)
	if mock.getMemberCount != 3 {
		t.Fatalf("expected user2 still in cache, got %d rest calls", mock.getMemberCount)
	}

	// 8. Request user1 -> was evicted by LRU! So it triggers a fresh REST call!
	_, _ = fallbackCaches.Member(guildID, user1)
	if mock.getMemberCount != 4 {
		t.Fatalf("expected user1 to trigger new rest call after eviction, got %d", mock.getMemberCount)
	}

	// 9. Verify channel fallback
	chID := snowflake.ID(888)
	ch, ok := fallbackCaches.Channel(chID)
	if !ok || ch.ID() != chID {
		t.Fatalf("expected channel fetched via rest")
	}
	if mock.getChannelCount != 1 {
		t.Fatalf("expected 1 channel rest call, got %d", mock.getChannelCount)
	}
	// Second channel call hits cache
	_, _ = fallbackCaches.Channel(chID)
	if mock.getChannelCount != 1 {
		t.Fatalf("expected channel cache hit, got %d calls", mock.getChannelCount)
	}

	// 10. Verify roles fallback
	role, ok := fallbackCaches.Role(guildID, 101)
	if !ok || role.Name != "AdminRole" {
		t.Fatalf("expected role fetched via rest")
	}
	if mock.getRoleCount != 1 {
		t.Fatalf("expected 1 role rest call, got %d", mock.getRoleCount)
	}
	// Other role in the guild was also populated
	modRole, ok := fallbackCaches.Role(guildID, 102)
	if !ok || modRole.Name != "ModRole" {
		t.Fatalf("expected mod role cached")
	}
	if mock.getRoleCount != 1 {
		t.Fatalf("expected 1 role rest call total, got %d", mock.getRoleCount)
	}
}

func TestTTL_Expiration_InCache(t *testing.T) {
	ttlCache := cache.NewTTLCache[string](cache.FlagsAll, cache.FlagsAll, nil, 50*time.Millisecond)
	id := snowflake.ID(12345)
	ttlCache.Put(id, "hello")

	val, ok := ttlCache.Get(id)
	if !ok || val != "hello" {
		t.Fatalf("expected item present immediately")
	}

	time.Sleep(60 * time.Millisecond)

	val, ok = ttlCache.Get(id)
	if ok {
		t.Fatalf("expected item to expire after 50ms, got: %v", val)
	}
}

func TestBot_WithRestFallback(t *testing.T) {
	opt := bot.WithRestFallback(cache.RestFallbackFlagMembers)
	if opt == nil {
		t.Fatalf("expected non-nil option")
	}
}

