package api

import (
	"net/http"

	"github.com/DisgoOrg/log"
)

// Options is the configuration used when creating the client
type Options struct {
	Logger                    log.Logger
	GatewayIntents                   GatewayIntents
	RestTimeout               int
	EnableWebhookInteractions bool
	ListenPort                int
	ListenURL                 string
	PublicKey                 string
	LargeThreshold            int
	RawGatewayEventsEnabled   bool
	HTTPClient                *http.Client
}
