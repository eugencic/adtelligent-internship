package main

import (
	"adtelligent-internship/db"
	"adtelligent-internship/db/data"
	"adtelligent-internship/db/printer"
	"adtelligent-internship/db/schema"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

func CampaignHandler(ctx *fasthttp.RequestCtx, db *sql.DB) {
	sourceIDStr := string(ctx.QueryArgs().Peek("source_id"))
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		log.Printf("invalid source_id: %s", sourceIDStr)
		ctx.Error("Invalid source_id", fasthttp.StatusBadRequest)
		return
	}

	query := "SELECT c.id, c.name FROM campaigns c INNER JOIN source_campaign sc ON c.id = sc.campaign_id WHERE sc.source_id = ?"
	rows, err := db.Query(query, sourceID)
	if err != nil {
		log.Printf("error querying campaigns: %v", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	var campaigns []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Printf("Error scanning campaign row: %v", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			return
		}
		campaigns = append(campaigns, map[string]interface{}{"id": id, "name": name})
	}

	response, err := json.Marshal(campaigns)
	if err != nil {
		log.Printf("Error marshalling campaigns to JSON: %v", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)
}

func main() {
	database, err := db.ConnectToDB()

	err = schema.CreateTables(database)
	if err != nil {
		return
	}

	err = data.PopulateData(database)
	if err != nil {
		return
	}

	printer.PrintSources(database)
	printer.PrintCampaigns(database)
	printer.PrintSourceCampaign(database)

	err = printer.PrintResults(database)
	if err != nil {
		return
	}

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/campaigns":
			CampaignHandler(ctx, database)
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}

	log.Println("Starting HTTP server on port 8080...")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
