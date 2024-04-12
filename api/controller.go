package api

import (
	"adtelligent-internship/api/handler"
	"github.com/valyala/fasthttp"
)

func NewRequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/campaigns":
			handler.CampaignHandler(ctx)
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}
