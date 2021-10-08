package core

import (
	"context"
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/internal/insecurerandstr"
)

var _ MemberChunkingManager = (*memberChunkingManagerImpl)(nil)

func NewMemberChunkingManager(bot *Bot, memberChunkingFilter MemberChunkingFilter) MemberChunkingManager {
	return &memberChunkingManagerImpl{
		bot:                  bot,
		memberChunkingFilter: memberChunkingFilter,
		chunkingRequests:     map[string]*chunkingRequest{},
	}
}

type MemberChunkingManager interface {
	Bot() *Bot
	MemberChunkingFilter() MemberChunkingFilter

	HandleChunk(payload discord.GuildMembersChunkGatewayEvent)

	LoadMembers(guildID discord.Snowflake, userIDs ...discord.Snowflake) ([]*Member, error)
	SearchMembers(guildID discord.Snowflake, query string, limit int) ([]*Member, error)
	LoadAllMembers(guildID discord.Snowflake) ([]*Member, error)
	FindMembers(guildID discord.Snowflake, memberFilterFunc func(member *Member) bool) ([]*Member, error)

	LoadMembersCtx(ctx context.Context, guildID discord.Snowflake, userIDs ...discord.Snowflake) ([]*Member, error)
	SearchMembersCtx(ctx context.Context, guildID discord.Snowflake, query string, limit int) ([]*Member, error)
	LoadAllMembersCtx(ctx context.Context, guildID discord.Snowflake) ([]*Member, error)
	FindMembersCtx(ctx context.Context, guildID discord.Snowflake, memberFilterFunc func(member *Member) bool) ([]*Member, error)

	LoadMembersChan(guildID discord.Snowflake, userIDs ...discord.Snowflake) (<-chan *Member, func(), error)
	SearchMembersChan(guildID discord.Snowflake, query string, limit int) (<-chan *Member, func(), error)
	LoadAllMembersChan(guildID discord.Snowflake) (<-chan *Member, func(), error)
	FindMembersChan(guildID discord.Snowflake, memberFilterFunc func(member *Member) bool) (<-chan *Member, func(), error)
}

type chunkingRequest struct {
	sync.Mutex
	nonce string

	memberChan chan<- *Member
	memberFilterFunc func(member *Member) bool

	chunks int
}

type memberChunkingManagerImpl struct {
	bot                  *Bot
	memberChunkingFilter MemberChunkingFilter

	chunkingRequestsMu sync.RWMutex
	chunkingRequests   map[string]*chunkingRequest
}

func (m *memberChunkingManagerImpl) Bot() *Bot {
	return m.bot
}

func (m *memberChunkingManagerImpl) MemberChunkingFilter() MemberChunkingFilter {
	return m.memberChunkingFilter
}

func (m *memberChunkingManagerImpl) HandleChunk(payload discord.GuildMembersChunkGatewayEvent) {
	m.chunkingRequestsMu.RLock()
	request, ok := m.chunkingRequests[payload.Nonce]
	m.chunkingRequestsMu.RUnlock()
	if !ok {
		m.Bot().Logger.Debug("received unknown member chunk event: ", payload)
		return
	}

	request.Lock()
	defer request.Unlock()

	for _, member := range payload.Members {
		coreMember := m.Bot().EntityBuilder.CreateMember(payload.GuildID, member, CacheStrategyYes)
		if request.memberFilterFunc != nil && !request.memberFilterFunc(coreMember) {
			continue
		}
		request.memberChan <- coreMember
	}

	// all chunks sent cleanup
	if request.chunks == payload.ChunkCount {
		cleanupRequest2(m, request)
		return
	}
	request.chunks++
}

func cleanupRequest2(m *memberChunkingManagerImpl, request *chunkingRequest) {
	close(request.memberChan)
	m.chunkingRequestsMu.Lock()
	delete(m.chunkingRequests, request.nonce)
	m.chunkingRequestsMu.Unlock()
}

