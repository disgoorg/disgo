package voice

import (
	"context"
	"log/slog"
	"sync"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/gateway"
)

type (
	// StateUpdateFunc is used to update the voice state of the bot from the Manager.
	StateUpdateFunc func(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error

	// Manager manages all voice connections.
	Manager interface {
		// HandleVoiceStateUpdate handles a gateway.EventVoiceStateUpdate
		HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate)

		// HandleVoiceServerUpdate handles a gateway.EventVoiceServerUpdate
		HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate)

		// CreateConn creates a new voice connection for the given guild.
		CreateConn(guildID snowflake.ID) Conn

		// GetConn returns the voice connection for the given guild.
		GetConn(guildID snowflake.ID) Conn

		// ForEachCon runs the given function for each voice connection. This is thread-safe.
		ForEachCon(f func(connection Conn))

		// RemoveConn removes the voice connection for the given guild.
		RemoveConn(guildID snowflake.ID)

		// Close closes all voice connections.
		Close(ctx context.Context)
	}
)

// NewManager creates a new Manager.
func NewManager(voiceStateUpdateFunc StateUpdateFunc, userID snowflake.ID, opts ...ManagerConfigOpt) Manager {
	config := DefaultManagerConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "voice"))

	return &managerImpl{
		config:               *config,
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
	m.config.Logger.Debug("new VoiceStateUpdate", slog.Int64("guild_id", int64(update.GuildID)))

	conn := m.GetConn(update.GuildID)
	if conn == nil {
		return
	}
	conn.HandleVoiceStateUpdate(update)
}

func (m *managerImpl) HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate) {
	m.config.Logger.Debug("new VoiceServerUpdate", slog.Int64("guild_id", int64(update.GuildID)))

	conn := m.GetConn(update.GuildID)
	if conn == nil {
		return
	}
	conn.HandleVoiceServerUpdate(update)
}

func (m *managerImpl) CreateConn(guildID snowflake.ID) Conn {
	m.config.Logger.Debug("Creating new voice conn", slog.Int64("guild_id", int64(guildID)))
	if conn := m.GetConn(guildID); conn != nil {
		return conn
	}

	m.connsMu.Lock()
	defer m.connsMu.Unlock()

	var once sync.Once
	removeFunc := func() { once.Do(func() { m.RemoveConn(guildID) }) }

	conn := m.config.ConnCreateFunc(guildID, m.userID, m.voiceStateUpdateFunc, removeFunc, append([]ConnConfigOpt{WithConnLogger(m.config.Logger)}, m.config.ConnOpts...)...)
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
	m.config.Logger.Debug("Removing voice conn", slog.Int64("guild_id", int64(guildID)))
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
