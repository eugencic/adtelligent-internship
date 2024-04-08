package printer

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

func PrintResults(database *sql.DB) error {
	queryBytes, err := ioutil.ReadFile("db/printer/ExtractFromTables.sql")
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

		rows, err := database.Query(trimmedQuery)
		if err != nil {
			return fmt.Errorf("error executing querry: %w", err)
		}

		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("error retrieving column names: %w", err)
		}

		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(scanArgs...)
			if err != nil {
				return fmt.Errorf("error scanning row: %w", err)
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
					return fmt.Errorf("unexpected data type: %w", v)
				}
			}
			fmt.Println(strings.Join(result, "\t"))
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
	}

	fmt.Println("Results printed successfully.")
	return nil
}
