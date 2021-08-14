package rate

import (
	"net/http"
	"time"

	"github.com/DisgoOrg/disgo/rest/route"
)

var MajorParameter = []string{"channel_id", "guild_id", "webhook_id/webhook_token"}

//goland:noinspection GoNameStartsWithPackageName
type RateLimiter interface {
	Close(force bool)
	IsRateLimited(route route.CompiledRoute) bool
	GetRateLimit(route route.CompiledRoute) time.Duration
	HandleResponse(response http.Response)
}
