package schema

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func CreateTables(db *sql.DB) error {
	queryBytes, err := ioutil.ReadFile("db/schema/CreateTables.sql")
	if err != nil {
		log.Fatal("Error reading file:", err)
		return err
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

	fmt.Println("Tables created successfully.")
	return nil
}
