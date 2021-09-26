package disgo

import (
	_ "github.com/DisgoOrg/disgo/core"
	_ "github.com/DisgoOrg/disgo/discord"
	_ "github.com/DisgoOrg/disgo/gateway"
	_ "github.com/DisgoOrg/disgo/httpserver"
	_ "github.com/DisgoOrg/disgo/info"
	_ "github.com/DisgoOrg/disgo/json"
	_ "github.com/DisgoOrg/disgo/oauth2"
	_ "github.com/DisgoOrg/disgo/rest"
	_ "github.com/DisgoOrg/disgo/rest/rate"
	_ "github.com/DisgoOrg/disgo/rest/route"
	_ "github.com/DisgoOrg/disgo/sharding"
	_ "github.com/DisgoOrg/disgo/sharding/rate"
	_ "github.com/DisgoOrg/disgo/webhook"
)
