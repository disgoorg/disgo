package events

import (
	"github.com/disgoorg/disgo/discord"
)

type ApplicationAuthorized struct {
	*Event
	*WebhookEvent
	IntegrationType discord.IntegrationType
	User            discord.User
	Scopes          []discord.OAuth2Scope
	Guild           discord.Guild
}

type ApplicationDeauthorized struct {
	*Event
	*WebhookEvent
	User discord.User
}
