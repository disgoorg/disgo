package rest

import (
	"context"
	"net/http"

	"github.com/DisgoOrg/disgo/rest/route"
)

func NewRequest(ctx context.Context, route *route.CompiledAPIRoute, rqBody interface{}, rsBody interface{}, headers http.Header) *Request {
	return &Request{
		Ctx:     ctx,
		Route:   route,
		RQBody:  rqBody,
		RSBody:  rsBody,
		Headers: headers,
	}
}

type Request struct {
	Ctx     context.Context
	Route   *route.CompiledAPIRoute
	RQBody  interface{}
	RSBody  interface{}
	Headers http.Header
}
