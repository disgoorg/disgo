package voice

import (
	"context"
	"sync"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

type (
	Manager interface {
		HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate)
		HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate)

		CreateConn(guildID snowflake.ID, channelID snowflake.ID) Conn
		GetConn(guildID snowflake.ID) Conn
		ForEachCon(f func(connection Conn))
		RemoveConn(guildID snowflake.ID)

		Close(ctx context.Context)
	}
	StateUpdateFunc func(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error
)

func NewManager(voiceStateUpdateFunc StateUpdateFunc, userID snowflake.ID, opts ...ManagerConfigOpt) Manager {
	config := DefaultManagerConfig()
	config.Apply(opts)
	return &managerImpl{
		config:               config,
		voiceStateUpdateFunc: voiceStateUpdateFunc,
		userID:               userID,
		conns:                map[snowflake.ID]Conn{},
	}
}

type managerImpl struct {
	config               ManagerConfig
	voiceStateUpdateFunc StateUpdateFunc
	userID               snowflake.ID

	conns   map[snowflake.ID]Conn
	connsMu sync.Mutex
}

func (m *managerImpl) HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate) {
	m.config.Logger.Debugf("VoiceStateUpdate for guild: %s", update.GuildID)

	conn := m.GetConn(update.GuildID)
	if conn == nil {
		return
	}
	conn.HandleVoiceStateUpdate(update)
}

func (m *managerImpl) HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate) {
	m.config.Logger.Debugf("VoiceServerUpdate for guild: %s", update.GuildID)

	conn := m.GetConn(update.GuildID)
	if conn == nil {
		return
	}
	conn.HandleVoiceServerUpdate(update)
}

func (m *managerImpl) CreateConn(guildID snowflake.ID, channelID snowflake.ID) Conn {
	m.config.Logger.Debugf("Creating new voice conn for guild: %s, channel: %s", guildID, channelID)
	if conn := m.GetConn(guildID); conn != nil {
		return conn
	}

	m.connsMu.Lock()
	defer m.connsMu.Unlock()

	var once sync.Once
	removeFunc := func() { once.Do(func() { m.RemoveConn(guildID) }) }

	conn := m.config.ConnCreateFunc(guildID, channelID, m.userID, m.voiceStateUpdateFunc, removeFunc, append([]ConnConfigOpt{WithConnLogger(m.config.Logger)}, m.config.ConnOpts...)...)
	m.conns[guildID] = conn

	return conn
}

func (m *managerImpl) GetConn(guildID snowflake.ID) Conn {
	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	return m.conns[guildID]
}

func (m *managerImpl) ForEachCon(f func(connection Conn)) {
	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	for _, connection := range m.conns {
		f(connection)
	}
}

func (m *managerImpl) RemoveConn(guildID snowflake.ID) {
	m.config.Logger.Debugf("Removing voice conn for guild: %s", guildID)
	conn := m.GetConn(guildID)
	if conn == nil {
		return
	}
	m.connsMu.Lock()
	defer m.connsMu.Unlock()
	delete(m.conns, guildID)
}

func (m *managerImpl) Close(ctx context.Context) {
	m.connsMu.Lock()
	conns := m.conns
	m.connsMu.Unlock()
	for i := range conns {
		conns[i].Close(ctx)
	}
	m.conns = map[snowflake.ID]Conn{}
}
