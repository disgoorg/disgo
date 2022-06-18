package voice

import (
	"context"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Manager interface {
	HandleVoiceStateUpdate(update discord.VoiceStateUpdate)
	HandleVoiceServerUpdate(update discord.VoiceServerUpdate)
	CreateConnection(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID) Connection
	GetConnection(guildID snowflake.ID) Connection
	ForConnections(f func(connection Connection))
	RemoveConnection(guildID snowflake.ID)

	Close(ctx context.Context)
}

func NewManager(opts ...ManagerConfigOpt) Manager {
	config := DefaultManagerConfig()
	config.Apply(opts)
	return &managerImpl{
		config:      *config,
		connections: map[snowflake.ID]Connection{},
	}
}

type managerImpl struct {
	config ManagerConfig

	connections   map[snowflake.ID]Connection
	connectionsMu sync.Mutex
}

func (m *managerImpl) HandleVoiceStateUpdate(update discord.VoiceStateUpdate) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	connection, ok := m.connections[update.GuildID]
	if !ok {
		return
	}

	m.config.Logger.Debugf("VoiceStateUpdate for guild: %s", update.GuildID)
	connection.HandleVoiceStateUpdate(update)
}

func (m *managerImpl) HandleVoiceServerUpdate(update discord.VoiceServerUpdate) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	connection, ok := m.connections[update.GuildID]
	if !ok {
		return
	}
	m.config.Logger.Debugf("VoiceServerUpdate for guild: %s", update.GuildID)
	connection.HandleVoiceServerUpdate(update)
}

func (m *managerImpl) CreateConnection(guildID snowflake.ID, channelID snowflake.ID, userID snowflake.ID) Connection {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	m.config.Logger.Debugf("Creating new voice connection for guild: %s, channel: %s, user: %s", guildID, channelID, userID)
	connection := m.config.ConnectionCreateFunc(guildID, channelID, userID, append([]ConnectionConfigOpt{WithConnectionLogger(m.config.Logger)}, m.config.ConnectionOpts...)...)
	m.connections[guildID] = connection
	return connection
}

func (m *managerImpl) GetConnection(guildID snowflake.ID) Connection {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	return m.connections[guildID]
}

func (m *managerImpl) ForConnections(f func(connection Connection)) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	for _, connection := range m.connections {
		f(connection)
	}
}

func (m *managerImpl) RemoveConnection(guildID snowflake.ID) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	m.config.Logger.Debugf("Removing voice connection for guild: %s", guildID)
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
