package handler

import (
	"adtelligent-internship/api/repository"
	"adtelligent-internship/api/util"
	"adtelligent-internship/model"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

func CampaignHandlerSlice(ctx *fasthttp.RequestCtx) {
	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
		return
	}

	requestedDomain := util.ExtractBaseDomain(strings.ToLower(string(ctx.QueryArgs().Peek("domain"))))

	cachedData, ok := repository.PreloadedSliceCache[sourceID]
	if !ok {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
		return
	}

	f, err := os.Create("profiles/mem_campaign_handler_slice.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()

	var filteredData []model.CampaignSlice
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

	campaignsWithPrices := make(chan model.CampaignSliceWithPrice, len(filteredData))

	for _, campaign := range filteredData {
		go func(c model.CampaignSlice) {
			campaignsWithPrices <- c.Call()
		}(campaign)
	}

	var campaignsWithPricesSlice []model.CampaignSliceWithPrice
	for i := 0; i < len(filteredData); i++ {
		cwp := <-campaignsWithPrices
		campaignsWithPricesSlice = append(campaignsWithPricesSlice, cwp)
	}

	close(campaignsWithPrices)

	sort.Slice(campaignsWithPricesSlice, func(i, j int) bool {
		return campaignsWithPricesSlice[i].Price < campaignsWithPricesSlice[j].Price
	})

	var sortedCampaigns []model.CampaignSlice
	for _, cwp := range campaignsWithPricesSlice {
		sortedCampaigns = append(sortedCampaigns, cwp.CampaignSlice)
	}

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("Could not write memory profile: ", err)
	}

	util.RespondWithJSON(ctx, filteredData)

}

func containsDomain(domains []string, domain string) bool {
	for _, d := range domains {
		if d == domain {
			return true
		}
	}
	return false
}
