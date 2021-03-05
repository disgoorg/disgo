package endpoints

import "fmt"

type CdnRoute string

func (r CdnRoute) Compile(args ...interface{}) string {
	return fmt.Sprintf(CDN+string(r), args...)
}
