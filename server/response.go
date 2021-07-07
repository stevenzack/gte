package server

import "net/http"

type Response struct {
	ctx *Context
	w   http.ResponseWriter
}

func NewResponse(ctx *Context, w http.ResponseWriter) *Response {
	return &Response{
		ctx: ctx,
		w:   w,
	}
}
