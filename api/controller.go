package api

import (
	"adtelligent-internship/api/handler"
	"github.com/valyala/fasthttp"
)

func NewRequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/campaigns_map":
			handler.CampaignHandlerMap(ctx)
		case "/campaigns_slice":
			handler.CampaignHandlerSlice(ctx)
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}
