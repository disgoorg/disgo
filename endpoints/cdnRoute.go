package endpoints

import "fmt"

// CdnRoute is a route for interacting with images hosted on discord's CDN
type CdnRoute string

// Compile builds the request URL
func (r CdnRoute) Compile(args ...interface{}) string {
	return fmt.Sprintf(CDN+string(r), args...)
}
