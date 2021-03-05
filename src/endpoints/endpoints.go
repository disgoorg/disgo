package endpoints

// All of the Discord endpoints used by the lib
var (
	Gateway    = NewRoute(GET, "gateway")
	GatewayBot = NewRoute(GET, "gateway/bot")

	UsersMe = NewRoute(GET, "users/@me")
	User    = NewRoute(GET, "users/%s")
)
