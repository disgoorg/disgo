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
	GetGateway    = NewRoute(GET, "gateway")
	GetGatewayBot = NewRoute(GET, "gateway/bot")

	GetUsersMe = NewRoute(GET, "users/@me")
	GetUser    = NewRoute(GET, "users/%s")

	PostUsersMeChannels = NewRoute(POST, "/users/@me/channels")

	PutReaction = NewRoute(PUT, "/channels/%s/messages/%s/reactions/%s/@me")

	GetMember   = NewRoute(GET, "guilds/%s/members/%s")
	PostMessage = NewRoute(POST, "channels/%s/messages")

	GuildIcon CdnRoute = "icons/%s/%s.%s"
)
