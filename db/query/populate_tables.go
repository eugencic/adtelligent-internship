package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func PopulateData(database *sql.DB) error {
	// List of domains to choose from
	domains := []string{
		"gmail.com",
		"yahoo.com",
		"hotmail.com",
		"aol.com",
		"msn.com",
		"mail.ru",
		"yandex.ru",
		"orange.fr",
	}

	// Insert sources
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Source %d", i+1)
		_, err := database.Exec("INSERT INTO sources (name) VALUES (?)", name)
		if err != nil {
			return fmt.Errorf("error inserting source: %w", err)
		}
	}

	// Insert campaigns with random domains
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Campaign %d", i+1)
		filterType := getRandomFilterType()

		// Select 5 random domains
		selectedDomains := getRandomDomains(domains, 5)

		// Convert selected domains to JSON format
		domainsJSON, err := json.Marshal(selectedDomains)
		if err != nil {
			return fmt.Errorf("error marshalling domains to JSON: %w", err)
		}

		// Insert campaign into database with domains JSON
		_, err = database.Exec("INSERT INTO campaigns (name, filter_type, domains) VALUES (?, ?, ?)",
			name, filterType, domainsJSON)
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

	// fmt.Println("Data populated successfully.")
	return nil
}

func getRandomDomains(allDomains []string, numDomains int) map[string]bool {
	rand.Seed(time.Now().UnixNano())
	selectedDomains := make(map[string]bool)

	for i := 0; i < numDomains; i++ {
		randIndex := rand.Intn(len(allDomains))
		selectedDomains[allDomains[randIndex]] = true
	}

	return selectedDomains
}

func getRandomFilterType() string {
	filterTypes := []string{
		"white",
		"black",
	}

	randIndex := rand.Intn(len(filterTypes))
	return filterTypes[randIndex]
}
