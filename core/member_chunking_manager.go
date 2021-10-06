package core

import (
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

	LoadAllMembers(guildID discord.Snowflake, presences bool) error
	LoadMembers(guildID discord.Snowflake, presences bool, userIDs ...discord.Snowflake) (<-chan *Member, func(), error)
	FindMembers(guildID discord.Snowflake, presences bool, memberFilterFunc func(member *Member) bool) (<-chan *Member, func(), error)
	SearchMembers(guildID discord.Snowflake, presences bool, query string, limit int) (<-chan *Member, func(), error)
}

type chunkingRequest struct {
	nonce            string
	memberFilterFunc func(member *Member) bool

	sync.Mutex
	returnChan chan<- *Member
}

type memberChunkingManagerImpl struct {
	bot                  *Bot
	memberChunkingFilter MemberChunkingFilter

	sync.RWMutex
	chunkingRequests map[string]*chunkingRequest
}

func (m *memberChunkingManagerImpl) Bot() *Bot {
	return m.bot
}

func (m *memberChunkingManagerImpl) MemberChunkingFilter() MemberChunkingFilter {
	return m.memberChunkingFilter
}

func (m *memberChunkingManagerImpl) HandleChunk(payload discord.GuildMembersChunkGatewayEvent) {
	request, ok := m.chunkingRequests[payload.Nonce]
	if !ok {
		m.Bot().Logger.Debug("received unknown member chunk event: ", payload)
		return
	}

	for _, member := range payload.Members {
		coreMember := m.Bot().EntityBuilder.CreateMember(payload.GuildID, member, CacheStrategyYes)
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
	delete(m.chunkingRequests, request.nonce)
	m.Unlock()
}

func (m *memberChunkingManagerImpl) requestGuildMembers(command discord.RequestGuildMembersCommand, memberFilterFunc func(member *Member) bool, returnChan chan<- *Member) (func(), error) {
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
	request := &chunkingRequest{
		nonce:            nonce,
		memberFilterFunc: memberFilterFunc,
		returnChan:       returnChan,
	}

	m.Lock()
	m.chunkingRequests[nonce] = request
	m.Unlock()

	command.Nonce = nonce
	shard, err := m.Bot().Shard(command.GuildID)
	if err != nil {
		return nil, err
	}
	err = shard.Send(discord.NewGatewayCommand(discord.GatewayOpcodeRequestGuildMembers, command))
	if err != nil {
		return nil, err
	}
	return func() {
		cleanupRequest(m, request)
	}, nil
}

func (m *memberChunkingManagerImpl) LoadAllMembers(guildID discord.Snowflake, presences bool) error {
	query := ""
	limit := 0
	_, err := m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, nil, nil)
	return err
}

func (m *memberChunkingManagerImpl) LoadMembers(guildID discord.Snowflake, presences bool, userIDs ...discord.Snowflake) (<-chan *Member, func(), error) {
	returnChan := make(chan *Member)
	cls, err := m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Presences: presences,
		UserIDs:   userIDs,
	}, nil, returnChan)
	return returnChan, cls, err
}

func (m *memberChunkingManagerImpl) FindMembers(guildID discord.Snowflake, presences bool, memberFindFunc func(member *Member) bool) (<-chan *Member, func(), error) {
	returnChan := make(chan *Member)
	query := ""
	limit := 0
	cls, err := m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, memberFindFunc, returnChan)
	return returnChan, cls, err
}

func (m *memberChunkingManagerImpl) SearchMembers(guildID discord.Snowflake, presences bool, query string, limit int) (<-chan *Member, func(), error) {
	returnChan := make(chan *Member)
	cls, err := m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, nil, returnChan)
	return returnChan, cls, err
}
