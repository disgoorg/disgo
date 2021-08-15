package rate

import (
	"net/http"
	"sync"
	"time"

	"github.com/DisgoOrg/disgo/rest/route"
)

const UnlimitedBucket = "unlimited"

//goland:noinspection GoNameStartsWithPackageName
type RateLimiterImpl struct {
	Lock sync.Mutex
	// Route -> Hash
	Hashes map[*route.APIRoute]string
	// Hash + Major Parameter -> Bucket
	Buckets map[string]Bucket
}

func (r *RateLimiterImpl) Close(force bool) {

}

func (r *RateLimiterImpl) GetRateLimit(route *route.CompiledAPIRoute) time.Duration {

}

func (r *RateLimiterImpl) GetRouteHash(route *route.APIRoute) string {
	hash, ok := r.Hashes[route]
	if !ok {
		hash = UnlimitedBucket + "+" + route.Method().String() + "/" + route.Route.Route()
		r.Hashes[route] = hash
	}
	return hash
}

func (r *RateLimiterImpl) GetBucket(route *route.CompiledAPIRoute) *Bucket {
	bucket :=
	return nil
}

func NewBucket(bucketID string) *Bucket {
	return &Bucket{
		ID:        bucketID,
		Remaining: 1,
		Limit:     1,
	}
}

type Bucket struct {
	sync.Mutex
	RateLimiter RateLimiter
	ID          string
	Reset       int
	Remaining   int
	Limit       int
}

func (b *Bucket) Free(response *http.Response) {
	remaining := response.Header.Get("X-RateLimit-Remaining")
	reset := response.Header.Get("X-RateLimit-Reset")
	global := response.Header.Get("X-RateLimit-Global")
	resetAfter := response.Header.Get("X-RateLimit-Reset-After")
}
