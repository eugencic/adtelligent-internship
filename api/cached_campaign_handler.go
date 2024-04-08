package api

import (
	"database/sql"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"sync"
)

type Campaign struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	cache      = make(map[int][]Campaign)
	cacheMutex sync.RWMutex
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
	rows, err := db.Query("SELECT c.id, c.name, sc.source_id FROM campaigns c INNER JOIN source_campaign sc ON c.id = sc.campaign_id")
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
		var name string
		if err := rows.Scan(&id, &name, &sourceID); err != nil {
			return err
		}

		cacheMutex.Lock()
		cache[sourceID] = append(cache[sourceID], Campaign{ID: id, Name: name})
		cacheMutex.Unlock()
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
	cacheMutex.RLock()
	cachedData, ok := cache[sourceID]
	cacheMutex.RUnlock()

	if ok {
		respondWithJSON(ctx, cachedData)
	} else {
		ctx.Error("Data not found", fasthttp.StatusNotFound)
	}
}
