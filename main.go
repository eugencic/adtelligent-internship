package main

import (
	"adtelligent-internship/api"
	"adtelligent-internship/db"
	"adtelligent-internship/db/data"
	"adtelligent-internship/db/printer"
	"adtelligent-internship/db/schema"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
		}
	}(database)

	err = schema.CreateTables(database)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	err = data.PopulateData(database)
	if err != nil {
		log.Fatalf("Failed to populate data: %v", err)
	}

	printer.PrintSources(database)
	printer.PrintCampaigns(database)
	printer.PrintSourceCampaign(database)

	err = api.PreloadData(database)
	if err != nil {
		log.Fatalf("Failed to preload data: %v", err)
	}

	requestHandler := api.NewRequestHandler(database)

	log.Println("Starting HTTP server on port 8080...")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
