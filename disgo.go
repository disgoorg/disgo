package disgo

import (
	_ "github.com/DisgoOrg/disgo/bot"
	_ "github.com/DisgoOrg/disgo/collectors"
	_ "github.com/DisgoOrg/disgo/core"
	_ "github.com/DisgoOrg/disgo/discord"
	_ "github.com/DisgoOrg/disgo/events"
	_ "github.com/DisgoOrg/disgo/gateway"
	_ "github.com/DisgoOrg/disgo/gateway/gatewayhandlers"
	_ "github.com/DisgoOrg/disgo/gateway/grate"
	_ "github.com/DisgoOrg/disgo/httpserver"
	_ "github.com/DisgoOrg/disgo/httpserver/httpserverhandlers"
	_ "github.com/DisgoOrg/disgo/info"
	_ "github.com/DisgoOrg/disgo/internal/insecurerandstr"
	_ "github.com/DisgoOrg/disgo/json"
	_ "github.com/DisgoOrg/disgo/oauth2"
	_ "github.com/DisgoOrg/disgo/rest"
	_ "github.com/DisgoOrg/disgo/rest/route"
	_ "github.com/DisgoOrg/disgo/rest/rrate"
	_ "github.com/DisgoOrg/disgo/sharding"
	_ "github.com/DisgoOrg/disgo/sharding/srate"
	_ "github.com/DisgoOrg/disgo/webhook"
)
