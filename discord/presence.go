package discord

// Presence (https://discord.com/developers/docs/topics/gateway#presence-update)
type Presence struct {
	User       User         `json:"user"`
	GuildID    Snowflake    `json:"guild_id"`
	Status     OnlineStatus `json:"status"`
	Activities []Activity   `json:"activities"`
}

// OnlineStatus (https://discord.com/developers/docs/topics/gateway#update-presence-status-types)
type OnlineStatus string

//goland:noinspection GoUnusedConst
const (
	OnlineStatusOnline    OnlineStatus = "online"
	OnlineStatusDND       OnlineStatus = "dnd"
	OnlineStatusIdle      OnlineStatus = "idle"
	OnlineStatusInvisible OnlineStatus = "invisible"
	OnlineStatusOffline   OnlineStatus = "offline"
)

// ClientStatus (https://discord.com/developers/docs/topics/gateway#client-status-object)
type ClientStatus struct {
	Desktop *OnlineStatus `json:"desktop,omitempty"`
	Mobile  *OnlineStatus `json:"mobile,omitempty"`
	Web     *OnlineStatus `json:"web,omitempty"`
}
