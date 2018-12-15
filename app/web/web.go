package web

import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/app/datastore"
	"github.com/ngalayko/highloadcup/app/logger"
)

// Web is a web server.
type Web struct {
	log       *logger.Logger
	datastore *datastore.Datastore
}

// New is a web constructor.
func New(
	log *logger.Logger,
	datastore *datastore.Datastore,
) *Web {
	return &Web{
		log:       log,
		datastore: datastore,
	}
}

// ListenAndServe starts the server.
func (w *Web) ListenAndServe(addr string) error {
	w.log.Info("starting server on %s", addr)

	s := &fasthttp.Server{
		Handler: w.handler,
	}
	return s.ListenAndServe(addr)
}

func (w *Web) handler(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	switch string(ctx.Method()) {
	case "GET":
		w.handlerGET(ctx)
	case "POST":
		w.handlerPOST(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
	w.log.Info("%s %s", ctx.URI(), time.Since(start))
}

func (w *Web) handlerPOST(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func (w *Web) handlerGET(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/accounts/filter/":
		w.accountsFilter()(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func (w *Web) error(ctx *fasthttp.RequestCtx, err error) {
	w.log.Error("%s: %s", ctx.Path(), err)

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func (w *Web) responseJSON(ctx *fasthttp.RequestCtx, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.error(ctx, err)
		return
	}

	ctx.Response.Header.Add("Connection", "keep-alive")
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetContentLength(ctx.Response.Header.ContentLength())

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(jsonData)
}
