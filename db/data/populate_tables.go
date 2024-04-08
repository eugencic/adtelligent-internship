package data

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func PopulateData(database *sql.DB) error {
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Source %d", i+1)
		_, err := database.Exec("INSERT INTO sources (name) VALUES (?)", name)
		if err != nil {
			return fmt.Errorf("error inserting source: %w", err)
		}
	}

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Campaign %d", i+1)
		_, err := database.Exec("INSERT INTO campaigns (name) VALUES (?)", name)
		if err != nil {
			return fmt.Errorf("error inserting campaign: %w", err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	// Simulate source and campaign associations
	for sourceID := 1; sourceID <= 100; sourceID++ {
		// Determine the number of campaigns for this source using a normal distribution
		numCampaigns := int(rand.NormFloat64()*2 + 5) // Mean of 5 campaigns with some variance

		// Ensure numCampaigns is within a reasonable range (1 to 10)
		if numCampaigns < 1 {
			numCampaigns = 1
		} else if numCampaigns > 10 {
			numCampaigns = 10
		}

		// Randomly select campaign IDs for this source
		campaignIDs := rand.Perm(100)[:numCampaigns]

		// Insert associations into source_campaign table
		for _, campaignID := range campaignIDs {
			// Ensure campaignID is within valid range (1 to 100)
			campaignID++ // Increment to match 1-based IDs in the database

			_, err := database.Exec(
				"INSERT INTO source_campaign (source_id, campaign_id) VALUES (?, ?)", sourceID, campaignID)
			if err != nil {
				return fmt.Errorf("error inserting source_campaign link: %w", err)
			}
		}
	}

	fmt.Println("Data populated successfully.")
	return nil
}
