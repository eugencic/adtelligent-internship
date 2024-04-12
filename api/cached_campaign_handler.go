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

var preloadedCache = make(map[int][]Campaign)

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

// extractBaseDomain extracts the base domain (root domain) from a full domain name
func extractBaseDomain(fullDomain string) string {
	// Convert domain to lowercase for case-insensitive comparison
	fullDomain = strings.ToLower(fullDomain)

	// Split domain into parts based on dot (.)
	parts := strings.Split(fullDomain, ".")

	// Determine the number of parts in the domain
	numParts := len(parts)

	// Extract base domain based on the number of parts
	var baseDomain string
	if numParts > 2 { // More than two parts (e.g., subdomain.domain.com)
		baseDomain = strings.Join(parts[numParts-2:], ".")
	} else { // Two parts or less (e.g., domain.com or subdomain.domain.com)
		baseDomain = fullDomain
	}

	return baseDomain
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

		var domains map[string]bool
		if err := json.Unmarshal([]byte(domainsJSON), &domains); err != nil {
			return fmt.Errorf("error unmarshalling domains JSON: %w", err)
		}

		campaign := Campaign{
			ID:         id,
			Name:       name,
			FilterType: filterType,
			Domains:    domains,
			SourceID:   sourceID,
		}

		preloadedCache[sourceID] = append(preloadedCache[sourceID], campaign)
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

	requestedDomain := extractBaseDomain(strings.ToLower(string(ctx.QueryArgs().Peek("domain"))))

	cachedData, ok := preloadedCache[sourceID]

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
