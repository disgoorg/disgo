package bot

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/internal/insecurerandstr"
)

var _ MemberChunkingManager = (*memberChunkingManagerImpl)(nil)

var ErrNoUserIDs = errors.New("no user ids to request")

// NewMemberChunkingManager returns a new MemberChunkingManager with the given MemberChunkingFilter.
func NewMemberChunkingManager(client Client, logger *slog.Logger, memberChunkingFilter MemberChunkingFilter) MemberChunkingManager {
	if memberChunkingFilter == nil {
		memberChunkingFilter = MemberChunkingFilterNone
	}
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("name", "bot_member_chunking_manager"))

	return &memberChunkingManagerImpl{
		client:               client,
		logger:               logger,
		memberChunkingFilter: memberChunkingFilter,
		chunkingRequests:     map[string]*chunkingRequest{},
	}
}

// MemberChunkingManager is used to request members for guilds from the discord gateway.
type MemberChunkingManager interface {
	// MemberChunkingFilter returns the configured MemberChunkingFilter used by this MemberChunkingManager.
	MemberChunkingFilter() MemberChunkingFilter

	// HandleChunk handles the discord.EventGuildMembersChunk event payloads from the discord gateway.
	HandleChunk(payload gateway.EventGuildMembersChunk)

	// RequestMembers requests members from the given guildID and userIDs.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembers(guildID snowflake.ID, userIDs ...snowflake.ID) ([]discord.Member, error)
	// RequestAllMembers requests all members from the given guildID.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestAllMembers(guildID snowflake.ID) ([]discord.Member, error)
	// RequestMembersWithQuery requests members from the given guildID and query.
	// query : string the username starts with
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersWithQuery(guildID snowflake.ID, query string, limit int) ([]discord.Member, error)
	// RequestMembersWithFilter requests members from the given guildID and userIDs. memberFilterFunc is used to filter all returned members.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersWithFilter(guildID snowflake.ID, memberFilterFunc func(member discord.Member) bool) ([]discord.Member, error)

	// RequestMembersCtx requests members from the given guildID and userIDs.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersCtx(ctx context.Context, guildID snowflake.ID, userIDs ...snowflake.ID) ([]discord.Member, error)
	// RequestAllMembersCtx requests all members from the given guildID.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestAllMembersCtx(ctx context.Context, guildID snowflake.ID) ([]discord.Member, error)
	// RequestMembersWithQueryCtx requests members from the given guildID and query.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersWithQueryCtx(ctx context.Context, guildID snowflake.ID, query string, limit int) ([]discord.Member, error)
	// RequestMembersWithFilterCtx requests members from the given guildID and userIDs. memberFilterFunc is used to filter all returned members.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersWithFilterCtx(ctx context.Context, guildID snowflake.ID, memberFilterFunc func(member discord.Member) bool) ([]discord.Member, error)

	// RequestMembersChan requests members from the given guildID and userIDs.
	// Returns a channel which will receive the members.
	// Returns a function which can be used to cancel the request and close the channel.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersChan(guildID snowflake.ID, userIDs ...snowflake.ID) (<-chan discord.Member, func(), error)
	// RequestAllMembersChan requests all members from the given guildID.
	// Returns a channel which will receive the members.
	// Returns a function which can be used to cancel the request and close the channel.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestAllMembersChan(guildID snowflake.ID) (<-chan discord.Member, func(), error)
	// RequestMembersWithQueryChan requests members from the given guildID and query.
	// Returns a channel which will receive the members.
	// Returns a function which can be used to cancel the request and close the channel.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersWithQueryChan(guildID snowflake.ID, query string, limit int) (<-chan discord.Member, func(), error)
	// RequestMembersWithFilterChan requests members from the given guildID and userIDs. memberFilterFunc is used to filter all returned members.
	// Returns a channel which will receive the members.
	// Returns a function which can be used to cancel the request and close the channel.
	// Notice: This action requires the gateway.IntentGuildMembers.
	RequestMembersWithFilterChan(guildID snowflake.ID, memberFilterFunc func(member discord.Member) bool) (<-chan discord.Member, func(), error)
}

type chunkingRequest struct {
	sync.Mutex
	nonce string

	memberChan       chan<- discord.Member
	memberFilterFunc func(member discord.Member) bool

	chunks int
}

type memberChunkingManagerImpl struct {
	client               Client
	logger               *slog.Logger
	memberChunkingFilter MemberChunkingFilter

	chunkingRequestsMu sync.RWMutex
	chunkingRequests   map[string]*chunkingRequest
}

func (m *memberChunkingManagerImpl) MemberChunkingFilter() MemberChunkingFilter {
	return m.memberChunkingFilter
}

func (m *memberChunkingManagerImpl) HandleChunk(payload gateway.EventGuildMembersChunk) {
	m.chunkingRequestsMu.RLock()
	request, ok := m.chunkingRequests[payload.Nonce]
	m.chunkingRequestsMu.RUnlock()
	if !ok {
		m.logger.Debug("received unknown member chunk event", slog.Any("payload", payload))
		return
	}

	request.Lock()
	defer request.Unlock()

	for _, member := range payload.Members {
		// try to cache member
		m.client.Caches().AddMember(member)
		if request.memberFilterFunc != nil && !request.memberFilterFunc(member) {
			continue
		}
		request.memberChan <- member
	}

	// all chunks sent cleanup
	if request.chunks == payload.ChunkCount-1 {
		cleanupRequest(m, request)
		return
	}
	request.chunks++
}

