package core

import (
	"sync"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/internal/insecurerandstr"
)

var _ MembersChunkingManager = (*membersChunkingManagerImpl)(nil)

func NewMembersChunkingManager(bot *Bot) MembersChunkingManager {
	return &membersChunkingManagerImpl{
		chunkingRequests: map[string]*chunkingRequest{},
		bot:              bot,
	}
}

type MembersChunkingManager interface {
	Bot() *Bot
	HandleChunk(payload discord.GuildMembersChunkGatewayEvent)

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

type membersChunkingManagerImpl struct {
	sync.RWMutex
	chunkingRequests map[string]*chunkingRequest

	bot *Bot
}

func (m *membersChunkingManagerImpl) Bot() *Bot {
	return m.bot
}
func (m *membersChunkingManagerImpl) HandleChunk(payload discord.GuildMembersChunkGatewayEvent) {
	request, ok := m.chunkingRequests[payload.Nonce]
	if !ok {
		m.Bot().Logger.Warn("received unknown member chunk event")
		return
	}

	for _, member := range payload.Members {
		coreMember := m.Bot().EntityBuilder.CreateMember(request.GuildID, member, CacheStrategyYes)
		if request.memberFilterFunc != nil && !request.memberFilterFunc(coreMember) {
			continue
		}
		request.Lock()
		if request.returnChan == nil {
			// channel is nil anyway abort all member parsing/sending
			request.Unlock()
			return
		}
		request.returnChan <- coreMember
		request.Unlock()
	}

	// all chunks sent cleanup
	if payload.ChunkIndex == payload.ChunkCount-1 {
		request.Lock()
		close(request.returnChan)
		request.returnChan = nil
		request.Unlock()
		m.Lock()
		delete(m.chunkingRequests, payload.Nonce)
		m.Unlock()
	}
}

func (m *membersChunkingManagerImpl) requestGuildMembers(command discord.RequestGuildMembersCommand, memberFilterFunc func(member *Member) bool) (<-chan *Member, func()) {
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
	returnChan := make(chan *Member)
	request := &chunkingRequest{
		RequestGuildMembersCommand: command,
		memberFilterFunc:           nil,
		returnChan:                 returnChan,
	}

	m.Lock()
	m.chunkingRequests[nonce] = request
	m.Unlock()
	return returnChan, func() {
		request.Lock()
		close(request.returnChan)
		request.returnChan = nil
		request.Unlock()
		m.Lock()
		delete(m.chunkingRequests, nonce)
		m.Unlock()
	}
}

func (m *membersChunkingManagerImpl) LoadMembers(guildID discord.Snowflake, presences bool, userIDs ...discord.Snowflake) (<-chan *Member, func()) {
	return m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Presences: presences,
		UserIDs:   userIDs,
	}, nil)
}

func (m *membersChunkingManagerImpl) FindMembers(guildID discord.Snowflake, presences bool, memberFindFunc func(member *Member) bool) (<-chan *Member, func()) {
	query := ""
	limit := 0
	return m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, memberFindFunc)
}

func (m *membersChunkingManagerImpl) SearchMembers(guildID discord.Snowflake, presences bool, query string, limit int) (<-chan *Member, func()) {
	return m.requestGuildMembers(discord.RequestGuildMembersCommand{
		GuildID:   guildID,
		Query:     &query,
		Limit:     &limit,
		Presences: presences,
	}, nil)
}
