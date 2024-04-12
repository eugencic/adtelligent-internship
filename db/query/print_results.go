package queries

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

func ExecuteSourceCampaignQuery(database *sql.DB) error {
	query, err := readQueryFromFile("db/sql/ExtractFromTables.sql", 1)
	if err != nil {
		return err
	}

	rows, err := database.Query(query)
	if err != nil {
		return fmt.Errorf("error executing query for SourceCampaign: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	for rows.Next() {
		var sourceName string
		var campaignCount int
		if err := rows.Scan(&sourceName, &campaignCount); err != nil {
			return fmt.Errorf("error scanning row for SourceCampaign: %w", err)
		}
		fmt.Printf("Source: %s, Campaign Count: %d\n", sourceName, campaignCount)
	}

	return nil
}

func ExecuteCampaignWithoutSourceQuery(database *sql.DB) error {
	query, err := readQueryFromFile("db/sql/ExtractFromTables.sql", 2)
	if err != nil {
		return err
	}

	rows, err := database.Query(query)
	if err != nil {
		return fmt.Errorf("error executing query for CampaignWithoutSource: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	for rows.Next() {
		var campaignName string
		if err := rows.Scan(&campaignName); err != nil {
			return fmt.Errorf("error scanning row for CampaignWithoutSource: %w", err)
		}
		fmt.Printf("Campaign without Source: %s\n", campaignName)
	}

	return nil
}

func ExecuteUniqueNamesQuery(database *sql.DB) error {
	query, err := readQueryFromFile("db/sql/ExtractFromTables.sql", 3)
	if err != nil {
		return err
	}

	rows, err := database.Query(query)
	if err != nil {
		return fmt.Errorf("error executing query for UniqueNames: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	var uniqueNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("error scanning row for UniqueNames: %w", err)
		}
		uniqueNames = append(uniqueNames, name)
	}

	fmt.Println("Unique Names:")
	for _, name := range uniqueNames {
		fmt.Println(name)
	}

	return nil
}

func readQueryFromFile(filePath string, queryIndex int) (string, error) {
	queryBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	queries := strings.Split(string(queryBytes), ";")
	if queryIndex <= 0 || queryIndex > len(queries) {
		return "", fmt.Errorf("invalid query index: %d", queryIndex)
	}

	trimmedQuery := strings.TrimSpace(queries[queryIndex-1])
	if trimmedQuery == "" {
		return "", fmt.Errorf("empty query at index %d", queryIndex)
	}

	return trimmedQuery, nil
}