func cleanupRequest(m *memberChunkingManagerImpl, request *chunkingRequest) {
	close(request.memberChan)
	m.chunkingRequestsMu.Lock()
	delete(m.chunkingRequests, request.nonce)
	m.chunkingRequestsMu.Unlock()
}

func (m *memberChunkingManagerImpl) requestGuildMembersChan(ctx context.Context, guildID snowflake.ID, query *string, limit *int, userIDs []snowflake.ID, memberFilterFunc func(member discord.Member) bool) (<-chan discord.Member, func(), error) {
	shard, err := m.client.Shard(guildID)
	if err != nil {
		return nil, nil, err
	}

	if shard.Intents().Missing(gateway.IntentGuildMembers) {
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
	memberChan := make(chan discord.Member)
	request := &chunkingRequest{
		nonce:            nonce,
		memberChan:       memberChan,
		memberFilterFunc: memberFilterFunc,
	}

	m.chunkingRequestsMu.Lock()
	m.chunkingRequests[nonce] = request
	m.chunkingRequestsMu.Unlock()

	command := gateway.MessageDataRequestGuildMembers{
		GuildID:   guildID,
		Query:     query,
		Limit:     limit,
		Presences: shard.Intents().Has(gateway.IntentGuildPresences),
		UserIDs:   userIDs,
		Nonce:     nonce,
	}

	return memberChan, func() {
		cleanupRequest(m, request)
	}, shard.Send(ctx, gateway.OpcodeRequestGuildMembers, command)
}

func (m *memberChunkingManagerImpl) requestGuildMembers(ctx context.Context, guildID snowflake.ID, query *string, limit *int, userIDs []snowflake.ID, memberFilterFunc func(member discord.Member) bool) ([]discord.Member, error) {
	var members []discord.Member
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

func (m *memberChunkingManagerImpl) RequestMembers(guildID snowflake.ID, userIDs ...snowflake.ID) ([]discord.Member, error) {
	return m.RequestMembersCtx(context.Background(), guildID, userIDs...)
}
func (m *memberChunkingManagerImpl) RequestMembersWithQuery(guildID snowflake.ID, query string, limit int) ([]discord.Member, error) {
	return m.RequestMembersWithQueryCtx(context.Background(), guildID, query, limit)
}
func (m *memberChunkingManagerImpl) RequestAllMembers(guildID snowflake.ID) ([]discord.Member, error) {
	return m.RequestAllMembersCtx(context.Background(), guildID)
}
func (m *memberChunkingManagerImpl) RequestMembersWithFilter(guildID snowflake.ID, memberFilterFunc func(member discord.Member) bool) ([]discord.Member, error) {
	return m.RequestMembersWithFilterCtx(context.Background(), guildID, memberFilterFunc)
}

func (m *memberChunkingManagerImpl) RequestMembersCtx(ctx context.Context, guildID snowflake.ID, userIDs ...snowflake.ID) ([]discord.Member, error) {
	if len(userIDs) == 0 {
		return nil, ErrNoUserIDs
	}
	return m.requestGuildMembers(ctx, guildID, nil, nil, userIDs, nil)
}

func (m *memberChunkingManagerImpl) RequestMembersWithQueryCtx(ctx context.Context, guildID snowflake.ID, query string, limit int) ([]discord.Member, error) {
	return m.requestGuildMembers(ctx, guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) RequestAllMembersCtx(ctx context.Context, guildID snowflake.ID) ([]discord.Member, error) {
	query := ""
	limit := 0
	return m.requestGuildMembers(ctx, guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) RequestMembersWithFilterCtx(ctx context.Context, guildID snowflake.ID, memberFilterFunc func(member discord.Member) bool) ([]discord.Member, error) {
	query := ""
	limit := 0
	return m.requestGuildMembers(ctx, guildID, &query, &limit, nil, memberFilterFunc)
}

func (m *memberChunkingManagerImpl) RequestMembersChan(guildID snowflake.ID, userIDs ...snowflake.ID) (<-chan discord.Member, func(), error) {
	if len(userIDs) == 0 {
		return nil, nil, ErrNoUserIDs
	}
	return m.requestGuildMembersChan(context.Background(), guildID, nil, nil, userIDs, nil)
}

func (m *memberChunkingManagerImpl) RequestMembersWithQueryChan(guildID snowflake.ID, query string, limit int) (<-chan discord.Member, func(), error) {
	return m.requestGuildMembersChan(context.Background(), guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) RequestAllMembersChan(guildID snowflake.ID) (<-chan discord.Member, func(), error) {
	query := ""
	limit := 0
	return m.requestGuildMembersChan(context.Background(), guildID, &query, &limit, nil, nil)
}

func (m *memberChunkingManagerImpl) RequestMembersWithFilterChan(guildID snowflake.ID, memberFilterFunc func(member discord.Member) bool) (<-chan discord.Member, func(), error) {
	query := ""
	limit := 0
	return m.requestGuildMembersChan(context.Background(), guildID, &query, &limit, nil, memberFilterFunc)
}
