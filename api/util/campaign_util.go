package util

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

func ExtractBaseDomain(fullDomain string) string {
	fullDomain = strings.ToLower(fullDomain)

	parts := strings.Split(fullDomain, ".")

	numParts := len(parts)

	var baseDomain string
	if numParts > 2 {
		baseDomain = strings.Join(parts[numParts-2:], ".")
	} else {
		baseDomain = fullDomain
	}

	return baseDomain
}

func RespondWithJSON(ctx *fasthttp.RequestCtx, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data to JSON: %v", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)
}
