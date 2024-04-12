package api

import (
	"database/sql"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

type Campaign struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Filter string `json:"filter"`
}

var (
	cache = make(map[int][]Campaign)
)

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
	rows, err := db.Query("SELECT c.id, c.name, c.domain, c.filter_type, sc.source_id FROM campaigns c INNER JOIN source_campaign sc ON c.id = sc.campaign_id")
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	for rows.Next() {
		var id, sourceID int
		var name, domain, filter string
		if err := rows.Scan(&id, &name, &domain, &filter, &sourceID); err != nil {
			return err
		}

		cache[sourceID] = append(cache[sourceID], Campaign{ID: id, Name: name, Domain: domain, Filter: filter})
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

	domain := strings.ToLower(string(ctx.QueryArgs().Peek("domain")))

	cachedData, ok := cache[sourceID]

	var filteredData []Campaign

	if ok {
		for _, campaign := range cachedData {
			match := (campaign.Filter == "white" && campaign.Domain == domain) ||
				(campaign.Filter == "black" && campaign.Domain != domain)
			if match {
				filteredData = append(filteredData, campaign)
			}
		}

		respondWithJSON(ctx, filteredData)
	} else {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
	}
}
