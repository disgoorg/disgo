package core

import (
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/internal/insecurerandstr"
)

var _ MemberChunkingManager = (*memberChunkingManagerImpl)(nil)

func NewMemberChunkingManager(bot *Bot) MemberChunkingManager {
	return &memberChunkingManagerImpl{
		chunkingRequests: map[string]*chunkingRequest{},
		bot:              bot,
	}
}

type MemberChunkingManager interface {
	Bot() *Bot
	HandleChunk(payload discord.GuildMembersChunkGatewayEvent)

	LoadAllMembers(guildID discord.Snowflake, presences bool)
	LoadMembers(guildID discord.Snowflake, presences bool, userIDs ...discord.Snowflake) (<-chan *Member, func())
	FindMembers(guildID discord.Snowflake, presences bool, memberFilterFunc func(member *Member) bool) (<-chan *Member, func())
	SearchMembers(guildID discord.Snowflake, presences bool, query string, limit int) (<-chan *Member, func())
}

type chunkingRequest struct {
	discord.RequestGuildMembersCommand
	memberFilterFunc func(member *Member) bool

	sync.Mutex
	returnChan chan<- *Member
}

type memberChunkingManagerImpl struct {
	sync.RWMutex
	chunkingRequests map[string]*chunkingRequest

	bot *Bot
}

func (m *memberChunkingManagerImpl) Bot() *Bot {
	return m.bot
}
func (m *memberChunkingManagerImpl) HandleChunk(payload discord.GuildMembersChunkGatewayEvent) {
	request, ok := m.chunkingRequests[payload.Nonce]
	if !ok {
		m.Bot().Logger.Debug("received unknown member chunk event: ", payload)
		return
	}

	for _, member := range payload.Members {
		coreMember := m.Bot().EntityBuilder.CreateMember(request.GuildID, member, CacheStrategyYes)
		if request.memberFilterFunc != nil && !request.memberFilterFunc(coreMember) {
			continue
		}
		request.Lock()
		if request.returnChan != nil {
			request.returnChan <- coreMember
		}
		request.Unlock()
	}

	// all chunks sent cleanup
	if payload.ChunkIndex == payload.ChunkCount-1 {
		cleanupRequest(m, request)
	}
}

func cleanupRequest(m *memberChunkingManagerImpl, request *chunkingRequest) {
	if request.returnChan != nil {
		request.Lock()
		close(request.returnChan)
		request.returnChan = nil
		request.Unlock()
	}
	m.Lock()
	delete(m.chunkingRequests, request.Nonce)
	m.Unlock()
}

func (m *memberChunkingManagerImpl) requestGuildMembers(command discord.RequestGuildMembersCommand, memberFilterFunc func(member *Member) bool, returnChan chan *Member) (<-chan *Member, func()) {
	var nonce string
	for {
		nonce = insecurerandstr.RandStr(32)
		m.RLock()
		_, ok := m.chunkingRequests[nonce]
		m.RUnlock()
		if !ok {
			break
		}
	}
	command.Nonce = nonce
	request := &chunkingRequest{
		RequestGuildMembersCommand: command,
		memberFilterFunc:           memberFilterFunc,
		returnChan:                 returnChan,
	}

	m.Lock()
	m.chunkingRequests[nonce] = request
	m.Unlock()
	return returnChan, func() {
		cleanupRequest(m, request)
	}
}

func (m *memberChunkingManagerImpl) LoadAllMembers(guildID discord.Snowflake, presences bool) {
	query := ""
	limit := 0
	m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, nil, nil)
}

func (m *memberChunkingManagerImpl) LoadMembers(guildID discord.Snowflake, presences bool, userIDs ...discord.Snowflake) (<-chan *Member, func()) {
	return m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Presences: presences,
		UserIDs:   userIDs,
	}, nil, make(chan *Member))
}

func (m *memberChunkingManagerImpl) FindMembers(guildID discord.Snowflake, presences bool, memberFindFunc func(member *Member) bool) (<-chan *Member, func()) {
	query := ""
	limit := 0
	return m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, memberFindFunc, make(chan *Member))
}

func (m *memberChunkingManagerImpl) SearchMembers(guildID discord.Snowflake, presences bool, query string, limit int) (<-chan *Member, func()) {
	return m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, nil, make(chan *Member))
}
