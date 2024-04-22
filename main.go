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
	"math/rand"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	runtime.GC()

	rand.Seed(time.Now().UnixNano())

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

	fmt.Println("Data is set.")

	requestHandler := api.NewRequestHandler()

	//go simulateRequests()

	log.Println("Starting HTTP server on port 8080...")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

//func simulateRequests() {
//	f, err := os.Create("profiles/mem_requests.pprof")
//	if err != nil {
//		log.Fatal("could not create memory profile: ", err)
//	}
//	defer f.Close()
//
//	url := "http://localhost:8080/campaigns_map?source_id=1&domain=gmail.com" // Adjust the URL based on your endpoint
//
//	client := &http.Client{}
//
//	numRequests := 10
//
//	for i := 0; i < numRequests; i++ {
//		req, err := http.NewRequest("GET", url, nil)
//		if err != nil {
//			log.Fatalf("Failed to create HTTP request: %v", err)
//		}
//
//		resp, err := client.Do(req)
//		if err != nil {
//			log.Fatalf("Failed to send HTTP request: %v", err)
//		}
//		defer resp.Body.Close()
//
//		fmt.Println("response Status:", resp.Status)
//	}
//
//	log.Printf("Simulated %d requests successfully", numRequests)
//
//	if err := pprof.WriteHeapProfile(f); err != nil {
//		log.Fatal("Could not write memory profile: ", err)
//	}
//}
