package schema

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

func CreateTables(db *sql.DB) error {
	queryBytes, err := ioutil.ReadFile("db/schema/CreateTables.sql")
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
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
			return fmt.Errorf("error executing query: %w", err)
		}
	}

	fmt.Println("Tables created successfully.")
	return nil
}
