package core

import (
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
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

	LoadMembers(guildID discord.Snowflake, userIDs []discord.Snowflake, opts ...gateway.TaskOpt) ([]*Member, error)

	SearchMembers(guildID discord.Snowflake, query string, limit int, opts ...gateway.TaskOpt) ([]*Member, error)
	LoadAllMembers(guildID discord.Snowflake, opts ...gateway.TaskOpt) ([]*Member, error)

	FindMembers(guildID discord.Snowflake, memberFilterFunc func(member *Member) bool, opts ...gateway.TaskOpt) ([]*Member, error)
}

type chunkingRequest struct {
	gateway.TaskConfig
	nonce string

	memberFilterFunc func(member *Member) bool
	memberFunc       func(member *Member)
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

	for _, member := range payload.Members {
		if request.Ctx.Err() != nil {
			cleanupRequest(m, request)
			return
		}
		coreMember := m.Bot().EntityBuilder.CreateMember(payload.GuildID, member, CacheStrategyYes)
		if request.memberFilterFunc != nil && !request.memberFilterFunc(coreMember) {
			continue
		}
		if request.memberFunc != nil {
			request.memberFunc(coreMember)
		}
	}

	// all chunks sent cleanup
	if payload.ChunkIndex == payload.ChunkCount-1 {
		cleanupRequest(m, request)
	}
}

func cleanupRequest(m *memberChunkingManagerImpl, request *chunkingRequest) {
	m.chunkingRequestsMu.Lock()
	delete(m.chunkingRequests, request.nonce)
	m.chunkingRequestsMu.Unlock()
}

func (m *memberChunkingManagerImpl) requestGuildMembers(guildID discord.Snowflake, query *string, limit *int, userIDs []discord.Snowflake, memberFilterFunc func(member *Member) bool, memberFunc func(member *Member)) error {
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
	request := &chunkingRequest{
		nonce:            nonce,
		memberFilterFunc: memberFilterFunc,
		memberFunc:       memberFunc,
	}

	m.chunkingRequestsMu.Lock()
	m.chunkingRequests[nonce] = request
	m.chunkingRequestsMu.Unlock()

	shard, err := m.Bot().Shard(guildID)
	if err != nil {
		return err
	}

	command := discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     query,
		Limit:     limit,
		Presences: shard.Config().GatewayIntents.Has(discord.GatewayIntentGuildPresences),
		UserIDs:   userIDs,
		Nonce:     nonce,
	}

	return shard.SendContext(discord.NewGatewayCommand(discord.GatewayOpcodeRequestGuildMembers, command))
}

func (m *memberChunkingManagerImpl) LoadAllMembers(guildID discord.Snowflake) error {
	query := ""
	limit := 0
	_, err := m.requestGuildMembers(guildID, &query, &limit, nil, nil, nil)
	return err
}

func (m *memberChunkingManagerImpl) LoadMembers(guildID discord.Snowflake, userIDs ...discord.Snowflake) (<-chan *Member, func(), error) {
	returnChan := make(chan *Member)
	cls, err := m.requestGuildMembers(guildID, nil, nil, userIDs, nil, returnChan)
	return returnChan, cls, err
}

func (m *memberChunkingManagerImpl) FindMembers(guildID discord.Snowflake, memberFindFunc func(member *Member) bool) (<-chan *Member, func(), error) {
	returnChan := make(chan *Member)
	query := ""
	limit := 0
	cls, err := m.requestGuildMembers(guildID, &query, &limit, nil, memberFindFunc, returnChan)
	return returnChan, cls, err
}

func (m *memberChunkingManagerImpl) SearchMembers(guildID discord.Snowflake, query string, limit int) (<-chan *Member, func(), error) {
	returnChan := make(chan *Member)
	cls, err := m.requestGuildMembers(guildID, &query, &limit, nil, nil, returnChan)
	return returnChan, cls, err
}
