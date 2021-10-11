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

func (r *Response) SetStatusCode(code int) error {
	r.w.WriteHeader(code)
	return nil
}

func (r *Response) Redirect(url string) error {
	http.Redirect(r.w, r.ctx.Request.Request, url, http.StatusTemporaryRedirect)
	return nil
}

func (r *Response) SetHeader(key, value string) error {
	r.w.Header().Set(key, value)
	return nil
}
