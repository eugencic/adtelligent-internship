package handler

import (
	"adtelligent-internship/api/repository"
	"adtelligent-internship/api/util"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool
	SourceID   int
}

type CampaignWithPrice struct {
	Campaign repository.Campaign
	Price    int
}

func (c Campaign) Call() CampaignWithPrice {
	price := rand.Intn(100) + 1
	fmt.Printf("Campaign %d generated price: %d\n", c.ID, price)

	delay := rand.Intn(5)
	fmt.Printf("Campaign %d simulating response delay: %d seconds\n", c.ID, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	return CampaignWithPrice{Campaign: repository.Campaign(c), Price: price}
}

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

	f, err := os.Create("profiles/mem_handler.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()

	var filteredData []repository.Campaign
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

	campaignsWithPrices := make(chan CampaignWithPrice, len(filteredData))

	for _, campaign := range filteredData {
		go func(c Campaign) {
			campaignsWithPrices <- c.Call()
		}(Campaign(campaign))
	}

	var campaignsWithPricesSlice []CampaignWithPrice
	for i := 0; i < len(filteredData); i++ {
		cwp := <-campaignsWithPrices
		campaignsWithPricesSlice = append(campaignsWithPricesSlice, cwp)
	}

	close(campaignsWithPrices)

	sort.Slice(campaignsWithPricesSlice, func(i, j int) bool {
		return campaignsWithPricesSlice[i].Price < campaignsWithPricesSlice[j].Price
	})

	var sortedCampaigns []repository.Campaign
	for _, cwp := range campaignsWithPricesSlice {
		sortedCampaigns = append(sortedCampaigns, cwp.Campaign)
	}

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("Could not write memory profile: ", err)
	}

	util.RespondWithJSON(ctx, sortedCampaigns)
}
