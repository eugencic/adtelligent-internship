package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

type DBConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
}

var dbConfig = DBConfig{
	host:     "localhost",
	port:     "3306",
	user:     "user",
	password: "password",
	database: "database",
}

func ConnectToDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.user, dbConfig.password, dbConfig.host, dbConfig.port, dbConfig.database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := ConnectToDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	fmt.Println("Connected to the database successfully!")

	queryBytes, err := ioutil.ReadFile("CreateTables.sql")
	if err != nil {
		log.Fatal("Error reading queries file:", err)
	}
	queries := string(queryBytes)

	queryList := strings.Split(queries, ";")

	for _, query := range queryList {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue
		}

		_, err := db.Exec(trimmedQuery)
		if err != nil {
			log.Fatal("Error executing query:", err)
		}
	}

	fmt.Println("Queries executed successfully.")
	
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Source%d", i+1)
		_, err := db.Exec("INSERT INTO sources (name) VALUES (?)", name)
		if err != nil {
			log.Fatal("Error inserting source:", err)
		}
	}

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Campaign%d", i+1)
		_, err := db.Exec("INSERT INTO campaigns (name) VALUES (?)", name)
		if err != nil {
			log.Fatal("Error inserting campaign:", err)
		}
	}

	// Populate source_campaign links
	for _, sourceID := range rand.Perm(100)[:50] {
		for _, campaignID := range rand.Perm(100)[:rand.Intn(11)] {
			_, err := db.Exec("INSERT INTO source_campaign (source_id, campaign_id) VALUES (?, ?)", sourceID+1, campaignID+1)
			if err != nil {
				log.Fatal("Error inserting source_campaign link:", err)
			}
		}
	}

	fmt.Println("Data populated successfully.")

	fmt.Println("Sources:")
	rows, err := db.Query("SELECT id, name FROM sources")
	if err != nil {
		log.Fatal("Error querying sources:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning sources row:", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nCampaigns:")
	rows, err = db.Query("SELECT id, name FROM campaigns")
	if err != nil {
		log.Fatal("Error querying campaigns:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning campaigns row:", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSource_campaign:")
	rows, err = db.Query("SELECT source_id, campaign_id FROM source_campaign")
	if err != nil {
		log.Fatal("Error querying source_campaign:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sourceID, campaignID int
		err := rows.Scan(&sourceID, &campaignID)
		if err != nil {
			log.Fatal("Error scanning source_campaign row:", err)
		}
		fmt.Printf("Source ID: %d, Campaign ID: %d\n", sourceID, campaignID)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database successfully!")

	queryBytes2, err := ioutil.ReadFile("SelectFromTables.sql")
	if err != nil {
		log.Fatal("Error reading queries2 file:", err)
	}
	queries2 := string(queryBytes2)

	queryList2 := strings.Split(queries2, ";")

	for _, query := range queryList2 {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue
		}

		rows, err := db.Query(trimmedQuery)
		if err != nil {
			log.Fatal("Error executing query:", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			log.Fatal("Error retrieving column names:", err)
		}

		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(scanArgs...)
			if err != nil {
				log.Fatal("Error scanning row:", err)
			}

			var result []string
			for _, col := range values {
				switch v := col.(type) {
				case nil:
					result = append(result, "NULL")
				case int64:
					result = append(result, fmt.Sprintf("%d", v))
				case []byte:
					result = append(result, string(v))
				default:
					log.Fatalf("Unexpected data type: %T", v)
				}
			}
			fmt.Println(strings.Join(result, "\t"))
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
