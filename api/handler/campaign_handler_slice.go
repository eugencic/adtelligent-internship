package handler

import (
	"adtelligent-internship/api/repository"
	"adtelligent-internship/api/util"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type CampaignSlice struct {
	ID         int
	Name       string
	FilterType string
	Domains    []string
	SourceID   int
}

func CampaignHandlerSlice(ctx *fasthttp.RequestCtx) {
	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
		return
	}

	requestedDomain := util.ExtractBaseDomain(strings.ToLower(string(ctx.QueryArgs().Peek("domain"))))

	cachedData, ok := repository.PreloadedSliceCache[sourceID]
	var filteredData []repository.CampaignSlice

	if ok {
		for _, campaign := range cachedData {
			switch campaign.FilterType {
			case "black":
				if !containsDomain(campaign.Domains, requestedDomain) {
					filteredData = append(filteredData, campaign)
				}
			case "white":
				if containsDomain(campaign.Domains, requestedDomain) {
					filteredData = append(filteredData, campaign)
				}
			}
		}

		util.RespondWithJSON(ctx, filteredData)
	} else {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
	}
}

func containsDomain(domains []string, domain string) bool {
	for _, d := range domains {
		if d == domain {
			return true
		}
	}
	return false
}
