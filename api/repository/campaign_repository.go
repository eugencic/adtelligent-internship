package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "net/http/pprof"
)

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool
	SourceID   int
}

var PreloadedCache = make(map[int][]Campaign, 100)

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

		//decoder := json.NewDecoder(strings.NewReader(domainsJSON))
		//if err := decoder.Decode(&domains); err != nil {
		//	return fmt.Errorf("error decoding domains JSON: %w", err)
		//}

		campaign := Campaign{
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

type CampaignSlice struct {
	ID         int
	Name       string
	FilterType string
	Domains    []string
	SourceID   int
}

var PreloadedSliceCache = make(map[int][]CampaignSlice)

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

		//decoder := json.NewDecoder(strings.NewReader(domainsJSON))
		//if err := decoder.Decode(&domainsMap); err != nil {
		//	return fmt.Errorf("error decoding domains JSON: %w", err)
		//}

		var domains []string
		for domain := range domainsMap {
			domains = append(domains, domain)
		}

		campaign := CampaignSlice{
			ID:         id,
			Name:       name,
			FilterType: filterType,
			Domains:    domains,
			SourceID:   sourceID,
		}

		if _, ok := PreloadedSliceCache[sourceID]; !ok {
			PreloadedSliceCache[sourceID] = make([]CampaignSlice, 0)
		}
		PreloadedSliceCache[sourceID] = append(PreloadedSliceCache[sourceID], campaign)
	}

	return nil
}
