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
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("profiles/mem_main.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()

	runtime.GC()

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

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

	err = repository.PreloadDataWithMap(database)
	if err != nil {
		log.Fatalf("Failed to preload data: %v", err)
	}

	err = repository.PreloadDataWithSlices(database)
	if err != nil {
		log.Fatalf("Failed to preload data: %v", err)
	}

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	fmt.Println("Data is set.")

	requestHandler := api.NewRequestHandler()

	log.Println("Starting HTTP server on port 8080...")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}

	//defer f.Close()
}
