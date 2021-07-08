package server

import (
	"net/http"

	"github.com/StevenZack/gte/config"
)

type Context struct {
	Config   config.Config
	route    config.Route
	Request  *Request
	Response *Response
}

func NewContext(cfg config.Config, route config.Route, w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Config: cfg,
		route:  route,
	}
	ctx.Request = NewRequest(ctx, r)
	ctx.Response = NewResponse(ctx, w)
	return ctx
}
