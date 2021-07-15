package server

import "net/http"

type Request struct {
	*http.Request
	ctx *Context
}

func NewRequest(ctx *Context, r *http.Request) *Request {
	return &Request{
		Request: r,
		ctx:     ctx,
	}
}

func (r *Request) GetParam(k string) string {
	m := r.ctx.route.Params(r.URL.Path)
	return m[k]
}
