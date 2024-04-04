package printer

import (
	"database/sql"
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

func PrintCampaigns(db *sql.DB) {
	fmt.Println("\nCampaigns:")
	rows, err := db.Query("SELECT id, name FROM campaigns")
	if err != nil {
		log.Fatal("Error querying campaigns:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning campaigns row:", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
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
