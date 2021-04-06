package api

// Options is the configuration used when creating the client
type Options struct {
	Intents                   Intents
	RestTimeout               int
	EnableWebhookInteractions bool
	ListenPort                int
	ListenURL                 string
	PublicKey                 string
	LargeThreshold            int
}
