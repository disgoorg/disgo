package endpoints


var (
	GatewayRoute = NewRoute(GET, "gateway")
	GatewayBotRoute = NewRoute(GET, "gateway/bot")
)
