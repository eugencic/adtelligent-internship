package db

import (
	"adtelligent-internship/db/config"
	"database/sql"
	"fmt"
)

func ConnectToDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DbConfig.User, config.DbConfig.Password, config.DbConfig.Host, config.DbConfig.Port, config.DbConfig.Database)

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
