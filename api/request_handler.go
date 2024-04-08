package api

import (
	"database/sql"
	"github.com/valyala/fasthttp"
)

func NewRequestHandler(db *sql.DB) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/cached-campaigns":
			CachedCampaignHandler(ctx)
		case "/campaigns":

			CampaignHandler(ctx, db)
		default:

			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}
