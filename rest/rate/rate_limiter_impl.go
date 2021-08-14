package rate

import (
	"sync"

	"github.com/DisgoOrg/disgo/rest/route"
)

//goland:noinspection GoNameStartsWithPackageName
type RateLimiterImpl struct {
	Lock sync.Mutex
	// Route -> Hash
	Hashes map[route.APIRoute]string
	// Hash + Major Parameter -> Bucket
	Buckets map[string]Bucket
}

type Bucket struct {
	sync.Mutex
	ID        string
	Reset     int
	Remaining int
	Limit     int
}
