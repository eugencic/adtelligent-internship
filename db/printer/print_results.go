package printer

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func PrintResults(database *sql.DB) {
	queryBytes, err := ioutil.ReadFile("db/printer/ExtractFromTables.sql")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	queries := string(queryBytes)

	queryList := strings.Split(queries, ";")

	for _, query := range queryList {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue
		}

		rows, err := database.Query(trimmedQuery)
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

	fmt.Println("Results printed successfully.")
}
