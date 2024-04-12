package utils

import (
	"adtelligent-internship/db/query"
	"database/sql"
)

func PrintData(database *sql.DB) {
	err := queries.PrintSources(database)
	if err != nil {
		return
	}

	err = queries.PrintCampaigns(database)
	if err != nil {
		return
	}

	err = queries.PrintSourceCampaign(database)
	if err != nil {
		return
	}

	err = queries.ExecuteSourceCampaignQuery(database)
	if err != nil {
		return
	}

	err = queries.ExecuteCampaignWithoutSourceQuery(database)
	if err != nil {
		return
	}

	err = queries.ExecuteUniqueNamesQuery(database)
	if err != nil {
		return
	}
}
