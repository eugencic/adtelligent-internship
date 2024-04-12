package main

import (
	"adtelligent-internship/api"
	"adtelligent-internship/api/repository"
	"adtelligent-internship/db"
	"adtelligent-internship/db/query"
	"database/sql"
	"fmt"
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

	err = queries.CreateTables(database)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	err = queries.PopulateData(database)
	if err != nil {
		log.Fatalf("Failed to populate data: %v", err)
	}

	// util.PrintData(database)

	err = repository.PreloadData(database)
	if err != nil {
		log.Fatalf("Failed to preload data: %v", err)
	}

	fmt.Println("Data is set.")

	requestHandler := api.NewRequestHandler()

	log.Println("Starting HTTP server on port 8080...")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
