package api

// GatewayRs contains the response for GET /gateway
type GatewayRs struct {
	URL string `json:"url"`
}

// GatewayBotRs contains the response for GET /gateway/bot
type GatewayBotRs struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}
