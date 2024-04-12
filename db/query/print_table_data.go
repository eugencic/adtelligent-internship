package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func PrintSources(db *sql.DB) error {
	fmt.Println("Sources:")
	rows, err := db.Query("SELECT id, name FROM sources")
	if err != nil {
		return fmt.Errorf("error querying sources: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return fmt.Errorf("error scanning sources row: %w", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
	}

	return nil
}

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool // Map to store domains
}

func PrintCampaigns(db *sql.DB) error {
	fmt.Println("\nCampaigns:")
	rows, err := db.Query("SELECT id, name, filter_type, domains FROM campaigns")
	if err != nil {
		return fmt.Errorf("error querying campaigns: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	for rows.Next() {
		var id int
		var name, filterType, domainsJSON string
		err := rows.Scan(&id, &name, &filterType, &domainsJSON)
		if err != nil {
			return fmt.Errorf("error scanning campaigns row: %w", err)
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
		}

		fmt.Printf("Campaign ID: %d, Name: %s, Filter Type: %s\n", campaign.ID, campaign.Name, campaign.FilterType)

		fmt.Println("Domains:")
		for domain := range campaign.Domains {
			fmt.Printf("- %s\n", domain)
		}

		fmt.Println()
	}

	if err := rows.Err(); err != nil {
	}

	return nil
}

func PrintSourceCampaign(db *sql.DB) error {
	fmt.Println("\nSource_campaign:")
	rows, err := db.Query("SELECT source_id, campaign_id FROM source_campaign")
	if err != nil {
		return fmt.Errorf("error querying source_campaign: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)
	for rows.Next() {
		var sourceID, campaignID int
		err := rows.Scan(&sourceID, &campaignID)
		if err != nil {
			return fmt.Errorf("error scanning source_campaign row: %w", err)
		}
		fmt.Printf("Source ID: %d, Campaign ID: %d\n", sourceID, campaignID)
	}
	if err := rows.Err(); err != nil {
	}

	return nil
}
