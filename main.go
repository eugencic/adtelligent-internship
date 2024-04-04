package main

import (
	"adtelligent-internship/db"
	"adtelligent-internship/db/data"
	"adtelligent-internship/db/printer"
	"adtelligent-internship/db/schema"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database, err := db.ConnectToDB()

	err = schema.CreateTables(database)
	if err != nil {
		return
	}

	data.PopulateData(database)

	printer.PrintSources(database)
	printer.PrintCampaigns(database)
	printer.PrintSourceCampaign(database)

	printer.PrintResults(database)
}
