package voice

import (
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func NewManager(opts ...ManagerConfigOpt) *Manager {
	config := DefaultManagerConfig()
	config.Apply(opts)
	return &Manager{
		config:      *config,
		connections: map[snowflake.ID]*Connection{},
	}
}

type Manager struct {
	config ManagerConfig

	connections   map[snowflake.ID]*Connection
	connectionsMu sync.Mutex
}

func (m *Manager) HandleVoiceStateUpdate(update discord.VoiceState) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	connection, ok := m.connections[update.GuildID]
	if !ok {
		return
	}
	connection.HandleVoiceStateUpdate(update)
}

func (m *Manager) HandleVoiceServerUpdate(update discord.VoiceServerUpdate) {
	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	connection, ok := m.connections[update.GuildID]
	if !ok {
		return
	}
	connection.HandleVoiceServerUpdate(update)
}

func (m *Manager) NewConnection(guildID snowflake.ID, userID snowflake.ID, opts ...ConnectionConfigOpt) *Connection {
	connection := m.config.ConnectionCreateFunc(guildID, userID, opts...)

	m.connectionsMu.Lock()
	defer m.connectionsMu.Unlock()
	m.connections[guildID] = connection

	return connection
}
