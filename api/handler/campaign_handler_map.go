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

func CampaignHandlerMap(ctx *fasthttp.RequestCtx) {
	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
		return
	}

	requestedDomain := util.ExtractBaseDomain(strings.ToLower(string(ctx.QueryArgs().Peek("domain"))))

	cachedData, ok := repository.PreloadedCache[sourceID]
	if !ok {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
		return
	}

	f, err := os.Create("profiles/mem_campaign_handler_map.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()

	var filteredData []model.Campaign
	for _, campaign := range cachedData {
		switch campaign.FilterType {
		case "black":
			if !campaign.Domains[requestedDomain] {
				filteredData = append(filteredData, campaign)
			}
		case "white":
			if campaign.Domains[requestedDomain] {
				filteredData = append(filteredData, campaign)
			}
		}
	}

	campaignsWithPrices := make(chan model.CampaignWithPrice, len(filteredData))

	for _, campaign := range filteredData {
		go func(c model.Campaign) {
			campaignsWithPrices <- c.Call()
		}(campaign)
	}

	var campaignsWithPricesSlice []model.CampaignWithPrice
	for i := 0; i < len(filteredData); i++ {
		cwp := <-campaignsWithPrices
		campaignsWithPricesSlice = append(campaignsWithPricesSlice, cwp)
	}

	close(campaignsWithPrices)

	sort.Slice(campaignsWithPricesSlice, func(i, j int) bool {
		return campaignsWithPricesSlice[i].Price < campaignsWithPricesSlice[j].Price
	})

	var sortedCampaigns []model.Campaign
	for _, cwp := range campaignsWithPricesSlice {
		sortedCampaigns = append(sortedCampaigns, cwp.Campaign)
	}

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("Could not write memory profile: ", err)
	}

	util.RespondWithJSON(ctx, sortedCampaigns)
}
