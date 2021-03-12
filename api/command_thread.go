package api

import (
	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/internal/events"
)

type CommandThread interface {
	Disgo() Disgo
	Event() events.SlashCommandEvent
	Ephemeral(ephemeral bool) CommandThread
	SendMessage() *promise.Promise
	EditOriginal(message string) *promise.Promise
	DeleteOriginal() * promise.Promise
}
