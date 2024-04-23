package main

import (
	"adtelligent-internship/api"
	"adtelligent-internship/api/repository"
	"adtelligent-internship/db"
	"adtelligent-internship/db/query"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/profile"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func main() {
	defer profile.Start(profile.ProfilePath("./profiles")).Stop()

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

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run server.go <ip_address> <port>")
		return
	}

	ipAddress := os.Args[1]
	port := os.Args[2]
	addr := fmt.Sprintf("%s:%s", ipAddress, port)

	log.Println("Starting HTTP server on port 8080...")
	if err := fasthttp.ListenAndServe(addr, requestHandler); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
