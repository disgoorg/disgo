package endpoints

// Discord Endpoint Constants
const (
	APIVersion = "8"
	Base       = "https://discord.com/"
	CDN        = "https://cdn.discordapp.com/"
	API        = Base + "api/v" + APIVersion + "/"
	WS         = "wss://gateway.discord.gg/"
)

// All of the Discord endpoints used by the lib
var (
	Gateway    = NewRoute(GET, "gateway")
	GatewayBot = NewRoute(GET, "gateway/bot")

	UsersMe = NewRoute(GET, "users/@me")
	User    = NewRoute(GET, "users/%s")

	GuildIcon CdnRoute = "icons/%s/%s.%s"
)
