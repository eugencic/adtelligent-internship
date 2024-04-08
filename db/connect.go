package db

import (
	"database/sql"
	"fmt"
)

func ConnectToDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.user, dbConfig.password, dbConfig.host, dbConfig.port, dbConfig.database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// fmt.Println("Connected to the database successfully.")
	return db, nil
}
