//package api
//
//import (
//	"database/sql"
//	"encoding/json"
//	"github.com/valyala/fasthttp"
//	"log"
//	"strconv"
//	"strings"
//)
//
//type Campaign struct {
//	ID   int    `json:"id"`
//	Name string `json:"name"`
//}
//
//var (
//	cache = make(map[int][]Campaign)
//)
//
//func respondWithJSON(ctx *fasthttp.RequestCtx, data interface{}) {
//	response, err := json.Marshal(data)
//	if err != nil {
//		log.Printf("Error marshalling data to JSON: %v", err)
//		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
//		return
//	}
//
//	ctx.SetContentType("application/json")
//	ctx.SetStatusCode(fasthttp.StatusOK)
//	ctx.SetBody(response)
//}
//
//func PreloadData(db *sql.DB) error {
//	rows, err := db.Query("SELECT c.id, c.name, sc.source_id FROM campaigns c INNER JOIN source_campaign sc ON c.id = sc.campaign_id")
//	if err != nil {
//		return err
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//		}
//	}(rows)
//
//	for rows.Next() {
//		var id, sourceID int
//		var name, domain, filter string
//		if err := rows.Scan(&id, &name, &domain, &filter, &sourceID); err != nil {
//			return err
//		}
//
//		cache[sourceID] = append(cache[sourceID], Campaign{ID: id, Name: name, Domain: domain, Filter: filter})
//	}
//
//	return nil
//}
//
//func CachedCampaignHandler(ctx *fasthttp.RequestCtx) {
//	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
//	sourceID, err := strconv.Atoi(sourceIDStr)
//	if err != nil {
//		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
//		return
//	}
//
//	domain := strings.ToLower(string(ctx.QueryArgs().Peek("domain")))
//
//	cachedData, ok := cache[sourceID]
//
//	var filteredData []Campaign
//
//	if ok {
//		for _, campaign := range cachedData {
//			match := (campaign.Filter == "white" && campaign.Domain == domain) ||
//				(campaign.Filter == "black" && campaign.Domain != domain)
//			if match {
//				filteredData = append(filteredData, campaign)
//			}
//		}
//
//		respondWithJSON(ctx, filteredData)
//	} else {
//		ctx.Error("Data not found", fasthttp.StatusNotFound)
//	}
//}

package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool
	SourceID   int
}

var cache = make(map[int][]Campaign)

func respondWithJSON(ctx *fasthttp.RequestCtx, data interface{}) {
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

func PreloadData(db *sql.DB) error {

	rows, err := db.Query("SELECT c.id, c.name, c.filter_type, c.domains, sc.source_id FROM campaigns c INNER JOIN source_campaign sc ON c.id = sc.campaign_id")
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}(rows)

	for rows.Next() {
		var id, sourceID int
		var name string
		var filterType string
		var domainsJSON string

		if err := rows.Scan(&id, &name, &filterType, &domainsJSON, &sourceID); err != nil {
			return err
		}

		// Unmarshal domains JSON into a map[string]bool
		var domains map[string]bool
		if err := json.Unmarshal([]byte(domainsJSON), &domains); err != nil {
			return fmt.Errorf("error unmarshalling domains JSON: %w", err)
		}

		// Create Campaign object with parsed data
		campaign := Campaign{
			ID:         id,
			Name:       name,
			FilterType: filterType,
			Domains:    domains,
			SourceID:   sourceID,
		}

		// Append campaign to cache map
		cache[sourceID] = append(cache[sourceID], campaign)
	}

	return nil
}

func CachedCampaignHandler(ctx *fasthttp.RequestCtx) {
	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
		return
	}

	requestedDomain := strings.ToLower(string(ctx.QueryArgs().Peek("domain")))

	cachedData, ok := cache[sourceID]

	var filteredData []Campaign

	if ok {
		for _, campaign := range cachedData {
			switch campaign.FilterType {
			case "black":
				// Exclude campaign if it's blacklisted and contains the requested domain
				if !campaign.Domains[requestedDomain] {
					filteredData = append(filteredData, campaign)
				}
			case "white":
				// Include campaign if it's whitelisted and contains the requested domain
				if campaign.Domains[requestedDomain] {
					filteredData = append(filteredData, campaign)
				}
			}
		}

		respondWithJSON(ctx, filteredData)
	} else {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
	}
}
