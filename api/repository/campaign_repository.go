package repository

import (
	"adtelligent-internship/model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "net/http/pprof"
)

var PreloadedCache = make(map[int][]model.Campaign, 100)
var PreloadedSliceCache = make(map[int][]model.CampaignSlice, 100)

func PreloadDataWithMap(db *sql.DB) error {
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

		campaign := model.Campaign{
			ID:         id,
			Name:       name,
			FilterType: filterType,
			Domains:    domains,
			SourceID:   sourceID,
		}

		PreloadedCache[sourceID] = append(PreloadedCache[sourceID], campaign)
	}

	return nil
}

func PreloadDataWithSlices(db *sql.DB) error {
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

		var domainsMap map[string]bool
		if err := json.Unmarshal([]byte(domainsJSON), &domainsMap); err != nil {
			return fmt.Errorf("error unmarshalling domains JSON: %w", err)
		}

		var domains []string
		for domain := range domainsMap {
			domains = append(domains, domain)
		}

		campaign := model.CampaignSlice{
			ID:         id,
			Name:       name,
			FilterType: filterType,
			Domains:    domains,
			SourceID:   sourceID,
		}

		if _, ok := PreloadedSliceCache[sourceID]; !ok {
			PreloadedSliceCache[sourceID] = make([]model.CampaignSlice, 0)
		}
		PreloadedSliceCache[sourceID] = append(PreloadedSliceCache[sourceID], campaign)
	}

	return nil
}
