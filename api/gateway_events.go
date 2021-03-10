package api


type ReadyEvent struct {
	GatewayCommand
	D ReadyEventData `json:"d"`
}

type ReadyEventData struct {
	User            User        `json:"user"`
	PrivateChannels []DMChannel `json:"channel"`
	Guilds          []Guild     `json:"guild"`
	SessionID       string      `json:"session_id"`
	Shard           [2]int      `json:"shard,omitempty"`
}