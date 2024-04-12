package printer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func PrintSources(db *sql.DB) {
	fmt.Println("Sources:")
	rows, err := db.Query("SELECT id, name FROM sources")
	if err != nil {
		log.Fatal("Error querying sources:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning sources row:", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

type Campaign struct {
	ID         int
	Name       string
	FilterType string
	Domains    map[string]bool // Map to store domains
}

func PrintCampaigns(db *sql.DB) {
	fmt.Println("\nCampaigns:")
	rows, err := db.Query("SELECT id, name, filter_type, domains FROM campaigns")
	if err != nil {
		log.Fatal("Error querying campaigns:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, filterType, domainsJSON string
		err := rows.Scan(&id, &name, &filterType, &domainsJSON)
		if err != nil {
			log.Fatal("Error scanning campaigns row:", err)
		}

		// Unmarshal domains JSON into a map[string]bool
		var domains map[string]bool
		if err := json.Unmarshal([]byte(domainsJSON), &domains); err != nil {
			log.Fatal("Error unmarshalling domains JSON:", err)
		}

		// Create Campaign object with parsed data
		campaign := Campaign{
			ID:         id,
			Name:       name,
			FilterType: filterType,
			Domains:    domains,
		}

		// Print campaign details
		fmt.Printf("Campaign ID: %d, Name: %s, Filter Type: %s\n", campaign.ID, campaign.Name, campaign.FilterType)

		// Print domains associated with the campaign
		fmt.Println("Domains:")
		for domain := range campaign.Domains {
			fmt.Printf("- %s\n", domain)
		}

		fmt.Println() // Print newline for readability
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func PrintSourceCampaign(db *sql.DB) {
	fmt.Println("\nSource_campaign:")
	rows, err := db.Query("SELECT source_id, campaign_id FROM source_campaign")
	if err != nil {
		log.Fatal("Error querying source_campaign:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sourceID, campaignID int
		err := rows.Scan(&sourceID, &campaignID)
		if err != nil {
			log.Fatal("Error scanning source_campaign row:", err)
		}
		fmt.Printf("Source ID: %d, Campaign ID: %d\n", sourceID, campaignID)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
