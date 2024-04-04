package data

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func PopulateData(database *sql.DB) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Source%d", i+1)
		_, err := database.Exec("INSERT INTO sources (name) VALUES (?)", name)
		if err != nil {
			log.Fatal("Error inserting source:", err)
		}
	}

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Campaign%d", i+1)
		_, err := database.Exec("INSERT INTO campaigns (name) VALUES (?)", name)
		if err != nil {
			log.Fatal("Error inserting campaign:", err)
		}
	}

	for _, sourceID := range rand.Perm(100)[:50] {
		for _, campaignID := range rand.Perm(100)[:rand.Intn(11)] {
			_, err := database.Exec(
				"INSERT INTO source_campaign (source_id, campaign_id) VALUES (?, ?)", sourceID+1, campaignID+1)
			if err != nil {
				log.Fatal("Error inserting source_campaign link:", err)
			}
		}
	}

	fmt.Println("Data populated successfully.")
}
