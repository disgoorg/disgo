package voice

import (
	"context"
	"sync"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

type Manager interface {
	HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate)
	HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate)

	CreateConn(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID) Conn
	GetConn(guildID snowflake.ID) Conn
	ForEachCon(f func(connection Conn))
	DeleteConn(guildID snowflake.ID)

	Close(ctx context.Context)
}

func NewManager(opts ...ManagerConfigOpt) Manager {
	config := DefaultManagerConfig()
	config.Apply(opts)
	return &managerImpl{
		config:      *config,
		connections: map[snowflake.ID]Conn{},
	}
}

type managerImpl struct {
	config ManagerConfig

	connections   map[snowflake.ID]Conn
	connectionsMu sync.Mutex
}

func (m *managerImpl) HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	connection, ok := m.connections[update.GuildID]
	if !ok {
		return
	}

	m.config.Logger.Debugf("VoiceStateUpdate for guild: %s", update.GuildID)
	connection.HandleVoiceStateUpdate(update)
}

func (m *managerImpl) HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	connection, ok := m.connections[update.GuildID]
	if !ok {
		return
	}
	m.config.Logger.Debugf("VoiceServerUpdate for guild: %s", update.GuildID)
	connection.HandleVoiceServerUpdate(update)
}

func (m *managerImpl) CreateConn(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID) Conn {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	m.config.Logger.Debugf("Creating new voice conn for guild: %s, channel: %s, user: %s", guildID, channelID, userID)
	connection := m.config.ConnCreateFunc(guildID, channelID, userID, append([]ConnConfigOpt{WithConnLogger(m.config.Logger)}, m.config.ConnOpts...)...)
	m.connections[guildID] = connection
	return connection
}

func (m *managerImpl) GetConn(guildID snowflake.ID) Conn {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	return m.connections[guildID]
}

func (m *managerImpl) ForEachCon(f func(connection Conn)) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	for _, connection := range m.connections {
		f(connection)
	}
}

func (m *managerImpl) DeleteConn(guildID snowflake.ID) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	m.config.Logger.Debugf("Removing voice conn for guild: %s", guildID)
	if conn, ok := m.connections[guildID]; ok {
		conn.Close()
		delete(m.connections, guildID)
	}
}

func (m *managerImpl) Close(ctx context.Context) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	for _, connection := range m.connections {
		connection.Close()
	}
}
