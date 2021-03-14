package api

// ReadyEvent is the event sent by discord when you successfully Identify
type ReadyEvent struct {
	GatewayCommand
	D ReadyEventData `json:"d"`
}

// ReadyEventData is the ReadyEvent.D payload
type ReadyEventData struct {
	User            User        `json:"user"`
	PrivateChannels []DMChannel `json:"channel"`
	Guilds          []Guild     `json:"guild_events"`
	SessionID       string      `json:"session_id"`
	Shard           [2]int      `json:"shard,omitempty"`
}