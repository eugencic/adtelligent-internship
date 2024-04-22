package handler

import (
	"adtelligent-internship/api/repository"
	"adtelligent-internship/api/util"
	"fmt"
	"github.com/valyala/fasthttp"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool
	SourceID   int
}

func CampaignHandlerMap(ctx *fasthttp.RequestCtx) {
	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
		return
	}

	requestedDomain := util.ExtractBaseDomain(strings.ToLower(string(ctx.QueryArgs().Peek("domain"))))
	//fmt.Println(repository.PreloadedCache)

	cachedData, ok := repository.PreloadedCache[sourceID]
	//fmt.Println(repository.PreloadedCache[sourceID])
	var filteredData []repository.Campaign

	if ok {
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

		//for _ = range cachedData {
		//	r := rand.Intn(5)
		//	var randTime = time.Duration(r) * time.Second
		//	fmt.Println(randTime)
		//	time.Sleep(randTime)
		//}

		var waitGroup sync.WaitGroup
		waitGroup.Add(len(filteredData))
		for _ = range filteredData {
			go func() {
				r := rand.Intn(5)
				var randTime = time.Duration(r) * time.Second
				fmt.Println(randTime)
				time.Sleep(randTime)
				waitGroup.Done()
			}()
		}

		waitGroup.Wait()

		util.RespondWithJSON(ctx, filteredData)
	} else {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
	}
}
