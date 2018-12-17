package web

import "github.com/valyala/fasthttp"

func (w *Web) healthcheck() func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
