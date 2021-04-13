package api

import "github.com/DisgoOrg/log"

// Options is the configuration used when creating the client
type Options struct {
	Logger                    log.Logger
	Intents                   Intents
	RestTimeout               int
	EnableWebhookInteractions bool
	ListenPort                int
	ListenURL                 string
	PublicKey                 string
	LargeThreshold            int
	RawGatewayEventsEnabled   bool
}
