package discord

// VoiceRegion (https://discord.com/developers/docs/resources/voice#voice-region-object)
type VoiceRegion struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Vip        bool         `json:"vip"`
	Optimal    bool         `json:"optimal"`
	Deprecated bool         `json:"deprecated"`
	Custom     bool         `json:"custom"`
}
