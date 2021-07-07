package server

import "net/http"

type Request struct {
	ctx *Context
	*http.Request
}

func NewRequest(ctx *Context, r *http.Request) *Request {
	return &Request{
		ctx:     ctx,
		Request: r,
	}
}

func (r *Request) GetPathParam(k string) string {
	m := r.ctx.route.Params(r.URL.Path)
	return m[k]
}

