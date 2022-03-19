package disgo

import (
	_ "github.com/DisgoOrg/disgo/core"
	_ "github.com/DisgoOrg/disgo/core/bot"
	_ "github.com/DisgoOrg/disgo/core/events"
	_ "github.com/DisgoOrg/disgo/core/handlers"
	_ "github.com/DisgoOrg/disgo/discord"
	_ "github.com/DisgoOrg/disgo/gateway"
	_ "github.com/DisgoOrg/disgo/gateway/grate"
	_ "github.com/DisgoOrg/disgo/gateway/sharding"
	_ "github.com/DisgoOrg/disgo/gateway/sharding/srate"
	_ "github.com/DisgoOrg/disgo/httpserver"
	_ "github.com/DisgoOrg/disgo/info"
	_ "github.com/DisgoOrg/disgo/internal/insecurerandstr"
	_ "github.com/DisgoOrg/disgo/internal/tokenhelper"
	_ "github.com/DisgoOrg/disgo/json"
	_ "github.com/DisgoOrg/disgo/oauth2"
	_ "github.com/DisgoOrg/disgo/rest"
	_ "github.com/DisgoOrg/disgo/rest/route"
	_ "github.com/DisgoOrg/disgo/rest/rrate"
	_ "github.com/DisgoOrg/disgo/webhook"
)