func (m *memberChunkingManagerImpl) requestGuildMembersChan(ctx context.Context, guildID discord.Snowflake, query *string, limit *int, userIDs []discord.Snowflake, memberFilterFunc func(member *Member) bool) (<-chan *Member, func(), error) {
	shard, err := m.Bot().Shard(guildID)
	if err != nil {
		return nil, nil, err
	}

	if shard.Config().GatewayIntents.Missing(discord.GatewayIntentGuildMembers) {
		return nil, nil, discord.ErrNoGuildMembersIntent
	}

	var nonce string
	for {
		nonce = insecurerandstr.RandStr(32)
		m.chunkingRequestsMu.RLock()
		_, ok := m.chunkingRequests[nonce]
		m.chunkingRequestsMu.RUnlock()
		if !ok {
			break
		}
	}
	memberChan := make(chan *Member)
	request := &chunkingRequest{
		nonce:      nonce,
		memberChan: memberChan,
		memberFilterFunc: memberFilterFunc,
	}

	m.chunkingRequestsMu.Lock()
	m.chunkingRequests[nonce] = request
	m.chunkingRequestsMu.Unlock()

	command := discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     query,
		Limit:     limit,
		Presences: shard.Config().GatewayIntents.Has(discord.GatewayIntentGuildPresences),
		UserIDs:   userIDs,
		Nonce:     nonce,
	}

	return memberChan, func() {
		cleanupRequest2(m, request)
	}, shard.SendCtx(ctx, discord.NewGatewayCommand(discord.GatewayOpcodeRequestGuildMembers, command))
}

func (m *memberChunkingManagerImpl) requestGuildMembers(ctx context.Context, guildID discord.Snowflake, query *string, limit *int, userIDs []discord.Snowflake, memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	var members []*Member
	memberChan, cls, err := m.requestGuildMembersChan(ctx, guildID, query, limit, userIDs, memberFilterFunc)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-ctx.Done():
			cls()
			return nil, ctx.Err()
		case member, ok := <-memberChan:
			if !ok {
				return members, nil
			}
			members = append(members, member)
		}
	}
}

func (m *memberChunkingManagerImpl) LoadMembers(guildID discord.Snowflake, userIDs ...discord.Snowflake) ([]*Member, error) {
	return m.LoadMembersCtx(context.Background(), guildID, userIDs...)
}
func (m *memberChunkingManagerImpl) SearchMembers(guildID discord.Snowflake, query string, limit int) ([]*Member, error) {
	return m.SearchMembersCtx(context.Background(), guildID, query, limit)
}
func (m *memberChunkingManagerImpl) LoadAllMembers(guildID discord.Snowflake) ([]*Member, error) {
	return m.LoadAllMembersCtx(context.Background(), guildID)
}
func (m *memberChunkingManagerImpl) FindMembers(guildID discord.Snowflake, memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	return m.FindMembersCtx(context.Background(), guildID, memberFilterFunc)
}

func (m *memberChunkingManagerImpl) LoadMembersCtx(ctx context.Context, guildID discord.Snowflake, userIDs ...discord.Snowflake) ([]*Member, error) {
	return m.requestGuildMembers(ctx, guildID, nil, nil, userIDs, nil)
}

func (m *memberChunkingManagerImpl) SearchMembersCtx(ctx context.Context, guildID discord.Snowflake, query string, limit int) ([]*Member, error) {
	return m.requestGuildMembers(ctx, guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) LoadAllMembersCtx(ctx context.Context, guildID discord.Snowflake) ([]*Member, error) {
	query := ""
	limit := 0
	return m.requestGuildMembers(ctx, guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) FindMembersCtx(ctx context.Context, guildID discord.Snowflake, memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	query := ""
	limit := 0
	return m.requestGuildMembers(ctx, guildID, &query, &limit, nil, memberFilterFunc)
}

func (m *memberChunkingManagerImpl) LoadMembersChan(guildID discord.Snowflake, userIDs ...discord.Snowflake) (<-chan *Member, func(), error) {
	return m.requestGuildMembersChan(context.Background(), guildID, nil, nil, userIDs, nil)
}

func (m *memberChunkingManagerImpl) SearchMembersChan(guildID discord.Snowflake, query string, limit int) (<-chan *Member, func(), error) {
	return m.requestGuildMembersChan(context.Background(), guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) LoadAllMembersChan(guildID discord.Snowflake) (<-chan *Member, func(), error) {
	query := ""
	limit := 0
	return m.requestGuildMembersChan(context.Background(), guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) FindMembersChan(guildID discord.Snowflake, memberFilterFunc func(member *Member) bool) (<-chan *Member, func(), error) {
	query := ""
	limit := 0
	return m.requestGuildMembersChan(context.Background(), guildID, &query, &limit, nil, nil)
}
